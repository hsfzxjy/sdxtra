package db

import (
	"context"
	"database/sql"

	"github.com/hsfzxjy/sdxtra/internal/isotime"
	"github.com/jmoiron/sqlx"
)

type Session struct {
	ID         SessionID
	Name       sql.NullString
	CreatedAt  isotime.String
	InstanceID InstanceID
	instance   *Instance
}

func (s *Session) SessionID() SessionID { return s.ID }

type SessionID uint64

func (id SessionID) SessionID() SessionID { return id }

type SessionRef interface {
	SessionID() SessionID
}

type SessionTemplate struct {
	Name      sql.NullString
	Instance  InstanceRef
	CreatedAt sql.Null[isotime.String]
}

func (t SessionTemplate) Insert(ctx context.Context, db sqlx.ExtContext) (*Session, error) {
	createdAt := isotime.OrNow(t.CreatedAt)
	instanceID := t.Instance.InstanceID()
	const q = `
	INSERT INTO sessions (name, created_at, instance_id)
	VALUES (?, ?, ?)
	ON CONFLICT (name, instance_id) DO FAIL`
	res, err := db.ExecContext(ctx, q, &t.Name, &createdAt, &instanceID)
	if err != nil {
		return nil, err
	}
	instance, _ := t.Instance.(*Instance)
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Session{
		ID:         SessionID(id),
		Name:       t.Name,
		CreatedAt:  createdAt,
		InstanceID: instance.ID,
		instance:   instance,
	}, nil
}
