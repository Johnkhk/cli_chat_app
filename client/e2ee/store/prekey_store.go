package store

import (
	"context"
	"database/sql"

	"github.com/RTann/libsignal-go/protocol/prekey"
)

type PreKeyStore struct {
	db *sql.DB
}

// NewSessionStore creates a new SQLite-backed session store.
func NewPreKeyStore(db *sql.DB) prekey.Store {
	return &PreKeyStore{db: db}
}

// Load retrieves a session record for the given address.
func (s *PreKeyStore) Load(ctx context.Context, id prekey.ID) (*prekey.PreKey, bool, error) {
	return nil, false, nil
}

// Store saves a session record for the given address.
func (s *PreKeyStore) Store(ctx context.Context, id prekey.ID, preKey *prekey.PreKey) error {

	return nil
}

// Delete removes a pre-key entry identified by the given ID from the store.
func (s *PreKeyStore) Delete(ctx context.Context, id prekey.ID) error {
	return nil
}

type SignedPreKeyStore struct {
	db *sql.DB
}

// NewSignedPreKeyStore creates a new SQLite-backed signed pre-key store.
func NewSignedPreKeyStore(db *sql.DB) prekey.SignedStore {
	return &SignedPreKeyStore{db: db}
}

// Load retrieves a signed pre-key record for the given ID.
func (s *SignedPreKeyStore) Load(ctx context.Context, id prekey.ID) (*prekey.SignedPreKey, bool, error) {
	return nil, false, nil
}

// Store saves a signed pre-key record for the given ID.
func (s *SignedPreKeyStore) Store(ctx context.Context, id prekey.ID, record *prekey.SignedPreKey) error {
	return nil
}
