package env

import (
	"fmt"
	"os"
)

func GetUsernamePassword() (string, string, error) {
	u, err := getEnvVar("ENPHOTO_USERNAME")
	if err != nil {
		return "", "", err
	}
	p, err := getEnvVar("ENPHOTO_PASSWORD")
	if err != nil {
		return "", "", err
	}
	return u, p, nil
}

func getEnvVar(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment variable %s is not defined", key)
	}
	return val, nil
}
