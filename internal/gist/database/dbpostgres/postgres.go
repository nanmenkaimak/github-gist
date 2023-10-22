package dbpostgres

import (
	"fmt"
	"github.com/nanmenkaimak/github-gist/internal/gist/config"
	"github.com/nanmenkaimak/github-gist/internal/gist/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	Db *gorm.DB
}

type Config config.DbNode

func (c Config) dsn() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Name, c.User, c.Password, c.SslMode)
}

func New(cfg config.DbNode) (*Db, error) {
	conf := Config(cfg)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  conf.dsn(),
		PreferSimpleProtocol: true,
	}))
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %w", err)
	}

	db.AutoMigrate(&entity.Gist{}, &entity.Commit{}, &entity.File{}, &entity.Comment{}, &entity.Fork{}, &entity.Star{})
	return &Db{
		Db: db,
	}, nil
}

func (d *Db) Close() error {
	sqlDB, _ := d.Db.DB()
	return sqlDB.Close()
}
