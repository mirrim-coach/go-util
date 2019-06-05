package goutils

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //Adds the postgres dialect to gorm
	mocket "github.com/selvatico/go-mocket"
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

var _DefaultConfig *DatabaseConfig

//GetDefaultConfig returns a default config option
func GetDefaultConfig() *DatabaseConfig {
	if _DefaultConfig == nil {
		_DefaultConfig = &DatabaseConfig{
			Host:         GetEnvVariable("DB_HOST", "localhost"),
			Port:         GetEnvVariable("DB_PORT", "5432"),
			User:         GetEnvVariable("DB_USER", "local"),
			Password:     GetEnvVariable("DB_PASS", ""),
			DatabaseName: GetEnvVariable("DB_NAME", "auth"),
			SSLMode:      GetEnvVariable("DB_SSL", "disable"),
		}
	}
	return _DefaultConfig
}

// GetDBWithConfig returns a new db instance with the specified config
func GetDBWithConfig(conf string) *gorm.DB {
	Logger().Debugf("Loading database with conf: \"%s\"", conf)
	db, err := gorm.Open("postgres", conf)
	if err != nil {
		Logger().Fatal(err.Error())
		return nil
	}
	return db
}

//GetMockDB returns the mock database
func GetMockDB() *gorm.DB {
	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
	mocket.Catcher.Logging = true
	// GORM
	db, err := gorm.Open(mocket.DriverName, "mock connection string") // Can be any connection string
	if err != nil {
		Logger().Fatal(err.Error())
	}
	return db
}

// GetDB returns an instance of the db
func GetDB() *gorm.DB {
	env := GetEnvVariable("ENV", "development")
	if env == "test" {
		return GetMockDB()
	}
	return GetDBWithConfig(GetDefaultConfig().toString())
}
