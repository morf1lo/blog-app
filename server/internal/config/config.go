package config

import "github.com/joho/godotenv"

func Init() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	return nil
}
