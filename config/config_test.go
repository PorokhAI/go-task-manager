package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("ALLOWED_ORIGINS", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "test-password")
	t.Setenv("DB_NAME", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	if cfg.Port != "8080" {
		t.Errorf("ожидался порт %q, получен %q", "8080", cfg.Port)
	}

	wantOrigins := []string{
		"http://localhost:3000",
		"http://localhost:8888",
	}

	if !reflect.DeepEqual(cfg.AllowedOrigins, wantOrigins) {
		t.Errorf(
			"ожидались источники %v, получены %v",
			wantOrigins,
			cfg.AllowedOrigins,
		)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("ожидался DB_HOST %q, получен %q", "localhost", cfg.Database.Host)
	}

	if cfg.Database.Port != "5432" {
		t.Errorf("ожидался DB_PORT %q, получен %q", "5432", cfg.Database.Port)
	}

	if cfg.Database.User != "task_manager_user" {
		t.Errorf("ожидался DB_USER %q, получен %q", "task_manager_user", cfg.Database.User)
	}

	if cfg.Database.Password != "test-password" {
		t.Errorf("неверное значение DB_PASSWORD")
	}

	if cfg.Database.Name != "task_manager" {
		t.Errorf("ожидалась DB_NAME %q, получено %q", "task_manager", cfg.Database.Name)
	}

}

func TestLoadFromEnvironment(t *testing.T) {
	t.Setenv("PORT", "9000")
	t.Setenv(
		"ALLOWED_ORIGINS",
		" https://app.example.com, http://localhost:3000 ",
	)
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "custom_user")
	t.Setenv("DB_PASSWORD", "test-password")
	t.Setenv("DB_NAME", "custom_database")

	cfg, err := Load()

	if err != nil {
		t.Errorf("Ошибка загрузки конфигурации")
	}

	if cfg.Port != "9000" {
		t.Errorf("ожидался порт %q, получен %q", "9000", cfg.Port)
	}

	wantOrigins := []string{
		"https://app.example.com",
		"http://localhost:3000",
	}

	if !reflect.DeepEqual(cfg.AllowedOrigins, wantOrigins) {
		t.Errorf(
			"ожидались источники %v, получены %v",
			wantOrigins,
			cfg.AllowedOrigins,
		)
	}

	if cfg.Database.Host != "db.example.com" {
		t.Errorf("неверный DB_HOST: %q", cfg.Database.Host)
	}

	if cfg.Database.Port != "5433" {
		t.Errorf("неверный DB_PORT: %q", cfg.Database.Port)
	}

	if cfg.Database.User != "custom_user" {
		t.Errorf("неверный DB_USER: %q", cfg.Database.User)
	}

	if cfg.Database.Name != "custom_database" {
		t.Errorf("неверный DB_NAME: %q", cfg.Database.Name)
	}
}

func TestFailedLoadEnvironment(t *testing.T) {
	tests := []struct {
		name    string
		port    string
		origins string
		dbPort  string
	}{
		{
			name:    "порт не является числом",
			port:    "abc",
			origins: "https://app.example.com",
		},
		{
			name:    "порт меньше 1",
			port:    "0",
			origins: "https://app.example.com",
		},
		{
			name:    "порт больше 65535",
			port:    "65536",
			origins: "https://app.example.com",
		},
		{
			name:    "пустой Origin",
			port:    "8080",
			origins: "https://app.example.com, ,http://localhost:3000",
		},
		{
			name:    "неподдерживаемый протокол в URL в Origins",
			port:    "8080",
			origins: "ftp://example.com",
		},
		{
			name:    "пустой Host в Origins",
			port:    "8080",
			origins: "https://",
		},
		{
			name:    "URL с путем в Origins",
			port:    "8080",
			origins: "https://app.example.com/dashboard",
		},
		{
			name:    "Origin с query",
			port:    "8080",
			origins: "https://example.com?token=123",
		},
		{
			name:    "Origin с fragment",
			port:    "8080",
			origins: "https://example.com#section",
		},
		{
			name:    "DB_PORT не является числом",
			port:    "8080",
			origins: "https://app.example.com",
			dbPort:  "abc",
		},
		{
			name:    "DB_PORT выходит за допустимый диапазон",
			port:    "8080",
			origins: "https://app.example.com",
			dbPort:  "70000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PORT", tt.port)
			t.Setenv("ALLOWED_ORIGINS", tt.origins)
			t.Setenv("DB_PASSWORD", "test-password")
			t.Setenv("DB_PORT", tt.dbPort)

			_, err := Load()

			if err == nil {
				t.Fatal("ожидалась ошибка конфигурации")
			}
		})
	}
}

