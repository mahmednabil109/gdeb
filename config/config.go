package config

import (
	"bufio"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	DirPrivateKey = "/home/husseljo/gdeb/private_key"
)

type Config struct {
	DataDir string
}

func New() *Config {
	return &Config{
		DirPrivateKey,
	}
}

func (c *Config) NodeKey() ed25519.PrivateKey {
	priv, err := loadPK(c.DataDir)
	if err != nil {
		priv, err = generateKey(c.DataDir)
		log.Println("NodeKey error: ", err)
	}
	return priv
}

func generateKey(file string) (ed25519.PrivateKey, error) {
	_, pk, _ := ed25519.GenerateKey(nil)
	f, err := os.Create(file)

	if err != nil {
		return nil, err
	}
	defer f.Close()
	f.WriteString(hex.EncodeToString(pk))

	return pk, nil
}

func loadPK(file string) (ed25519.PrivateKey, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	nonEmpty := scanner.Scan()
	var privKey string
	if nonEmpty {
		privKey = scanner.Text()
	} else {
		return nil, fmt.Errorf("Private Key file is empty!")
	}

	privKey = strings.TrimSpace(privKey)
	if len(privKey) != 128 {
		return nil, fmt.Errorf("Invalid Private Key")
	}

	sk, _ := hex.DecodeString(privKey)

	return ed25519.PrivateKey(sk), nil

}
