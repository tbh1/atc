package dbng

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
)

var (
	ErrWorkerNotPresent         = errors.New("worker-not-present-in-db")
	ErrCannotPruneRunningWorker = errors.New("worker-not-stalled-for-pruning")
)

type Worker interface {
	Model() atc.Worker
	Refresh() error // TODO: implement this

	Land() error
	Retire() error
	Prune() error
	Delete() error
}

type worker struct {
	name        string
	cachedModel *atc.Worker
	conn        Conn
}

func (worker *worker) Model() atc.Worker {
	return *worker.cachedModel
	// if worker.cachedModel != nil {
	// 	return *worker.cachedModel
	// }
	//
	// row := workersQuery.Where(sq.Eq{"w.name": worker.name}).
	// 	RunWith(worker.conn).
	// 	QueryRow()
	//
	// selectedWorker, err := scanWorker(row)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return atc.Worker{}
	// 	}
	// 	return atc.Worker{}
	// }
	//
	// return *selectedWorker
}

func (worker *worker) Land() error {
	tx, err := worker.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	cSql, _, err := sq.Case("state").
		When("'landed'::worker_state", "'landed'::worker_state").
		Else("'landing'::worker_state").
		ToSql()
	if err != nil {
		return err
	}

	_, err = psql.Update("workers").
		Set("state", sq.Expr("("+cSql+")")).
		Where(sq.Eq{"name": worker.name}).
		RunWith(tx).
		Exec()
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrWorkerNotPresent
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (worker *worker) Retire(name string) error {
	tx, err := worker.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = psql.Update("workers").
		SetMap(map[string]interface{}{
			"state": string(atc.WorkerStateRetiring),
		}).
		Where(sq.Eq{"name": name}).
		RunWith(tx).
		Exec()
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrWorkerNotPresent
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (worker *worker) Prune(name string) error {
	tx, err := worker.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := sq.Delete("workers").
		Where(sq.Eq{
			"name": name,
		}).
		Where(sq.NotEq{
			"state": string(atc.WorkerStateRunning),
		}).
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

	affected, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		atcWorker := worker.Model()

		if atcWorker == (atc.Worker{}) {
			return ErrWorkerNotPresent
		}

		return ErrCannotPruneRunningWorker
	}

	return nil
}

func (worker *worker) Delete(name string) error {
	tx, err := worker.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = sq.Delete("workers").
		Where(sq.Eq{
			"name": name,
		}).
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
