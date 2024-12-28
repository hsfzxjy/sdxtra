package db

import (
	"context"
	"errors"
	"iter"

	"github.com/jmoiron/sqlx"
)

type File struct {
	TaskID TaskID
	MetaID FileMetaID
	meta   *FileMeta
	task   *Task
}

type FileID uint64
type FileMetaID uint64

type FileMeta struct {
	ID   FileMetaID
	Path string `db:"path"`
	Meta []byte `db:"meta"`
}

func FileCreateForTask(ctx context.Context, db sqlx.ExtContext, task TaskRef, metas []FileMeta) ([]File, error) {
	const q = `
	INSERT INTO file_metas (path, meta)
	VALUES (:path, :meta)
	ON CONFLICT (path) DO UPDATE SET meta=excluded.meta
	RETURNING id`
	rows, err := sqlx.NamedQueryContext(ctx, db, q, metas)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for i := range metas {
		if !rows.Next() {
			return nil, errors.New("not enough rows")
		}
		meta := &metas[i]
		err := rows.Scan(&meta.ID)
		if err != nil {
			return nil, err
		}
	}
	task_id := task.TaskID()
	taskPtr, _ := task.(*Task)
	files := make([]File, len(metas))
	for i := range metas {
		meta := &metas[i]
		files[i] = File{
			TaskID: task_id,
			MetaID: meta.ID,
			meta:   meta,
			task:   taskPtr,
		}
	}
	const q2 = `
	INSERT INTO file_task (meta_id, task_id)
	VALUES (:meta_id, :task_id)`
	_, err = sqlx.NamedExecContext(ctx, db, q2, files)
	if err != nil {
		return nil, err
	}
	return files, nil
}

type FileDeleteCondition struct {
	Expr string
	Args []any
}

func FileDelete(ctx context.Context, db sqlx.ExtContext, cond FileDeleteCondition) error {
	var q = `
	DELETE FROM file_task
	WHERE (` + cond.Expr + `)`
	_, err := db.ExecContext(ctx, q, cond.Args...)
	return err
}

func FileCleanOrphans(ctx context.Context, db sqlx.ExtContext) (iter.Seq2[string, error], error) {
	rows, err := db.QueryContext(ctx, `
	DELETE FROM file_metas
	WHERE NOT EXISTS (SELECT 1 FROM file_task WHERE file_metas.id = file_task.meta_id)
	RETURNING path`)
	if err != nil {
		return nil, err
	}
	return func(yield func(string, error) bool) {
		defer rows.Close()
		var p string
		for rows.Next() {
			err := rows.Scan(&p)
			if !yield(p, err) {
				return
			}
			if err != nil {
				return
			}
		}
	}, nil
}
