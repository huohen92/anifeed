package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

const encryptedCookiesPrefix = "afenc:v1:"

func cookiesKeyPath(path string) string {
	return path + ".key"
}

func loadCookiesKey(path string, create bool) ([]byte, error) {
	keyPath := cookiesKeyPath(path)
	key, err := os.ReadFile(keyPath)
	if err == nil && len(key) == 32 {
		return key, nil
	}
	if !create {
		return nil, fmt.Errorf("cookies 密钥不存在或长度不正确: %s", keyPath)
	}
	key = make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	if err := os.WriteFile(keyPath, key, 0o600); err != nil {
		return nil, err
	}
	return key, nil
}

// ReadCookiesFile reads a plaintext legacy cookies file or decrypts the
// afenc:v1 format used by the Web UI.
func ReadCookiesFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	text := strings.TrimSpace(string(data))
	if !strings.HasPrefix(text, encryptedCookiesPrefix) {
		return text, nil
	}
	key, err := loadCookiesKey(path, false)
	if err != nil {
		return "", err
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(text, encryptedCookiesPrefix))
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	if len(raw) < gcm.NonceSize() {
		return "", fmt.Errorf("cookies 密文长度不正确")
	}
	plain, err := gcm.Open(nil, raw[:gcm.NonceSize()], raw[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}

// WriteCookiesFile atomically writes cookies. Non-empty values are encrypted
// with AES-GCM; the per-file key is created on demand with mode 0600.
func WriteCookiesFile(path, cookies string) error {
	cookies = strings.TrimSpace(cookies)
	if cookies == "" {
		return os.WriteFile(path, nil, 0o600)
	}
	key, err := loadCookiesKey(path, true)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(cookies), nil)
	return os.WriteFile(path, []byte(encryptedCookiesPrefix+base64.StdEncoding.EncodeToString(sealed)+"\n"), 0o600)
}
