package store

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Johnkhk/libsignal-go/protocol/curve"
	"github.com/Johnkhk/libsignal-go/protocol/identity"
	"github.com/Johnkhk/libsignal-go/protocol/prekey"
	"github.com/Johnkhk/libsignal-go/protocol/protocol"
	"github.com/Johnkhk/libsignal-go/protocol/session"
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore represents a complete implementation of the libsignal-go Store interface using SQLite.
var _ protocol.Store = (*SQLiteStore)(nil)

type SQLiteStore struct {
	DB                *sql.DB
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

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
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

	// Initialize store components
	return &SQLiteStore{
		DB:                db,
		sessionStore:      NewSessionStore(db),
		preKeyStore:       NewPreKeyStore(db),
		signedPreKeyStore: NewSignedPreKeyStore(db),
		identityStore:     NewIdentityStore(db),
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
		registration_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		device_id INTEGER NOT NULL
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

type LocalIdentity struct {
	IdentityPublicKey     []byte
	PreKeyID              uint32
	PreKeyPublicKey       []byte
	SignedPreKeyID        uint32
	SignedPreKeyPublicKey []byte
	Signature             []byte
}

// func (s *SQLiteStore) CreateLocalIdentity(registrationID uint32) (*LocalIdentity, error) {
// 	db := s.DB
// 	ctx := context.Background()

// 	// 1. Generate a new identity key pair
// 	identityKeyPair, err := identity.GenerateKeyPair(rand.Reader)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate identity key pair: %v", err)
// 	}

// 	// Serialize the key pair for local storage
// 	keyPairData, err := SerializeKeyPairAndEncode(identityKeyPair)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to encode identity key pair: %v", err)
// 	}

// 	// Insert the local identity into the database
// 	insertQuery := "INSERT INTO local_identity (key_pair, registration_id) VALUES (?, ?)"
// 	_, err = db.Exec(insertQuery, keyPairData, registrationID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to store local identity: %v", err)
// 	}

// 	// Schema to generate unique IDs (incremental IDs)
// 	preKeyID := uint32(1)       // Example ID schema
// 	signedPreKeyID := uint32(1) // Example ID schema for signed prekey

// 	// 2. Generate the PreKey
// 	preKeyPair, err := curve.GenerateKeyPair(rand.Reader)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate pre-key pair: %v", err)
// 	}

// 	// Store PreKey locally (private part)
// 	preKey := prekey.NewPreKey(prekey.ID(preKeyID), preKeyPair)
// 	err = s.preKeyStore.Store(ctx, prekey.ID(preKeyID), preKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to store pre-key: %v", err)
// 	}

// 	// 3. Generate the Signed PreKey
// 	signedPreKeyPair, err := curve.GenerateKeyPair(rand.Reader)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to generate signed pre-key pair: %v", err)
// 	}

// 	// Store Signed PreKey locally (private part)
// 	signedPreKey := prekey.NewSigned(prekey.ID(signedPreKeyID), uint64(time.Now().Unix()), signedPreKeyPair, nil)
// 	err = s.signedPreKeyStore.Store(ctx, prekey.ID(signedPreKeyID), signedPreKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to store signed pre-key: %v", err)
// 	}

// 	// 4. Sign the public part of the Signed PreKey with the private Identity Key
// 	signature, err := identityKeyPair.PrivateKey().Sign(rand.Reader, signedPreKeyPair.PublicKey().Bytes())
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to sign signed pre-key: %v", err)
// 	}

// 	// Return a struct with all the relevant fields
// 	return &LocalIdentity{
// 		IdentityPublicKey:     identityKeyPair.PublicKey().Bytes(),
// 		PreKeyID:              preKeyID,
// 		PreKeyPublicKey:       preKeyPair.PublicKey().Bytes(),
// 		SignedPreKeyID:        signedPreKeyID,
// 		SignedPreKeyPublicKey: signedPreKeyPair.PublicKey().Bytes(),
// 		Signature:             signature,
// 	}, nil
// }

func (s *SQLiteStore) CreateLocalIdentity(registrationID uint32) (*LocalIdentity, error) {
	db := s.DB
	ctx := context.Background()

	// 1. Generate a new identity key pair
	identityKeyPair, err := identity.GenerateKeyPair(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate identity key pair: %v", err)
	}

	// Serialize the key pair for local storage
	keyPairData, err := SerializeKeyPairAndEncode(identityKeyPair)
	if err != nil {
		return nil, fmt.Errorf("failed to encode identity key pair: %v", err)
	}

	// Insert the local identity into the database
	extractedUserID, deviceID := ExtractUserIDAndDeviceID(registrationID)
	insertQuery := "INSERT INTO local_identity (key_pair, registration_id, user_id, device_id) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(insertQuery, keyPairData, registrationID, extractedUserID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to store local identity: %v", err)
	}

	// Schema to generate unique IDs (incremental IDs)
	preKeyID := uint32(1)       // Example ID schema
	signedPreKeyID := uint32(1) // Example ID schema for signed prekey
	// Generate unique PreKeyID and SignedPreKeyID
	// preKeyID, err := generateRandomID()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to generate PreKey ID: %v", err)
	// }
	// signedPreKeyID, err := generateRandomID()
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to generate Signed PreKey ID: %v", err)
	// }

	// 2. Generate the PreKey
	preKeyPair, err := curve.GenerateKeyPair(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate pre-key pair: %v", err)
	}

	// Store PreKey locally (private part)
	preKey := prekey.NewPreKey(prekey.ID(preKeyID), preKeyPair)
	err = s.preKeyStore.Store(ctx, prekey.ID(preKeyID), preKey)
	if err != nil {
		return nil, fmt.Errorf("failed to store pre-key: %v", err)
	}

	// 3. Generate the Signed PreKey
	signedPreKeyPair, err := curve.GenerateKeyPair(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate signed pre-key pair: %v", err)
	}

	// 4. Sign the public part of the Signed PreKey with the private Identity Key
	signature, err := identityKeyPair.PrivateKey().Sign(rand.Reader, signedPreKeyPair.PublicKey().Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to sign signed pre-key: %v", err)
	}

	// 5. Store Signed PreKey locally (with the public key signature)
	signedPreKey := prekey.NewSigned(prekey.ID(signedPreKeyID), uint64(time.Now().Unix()), signedPreKeyPair, signature)
	err = s.signedPreKeyStore.Store(ctx, prekey.ID(signedPreKeyID), signedPreKey)
	if err != nil {
		return nil, fmt.Errorf("failed to store signed pre-key: %v", err)
	}

	// Return a struct with all the relevant fields
	return &LocalIdentity{
		IdentityPublicKey:     identityKeyPair.PublicKey().Bytes(),
		PreKeyID:              preKeyID,
		PreKeyPublicKey:       preKeyPair.PublicKey().Bytes(),
		SignedPreKeyID:        signedPreKeyID,
		SignedPreKeyPublicKey: signedPreKeyPair.PublicKey().Bytes(),
		Signature:             signature,
	}, nil
}

