package db

import (
	"context"
	"fmt"

	"github.com/hsfzxjy/sdxtra/internal/isotime"
	"github.com/hsfzxjy/sdxtra/internal/log"
	"github.com/jmoiron/sqlx"
)

type LogKind int

// NOTE: New kinds must be added to the end of the list.
const (
	LOGKIND_UNKNOWN  LogKind = 0
	LOGKIND_INSTANCE LogKind = 1
	LOGKIND_TASK     LogKind = 2
)

type LogOwner interface {
	LogOwnerInfo() (LogKind, uint64)
}

type LogEntry struct {
	AssocID   uint64
	Kind      LogKind
	Lvl       log.Level
	Message   string
	CreatedAt isotime.String
}

type ErrLogEntryOwnerTypeMismatch struct {
	Idx   int
	Entry *log.Entry
}

func (e ErrLogEntryOwnerTypeMismatch) Error() string {
	return fmt.Sprintf("log entry owner at index %d is not a LogOwner: %v", e.Idx, e.Entry.Owner)
}

func LogEntryBulkInsert(ctx context.Context, db sqlx.ExtContext, rawEntries []*log.Entry) error {
	entries := make([]LogEntry, 0, len(rawEntries))
	for i, e := range rawEntries {
		owner, ok := e.Owner.(LogOwner)
		if !ok {
			return ErrLogEntryOwnerTypeMismatch{i, e}
		}
		kind, assocID := owner.LogOwnerInfo()
		entries = append(entries, LogEntry{
			AssocID:   assocID,
			Kind:      kind,
			Lvl:       e.Level,
			Message:   e.Message,
			CreatedAt: isotime.Encode(e.Time),
		})
	}
	const q = `
	INSERT INTO logs (assoc_id, kind, lvl, message, created_at)
	VALUES (:assoc_id, :kind, :lvl, :message, :created_at)`
	_, err := sqlx.NamedExecContext(ctx, db, q, entries)
	return err
}
