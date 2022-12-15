package db

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"messenger/config"
	"messenger/graphql"
)

type Database interface {
	Create(message *graphql.Message) error
	LoadUserMessages(userId string) (*[]graphql.Message, error)
	Close() error
}

type GormDB struct {
	cfg    config.Config
	dbConn *gorm.DB
}

func Init(cfg config.Config) (Database, error) {
	dbCfg := &gorm.Config{}
	url := GetDbUrl(cfg)
	var dbConn *gorm.DB
	var err error
	switch cfg.Database.Dialect {
	case "postgres":
		if dbConn, err = gorm.Open(postgres.Open(url), dbCfg); err != nil {
			log.Fatalln(err)
		}
		break
	}

	models := []interface{}{
		&graphql.User{},
		&graphql.Message{},
	}
	for _, m := range models {
		log.Printf("Auto migration of m")
		err = dbConn.AutoMigrate(m)
		if err != nil {
			log.Printf("Error migrating m: %s", err.Error())
			return nil, errors.Wrap(dbConn.Error, "Failed to migrate m")
		}
	}

	database := &GormDB{
		dbConn: dbConn,
		cfg:    cfg,
	}

	return database, nil
}

func GetDbUrl(cfg config.Config) string {
	//dbURL := "postgres://pg:pass@localhost:5432/messenger"
	db := cfg.Database
	return fmt.Sprintf("%s://%s:%s@%s/%s", db.Dialect, db.User, db.Password, db.IP, db.Name)
}

func (d *GormDB) Create(message *graphql.Message) error {
	result := d.dbConn.Create(message)
	if result.Error != nil {
		return errors.Wrap(result.Error, "Failed to persist notification in DB")
	}
	return nil
}

func (d *GormDB) LoadUserMessages(userId string) (*[]graphql.Message, error) {
	var messages *[]graphql.Message

	if result := d.dbConn.Find(&messages, ""); result.Error != nil {
		fmt.Println(result.Error)
	}

	return messages, nil
}

func (d *GormDB) Close() error {
	db, err := d.dbConn.DB()
	if err != nil {
		return err
	}

	return db.Close()
}
