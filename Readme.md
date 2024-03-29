# Starenv

[![Go
Reference](https://pkg.go.dev/badge/github.com/oxplot/starenv.svg)](https://pkg.go.dev/github.com/oxplot/starenv)

**starenv** (`*env`) allows populating environmental variables from
variety of sources, such as AWS Parameter Store, GPG encrypted files
and more, with extreme ease.

For the impatient, import the `autoload` package and you're set:

```go
package main

import (
  "fmt"
  "os"

  _ "github.com/oxplot/starenv/autoload"
)

func main() {
  fmt.Printf("GITHUB_TOKEN=%s\n", os.Getenv("GITHUB_TOKEN"))
}
```

and set the value of the environmental variable to load from Parameter
Store:

```sh
$ export GITHUB_TOKEN=*ssm:/github_token
$ go run main.go
GITHUB_TOKEN=abcdef-1235143-abcdef-123-abcdef-12314
```

or from a GPG encrypted file:

```sh
$ export GITHUB_TOKEN=*gpg*file:github_token.gpg
$ go run main.go
GITHUB_TOKEN=abcdef-1235143-abcdef-123-abcdef-12314
```

why not ditch the file and embed its content:

```sh
$ export GITHUB_TOKEN=*gpg*b64:eNeO7D2rBrBOOcW6TuETyHdyPEOaAfdgaTzgOTSvROI=
$ go run main.go
GITHUB_TOKEN=abcdef-1235143-abcdef-123-abcdef-12314
```

and thanks to the amazing [godotenv](https://github.com/joho/godotenv)
which is run as part of starenv's `autoload` package, you can even do:

```sh
$ echo 'GITHUB_TOKEN=*keyring:awesome_app/github_token' > .env
$ go run main.go
GITHUB_TOKEN=abcdef-1235143-abcdef-123-abcdef-12314
```

For a full list, see [the
docs](https://pkg.go.dev/github.com/oxplot/starenv/derefer#NewDefault).
