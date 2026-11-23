package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           string
	AllowedOrigins []string
	Database       DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func CreateDefaultConfig() Config {
	return Config{
		Port:           "8080",
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8888"},
		Database: DatabaseConfig{
			Host: "localhost",
			Port: "5432",
			User: "task_manager_user",
			Name: "task_manager",
		},
	}
}

func validateOrigin(origin string) error {
	if origin == "" {
		return fmt.Errorf("Обнаружен пустой Origin")
	}
	parsedURL, err := url.Parse(origin)
	if err != nil {
		return fmt.Errorf("Обнаружен некорректный URL в Origins")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("Обнаружен недопустимый протокол в Origins")
	}
	if parsedURL.Host == "" {
		return fmt.Errorf("Обнаружен пустой Host в Origins")
	}
	port := parsedURL.Port()
	if port != "" {
		portNumber, err := strconv.Atoi(port)
		if err != nil || portNumber < 1 || portNumber > 65535 {
			return fmt.Errorf("Origin содержит некорректный порт")
		}
	}
	if parsedURL.Path != "" {
		return fmt.Errorf("Обнаружен URL с путем в Origins")
	}
	if parsedURL.RawQuery != "" || parsedURL.ForceQuery {
		return fmt.Errorf("Origin не должен содержать query-параметры")
	}
	if parsedURL.Fragment != "" {
		return fmt.Errorf("Origin не должен содержать fragment")
	}
	if parsedURL.User != nil {
		return fmt.Errorf("Origin не должен содержать учетные данные")
	}
	return nil
}

func Load() (Config, error) {

	c := CreateDefaultConfig()

	if port := os.Getenv("PORT"); port != "" {
		c.Port = port
	}

	port, err := strconv.Atoi(c.Port)
	if err != nil {
		return CreateDefaultConfig(), fmt.Errorf("PORT должен быть числом")
	}

	if port > 65535 || port < 1 {
		return CreateDefaultConfig(), fmt.Errorf("PORT должен быть числом от 1 до 65535")
	}

	if originsString := os.Getenv("ALLOWED_ORIGINS"); originsString == "" {
		c.AllowedOrigins = []string{"http://localhost:3000", "http://localhost:8888"}
	} else {
		c.AllowedOrigins = strings.Split(originsString, ",")
		for i := range c.AllowedOrigins {
			c.AllowedOrigins[i] = strings.TrimSpace(c.AllowedOrigins[i])
			err := validateOrigin(c.AllowedOrigins[i])
			if err != nil {
				return CreateDefaultConfig(), err
			}
		}
	}

	c.Database.Host = os.Getenv("DB_HOST")
	if c.Database.Host == "" {
		c.Database.Host = "localhost"
	}

	c.Database.Port = os.Getenv("DB_PORT")
	if c.Database.Port == "" {
		c.Database.Port = "5432"
	}

	dbPort, err := strconv.Atoi(c.Database.Port)
	if err != nil || dbPort < 1 || dbPort > 65535 {
		return Config{}, fmt.Errorf("DB_PORT должен быть числом от 1 до 65535")
	}

	c.Database.User = os.Getenv("DB_USER")
	if c.Database.User == "" {
		c.Database.User = "task_manager_user"
	}

	passwordFile := os.Getenv("DB_PASSWORD_FILE")

	if passwordFile != "" {
		password, err := os.ReadFile(passwordFile)
		if err != nil {
			return Config{}, fmt.Errorf("не удалось прочитать файл DB_PASSWORD_FILE: %w", err)
		}

		c.Database.Password = strings.TrimRight(string(password), "\r\n")
	} else {
		c.Database.Password = os.Getenv("DB_PASSWORD")
	}

	if c.Database.Password == "" {
		return Config{}, fmt.Errorf("DB_PASSWORD или DB_PASSWORD_FILE не задан")
	}

	c.Database.Name = os.Getenv("DB_NAME")
	if c.Database.Name == "" {
		c.Database.Name = "task_manager"
	}

	return c, nil
}
