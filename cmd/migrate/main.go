package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"example.com/go-backend/config"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var dirtyErr migrate.ErrDirty

func main() {
	if err := run(); err != nil {

		if errors.As(err, &dirtyErr) {
			log.Fatal(fmt.Errorf("База находится в dirty-состоянии: %w", err))
		}

		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("Ошибка конфигурации: %w", err)
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.Database.DSN(),
	)
	if err != nil {
		return fmt.Errorf("Ошибка инициализации миграций: %w", err)
	}

	defer m.Close()

	if len(os.Args) < 2 {
		return errors.New("Не указано направление миграции: Используйте 'up' или 'down'")
	}

	direction := os.Args[1]

	if direction != "up" && direction != "down" {
		return fmt.Errorf("Неверное направление миграции: %s. Используйте 'up' или 'down'", direction)
	}

	var migrationErr error

	switch direction {
	case "up":
		migrationErr = m.Up()
	case "down":
		migrationErr = m.Steps(-1)
	}

	version, dirty, versionErr := m.Version()

	if errors.Is(versionErr, migrate.ErrNilVersion) {
		if migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNoChange) {
			return fmt.Errorf("ошибка миграции: %w; активная версия БД отсутствует", migrationErr)
		}

		if direction == "down" && migrationErr == nil {
			log.Println("Миграция успешно откатана")
		}

		if direction == "down" && errors.Is(migrationErr, migrate.ErrNoChange) {
			log.Println("Нет миграций для отката")
		}

		log.Println("Активная версия миграций отсутствует")
		return nil
	}

	if versionErr != nil {
		if migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNoChange) {
			return fmt.Errorf(
				"ошибка миграции: %w; не удалось получить состояние БД: %v",
				migrationErr,
				versionErr,
			)
		}

		return fmt.Errorf("не удалось получить состояние БД: %w", versionErr)
	}

	if migrationErr != nil {
		if errors.Is(migrationErr, migrate.ErrNoChange) {
			switch direction {
			case "up":
				log.Println("Миграции уже актуальны")
			case "down":
				log.Println("Нет миграций для отката")
			}
		} else {
			return fmt.Errorf(
				"ошибка миграции: %w; состояние БД: version=%d, dirty=%t",
				migrationErr,
				version,
				dirty,
			)
		}
	}

	if !errors.Is(migrationErr, migrate.ErrNoChange) {
		switch direction {
		case "up":
			log.Println("Миграции успешно применены")
		case "down":
			log.Println("Миграции успешно откатаны")
		}
	}

	log.Printf("Версия БД: %d, Dirty: %t", version, dirty)
	return nil
}
