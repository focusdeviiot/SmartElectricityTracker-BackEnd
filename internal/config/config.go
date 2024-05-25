package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Dbname   string
	}
	PowerMeter struct {
		Device string
	}
	JWTSecret              string
	JWTExpiration          time.Duration
	RefreshTokenExpiration time.Duration
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../configs")

	// อ่านค่าจาก environment variables
	viper.AutomaticEnv()

	// Binding environment variables to Viper keys
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("JWTSecret", "JWT_SECRET")
	viper.BindEnv("power_meter.device", "POWER_METER_DEVICE")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
