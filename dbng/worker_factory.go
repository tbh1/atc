package dbng

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
)

//go:generate counterfeiter . WorkerFactory

type WorkerFactory interface {
	GetWorker(name string) (Worker, bool, error)
	SaveWorker(worker atc.Worker, ttl time.Duration) (Worker, error)
	HeartbeatWorker(worker atc.Worker, ttl time.Duration) (Worker, error)
	Workers() ([]*Worker, error)

	// move to *Team
	WorkersForTeam(teamName string) ([]*atc.Worker, error)
}

type workerFactory struct {
	conn Conn
}

func NewWorkerFactory(conn Conn) WorkerFactory {
	return &workerFactory{
		conn: conn,
	}
}

var workersQuery = psql.Select(`
		w.name,
		w.addr,
		w.state,
		w.baggageclaim_url,
		w.http_proxy_url,
		w.https_proxy_url,
		w.no_proxy,
		w.active_containers,
		w.resource_types,
		w.platform,
		w.tags,
		w.start_time,
		t.name,
		EXTRACT(epoch FROM w.expires - NOW())
	`).
	From("workers w").
	LeftJoin("teams t ON w.team_id = t.id")

func (f *workerFactory) GetWorker(name string) (Worker, bool, error) {
	row := workersQuery.Where(sq.Eq{"name": name}).
		RunWith(f.conn).
		QueryRow()

	model, err := scanWorker(row)
	if err != nil {
		return nil, false, err
	}

	return &worker{
		name:        name,
		cachedModel: &model,
		conn:        f.conn,
	}, true, nil
}

func (f *workerFactory) Workers() ([]Worker, error) {
	return f.getWorkers(workersQuery)
}

func (f *workerFactory) WorkersForTeam(teamName string) ([]Worker, error) {
	return f.getWorkers(workersQuery.Where(sq.Or{
		sq.Eq{"t.name": teamName},
		sq.Eq{"w.team_id": nil},
	}))
}

func getWorkers(conn Conn, query sq.SelectBuilder) ([]Worker, error) {
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	workers := []Worker{}

	for rows.Next() {
		model, err := scanWorker(rows)
		if err != nil {
			return nil, err
		}

		workers = append(workers, &worker{
			name:        model.Name,
			cachedModel: &model,
			conn:        conn,
		})
	}

	return workers, nil
}

func scanWorker(row scannable) (*atc.Worker, error) {
	var (
		name          string
		addStr        sql.NullString
		state         string
		bcURLStr      sql.NullString
		httpProxyURL  sql.NullString
		httpsProxyURL sql.NullString
		noProxy       sql.NullString

		activeContainers int
		resourceTypes    []byte
		platform         sql.NullString
		tags             []byte
		teamName         sql.NullString
		startTime        int64

		expiresIn *float64
	)

	err := row.Scan(
		&name,
		&addStr,
		&state,
		&bcURLStr,
		&httpProxyURL,
		&httpsProxyURL,
		&noProxy,
		&activeContainers,
		&resourceTypes,
		&platform,
		&tags,
		&teamName,
		&startTime,
		&expiresIn,
	)
	if err != nil {
		return nil, err
	}

	var addr *string
	if addStr.Valid {
		addr = &addStr.String
	}

	var bcURL *string
	if bcURLStr.Valid {
		bcURL = &bcURLStr.String
	}

	worker := atc.Worker{
		Name:            name,
		GardenAddr:      addr,
		BaggageclaimURL: bcURL,
		State:           WorkerState(state),

		ActiveContainers: activeContainers,
		StartTime:        startTime,
	}

	if expiresIn != nil {
		worker.ExpiresIn = time.Duration(*expiresIn) * time.Second
	}

	if httpProxyURL.Valid {
		worker.HTTPProxyURL = httpProxyURL.String
	}

	if httpsProxyURL.Valid {
		worker.HTTPSProxyURL = httpsProxyURL.String
	}

	if noProxy.Valid {
		worker.NoProxy = noProxy.String
	}

	if teamName.Valid {
		worker.TeamName = teamName.String
	}

	if platform.Valid {
		worker.Platform = platform.String
	}

	err = json.Unmarshal(resourceTypes, &worker.ResourceTypes)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(tags, &worker.Tags)
	if err != nil {
		return nil, err
	}
	return &worker, nil
}

func (f *workerFactory) HeartbeatWorker(worker atc.Worker, ttl time.Duration) (Worker, error) {
	// In order to be able to calculate the ttl that we return to the caller
	// we must compare time.Now() to the worker.expires column
	// However, workers.expires column is a "timestamp (without timezone)"
	// So we format time.Now() without any timezone information and then
	// parse that using the same layout to strip the timezone information
	layout := "Jan 2, 2006 15:04:05"
	nowStr := time.Now().Format(layout)
	now, err := time.Parse(layout, nowStr)
	if err != nil {
		return nil, err
	}

	tx, err := f.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	expires := "NULL"
	if ttl != 0 {
		expires = fmt.Sprintf(`NOW() + '%d second'::INTERVAL`, int(ttl.Seconds()))
	}

	cSql, _, err := sq.Case("state").
		When("'landing'::worker_state", "'landing'::worker_state").
		When("'landed'::worker_state", "'landed'::worker_state").
		When("'retiring'::worker_state", "'retiring'::worker_state").
		Else("'running'::worker_state").
		ToSql()

	if err != nil {
		return nil, err
	}

	addrSql, _, err := sq.Case("state").
		When("'landed'::worker_state", "NULL").
		Else("'" + worker.GardenAddr + "'").
		ToSql()
	if err != nil {
		return nil, err
	}

	bcSql, _, err := sq.Case("state").
		When("'landed'::worker_state", "NULL").
		Else("'" + worker.BaggageclaimURL + "'").
		ToSql()
	if err != nil {
		return nil, err
	}

	var (
		workerName       string
		workerStateStr   string
		activeContainers int
		expiresAt        time.Time
		addrStr          sql.NullString
		bcURLStr         sql.NullString
	)

	err = psql.Update("workers").
		Set("expires", sq.Expr(expires)).
		Set("addr", sq.Expr("("+addrSql+")")).
		Set("baggageclaim_url", sq.Expr("("+bcSql+")")).
		Set("active_containers", worker.ActiveContainers).
		Set("state", sq.Expr("("+cSql+")")).
		Where(sq.Eq{"name": worker.Name}).
		RunWith(tx).
		QueryRow().
		Scan(&workerName, &addrStr, &bcURLStr, &workerStateStr, &expiresAt, &activeContainers)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrWorkerNotPresent
		}
		return nil, err
	}

	row := workersQuery.Where(sq.Eq{"name": worker.Name}).
		RunWith(tx).
		QueryRow()

	model, err := scanWorker(rows)
	if err != nil {
		return nil, false, err
	}

	return &worker{
		name:        name,
		cachedModel: &model,
		conn:        f.conn,
	}, true, nil

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

}

