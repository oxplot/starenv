package autoload

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/oxplot/starenv"
)

func init() {
	if os.Getenv("DOTENV_ENABLED") != "0" {
		if err := godotenv.Load(); err != nil {
			log.Printf("starenv.autoload: %s", err)
		}
	}

	if errs := starenv.DefaultLoader.Load(); errs != nil {
		for _, err := range errs {
			log.Printf("starenv.autoload: %s", err)
		}
	}
}