// LoadLocalIdentity loads the identity key pair and registration ID based on the current device's MAC address.
func LoadLocalIdentity(db *sql.DB, userID uint32) (identity.KeyPair, uint32, error) {
	// Get MAC address
	macAddress, err := GetMACAddress()
	if err != nil {
		return identity.KeyPair{}, 0, fmt.Errorf("failed to get MAC address: %v", err)
	}

	// Convert MAC address to uint32 for deviceID
	deviceID, err := MacToUint32(macAddress)
	if err != nil {
		return identity.KeyPair{}, 0, fmt.Errorf("failed to convert MAC address to uint32: %v", err)
	}

	// Generate registration ID using userID and deviceID
	registrationID := GenerateRegistrationID(userID, deviceID)

	// Query the database to load the key pair and registration ID
	var keyPairData []byte
	var storedRegistrationID uint32

	query := "SELECT key_pair, registration_id FROM local_identity WHERE registration_id = ? LIMIT 1"
	err = db.QueryRow(query, registrationID).Scan(&keyPairData, &storedRegistrationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return identity.KeyPair{}, 0, fmt.Errorf("local identity not found for registration ID: %d", registrationID)
		}
		return identity.KeyPair{}, 0, fmt.Errorf("failed to load local identity: %v", err)
	}

	// Decode the key pair
	var keyPair identity.KeyPair
	err = DecodeAndDeserializeKeyPair(keyPairData, &keyPair)
	if err != nil {
		return identity.KeyPair{}, 0, fmt.Errorf("failed to decode key pair: %v", err)
	}

	return keyPair, storedRegistrationID, nil
}

func GetMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		mac := i.HardwareAddr.String()
		if len(mac) > 0 {
			return mac, nil
		}
	}

	return "", fmt.Errorf("no network interfaces found")
}

// Convert a MAC address string into a uint32 with a limit to 10 bits
func MacToUint32(macAddr string) (uint32, error) {
	// Hash the MAC address to get a larger value
	hash := sha256.Sum256([]byte(macAddr))

	// Take the first 4 bytes of the hash to form a uint32, then limit it to 10 bits
	return binary.BigEndian.Uint32(hash[:4]) & 0x3FF, nil // Mask to get the lower 10 bits
}

////////////////////////////// Session Management //////////////////////////////////
