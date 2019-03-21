package goutils

import (
	"encoding/base32"
	"github.com/joho/godotenv"
	"math/rand"
	"os"
	"strconv"
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
			Logger().Debugf(".Env file failed to load: %s. This is ok!", err.Error())
		}
	}

	val, exist := os.LookupEnv(env)

	if !exist {
		return d
	}

	return val
}

//GetEnvVariableInt64 returns the environment variable as an int64
func GetEnvVariableInt64(env string, d string) (int64, error) {
	v := GetEnvVariable(env, d)
	i, err := strconv.ParseInt(v, 10, 64)

	if err != nil {
		return 0, err
	}
	return i, nil
}
