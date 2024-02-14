package config

import "os"

func DBPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func DBName() string {
	return os.Getenv("DB_NAME")
}
