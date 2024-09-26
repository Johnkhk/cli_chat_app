package app

import (
	"fmt"
	"net"
)

func getMACAddress() (string, error) {
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

// func generateOneTimePreKeys(batchSize int) ([][]byte, error) {
//     var oneTimePreKeys [][]byte
//     for i := 0; i < batchSize; i++ {
//         oneTimePreKeyPair, err := curve.GenerateKeyPair(rand.Reader)
//         if err != nil {
//             return nil, fmt.Errorf("failed to generate one-time pre-key pair: %v", err)
//         }
//         oneTimePreKeys = append(oneTimePreKeys, oneTimePreKeyPair.PublicKey().Bytes())
//         // Store the private One-Time PreKey locally for session use
//         err = storeOneTimePreKeyLocally(oneTimePreKeyPair)
//         if err != nil {
//             return nil, err
//         }
//     }
//     return oneTimePreKeys, nil
// }