func (f *workerFactory) SaveWorker(worker atc.Worker, ttl time.Duration) (Worker, error) {
	tx, err := f.conn.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	savedWorker, err := saveWorker(tx, worker, nil, ttl)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &worker{
		name:        worker.Name,
		cachedModel: &savedWorker,
		conn:        f.conn,
	}, nil
}

func saveWorker(tx Tx, worker atc.Worker, teamID *int, ttl time.Duration) (atc.Worker, error) {
	resourceTypes, err := json.Marshal(worker.ResourceTypes)
	if err != nil {
		return nil, err
	}

	tags, err := json.Marshal(worker.Tags)
	if err != nil {
		return nil, err
	}

	expires := "NULL"
	if ttl != 0 {
		expires = fmt.Sprintf(`NOW() + '%d second'::INTERVAL`, int(ttl.Seconds()))
	}

	var oldTeamID sql.NullInt64

	err = psql.Select("team_id").From("workers").Where(sq.Eq{
		"name": worker.Name,
	}).RunWith(tx).QueryRow().Scan(&oldTeamID)

	var workerState WorkerState
	if worker.State != "" {
		workerState = WorkerState(worker.State)
	} else {
		workerState = WorkerStateRunning
	}

	if err != nil {
		if err == sql.ErrNoRows {
			_, err = psql.Insert("workers").
				Columns(
					"addr",
					"expires",
					"active_containers",
					"resource_types",
					"tags",
					"platform",
					"baggageclaim_url",
					"http_proxy_url",
					"https_proxy_url",
					"no_proxy",
					"name",
					"start_time",
					"team_id",
					"state",
				).
				Values(
					worker.GardenAddr,
					sq.Expr(expires),
					worker.ActiveContainers,
					resourceTypes,
					tags,
					worker.Platform,
					worker.BaggageclaimURL,
					worker.HTTPProxyURL,
					worker.HTTPSProxyURL,
					worker.NoProxy,
					worker.Name,
					worker.StartTime,
					teamID,
					string(workerState),
				).
				RunWith(tx).
				Exec()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		if (oldTeamID.Valid == (teamID == nil)) ||
			(oldTeamID.Valid && (*teamID != int(oldTeamID.Int64))) {
			return nil, errors.New("update-of-other-teams-worker-not-allowed")
		}

		_, err = psql.Update("workers").
			Set("addr", worker.GardenAddr).
			Set("expires", sq.Expr(expires)).
			Set("active_containers", worker.ActiveContainers).
			Set("resource_types", resourceTypes).
			Set("tags", tags).
			Set("platform", worker.Platform).
			Set("baggageclaim_url", worker.BaggageclaimURL).
			Set("http_proxy_url", worker.HTTPProxyURL).
			Set("https_proxy_url", worker.HTTPSProxyURL).
			Set("no_proxy", worker.NoProxy).
			Set("name", worker.Name).
			Set("start_time", worker.StartTime).
			Set("state", string(workerState)).
			Where(sq.Eq{
				"name": worker.Name,
			}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, err
		}
	}

	savedWorker := atc.Worker{
		Name:       worker.Name,
		GardenAddr: &worker.GardenAddr,
		State:      workerState,
	}

	workerBaseResourceTypeIDs := []int{}
	for _, resourceType := range worker.ResourceTypes {
		workerResourceType := WorkerResourceType{
			Worker:  savedWorker,
			Image:   resourceType.Image,
			Version: resourceType.Version,
			BaseResourceType: &BaseResourceType{
				Name: resourceType.Type,
			},
		}

		brt := BaseResourceType{
			Name: resourceType.Type,
		}

		ubrt, err := brt.FindOrCreate(tx)
		if err != nil {
			return nil, err
		}

		_, err = psql.Delete("worker_base_resource_types").
			Where(sq.Eq{
				"worker_name":           worker.Name,
				"base_resource_type_id": ubrt.ID,
			}).
			Where(sq.NotEq{
				"version": resourceType.Version,
			}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, err
		}
		uwrt, err := workerResourceType.FindOrCreate(tx)
		if err != nil {
			return nil, err
		}

		workerBaseResourceTypeIDs = append(workerBaseResourceTypeIDs, uwrt.ID)
	}

	_, err = psql.Delete("worker_base_resource_types").
		Where(sq.Eq{
			"worker_name": worker.Name,
		}).
		Where(sq.NotEq{
			"id": workerBaseResourceTypeIDs,
		}).
		RunWith(tx).
		Exec()
	if err != nil {
		return nil, err
	}

	return savedWorker, nil
}
