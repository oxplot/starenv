package autoload

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"github.com/oxplot/starenv"
	"github.com/oxplot/starenv/derefer"
)

const blacklistEnv = "STARENV_AUTOLOAD_BLACKLIST"
const failEnv = "STARENV_AUTOLOAD_FAIL"

func init() {
	_ = godotenv.Load()

	fail := len(os.Getenv(failEnv)) > 0
	s := strings.Split(os.Getenv(blacklistEnv), ",")
	os.Unsetenv(blacklistEnv)
	os.Unsetenv(failEnv)
	blacklist := make(map[string]struct{}, len(s))
	for _, t := range s {
		blacklist[t] = struct{}{}
	}
	for t, init := range derefer.NoConfigInit {
		if _, ok := blacklist[t]; !ok {
			d, err := init()
			if err != nil {
				if fail {
					log.Fatal("starenv.autoload: ", err)
				}
				continue
			}
			starenv.Register(t, d)
		}
	}
	if err := starenv.Load(); err != nil && fail {
		log.Fatal("starenv.autoload: ", err)
	}
}
