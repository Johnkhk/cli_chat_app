package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"testing"

	"github.com/Johnkhk/libsignal-go/protocol/address"
	"github.com/Johnkhk/libsignal-go/protocol/curve"
	"github.com/Johnkhk/libsignal-go/protocol/identity"
	"github.com/Johnkhk/libsignal-go/protocol/prekey"
	"github.com/Johnkhk/libsignal-go/protocol/ratchet"
	"github.com/Johnkhk/libsignal-go/protocol/session"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test SQLiteStore
func createTestSQLiteStore(t *testing.T) *SQLiteStore {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err, "should create in-memory SQLite DB without error")

	err = CreateTables(db)
	assert.NoError(t, err, "should create tables without error")

	store := &SQLiteStore{
		DB:                db,
		sessionStore:      NewSessionStore(db),
		preKeyStore:       NewPreKeyStore(db),
		signedPreKeyStore: NewSignedPreKeyStore(db),
		identityStore:     NewIdentityStore(db),
	}
	return store
}

// Session Store Tests
func TestSessionStore(t *testing.T) {
	store := createTestSQLiteStore(t)
	ctx := context.Background()
	aliceIdentity, err := identity.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err, "should generate alice identity key without error")
	bobIdentity, err := identity.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err, "should generate bob identity key without error")
	aliceBaseKey, err := curve.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err, "should generate alice base key without error")
	bobBaseKey, err := curve.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err, "should generate bob base key without error")
	bobEphemeralKey := bobBaseKey

	aliceParams := &ratchet.AliceParameters{
		OurIdentityKeyPair: aliceIdentity,
		OurBaseKeyPair:     aliceBaseKey,
		TheirIdentityKey:   bobIdentity.IdentityKey(),
		TheirSignedPreKey:  bobBaseKey.PublicKey(),
		TheirOneTimePreKey: nil,
		TheirRatchetKey:    bobEphemeralKey.PublicKey(),
	}
	aliceSession, err := session.InitializeAliceSessionRecord(rand.Reader, aliceParams)

	// Store Alice's session
	err = store.SessionStore().Store(ctx, address.Address{Name: "Alice", DeviceID: 1}, aliceSession)
	assert.NoError(t, err, "should store Alice's session without error")

	// Load Alice's session
	_, found, err := store.SessionStore().Load(ctx, address.Address{Name: "Alice", DeviceID: 1})
	assert.NoError(t, err, "should load Alice's session without error")
	assert.True(t, found, "Alice's session should be found")

	// Compare loaded session with stored session
	// assert.Equal(t, aliceSession.previousSessions, loadedAliceSession.previousSessions, "loaded session should match stored session")
	// assert.Equal(t, aliceSession., loadedAliceSession.previousSessions, "loaded session should match stored session")

}

// Identity Store Tests
func TestIdentityStore(t *testing.T) {
	store := createTestSQLiteStore(t)

	ctx := context.Background()
	addr := address.Address{
		Name:     "testUser",
		DeviceID: 1,
	}

	// Create and store a new identity key
	newIdentityKey, _ := identity.GenerateKeyPair(rand.Reader)
	_, err := store.IdentityStore().Store(ctx, addr, newIdentityKey.IdentityKey())
	assert.NoError(t, err, "should store identity key without error")

	// Load identity key
	loadedKey, found, err := store.IdentityStore().Load(ctx, addr)
	assert.NoError(t, err, "should load identity key without error")
	assert.True(t, found, "identity key should be found")
	assert.Equal(t, newIdentityKey.IdentityKey().Bytes(), loadedKey.Bytes(), "loaded identity key should match stored key")
}

// PreKey Store Tests
func TestPreKeyStore(t *testing.T) {
	store := createTestSQLiteStore(t)

	ctx := context.Background()
	preKeyID := prekey.ID(1)
	preKeyPair, err := curve.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err, "should generate key pair without error")
	preKey := prekey.NewPreKey(preKeyID, preKeyPair)

	// Store pre-key
	err = store.PreKeyStore().Store(ctx, preKeyID, preKey)
	assert.NoError(t, err, "should store pre-key without error")

	// Load pre-key
	loadedPreKey, found, err := store.PreKeyStore().Load(ctx, preKeyID)
	assert.NoError(t, err, "should load pre-key without error")
	fmt.Printf("SUCC 3 %v\n", err)
	assert.True(t, found, "pre-key should be found")
	fmt.Println("SUCC 4")
	assert.Equal(t, preKey, loadedPreKey, "loaded pre-key should match stored pre-key")

	// Delete pre-key
	err = store.PreKeyStore().Delete(ctx, preKeyID)
	assert.NoError(t, err, "should delete pre-key without error")

	// Verify deletion
	_, found, err = store.PreKeyStore().Load(ctx, preKeyID)
	assert.NoError(t, err, "should not error when loading deleted pre-key")
	assert.False(t, found, "deleted pre-key should not be found")
}

// Signed PreKey Store Tests
func TestSignedPreKeyStore(t *testing.T) {
	store := createTestSQLiteStore(t)

	ctx := context.Background()
	signedPreKeyID := prekey.ID(1)
	signedPreKeyPair, err := curve.GenerateKeyPair(rand.Reader)
	assert.NoError(t, err, "should generate key pair without error")
	signedPreKey := prekey.NewSigned(signedPreKeyID, 123, signedPreKeyPair, []byte("signature"))

	// Store signed pre-key
	err = store.SignedPreKeyStore().Store(ctx, signedPreKeyID, signedPreKey)
	assert.NoError(t, err, "should store signed pre-key without error")

	// Load signed pre-key
	loadedSignedPreKey, found, err := store.SignedPreKeyStore().Load(ctx, signedPreKeyID)
	assert.NoError(t, err, "should load signed pre-key without error")
	assert.True(t, found, "signed pre-key should be found")
	// fmt.Printf(("loadedSignedPreKey: %v\n"), loadedSignedPreKey.GetSigned().PrivateKey().Bytes())
	assert.Equal(t, signedPreKey.GetSigned().PrivateKey, loadedSignedPreKey.GetSigned().PrivateKey, "loaded signed private key should match stored signed private key")
	assert.Equal(t, signedPreKey.GetSigned().PublicKey, loadedSignedPreKey.GetSigned().PublicKey, "loaded signed public key should match stored signed public key")
	assert.Equal(t, signedPreKey.GetSigned().Signature, loadedSignedPreKey.GetSigned().Signature, "loaded signed signature should match stored signed signature")

	// // Delete signed pre-key
	// err = store.SignedPreKeyStore().Delete(ctx, signedPreKeyID)
	// assert.NoError(t, err, "should delete signed pre-key without error")

	// // Verify deletion
	// _, found, err = store.SignedPreKeyStore().Load(ctx, signedPreKeyID)
	// assert.NoError(t, err, "should not error when loading deleted signed pre-key")
	// assert.False(t, found, "deleted signed pre-key should not be found")
}
