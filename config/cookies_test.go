package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCookiesEncryptedRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".cookies")
	const cookies = "UID=1; CID=2; SEID=3; KID=4"
	if err := WriteCookiesFile(path, cookies); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(raw), encryptedCookiesPrefix) {
		t.Fatalf("cookies file is not encrypted: %q", raw)
	}
	if info, err := os.Stat(path); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("cookies permissions = %v, err = %v", info.Mode().Perm(), err)
	}
	if info, err := os.Stat(path + ".key"); err != nil || info.Mode().Perm() != 0o600 {
		t.Fatalf("key permissions = %v, err = %v", info.Mode().Perm(), err)
	}
	got, err := ReadCookiesFile(path)
	if err != nil || got != cookies {
		t.Fatalf("round trip = %q, err = %v", got, err)
	}
}

func TestCookiesLegacyPlaintextStillReadable(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".cookies")
	if err := os.WriteFile(path, []byte("UID=legacy\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ReadCookiesFile(path)
	if err != nil || got != "UID=legacy" {
		t.Fatalf("legacy read = %q, err = %v", got, err)
	}
}

func TestLoadDecryptsCookiesFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", dir)
	const cookies = "UID=1; CID=2; SEID=3; KID=4"
	if err := WriteCookiesFile(filepath.Join(dir, ".cookies"), cookies); err != nil {
		t.Fatal(err)
	}
	content := "[auth]\ncookies_file = \".cookies\"\n\n[server]\nport = 8115\n"
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, _, err := LoadWithOptions(CLIParams{}, LoadOptions{Auth: true})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Auth.Cookies != cookies {
		t.Fatalf("Cookies = %q, want plaintext", cfg.Auth.Cookies)
	}
}
