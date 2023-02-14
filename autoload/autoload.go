package autoload

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/oxplot/starenv"
	"github.com/oxplot/starenv/derefer"
)

func init() {
	if os.Getenv("DOTENV_ENABLED") != "0" {
		_ = godotenv.Load()
	}

	for t, n := range derefer.NewDefault {
		starenv.Register(t, &derefer.Lazy{New: n})
	}

	if err := starenv.Load(); err != nil {
		log.Fatal("starenv.autoload: ", err)
	}
}
