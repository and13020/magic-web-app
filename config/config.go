package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const badConfig string = "Failed to read configs:\nEnv var: %v\nExists: %v\n"

type Config struct {
	DbConfig       DbConfig
	ServerConfig   ServerConfig
	SessionAuthKey string
	SessionEncrKey string
	SessionName    string
}

// DbConfig carries db related config details
type DbConfig struct {
	DbPath string
}

// ServerConfig carries server related config details
type ServerConfig struct {
	Address  string
	RTimeout int
	WTimeout int
}

func getServerConfig() *ServerConfig {
	var sConfig ServerConfig

	enVar := []string{
		"SERVER_ADDRESS",
		"SERVER_READ_TIMEOUT",
		"SERVER_WRITE_TIMEOUT",
	}

	for _, v := range enVar {
		str, exists := os.LookupEnv(v)
		if str == "" || !exists {
			log.Fatalf(badConfig, str, exists)
		}
		switch v {
		case "SERVER_ADDRESS":
			sConfig.Address = str
		case "SERVER_READ_TIMEOUT":
			i, err := strconv.Atoi(str)
			if err != nil {
				log.Fatalf("Failed to read configs:\nEnv var: %v\nExists: %v\n", str, exists)
			}
			sConfig.RTimeout = i
		case "SERVER_WRITE_TIMEOUT":
			i, err := strconv.Atoi(str)
			if err != nil {
				log.Fatalf("Failed to read configs:\nEnv var: %v\nExists: %v\n", str, exists)
			}
			sConfig.WTimeout = i

		}
	}
	return &sConfig
}

func getDBConfig() *DbConfig {
	var d DbConfig

	str, exists := os.LookupEnv("DB_PATH")
	if str == "" || !exists {
		log.Fatalf(badConfig, str, exists)
	}
	d.DbPath = str

	return &d
}

func LoadEnv() *Config {
	// Read from .env into environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sAuthKey := os.Getenv("SESSION_AUTH_KEY")
	if sAuthKey == "" {
		log.Fatal("SESSION_AUTH_KEY is required")
	}

	sEncrKey := os.Getenv("SESSION_ENCRYPTION_KEY")
	if sEncrKey == "" {
		log.Fatal("SESSION_ENCRYPTION_KEY is required")
	}

	sessionName := os.Getenv("SESSION_NAME")
	if sessionName == "" {
		sessionName = "mtg_app_session"
	}

	config := Config{
		ServerConfig:   *getServerConfig(),
		DbConfig:       *getDBConfig(),
		SessionAuthKey: sAuthKey,
		SessionEncrKey: sEncrKey,
		SessionName:    sessionName,
	}

	return &config
}
