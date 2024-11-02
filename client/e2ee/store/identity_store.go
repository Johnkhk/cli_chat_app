package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Johnkhk/libsignal-go/protocol/address"
	"github.com/Johnkhk/libsignal-go/protocol/direction"
	"github.com/Johnkhk/libsignal-go/protocol/identity"
)

// IdentityStore represents a SQLite-backed identity store.

var _ identity.Store = (*IdentityStore)(nil)

type IdentityStore struct {
	db *sql.DB
}

// NewIdentityStore creates a new SQLite-backed identity store.
func NewIdentityStore(db *sql.DB) identity.Store {
	return &IdentityStore{
		db: db,
	}
}

// KeyPair returns the associated identity key pair from the database.
func (s *IdentityStore) KeyPair(_ context.Context) identity.KeyPair {
	var keyPair identity.KeyPair
	var keyPairData []byte

	// Query for the key_pair from the database
	query := "SELECT key_pair FROM local_identity LIMIT 1"
	err := s.db.QueryRow(query).Scan(&keyPairData)
	if err != nil {
		if err == sql.ErrNoRows {
			// Log that no identity key pair was found
			log.Printf("No identity key pair found: %v", err)
		} else {
			// Log other errors
			log.Printf("Failed to retrieve identity key pair: %v", err)
		}
		return keyPair // Return zero-value key pair if error occurs
	}

	// Decode the key pair from gob format
	err = DecodeAndDeserializeKeyPair(keyPairData, &keyPair)
	if err != nil {
		// Log decoding error
		log.Printf("Failed to decode identity key pair: %v", err)
		return identity.KeyPair{} // Return zero-value key pair in case of decoding error
	}

	return keyPair
}

// LocalRegistrationID returns the associated registration ID from the database.
func (s *IdentityStore) LocalRegistrationID(_ context.Context) uint32 {
	var registrationID uint32

	// Query for the registration_id from the database
	query := "SELECT registration_id FROM local_identity LIMIT 1"
	err := s.db.QueryRow(query).Scan(&registrationID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Log that no registration ID was found
			log.Printf("No registration ID found: %v", err)
		} else {
			// Log other errors
			log.Printf("Failed to retrieve registration ID: %v", err)
		}
		return 0 // Return default value of 0 if error occurs
	}

	return registrationID
}

// Load loads the identity key associated with the remote address.
func (s *IdentityStore) Load(ctx context.Context, addr address.Address) (identity.Key, bool, error) {
	var keyData []byte
	query := "SELECT key_data FROM identities WHERE address = ?"
	err := s.db.QueryRowContext(ctx, query, addr.String()).Scan(&keyData)

	if err != nil {
		if err == sql.ErrNoRows {
			return identity.Key{}, false, nil // No record found
		}
		return identity.Key{}, false, fmt.Errorf("failed to load identity key: %w", err)
	}

	// Rebuild the identity key from the stored bytes
	key, err := identity.NewKey(keyData)
	if err != nil {
		return identity.Key{}, false, fmt.Errorf("failed to create identity key from stored data: %w", err)
	}

	return key, true, nil
}

// Store stores the identity key associated with the remote address and returns
// "true" if there is already an entry for the address that is overwritten
// with a new identity key.
func (s *IdentityStore) Store(ctx context.Context, addr address.Address, identityKey identity.Key) (bool, error) {
	// Check if the identity already exists
	var existingKeyData []byte
	query := "SELECT key_data FROM identities WHERE address = ?"
	err := s.db.QueryRowContext(ctx, query, addr.String()).Scan(&existingKeyData)

	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("failed to check existing identity key: %w", err)
	}

	// Store the new identity key (public key bytes)
	// fmt.Printf("Storing identity key for address: %v\n", identityKey.PublicKey().Bytes())
	newKeyData := identityKey.Bytes()

	// Insert or update the identity key in the SQLite database
	insertQuery := "INSERT OR REPLACE INTO identities (address, key_data, trust_level) VALUES (?, ?, ?)"
	_, err = s.db.ExecContext(ctx, insertQuery, addr.String(), newKeyData, 1) // Trust level can be set as 1 (trusted) initially
	if err != nil {
		return false, fmt.Errorf("failed to store identity key: %w", err)
	}

	// Return true if the key already existed
	return len(existingKeyData) > 0, nil
}

// Clear removes all items from the store.
func (s *IdentityStore) Clear() error {
	_, err := s.db.Exec("DELETE FROM identities")
	if err != nil {
		return fmt.Errorf("failed to clear identity store: %w", err)
	}
	return nil
}

// IsTrustedIdentity returns "true" if the given identity key for the given address is already trusted.
// If there is no entry for the given address, the given identity key is trusted.
func (s *IdentityStore) IsTrustedIdentity(ctx context.Context, addr address.Address, identityKey identity.Key, _ direction.Direction) (bool, error) {
	// var keyData []byte
	// query := "SELECT key_data FROM identities WHERE address = ?"
	// err := s.db.QueryRowContext(ctx, query, addr.String()).Scan(&keyData)

	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return true, nil // No record found, trust the new identity key
	// 	}
	// 	return false, fmt.Errorf("failed to load identity key: %w", err)
	// }

	// // Decode the stored identity key
	// var storedKeyPair identity.KeyPair
	// fmt.Println("WHATTTTTTTTTTTT", keyData)
	// err = DecodeAndDeserializeKeyPair(keyData, &storedKeyPair)
	// if err != nil {
	// 	return false, fmt.Errorf("failed to decode identity key: %w", err)
	// }

	// // Compare the provided identity key with the stored one
	// return identityKey.Equal(storedKeyPair.IdentityKey()), nil

	return true, nil
}
