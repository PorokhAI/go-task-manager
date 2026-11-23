package config

import (
	"net/url"
	"strings"
	"testing"
)

func TestDatabaseConfigDSN(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "task_manager_user",
		Password: "password",
		Name:     "task_manager",
	}

	dsn := cfg.DSN()

	want := "postgres://task_manager_user:password@localhost:5432/task_manager?sslmode=disable"

	if dsn != want {
		t.Errorf("ожидался DSN %q, получен %q", want, dsn)
	}
}

func TestDatabaseConfigDSNPasswordParse(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "task_manager_user",
		Password: "my@pass:word#123",
		Name:     "task_manager",
	}

	dsn := cfg.DSN()

	if !strings.Contains(dsn, "my%40pass%3Aword%23123") {
		t.Errorf("пароль некорректно закодирован в DSN: %q", dsn)
	}

	parsedURL, err := url.Parse(dsn)

	if err != nil {
		t.Fatalf("не удалось разобрать DSN: %v", err)
	}

	password, hasPassword := parsedURL.User.Password()
	if !hasPassword {
		t.Error("ожидался пароль в DSN")
	}

	if password != "my@pass:word#123" {
		t.Errorf("ожидался пароль %q, получен %q", "my@pass:word#123", password)
	}
}

func TestDatabaseConfigDSNUserParse(t *testing.T) {
	cfg := DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "task_manager_user",
		Password: "password",
		Name:     "task_manager#test",
	}

	dsn := cfg.DSN()

	parsedURL, err := url.Parse(dsn)

	if err != nil {
		t.Fatalf("не удалось разобрать DSN: %v", err)
	}

	dbName := parsedURL.Path
	if dbName != "/task_manager#test" {
		t.Errorf("ожидалось имя базы данных %q, получено %q", "/task_manager#test", dbName)
	}
	if !strings.Contains(dsn, "task_manager%23test") {
		t.Errorf("символ # не закодирован в DSN: %q", dsn)
	}
}
