package store

import (
	"database/sql"

	"github.com/RTann/libsignal-go/protocol/identity"
	"github.com/RTann/libsignal-go/protocol/prekey"
	"github.com/RTann/libsignal-go/protocol/protocol"
	"github.com/RTann/libsignal-go/protocol/session"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore represents a complete implementation of the libsignal-go Store interface using SQLite.
type SQLiteStore struct {
	db                *sql.DB
	sessionStore      session.Store
	preKeyStore       prekey.Store
	signedPreKeyStore prekey.SignedStore
	identityStore     identity.Store
	groupStore        session.GroupStore
}

func (s *SQLiteStore) SessionStore() session.Store {
	return s.sessionStore
}
func (s *SQLiteStore) IdentityStore() identity.Store {
	return s.identityStore
}
func (s *SQLiteStore) PreKeyStore() prekey.Store {
	return s.preKeyStore
}
func (s *SQLiteStore) SignedPreKeyStore() prekey.SignedStore {
	return s.signedPreKeyStore
}
func (s *SQLiteStore) GroupStore() session.GroupStore {
	return s.groupStore
}

func NewSQLiteStore(keyPair identity.KeyPair, registrationID uint32, db *sql.DB) protocol.Store {
	// Make DB
	return &SQLiteStore{
		db:                db,
		sessionStore:      NewSessionStore(db),
		preKeyStore:       NewPreKeyStore(db),
		signedPreKeyStore: NewSignedPreKeyStore(db),
		identityStore:     NewIdentityStore(keyPair, registrationID, db),
		groupStore:        NewGroupStore(), // Not supported yet
	}
}
