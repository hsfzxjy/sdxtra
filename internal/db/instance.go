package db

import (
	"context"
	"database/sql"

	"github.com/hsfzxjy/sdxtra/internal/isotime"
	"github.com/jmoiron/sqlx"
)

type InstanceParam struct {
	ID     InstanceParamID
	Sha256 Sha256Digest
	Params []byte
}

func (p *InstanceParam) InstanceParamID() InstanceParamID {
	return p.ID
}

type InstanceParamTemplate struct {
	Sha256 Sha256Digest
	Params []byte
}

func (t InstanceParamTemplate) Upsert(ctx context.Context, db sqlx.ExtContext) (*InstanceParam, error) {
	const q = `
	INSERT INTO instance_params (sha256, params)
	VALUES (?, ?)
	ON CONFLICT (sha256) DO NOTHING
	RETURNING id`
	row := db.QueryRowxContext(ctx, q, &t.Sha256, &t.Params)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var id InstanceParamID
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}
	return &InstanceParam{
		ID:     id,
		Sha256: t.Sha256,
		Params: t.Params,
	}, nil
}

type InstanceParamID uint64

func (id InstanceParamID) InstanceParamID() InstanceParamID {
	return id
}

type InstanceParamRef interface {
	InstanceParamID() InstanceParamID
}

type Instance struct {
	ID        InstanceID
	Name      sql.NullString
	CreatedAt isotime.String
	ParamID   InstanceParamID
	param     *InstanceParam
}

func (i *Instance) InstanceID() InstanceID {
	return i.ID
}
func (i *Instance) LogOwnerInfo() (LogKind, uint64) {
	return i.ID.LogKindAssocID()
}

type InstanceTemplate struct {
	Name      sql.NullString
	Param     InstanceParamRef
	CreatedAt sql.Null[isotime.String]
}

func (t InstanceTemplate) Insert(ctx context.Context, db sqlx.ExtContext) (*Instance, error) {
	createdAt := isotime.OrNow(t.CreatedAt)
	paramID := t.Param.InstanceParamID()
	const q = `
	INSERT INTO instances (name, created_at, param_id)
	VALUES (?, ?, ?)`
	res, err := db.ExecContext(ctx, q, &t.Name, &createdAt, &paramID)
	if err != nil {
		return nil, err
	}
	param, _ := t.Param.(*InstanceParam)
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &Instance{
		ID:        InstanceID(id),
		Name:      t.Name,
		CreatedAt: createdAt,
		ParamID:   param.ID,
		param:     param,
	}, nil
}

type InstanceID uint64

func (id InstanceID) InstanceID() InstanceID {
	return id
}
func (id InstanceID) LogKindAssocID() (LogKind, uint64) {
	return LOGKIND_INSTANCE, uint64(id)
}

func (id InstanceID) Rename(ctx context.Context, db sqlx.ExtContext, newName sql.NullString) (success bool, err error) {
	const q = `
	UPDATE instances
	SET name = ?
	WHERE id = ?`
	res, err := db.ExecContext(ctx, q, &newName, &id)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

type InstanceRef interface {
	InstanceID() InstanceID
	LogOwner
}
