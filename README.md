# ⚡ G3

### microframework

## نصب

```bash
go get github.com/zhitoo/g3
```

```go
package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/zhitoo/g3"
)

type User struct {
	Name []string `form:"name"`
	Age  int      `form:"age"`
}

func main() {
	g := g3.New(":5500")

	g.Group("/users/{id}", func() {
		g.Use(func(next g3.Controller) g3.Controller {
			return func(r *g3.Request) (g3.Response, error) {
				fmt.Println("👉 New Request received From " + r.PathParams["id"])
				return next(r)
			}
		})
		g.Get("/", func(r *g3.Request) (g3.Response, error) {
			response := g3.Response{}
			response.Body = []byte("All Users")
			return response, nil
		})
	})

	g.Get("/", func(r *g3.Request) (g3.Response, error) {
		response := g3.Response{}
		response.Body = []byte("Hello World")
		fmt.Println("Query Param", r.Get("name"))
		return response, nil
	}).Post("/", func(r *g3.Request) (g3.Response, error) {
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
		}).AddValidation("age", func(r *g3.Request) (bool, string) {
			age, ok := r.PostParams["age"]
			if ok {
				ageNum, err := strconv.Atoi(age[0])
				if err != nil {
					return false, "age must be a valid number"
				}
				if ageNum < 12 {
					return false, "age must gt than 12"
				}
				if ageNum > 75 {
					return false, "age must st than 75"
				}
			} else {
				return false, "age is required"
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

	g.Use(func(next g3.Controller) g3.Controller {
		return func(r *g3.Request) (g3.Response, error) {
			fmt.Println("👉 New Request received")
			return next(r)
		}
	})

	g.Get("/hello/{name}", func(r *g3.Request) (g3.Response, error) {
		response := g3.Response{}

		response.Redirect("/users/"+r.PathParams["name"], 301)

		return response, nil
	})
	log.Fatal(g.Serve())

}


```
