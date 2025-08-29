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

	//middleware
	g.Use(func(next g3.Controller) g3.Controller {
		return func(r *g3.Request) (g3.Response, error) {
			fmt.Println("👉 New Request received")
			return next(r)
		}
	})

	//route group
	g.Group("/users", func() {
		g.Get("/", func(r *g3.Request) (g3.Response, error) {
			response := g3.Response{}
			response.Body = []byte("All Users")
			return response, nil
		})
	})

	g.Get("/g3/{id?:^[0-9]*$}", func(r *g3.Request) (g3.Response, error) {
		response := g3.Response{}
		response.Body = []byte("Hello, G3!")
		return response, nil

	}).Get("/", func(r *g3.Request) (g3.Response, error) {
		response := g3.Response{}
		response.Body = []byte("Hello, World!")
		return response, nil

	})

	//redirect
	g.Get("/hello/{name}", func(r *g3.Request) (g3.Response, error) {
		response := g3.Response{}

		response.Redirect("/" + r.PathParams["name"], 301)

		return response, nil
	})

	log.Fatal(server.Serve())
}

```
