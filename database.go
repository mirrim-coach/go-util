package goutils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //Adds the postgres dialect to gorm
)

var (
	defaultConfig = DatabaseConfig{
		Host:         GetEnvVariable("DB_HOST", "localhost"),
		Port:         GetEnvVariable("DB_PORT", "5432"),
		User:         GetEnvVariable("DB_USER", "local"),
		Password:     GetEnvVariable("DB_PASS", ""),
		DatabaseName: GetEnvVariable("DB_NAME", "auth"),
		SSLMode:      GetEnvVariable("DB_SSL", "disable"),
	}
)

// DatabaseConfig is a struct containing the database configuration
type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	SSLMode      string
}

func (d DatabaseConfig) toString() string {
	return fmt.Sprintf(`host=%s port=%s user=%s dbname=%s password="%s" sslmode=%s`, d.Host, d.Port, d.User, d.DatabaseName, d.Password, d.SSLMode)
}

// GetDBWithConfig returns a new db instance with the specified config
func GetDBWithConfig(conf string) *gorm.DB {
	db, err := gorm.Open("postgres", conf)
	if err != nil {
		Logger().Fatal(err.Error())
	}

	return db
}

// GetDB returns an instance of the db
func GetDB() *gorm.DB {
	return GetDBWithConfig(defaultConfig.toString())
}
