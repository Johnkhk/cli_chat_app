package chat

import (
	"context"
	"crypto/rand"
	"testing"
)

var (
	ctx    = context.Background()
	random = rand.Reader
)

func TestSessionManager(t *testing.T) {
	// ctx := context.Background()

	// // Define Alice's and Bob's addresses
	// aliceAddress := address.Address{
	// 	Name:     "+14151111111",
	// 	DeviceID: address.DeviceID(1),
	// }
	// bobAddress := address.Address{
	// 	Name:     "+14151111112",
	// 	DeviceID: address.DeviceID(1),
	// }

	// // Create in-memory stores for testing
	// identityKeyPair, err := identity.GenerateKeyPair(random)
	// require.NoError(t, err)
	// registrationID := uint32(5)
	// aliceIdentityStore := protocol.NewInMemStore(identityKeyPair, registrationID)
	// bobIdentityStore := protocol.NewInMemStore(identityKeyPair, registrationID)

	// aliceSessionStore := session.NewMemoryStore()
	// bobSessionStore := session.NewMemoryStore()

	// alicePreKeyStore := prekey.NewMemoryStore()
	// aliceSignedPreKeyStore := prekey.NewMemoryStore()

	// // Generate keys for Alice and Bob
	// aliceKeyPair, err := identity.GenerateKeyPair(rand.Reader)
	// require.NoError(t, err)

	// bobKeyPair, err := identity.GenerateKeyPair(rand.Reader)
	// require.NoError(t, err)

	// // Store keys in Alice's and Bob's identity stores
	// require.NoError(t, aliceIdentityStore.Store(ctx, aliceAddress, aliceKeyPair.IdentityKey()))
	// require.NoError(t, bobIdentityStore.Store(ctx, bobAddress, bobKeyPair.IdentityKey()))

	// // Initialize Alice's SessionManager
	// aliceSessionManager, err := NewSessionManager(bobAddress, aliceIdentityStore, aliceSessionStore, alicePreKeyStore, aliceSignedPreKeyStore)
	// require.NoError(t, err)

	// // Create a pre-key bundle for Bob to be used by Alice
	// bobPreKeyPair, err := prekey.GenerateKeyPair(rand.Reader)
	// require.NoError(t, err)

	// bobSignedPreKeyPair, err := prekey.GenerateKeyPair(rand.Reader)
	// require.NoError(t, err)

	// bobSignedPreKeySignature, err := bobKeyPair.PrivateKey().Sign(rand.Reader, bobSignedPreKeyPair.PublicKey().Bytes())
	// require.NoError(t, err)

	// preKeyID := prekey.ID(31337)
	// signedPreKeyID := prekey.ID(22)

	// bobPreKeyBundle := &prekey.Bundle{
	// 	RegistrationID:        bobIdentityStore.LocalRegistrationID(ctx),
	// 	DeviceID:              1,
	// 	PreKeyID:              &preKeyID,
	// 	PreKeyPublic:          bobPreKeyPair.PublicKey(),
	// 	SignedPreKeyID:        signedPreKeyID,
	// 	SignedPreKeyPublic:    bobSignedPreKeyPair.PublicKey(),
	// 	SignedPreKeySignature: bobSignedPreKeySignature,
	// 	IdentityKey:           bobKeyPair.IdentityKey(),
	// }

	// // Alice establishes a session with Bob using Bob's pre-key bundle
	// require.NoError(t, aliceSessionManager.EstablishSession(ctx, rand.Reader, bobPreKeyBundle))

	// // Verify session creation
	// aliceSession, exists, err := aliceSessionStore.Load(ctx, bobAddress)
	// assert.NoError(t, err)
	// assert.True(t, exists)

	// // Encrypt a message from Alice to Bob
	// plaintext := []byte("Hello, Bob!")
	// ciphertext, err := aliceSessionManager.EncryptMessage(ctx, plaintext)
	// require.NoError(t, err)

	// // Initialize Bob's SessionManager
	// bobSessionManager, err := NewSessionManager(aliceAddress, bobIdentityStore, bobSessionStore, alicePreKeyStore, aliceSignedPreKeyStore)
	// require.NoError(t, err)

	// // Process the incoming pre-key message at Bob's side
	// incomingMsg, err := message.NewPreKeyFromBytes(ciphertext.Bytes())
	// require.NoError(t, err)

	// require.NoError(t, bobSessionManager.ProcessIncomingMessage(ctx, incomingMsg))

	// // Decrypt the message at Bob's side
	// decryptedMessage, err := bobSessionManager.DecryptMessage(ctx, ciphertext)
	// require.NoError(t, err)
	// assert.Equal(t, plaintext, decryptedMessage)

	// // Verify that Bob has established a session with Alice
	// bobSession, exists, err := bobSessionStore.Load(ctx, aliceAddress)
	// assert.NoError(t, err)
	// assert.True(t, exists)

	// // Further message exchange can be added here for more extensive testing...
}
