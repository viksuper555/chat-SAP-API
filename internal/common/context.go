package common

import (
	"gorm.io/gorm"
)

type Context struct {
	Database *gorm.DB
}
