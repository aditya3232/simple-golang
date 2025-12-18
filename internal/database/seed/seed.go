package seed

import (
	"github.com/labstack/gommon/log"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) {
	log.Infof("Running database seeds...")
	UserSeed(db)
}
