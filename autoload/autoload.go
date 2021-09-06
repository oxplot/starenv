package autoload

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/oxplot/starenv"
	"github.com/oxplot/starenv/derefer"
)

const ignoreErrEnv = "STARENV_IGNORE_ERR"

func init() {
	_ = godotenv.Load()

	for t, n := range derefer.NewDefault {
		starenv.Register(t, &derefer.Lazy{New: n})
	}

	ignoreErr := len(os.Getenv(ignoreErrEnv)) > 0
	os.Unsetenv(ignoreErrEnv)

	if err := starenv.Load(); err != nil && !ignoreErr {
		log.Fatal("starenv.autoload: ", err)
	}
}
