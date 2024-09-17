package store

import (
	"context"
	"database/sql"

	"github.com/RTann/libsignal-go/protocol/address"
	"github.com/RTann/libsignal-go/protocol/session"
)

type SessionStore struct {
	db *sql.DB
}

// NewSessionStore creates a new SQLite-backed session store.
func NewSessionStore(db *sql.DB) session.Store {
	return &SessionStore{db: db}
}

// Load retrieves a session record for the given address.
func (s *SessionStore) Load(ctx context.Context, addr address.Address) (*session.Record, bool, error) {

	record := &session.Record{}
	return record, true, nil
}

// Store saves a session record for the given address.
func (s *SessionStore) Store(ctx context.Context, addr address.Address, record *session.Record) error {

	return nil
}
