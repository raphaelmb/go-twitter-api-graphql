package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type database struct {
	URL string
}

type jwt struct {
	Secret string
	Issuer string
}

type Config struct {
	Database database
	JWT      jwt
}

func LoadEnv(fileName string) {
	dir := "go-twitter-api-graphql"
	re := regexp.MustCompile(`^(.*` + dir + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	fmt.Println(string(rootPath) + `/` + fileName)
	err := godotenv.Load(string(rootPath) + `/` + fileName)
	if err != nil {
		godotenv.Load()
	}
}

func New() *Config {
	return &Config{Database: database{
		URL: os.Getenv("DATABASE_URL"),
	}, JWT: jwt{
		Secret: os.Getenv(("JWT_SECRET")),
		Issuer: os.Getenv("DOMAIN"),
	}}
}
