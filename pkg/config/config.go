package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
)

type Config struct {
	TestField     string `json:"test_field"`
	ListenAddress string `json:"listen_address"`
}

func (config Config) String() string {
	return fmt.Sprintf("{TestField: %q, ListenAddress: %q}",
		config.TestField, config.ListenAddress)
}

func GeneratePreSharedKey() (string, error) {
	// 生成 256 位密钥
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func LoadConfig(configPath string) (*Config, error) {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadPeerKey(filePath string) (crypto.PrivKey, crypto.PubKey, error) {
	privBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}
	priv, err := crypto.UnmarshalPrivateKey(privBytes)
	return priv, priv.GetPublic(), err
}

func SavePeerKey(filePath string, priv crypto.PrivKey) error {
	privBytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, privBytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

func GeneratePeerKey(filePath string) (crypto.PrivKey, crypto.PubKey, error) {
	priv, pub, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, nil, err
	} else {
		err := SavePeerKey(filePath, priv)
		if err != nil {
			return nil, nil, err
		} else {
			return priv, pub, err
		}
	}
}
