package autoload

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/oxplot/starenv"
	"github.com/oxplot/starenv/derefer"
)

const failEnv = "STARENV_AUTOLOAD_FAIL"

func init() {
	_ = godotenv.Load()

	for t, n := range derefer.NewDefault {
		starenv.Register(t, &derefer.Lazy{New: n})
	}

	fail := len(os.Getenv(failEnv)) > 0
	os.Unsetenv(failEnv)

	if err := starenv.Load(); err != nil && fail {
		log.Fatal("starenv.autoload: ", err)
	}
}
