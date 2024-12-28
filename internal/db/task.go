package db

import (
	"context"
	"database/sql"

	"github.com/hsfzxjy/sdxtra/internal/isotime"
	"github.com/jmoiron/sqlx"
)

type TaskParam struct {
	ID     TaskParamID
	Sha256 Sha256Digest
	Param  []byte
}

func (p *TaskParam) TaskParamID() TaskParamID { return p.ID }

type TaskParamID uint64

func (id TaskParamID) TaskParamID() TaskParamID { return id }

type TaskParamRef interface {
	TaskParamID() TaskParamID
}

type TaskParamTemplate struct {
	Sha256 Sha256Digest
	Param  []byte
}

func (t TaskParamTemplate) Upsert(ctx context.Context, db sqlx.ExtContext) (*TaskParam, error) {
	const q = `
	INSERT INTO task_params (sha256, param)
	VALUES (?, ?)
	ON CONFLICT (sha256) DO NOTHING
	RETURNING id`
	row := db.QueryRowxContext(ctx, q, &t.Sha256, &t.Param)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var id TaskParamID
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}
	return &TaskParam{
		ID:     id,
		Sha256: t.Sha256,
		Param:  t.Param,
	}, nil
}

type Task struct {
	ID        TaskID
	Name      sql.NullString
	CreatedAt isotime.String

	ParamID TaskParamID
	param   *TaskParam

	Result []byte

	SessionID SessionID
	session   *Session
}

func (t *Task) TaskID() TaskID                    { return t.ID }
func (t *Task) LogOwnerInfo() (LogKind, uint64) { return (t.ID).LogKindAssocID() }

type TaskID uint64

func (id TaskID) TaskID() TaskID                    { return id }
func (id TaskID) LogKindAssocID() (LogKind, uint64) { return LOGKIND_TASK, uint64(id) }

type TaskTemplate struct {
	Name      sql.NullString
	Param     TaskParamRef
	CreatedAt sql.Null[isotime.String]
	Result    []byte
}

func (t TaskTemplate) Insert(ctx context.Context, db sqlx.ExtContext, session SessionRef) (*Task, error) {
	createdAt := isotime.OrNow(t.CreatedAt)
	paramID := t.Param.TaskParamID()
	sessionID := session.SessionID()
	const q = `
	INSERT INTO tasks (name, created_at, param_id, session_id, result)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT (name, session_id) DO NOTHING
	RETURNING id`
	res, err := db.ExecContext(ctx, q, &t.Name, &createdAt, &paramID, &sessionID, &t.Result)
	if err != nil {
		return nil, err
	}
	sessionPtr, _ := session.(*Session)
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	paramPtr := t.Param.(*TaskParam)
	return &Task{
		ID:        TaskID(id),
		Name:      t.Name,
		CreatedAt: createdAt,
		ParamID:   paramID,
		param:     paramPtr,
		Result:    t.Result,
		SessionID: sessionID,
		session:   sessionPtr,
	}, nil
}

type TaskRef interface {
	TaskID() TaskID
	LogOwner
}
