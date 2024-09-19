package store

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/Johnkhk/libsignal-go/protocol/curve"
	"github.com/Johnkhk/libsignal-go/protocol/identity"
)

// serializeKeyPairAndEncode serializes a struct to JSON format (mimics gobEncode).
func serializeKeyPairAndEncode(keyPair interface{}) ([]byte, error) {
	// Serialize the key pair to privateKeyBytes and publicKeyBytes
	fmt.Printf("Type of keyPair: %T\n", keyPair) // Debugging line
	privateKeyBytes, publicKeyBytes, err := serializeKeyPair(keyPair)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize key pair: %v", err)
	}

	// Create a struct to hold the serialized bytes
	serialized := struct {
		PrivateKeyBytes []byte
		PublicKeyBytes  []byte
	}{
		PrivateKeyBytes: privateKeyBytes,
		PublicKeyBytes:  publicKeyBytes,
	}

	// Use gob to encode the serialized byte slices
	encodedData, err := gobEncode(serialized)
	if err != nil {
		return nil, fmt.Errorf("failed to encode serialized key pair: %v", err)
	}

	return encodedData, nil
}

// decodeAndDeserializeKeyPair deserializes a JSON-encoded byte slice into the provided struct (mimics gobDecode).
func decodeAndDeserializeKeyPair(data []byte, value interface{}) error {
	// Define a struct to hold the decoded private and public key bytes
	var serialized struct {
		PrivateKeyBytes []byte
		PublicKeyBytes  []byte
	}

	// Use gob to decode the data into the serialized struct
	err := gobDecode(data, &serialized)
	if err != nil {
		return fmt.Errorf("failed to decode key pair data: %v", err)
	}

	// Deserialize the key pair from the byte slices
	var isIdentityKey bool
	switch value.(type) {
	case *identity.KeyPair:
		isIdentityKey = true
	case *curve.KeyPair:
		isIdentityKey = false
	default:
		return fmt.Errorf("unsupported key pair type %T", value)
	}

	keyPair, err := deserializeKeyPair(serialized.PrivateKeyBytes, serialized.PublicKeyBytes, isIdentityKey)
	if err != nil {
		return fmt.Errorf("failed to deserialize key pair: %v", err)
	}

	// Assign the deserialized key pair back to the provided value interface
	switch v := value.(type) {
	case *identity.KeyPair:
		identityKeyPair, ok := keyPair.(identity.KeyPair)
		if !ok {
			return fmt.Errorf("failed to assert key pair to identity.KeyPair")
		}
		*v = identityKeyPair

	case *curve.KeyPair:
		curveKeyPair, ok := keyPair.(*curve.KeyPair)
		if !ok {
			return fmt.Errorf("failed to assert key pair to curve.KeyPair")
		}
		*v = *curveKeyPair

	default:
		return fmt.Errorf("unsupported key pair type")
	}

	return nil
}

func generateRegistrationID(userID, deviceID uint32) uint32 {
	return (userID << 10) | deviceID // Bit-shifting to combine userID and deviceID
}

type SerializableKeyPair struct {
	PrivateKey []byte `json:"privateKey"`
	PublicKey  []byte `json:"publicKey"`
}

func deserializeKeyPair(privateKeyBytes, publicKeyBytes []byte, isIdentityKey bool) (interface{}, error) {
	// Deserialize curve keys first
	privateKey, err := curve.NewPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize private key: %v", err)
	}

	publicKey, err := curve.NewPublicKey(publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize public key: %v", err)
	}

	if isIdentityKey {
		// Create identity.Key from publicKeyBytes
		identityKey, err := identity.NewKey(publicKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to create identity key: %v", err)
		}
		// Return identity.KeyPair
		return identity.NewKeyPair(privateKey, identityKey), nil
	}

	// Return curve.KeyPair (no need to use Bytes() since privateKey and publicKey are already curve keys)
	return curve.NewKeyPair(privateKey.Bytes(), publicKey.Bytes())
}

func serializeKeyPair(keyPair interface{}) (privateKeyBytes, publicKeyBytes []byte, err error) {
	fmt.Printf("Type of keyPair2222: %T\n", keyPair) // Debugging line
	switch k := keyPair.(type) {
	case identity.KeyPair:
		fmt.Println("XDDDDDDDDDDDDDDd")
		// Serialize the identity.KeyPair
		privateKeyBytes = k.PrivateKey().Bytes() // Get the byte representation of the private key
		publicKeyBytes = k.IdentityKey().Bytes() // Get the byte representation of the identity public key
		return privateKeyBytes, publicKeyBytes, nil

	case *curve.KeyPair:
		fmt.Println("YYYYYYYYYYYYYYYYY")
		// Serialize the curve.KeyPair
		privateKeyBytes = k.PrivateKey().Bytes() // Get the byte representation of the private key
		publicKeyBytes = k.PublicKey().Bytes()   // Get the byte representation of the public key
		return privateKeyBytes, publicKeyBytes, nil
	// case *curve.KeyPair:
	// 	// Handle pointer to curve.KeyPair
	// 	privateKeyBytes = k.PrivateKey().Bytes()
	// 	publicKeyBytes = k.PublicKey().Bytes()
	// 	return privateKeyBytes, publicKeyBytes, nil

	default:
		// Unknown key type
		return nil, nil, fmt.Errorf("unsupported key pair type %T", keyPair)
	}
}

// gobEncode serializes a struct to gob format.
func gobEncode(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode gob: %v", err)
	}
	return buf.Bytes(), nil
}

// gobDecode deserializes a gob-encoded byte slice into the provided struct.
func gobDecode(data []byte, value interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(value)
	if err != nil {
		return fmt.Errorf("failed to decode gob: %v", err)
	}
	return nil
}
