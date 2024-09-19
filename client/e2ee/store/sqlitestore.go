package store

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"

	"github.com/Johnkhk/libsignal-go/protocol/identity"
	"github.com/Johnkhk/libsignal-go/protocol/prekey"
	"github.com/Johnkhk/libsignal-go/protocol/protocol"
	"github.com/Johnkhk/libsignal-go/protocol/session"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore represents a complete implementation of the libsignal-go Store interface using SQLite.
var _ protocol.Store = (*SQLiteStore)(nil)

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

func NewSQLiteStore(dbPath string, userID uint32, isSignup bool) (protocol.Store, error) {
	// Initialize SQLite database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Create tables if they do not exist
	err = CreateTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	// Generate registration ID (using default device ID = 1 for now)
	deviceID := uint32(1)
	registrationID := generateRegistrationID(userID, deviceID)

	// If this is a signup, create and store the local identity
	if isSignup {
		err := createLocalIdentity(db, registrationID)
		if err != nil {
			return nil, fmt.Errorf("failed to create local identity: %v", err)
		}
	}

	// Load the local identity for subsequent app runs
	keyPair, existingRegistrationID, err := loadLocalIdentity(db)
	if err != nil {
		return nil, fmt.Errorf("failed to load local identity: %v", err)
	}

	// Initialize store components
	return &SQLiteStore{
		db:                db,
		sessionStore:      NewSessionStore(db),
		preKeyStore:       NewPreKeyStore(db),
		signedPreKeyStore: NewSignedPreKeyStore(db),
		identityStore:     NewIdentityStore(keyPair, existingRegistrationID, db),
		groupStore:        NewGroupStore(), // Group store is not yet supported
	}, nil
}

// CreateTables ensures that the necessary tables for the Signal Protocol stores exist.
func CreateTables(db *sql.DB) error {
	// Create session table
	sessionTable := `
	CREATE TABLE IF NOT EXISTS sessions (
    address TEXT NOT NULL,
    device_id INTEGER NOT NULL,
    record BLOB NOT NULL,
    PRIMARY KEY (address, device_id)
	);
	`

	// Create pre-key table
	preKeyTable := `
	CREATE TABLE IF NOT EXISTS prekeys (
		id INTEGER PRIMARY KEY,
		record BLOB NOT NULL
	);`

	// Create signed pre-key table
	signedPreKeyTable := `
	CREATE TABLE IF NOT EXISTS signed_prekeys (
		id INTEGER PRIMARY KEY,
		record BLOB NOT NULL
	);`

	// Create identity table
	identityTable := `
	CREATE TABLE IF NOT EXISTS identities (
		address TEXT PRIMARY KEY,
		key_data BLOB NOT NULL,
		trust_level INTEGER NOT NULL
	);`

	localIdentityTable := `
	CREATE TABLE IF NOT EXISTS local_identity (
		key_pair BLOB NOT NULL,
		registration_id INTEGER NOT NULL
	);
	`

	// Create table queries in a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Execute table creation queries
	_, err = tx.Exec(sessionTable)
	if err != nil {
		return fmt.Errorf("failed to create sessions table: %v", err)
	}
	_, err = tx.Exec(preKeyTable)
	if err != nil {
		return fmt.Errorf("failed to create prekeys table: %v", err)
	}
	_, err = tx.Exec(signedPreKeyTable)
	if err != nil {
		return fmt.Errorf("failed to create signed_prekeys table: %v", err)
	}
	_, err = tx.Exec(identityTable)
	if err != nil {
		return fmt.Errorf("failed to create identities table: %v", err)
	}

	_, err = tx.Exec(localIdentityTable)
	if err != nil {
		return fmt.Errorf("failed to create local identity table table: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Println("All necessary tables created successfully.")
	return nil
}

func createLocalIdentity(db *sql.DB, registrationID uint32) error {
	// Generate a new key pair
	keyPair, err := identity.GenerateKeyPair(rand.Reader) // Assume GenerateKeyPair exists
	if err != nil {
		return fmt.Errorf("failed to generate key pair: %v", err)
	}

	// Serialize the key pair
	keyPairData, err := serializeKeyPairAndEncode(keyPair)

	if err != nil {
		return fmt.Errorf("failed to encode key pair: %v", err)
	}

	// Insert the local identity into the database
	insertQuery := "INSERT INTO local_identity (key_pair, registration_id) VALUES (?, ?)"
	_, err = db.Exec(insertQuery, keyPairData, registrationID)
	if err != nil {
		return fmt.Errorf("failed to store local identity: %v", err)
	}

	return nil
}

func loadLocalIdentity(db *sql.DB) (identity.KeyPair, uint32, error) {
	var keyPairData []byte
	var registrationID uint32

	// Load key pair and registration ID from the database
	query := "SELECT key_pair, registration_id FROM local_identity LIMIT 1"
	err := db.QueryRow(query).Scan(&keyPairData, &registrationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return identity.KeyPair{}, 0, fmt.Errorf("local identity not found")
		}
		return identity.KeyPair{}, 0, fmt.Errorf("failed to load local identity: %v", err)
	}

	// Decode the key pair
	var keyPair identity.KeyPair
	err = decodeAndDeserializeKeyPair(keyPairData, &keyPair)
	if err != nil {
		return identity.KeyPair{}, 0, fmt.Errorf("failed to decode key pair: %v", err)
	}

	return keyPair, registrationID, nil
}
