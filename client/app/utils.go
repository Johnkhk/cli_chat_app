package app

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
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

func GetAppDirPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	// Check if APP_DIR_PATH is set in the environment
	if os.Getenv("APP_DIR_PATH") != "" {
		return os.Getenv("APP_DIR_PATH"), nil
	}

	var appDir string
	switch runtime.GOOS {
	case "windows":
		appDir = filepath.Join(usr.HomeDir, "AppData", "Local", "cli_chat_app")
	case "darwin":
		appDir = filepath.Join(usr.HomeDir, "Library", "Application Support", "cli_chat_app")
	case "linux":
		appDir = filepath.Join(usr.HomeDir, ".cli_chat_app")
	default:
		return "", fmt.Errorf("unsupported platform")
	}

	return appDir, nil
}
