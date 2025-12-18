package util

import (
	"context"
	"errors"
	"simple-golang/internal/adapter/outbound/postgres/model"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type Validator struct {
	Validator  *validator.Validate
	Translator ut.Translator
	DB         *gorm.DB
}

func NewValidator(db *gorm.DB) *Validator {
	enLocale := en.New()
	uni := ut.New(enLocale, enLocale)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatalf("[NewValidator-1] Translator not found")
	}

	validate := validator.New()

	v := &Validator{
		Validator:  validate,
		Translator: trans,
		DB:         db,
	}

	// Register custom validation
	v.registerCustomValidations()

	return v
}

func (v *Validator) Validate(i interface{}) error {
	err := v.Validator.Struct(i)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				log.Infof("[Validate-1] %s: %s", e.Field(), e.Translate(v.Translator))
				return errors.New(e.Translate(v.Translator))
			}
		}
		// fallback kalau bukan ValidationErrors biasa
		return err
	}
	return nil
}

// ================================================================
// Custom Validation Section
// ================================================================
func (v *Validator) registerCustomValidations() {
	err := v.Validator.RegisterValidation("uniqueEmail", v.uniqueEmail)
	if err != nil {
		log.Errorf("[Validator] failed to register uniqueEmail: %v", err)
	}
}

func (v *Validator) uniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	var user model.User
	err := v.DB.WithContext(context.Background()).
		Where("email = ?", email).
		First(&user).Error

	log.Infof("[uniqueEmail] Checking email: %s", email)

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		log.Infof("[uniqueEmail] Email is unique (not found in database)")
		return true

	case err != nil:
		log.Errorf("[uniqueEmail] Database error while checking email: %v", err)
		return false

	default:
		log.Infof("[uniqueEmail] Email already exists in database as: %s", user.Email)
		return false
	}
}
