package data

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	ErrSavingSecrets      = errors.New("failed to save secrets")
	ErrSecretDoesNotExist = errors.New("secret does not exist")
)

type SecretStore struct {
	filePath string
	data     map[string]string
	mutex    *sync.Mutex
}

func CreateSecretStore(filePath string) (*SecretStore, error) {
	store := &SecretStore{
		data:     make(map[string]string),
		filePath: filePath,
		mutex:    &sync.Mutex{},
	}

	_, err := os.Stat(store.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = store.save(); err != nil {
				return nil, err
			}
			return store, nil
		} else {
			return nil, err
		}
	}

	if err = store.load(); err != nil {
		return store, err
	}

	return store, nil
}

func (s *SecretStore) Get(id string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	secret, ok := s.data[id]
	if !ok {
		return "", ErrSecretDoesNotExist
	}

	delete(s.data, id)
	if err := s.save(); err != nil {
		return "", ErrSavingSecrets
	}

	return secret, nil
}

func (s *SecretStore) Put(secret string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := secretId(secret)
	s.data[id] = secret
	if err := s.save(); err != nil {
		delete(s.data, id)
		return "", ErrSavingSecrets
	}

	return id, nil
}

func (s *SecretStore) save() error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s.data)
	if err != nil {
		return err
	}

	return nil
}

func (s *SecretStore) load() error {
	file, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s.data)
	if err != nil {
		return err
	}

	return nil
}

func secretId(secret string) string {
	secretIdHex := md5.Sum([]byte(secret))
	return hex.EncodeToString(secretIdHex[:])
}
