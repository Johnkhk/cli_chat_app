package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Johnkhk/libsignal-go/protocol/address"
	"github.com/Johnkhk/libsignal-go/protocol/session"
	"google.golang.org/protobuf/proto"
)

var _ session.Store = (*SessionStore)(nil)

type SessionStore struct {
	db *sql.DB
}

// NewSessionStore creates a new SQLite-backed session store.
func NewSessionStore(db *sql.DB) session.Store {
	// Ensure that the sessions table exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS sessions (
		id TEXT PRIMARY KEY,
		record BLOB
	);
	`
	_, err := db.Exec(createTableQuery)
	if err != nil {
		panic(fmt.Sprintf("failed to create sessions table: %v", err))
	}

	return &SessionStore{db: db}
}

// Load retrieves a session record for the given address.
func (s *SessionStore) Load(ctx context.Context, addr address.Address) (*session.Record, bool, error) {
	var recordData []byte
	query := "SELECT record FROM sessions WHERE address = ? AND device_id = ?"
	err := s.db.QueryRowContext(ctx, query, addr.Name, addr.DeviceID).Scan(&recordData)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil // No record found
		}
		return nil, false, fmt.Errorf("failed to load session record: %w", err)
	}

	// Unmarshal the session record
	record, err := session.NewRecordBytes(recordData)
	if err != nil {
		return nil, false, fmt.Errorf("failed to load new record bytes: %w", err)
	}

	return record, true, nil
}

// Store saves a session record for the given address.
func (s *SessionStore) Store(ctx context.Context, addr address.Address, record *session.Record) error {
	// Marshal the session state to bytes
	v1SessionState := record.State().GetState()
	recordData, err := proto.Marshal(v1SessionState)
	if err != nil {
		return fmt.Errorf("failed to marshal session record: %w", err)
	}

	// Insert or replace the session record in the SQLite database
	query := "INSERT OR REPLACE INTO sessions (address, device_id, record) VALUES (?, ?, ?)"
	_, err = s.db.ExecContext(ctx, query, addr.Name, addr.DeviceID, recordData)
	if err != nil {
		return fmt.Errorf("failed to store session record: %w", err)
	}

	return nil
}
