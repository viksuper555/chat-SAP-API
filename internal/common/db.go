package common

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"messenger/config"
	"messenger/model"
)

var Db *gorm.DB

func InitDb(cfg config.Config) (*gorm.DB, error) {
	dbCfg := &gorm.Config{}
	url := GetDbUrl(cfg)
	var err error
	switch cfg.Database.Dialect {
	case "postgres":
		if Db, err = gorm.Open(postgres.Open(url), dbCfg); err != nil {
			log.Fatalln(err)
		}
		break
	}
	if err != nil {
		return nil, err
	}

	models := []interface{}{
		&model.User{},
		&model.Message{},
	}
	for _, m := range models {
		log.Printf("Auto migration of m")
		err = Db.AutoMigrate(m)
		if err != nil {
			log.Printf("Error migrating m: %s", err.Error())
			return nil, errors.Wrap(Db.Error, "Failed to migrate m")
		}
	}
	if err != nil {
		return nil, err
	}
	return Db, nil
}

func GetDbUrl(cfg config.Config) string {
	//dbURL := "postgres://pg:pass@localhost:5432/messenger"
	db := cfg.Database
	return fmt.Sprintf("%s://%s:%s@%s/%s", db.Dialect, db.User, db.Password, db.IP, db.Name)
}
