package config

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
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
	AdminUser struct {
		Username string
		Name     string
		Password string
	}
	JWTSecret              string
	JWTExpiration          time.Duration
	RefreshTokenExpiration time.Duration
	Devices                struct {
		USB               string
		BaudRate          int
		DataBits          int
		StopBits          int
		Parity            string
		TimeOut           time.Duration
		LoopReadTime      time.Duration
		LoopbroadcastTime time.Duration
		LoopRecordTime    time.Duration
		DEVICE01          struct {
			SlaveId  int
			DeviceId string
			Name     string
		}
		DEVICE02 struct {
			SlaveId  int
			DeviceId string
			Name     string
		}
		DEVICE03 struct {
			SlaveId  int
			DeviceId string
			Name     string
		}
	}
}

func LoadConfig() (*Config, error) {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		log.Info("CONFIG_FILE environment variable is not set")
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

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

	log.Info("Config loaded successfully", config)

	return &config, nil
}
