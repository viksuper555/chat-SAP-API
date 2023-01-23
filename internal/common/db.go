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
		&model.Room{},
	}
	log.Printf("Auto migration")
	for _, m := range models {
		err = Db.AutoMigrate(m)
		if err != nil {
			log.Printf("Error migrating m: %s", err.Error())
			return nil, errors.Wrap(Db.Error, "Failed to migrate m")
		}
	}
	err = AutoInitRooms()
	if err != nil {
		return nil, err
	}
	return Db, nil
}
func AutoInitRooms() error {
	var rm model.Room
	if err := Db.Where("id = ?", "global").Find(&rm).Error; err == nil {
		return nil
	}
	log.Printf("Creating default room")
	grm := &model.Room{
		ID:    "global",
		Name:  "Global Room",
		Users: []*model.User{},
	}

	if err := Db.Create(&grm).Error; err != nil {
		return err
	}

	return nil
}
func GetDbUrl(cfg config.Config) string {
	//dbURL := "postgres://pg:pass@localhost:5432/messenger"
	db := cfg.Database
	return fmt.Sprintf("%s://%s:%s@%s/%s", db.Dialect, db.User, db.Password, db.IP, db.Name)
}

func GetUserRoomIds(db *gorm.DB, userId int) ([]string, error) {
	roomIds := make([]string, 0)
	if err := db.Table("user_room").Distinct("room_id").Where("user_id = ?", userId).Find(&roomIds).Error; err != nil {
		return nil, err
	}
	return roomIds, nil
}
func GetRooms(db *gorm.DB) ([]*model.Room, error) {
	var rooms []*model.Room
	if err := db.Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func AddUserToRoom(db *gorm.DB, userId int, roomId string) error {
	var u model.User
	err := db.Where(&model.User{ID: userId}).First(&u).Error
	if err != nil {
		return err
	}
	var r model.Room
	err = db.Where(&model.Room{ID: roomId}).First(&r).Error
	if err != nil {
		return err
	}
	err = db.Model(&r).Association("Users").Append(&u)
	if err != nil {
		return err
	}
	return nil
}

func RemoveUserFromRoom(db *gorm.DB, userId int, roomId string) error {
	var u model.User
	err := db.Where(&model.User{ID: userId}).First(&u).Error
	if err != nil {
		return err
	}
	var r model.Room
	err = db.Preload("Users").Where(&model.Room{ID: roomId}).First(&r).Error
	if err != nil {
		return err
	}
	if len(r.Users) > 1 {
		err = db.Model(&r).Association("Users").Delete(&u)
		if err != nil {
			return err
		}
	} else {
		err = db.Delete(&r).Error
		if err != nil {
			return err
		}
	}
	return nil
}
