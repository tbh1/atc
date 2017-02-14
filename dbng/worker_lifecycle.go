package dbng

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
)

//go:generate counterfeiter . WorkerLifecycle

type WorkerLifecycle interface {
	StallUnresponsiveWorkers() ([]string, error)
	LandFinishedLandingWorkers() ([]string, error)
	DeleteFinishedRetiringWorkers() ([]string, error)
}

type workerLifecycle struct {
	conn Conn
}

func (lifecycle *workerLifecycle) StallUnresponsiveWorkers() ([]*atc.Worker, error) {
	query, args, err := psql.Update("workers").
		SetMap(map[string]interface{}{
			"state":            string(WorkerStateStalled),
			"addr":             nil,
			"baggageclaim_url": nil,
			"expires":          nil,
		}).
		Where(sq.Eq{"state": string(WorkerStateRunning)}).
		Where(sq.Expr("expires < NOW()")).
		Suffix("RETURNING name, addr, baggageclaim_url, state").
		ToSql()
	if err != nil {
		return []*atc.Worker{}, err
	}

	rows, err := lifecycle.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workers := []*atc.Worker{}

	for rows.Next() {
		var (
			name     string
			addrStr  sql.NullString
			bcURLStr sql.NullString
			state    string
		)

		err = rows.Scan(&name, &addrStr, &bcURLStr, &state)
		if err != nil {
			return nil, err
		}

		var addr *string
		if addrStr.Valid {
			addr = &addrStr.String
		}

		var bcURL *string
		if bcURLStr.Valid {
			bcURL = &bcURLStr.String
		}

		workers = append(workers, &atc.Worker{
			Name:            name,
			GardenAddr:      addr,
			BaggageclaimURL: bcURL,
			State:           WorkerState(state),
		})
	}

	return workers, nil
}

func (lifecycle *workerLifecycle) DeleteFinishedRetiringWorkers() error {
	tx, err := lifecycle.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Squirrel does not have default support for subqueries in where clauses.
	// We hacked together a way to do it
	//
	// First we generate the subquery's SQL and args using
	// sq.Select instead of psql.Select so that we get
	// unordered placeholders instead of psql's ordered placeholders
	subQ, subQArgs, err := sq.Select("w.name").
		Distinct().
		From("builds b").
		Join("containers c ON b.id = c.build_id").
		Join("workers w ON w.name = c.worker_name").
		LeftJoin("jobs j ON j.id = b.job_id").
		Where(sq.Or{
			sq.Eq{
				"b.status": string(BuildStatusStarted),
			},
			sq.Eq{
				"b.status": string(BuildStatusPending),
			},
		}).
		Where(sq.Or{
			sq.Eq{
				"j.interruptible": false,
			},
			sq.Eq{
				"b.job_id": nil,
			},
		}).ToSql()

	if err != nil {
		return err
	}

	// Then we inject the subquery sql directly into
	// the where clause, and "add" the args from the
	// first query to the second query's args
	//
	// We use sq.Delete instead of psql.Delete for the same reason
	// but then change the placeholders using .PlaceholderFormat(sq.Dollar)
	// to go back to postgres's format
	_, err = sq.Delete("workers").
		Where(sq.Eq{
			"state": string(WorkerStateRetiring),
		}).
		Where("name NOT IN ("+subQ+")", subQArgs...).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (lifecycle *workerLifecycle) LandFinishedLandingWorkers() error {
	tx, err := lifecycle.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	subQ, subQArgs, err := sq.Select("w.name").
		Distinct().
		From("builds b").
		Join("containers c ON b.id = c.build_id").
		Join("workers w ON w.name = c.worker_name").
		LeftJoin("jobs j ON j.id = b.job_id").
		Where(sq.Or{
			sq.Eq{
				"b.status": string(BuildStatusStarted),
			},
			sq.Eq{
				"b.status": string(BuildStatusPending),
			},
		}).
		Where(sq.Or{
			sq.Eq{
				"j.interruptible": false,
			},
			sq.Eq{
				"b.job_id": nil,
			},
		}).ToSql()

	if err != nil {
		return err
	}

	_, err = sq.Update("workers").
		Set("state", string(WorkerStateLanded)).
		Set("addr", nil).
		Set("baggageclaim_url", nil).
		Where(sq.Eq{
			"state": string(WorkerStateLanding),
		}).
		Where("name NOT IN ("+subQ+")", subQArgs...).
		PlaceholderFormat(sq.Dollar).
		RunWith(tx).
		Exec()

	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
