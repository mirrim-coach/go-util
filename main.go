package goutils

import (
	"encoding/base32"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"time"
)

var (
	envIsLoaded = false
)

// GenRandomBase32Bytes generates random number of bytes in base32 format
func GenRandomBase32Bytes(num int) string {
	rand.Seed(time.Now().UnixNano())

	token := make([]byte, num)
	rand.Read(token)

	return base32.StdEncoding.EncodeToString(token)
}

// GetEnvVariable get environment variable with default
func GetEnvVariable(env string, d string) string {
	if !envIsLoaded {
		envIsLoaded = true
		err := godotenv.Load()
		if err != nil {
			log.Info(".Env file failed to load: " + err.Error())
		}
	}

	val, exist := os.LookupEnv(env)

	if !exist {
		return d
	}

	return val
}
