package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Johnkhk/libsignal-go/protocol/curve"
	v1 "github.com/Johnkhk/libsignal-go/protocol/generated/v1"
	"github.com/Johnkhk/libsignal-go/protocol/prekey"
	"google.golang.org/protobuf/proto"
)

// PreKeyStore represents the SQLite-backed pre-key store.
// Stores user's own pre-keys for later use.
var _ prekey.Store = (*PreKeyStore)(nil)

type PreKeyStore struct {
	db *sql.DB
}

// NewPreKeyStore creates a new SQLite-backed pre-key store.
func NewPreKeyStore(db *sql.DB) prekey.Store {
	return &PreKeyStore{db: db}
}

// Load retrieves a pre-key record for the given ID.
func (s *PreKeyStore) Load(ctx context.Context, id prekey.ID) (*prekey.PreKey, bool, error) {
	// Prepare the SQL query
	query := `SELECT record FROM prekeys WHERE id = ?`
	var blob []byte

	// Execute the query to get the pre-key record
	fmt.Printf("Loading pre-key with ID: %v\n", id)
	err := s.db.QueryRowContext(ctx, query, id).Scan(&blob)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("failed to load pre-key: %v", err)
	}

	// Deserialize the pre-key key pair from the blob
	var keyPair curve.KeyPair
	err = decodeAndDeserializeKeyPair(blob, &keyPair)
	if err != nil {
		return nil, false, fmt.Errorf("failed to decode pre-key key pair: %v", err)
	}

	// Recreate the PreKey using the ID and the deserialized KeyPair
	preKey := prekey.NewPreKey(id, &keyPair)

	return preKey, true, nil
}

// Store saves a pre-key record for the given ID.
func (s *PreKeyStore) Store(ctx context.Context, id prekey.ID, preKey *prekey.PreKey) error {
	// Serialize the pre-key record (KeyPair)
	preKeyPair, err := preKey.KeyPair()
	if err != nil {
		return fmt.Errorf("failed to get key pair from pre-key: %v", err)
	}

	// Serialize the KeyPair
	blob, err := serializeKeyPairAndEncode(preKeyPair)
	if err != nil {
		return fmt.Errorf("failed to encode pre-key key pair: %v", err)
	}

	// Insert or replace the pre-key record in the database
	query := `INSERT OR REPLACE INTO prekeys (id, record) VALUES (?, ?)`
	_, err = s.db.ExecContext(ctx, query, id, blob)
	if err != nil {
		return fmt.Errorf("failed to store pre-key: %v", err)
	}

	return nil
}

// Delete removes a pre-key entry identified by the given ID from the store.
func (s *PreKeyStore) Delete(ctx context.Context, id prekey.ID) error {
	// Prepare the SQL query to delete the pre-key
	query := `DELETE FROM prekeys WHERE id = ?`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pre-key: %v", err)
	}
	return nil
}

// SignedPreKeyStore represents the SQLite-backed signed pre-key store.
// Stores user's own signed pre-keys for later use.
var _ prekey.SignedStore = (*SignedPreKeyStore)(nil)

type SignedPreKeyStore struct {
	db *sql.DB
}

// NewSignedPreKeyStore creates a new SQLite-backed signed pre-key store.
func NewSignedPreKeyStore(db *sql.DB) prekey.SignedStore {
	return &SignedPreKeyStore{db: db}
}

// Load retrieves a signed pre-key record for the given ID.
func (s *SignedPreKeyStore) Load(ctx context.Context, id prekey.ID) (*prekey.SignedPreKey, bool, error) {
	var recordData []byte
	query := "SELECT record FROM signed_prekeys WHERE id = ?"
	err := s.db.QueryRowContext(ctx, query, id).Scan(&recordData)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil // No record found
		}
		return nil, false, fmt.Errorf("failed to load signed pre-key record: %w", err)
	}

	// Decode the record data using gob
	// var signedPreKey prekey.SignedPreKey
	// err = decodeAndDeserializeKeyPair(recordData, &signedPreKey)
	var tmp v1.SignedPreKeyRecordStructure
	if err := proto.Unmarshal(recordData, &tmp); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal signed pre-key record: %w", err)
	}

	newKeyPair, err := curve.NewKeyPair(tmp.GetPrivateKey(), tmp.GetPublicKey())
	if err != nil {
		return nil, false, fmt.Errorf("failed to create key pair from bytes: %w", err)
	}
	signedPreKey := prekey.NewSigned(prekey.ID(tmp.GetId()), tmp.GetTimestamp(), newKeyPair, tmp.GetSignature())

	return signedPreKey, true, nil
}

// Store saves a signed pre-key record for the given ID.
func (s *SignedPreKeyStore) Store(ctx context.Context, id prekey.ID, record *prekey.SignedPreKey) error {
	// Get the internal SignedPreKeyRecordStructure from the SignedPreKey object
	signedPreKeyRecord := record.GetSigned()
	if len(signedPreKeyRecord.GetPublicKey()) == 0 || len(signedPreKeyRecord.GetPrivateKey()) == 0 {
		return fmt.Errorf("invalid signed pre-key record: missing public or private key")
	}
	if len(signedPreKeyRecord.GetSignature()) == 0 {
		return fmt.Errorf("invalid signed pre-key record: missing signature")
	}

	// Serialize the SignedPreKeyRecordStructure using protobuf
	recordData, err := proto.Marshal(signedPreKeyRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal signed pre-key record: %w", err)
	}

	// Insert or replace the serialized SignedPreKey in the database
	query := "INSERT OR REPLACE INTO signed_prekeys (id, record) VALUES (?, ?)"
	_, err = s.db.ExecContext(ctx, query, id, recordData)
	if err != nil {
		return fmt.Errorf("failed to store signed pre-key record: %w", err)
	}

	return nil
}
