**starenv** (`*env`) allows populating environmental variables from
variety of sources, such as AWS Parameter Store, Google Cloud Secrets
Manager and more, with extreme ease.

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

```
$ export GITHUB_TOKEN=*ssm:arn:aws:ssm:us-west-2:123456123456:parameter/github_token
$ go run main.go
GITHUB_TOKEN=abcdef-1235143-abcdef-123-abcdef-12314
```

or from a GPG encrypted file:

```
$ export GITHUB_TOKEN=*gpg:*file:github_token.gpg
$ go run main.go
GITHUB_TOKEN=abcdef-1235143-abcdef-123-abcdef-12314
```

why not ditch the file and embed the file:

```
$ export GITHUB_TOKEN=*gpg:*b64:eNeO7D2rBrBOOcW6TuETyHdyPEOaAfdgaTzgOTSvROI=
$ go run main.go
GITHUB_TOKEN=abcdef-1235143-abcdef-123-abcdef-12314
```

and thanks to the amazing [godotenv](https://github.com/joho/godotenv)
which is run as part of starenv's `autoload` package, you can even do:

```
$ echo 'GITHUB_TOKEN=*gpg:*file:github_token.gpg' > .env
$ go run main.go
GITHUB_TOKEN=abcdef-1235143-abcdef-123-abcdef-12314
```
