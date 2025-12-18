package config

import "github.com/spf13/viper"

type App struct {
	AppPort      string `json:"app_port"`
	AppEnv       string `json:"app_env"`
	JwtSecretKey string `json:"jwt_secret_key"`
	JwtIssuer    string `json:"jwt_issuer"`
}

type PsqlDB struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	User      string `json:"user"`
	Password  string `json:"password"`
	DBName    string `json:"db_name"`
	DBMaxOpen int    `json:"db_max_open"`
	DBMaxIdle int    `json:"db_max_idle"`
}

type Redis struct {
	Addr     string `json:"addr"`
	DB       int    `json:"db"`
	Password string `json:"password"`
}

type Config struct {
	App    App    `json:"app"`
	PsqlDB PsqlDB `json:"psql_db"`
	Redis  Redis  `json:"redis"`
}

func NewConfig() *Config {
	return &Config{
		App: App{
			AppPort:      viper.GetString("APP_PORT"),
			AppEnv:       viper.GetString("APP_ENV"),
			JwtSecretKey: viper.GetString("JWT_SECRET_KEY"),
			JwtIssuer:    viper.GetString("JWT_ISSUER"),
		},
		PsqlDB: PsqlDB{
			Host:      viper.GetString("DATABASE_HOST"),
			Port:      viper.GetString("DATABASE_PORT"),
			User:      viper.GetString("DATABASE_USER"),
			Password:  viper.GetString("DATABASE_PASSWORD"),
			DBName:    viper.GetString("DATABASE_NAME"),
			DBMaxOpen: viper.GetInt("DATABASE_MAX_OPEN_CONNECTION"),
			DBMaxIdle: viper.GetInt("DATABASE_MAX_IDLE_CONNECTION"),
		},
		Redis: Redis{
			Addr:     viper.GetString("REDIS_ADDR"),
			DB:       viper.GetInt("REDIS_DB"),
			Password: viper.GetString("REDIS_PASSWORD"),
		},
	}
}
