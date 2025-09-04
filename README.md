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
	g.Post("/", func(r *g3.Request) (g3.Response, error) {
		res := g3.Response{}

		r.AddValidation("name", func(r *g3.Request) (bool, string) {
			value, ok := r.PostParams["name"]
			if ok {
				if len(value[0]) < 5 {
					return false, "name must at least be 5 chars"
				}
			} else {
				return false, "name is required"
			}
			return true, ""
		})

		err := r.Validate()
		if err != nil {
			res.SetStatusCode(422)
			return res, err
		}

		user := User{}
		err = r.Bind(&user)
		if err != nil {
			res.Body = []byte(err.Error())
			res.SetStatusCode(400)
			return res, nil
		}
		fmt.Printf("User: %+v\n", user)

		res.Body = []byte("Hello World!!! POST")
		res.SetStatusCode(201)
		return res, nil
	})

	log.Fatal(server.Serve())
}

```
