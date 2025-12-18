package seed

import (
	"simple-golang/internal/adapter/outbound/postgres/model"
	"simple-golang/util"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func UserSeed(db *gorm.DB) {
	bytes, err := util.HashPassword("admin123")
	if err != nil {
		log.Fatalf("[UserSeed-1]: %v", err)
	}

	admin := model.User{
		Name:     "super admin",
		Email:    "superadmin@mail.com",
		Password: bytes,
		Phone:    "085162665063",
		Address:  "Kelurahan Cimanggis, BOJONGGEDE, KAB. BOGOR, JAWA BARAT, ID, 16920",
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "superadmin@mail.com"}).Error; err != nil {
		log.Errorf("[SeedAdmin-2]: %v", err)
	} else {
		log.Infof("User %s created", admin.Name)
	}
}