func TestValidateOrigin(t *testing.T) {
	tests := []struct {
		name      string
		origin    string
		wantError bool
	}{
		{
			name:      "корректный HTTPS Origin",
			origin:    "https://app.example.com",
			wantError: false,
		},
		{
			name:      "корректный HTTP Origin с портом",
			origin:    "http://localhost:3000",
			wantError: false,
		},
		{
			name:      "пустой Origin",
			origin:    "",
			wantError: true,
		},
		{
			name:      "недопустимый протокол",
			origin:    "ftp://example.com",
			wantError: true,
		},
		{
			name:      "отсутствует Host",
			origin:    "https://",
			wantError: true,
		},
		{
			name:      "Origin содержит путь",
			origin:    "https://app.example.com/dashboard",
			wantError: true,
		},
		{
			name:      "Origin содержит query",
			origin:    "https://example.com?token=123",
			wantError: true,
		},
		{
			name:      "Origin содержит fragment",
			origin:    "https://example.com#section",
			wantError: true,
		},
		{
			name:      "Origin содержит учетные данные",
			origin:    "https://user:pass@example.com",
			wantError: true,
		},
		{
			name:      "Порт Origin не является числом",
			origin:    "https://example.com:abc",
			wantError: true,
		},
		{
			name:      "Порт Origin больше 65535",
			origin:    "https://example.com:70000",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOrigin(tt.origin)

			if (err != nil) != tt.wantError {
				t.Errorf(
					"validateOrigin(%q): error = %v, wantError = %v",
					tt.origin,
					err,
					tt.wantError,
				)
			}
		})
	}
}

func TestLoadWithoutDatabasePassword(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("ALLOWED_ORIGINS", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_PASSWORD_FILE", "")

	_, err := Load()

	if err == nil {
		t.Fatal("ожидалась ошибка при отсутствии DB_PASSWORD")
	}
}

func TestLoadEnvironmentPassword(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("ALLOWED_ORIGINS", "")
	t.Setenv("DB_PASSWORD", "test-password")
	t.Setenv("DB_PASSWORD_FILE", "")

	c, err := Load()

	if err != nil {
		t.Fatal("Не ожидалась ошибка при загрузке конфигурации с DB_PASSWORD")
	}

	if c.Database.Password != "test-password" {
		t.Errorf("неверное значение DB_PASSWORD: получено %q, ожидалось %q", c.Database.Password, "test-password")
	}
}

func TestLoadFilePassword(t *testing.T) {
	direction := t.TempDir()
	passwordFile := filepath.Join(direction, "db_password.txt")

	if err := os.WriteFile(passwordFile, []byte("test-file-password"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("PORT", "")
	t.Setenv("ALLOWED_ORIGINS", "")
	t.Setenv("DB_PASSWORD", "test-password")
	t.Setenv("DB_PASSWORD_FILE", passwordFile)

	c, err := Load()

	if err != nil {
		t.Fatal("Не ожидалась ошибка при загрузке конфигурации с DB_PASSWORD")
	}

	if c.Database.Password != "test-file-password" {
		t.Errorf("неверное значение DB_PASSWORD: получено %q, ожидалось %q", c.Database.Password, "test-file-password")
	}
}

func TestFailedLoadFilePassword(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("ALLOWED_ORIGINS", "")
	t.Setenv("DB_PASSWORD", "test-password")
	t.Setenv("DB_PASSWORD_FILE", "secret/db_password.txt")

	_, err := Load()

	if err == nil {
		t.Fatal("Ожидалась ошибка при загрузке конфигурации с несуществующим файлом DB_PASSWORD_FILE")
	}
}

func TestLoadFilePasswordWithNewline(t *testing.T) {
	dir := t.TempDir()
	passwordFile := filepath.Join(dir, "db_password.txt")

	if err := os.WriteFile(passwordFile, []byte("test-file-password\r\n"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_PASSWORD_FILE", passwordFile)

	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if c.Database.Password != "test-file-password" {
		t.Errorf("получен пароль %q, ожидался %q", c.Database.Password, "test-file-password")
	}
}
