package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TokenStorage defines the interface for storing and retrieving JWT and refresh tokens.
type TokenStorage interface {
	StoreTokens(accessToken, refreshToken string) error
	ReadTokens() (string, string, error)
}

// FileTokenStorage is a concrete implementation of TokenStorage using a local file.
type FileTokenStorage struct {
	filePath string
}

// NewFileTokenStorage creates a new FileTokenStorage instance with a specified file path.
func NewFileTokenStorage(filePath string) *FileTokenStorage {
	return &FileTokenStorage{filePath: filePath}
}

// StoreTokens stores both access and refresh tokens in a local file.
func (fts *FileTokenStorage) StoreTokens(accessToken, refreshToken string) error {
	err := os.MkdirAll(filepath.Dir(fts.filePath), 0700)
	if err != nil {
		return err
	}

	data := fmt.Sprintf("access_token:%s\nrefresh_token:%s", accessToken, refreshToken)
	return ioutil.WriteFile(fts.filePath, []byte(data), 0600)
}

// ReadTokens retrieves the access and refresh tokens from a local file.
func (fts *FileTokenStorage) ReadTokens() (string, string, error) {
	data, err := ioutil.ReadFile(fts.filePath)
	if err != nil {
		return "", "", err
	}

	var accessToken, refreshToken string
	_, err = fmt.Sscanf(string(data), "access_token:%s\nrefresh_token:%s", &accessToken, &refreshToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
