package cmd

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"strings"
	"testing"
)

func TestHashSHA512(t *testing.T) {
	root := NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"hash", "testpassword", "sha512"})

	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := strings.TrimSpace(buf.String())
	expected := fmt.Sprintf("%x", sha512.Sum512([]byte("testpassword")))
	if result != expected {
		t.Errorf("expected SHA-512 hash '%s', got '%s'", expected, result)
	}
}

func TestHashBcrypt(t *testing.T) {
	root := NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"hash", "testpassword", "bcrypt"})

	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(result, "$2a$") && !strings.HasPrefix(result, "$2b$") {
		t.Errorf("expected bcrypt hash starting with $2a$12$ or $2b$12$, got '%s'", result)
	}
}

func TestHashPlain(t *testing.T) {
	root := NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"hash", "mysecret", "plain"})

	if err := root.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result := strings.TrimSpace(buf.String())
	if result != "mysecret" {
		t.Errorf("expected plain output 'mysecret', got '%s'", result)
	}
}

func TestHashInvalidType(t *testing.T) {
	root := NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs([]string{"hash", "password", "invalid_type"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("expected no error from cobra, got %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "unsupported hash type") {
		t.Errorf("expected error message about unsupported hash type in output, got: %s", output)
	}
}

func TestHashExactArgs(t *testing.T) {
	root := NewRootCmd()
	root.SetArgs([]string{"hash", "password"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for missing argument, got nil")
	}

	root.SetArgs([]string{"hash", "password", "sha512", "extra"})
	err = root.Execute()
	if err == nil {
		t.Fatal("expected error for extra argument, got nil")
	}
}

func TestNewHashCmd(t *testing.T) {
	cmd := NewHashCmd()
	if cmd == nil {
		t.Fatal("expected non-nil command")
	}
	if cmd.Use != "hash <password> <type>" {
		t.Errorf("expected use 'hash <password> <type>', got '%s'", cmd.Use)
	}
	if cmd.Short != "Generate a password hash" {
		t.Errorf("expected short 'Generate a password hash', got '%s'", cmd.Short)
	}
}

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	if cmd == nil {
		t.Fatal("expected non-nil root command")
	}
	if cmd.Use != "dedctl" {
		t.Errorf("expected use 'dedctl', got '%s'", cmd.Use)
	}
}
