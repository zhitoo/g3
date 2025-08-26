# ⚡ G3

### microframework

## نصب

```bash
go get github.com/zhitoo/g3
```

```go
package main

import (
	"log"

	"github.com/zhitoo/g3"
)

func main() {
	server := g3.New(":5500")
	server.Get("/g3/{id?:^[0-9]*$}", func(r *g3.Request) (g3.Response, error) {
		response := g3.Response{}
		response.Body = []byte("Hello, G3!")
		return response, nil

	}).Get("/", func(r *g3.Request) (g3.Response, error) {
		response := g3.Response{}
		response.Body = []byte("Hello, World!")
		return response, nil

	})
	log.Fatal(server.Serve())
}

```
