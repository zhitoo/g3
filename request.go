package g3

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Method      string
	Path        string
	PathParams  map[string]string
	QueryParams map[string][]string
	PostParams  map[string][]string
}

func (r *Request) Get(name string) string {
	if value, ok := r.QueryParams[name]; ok {
		return value[0]
	}
	return ""
}

func (r *Request) Post(name string) string {
	if value, ok := r.PostParams[name]; ok {
		return value[0]
	}
	return ""
}

func (r *Request) Input(name string) string {
	value := r.Post(name)
	if value == "" {
		value = r.Get(name)
	}
	return value
}

func (r *Request) Array(name string) []string {
	value := r.PostParams[name]
	if value == nil {
		value = r.QueryParams[name]
	}
	return value
}

func (r *Request) Has(key string) bool {
	_, ok := r.PostParams[key]
	if !ok {
		_, ok = r.QueryParams[key]
	}
	return ok
}

func (gr *Request) setQueryParams(r *http.Request) error {
	query := r.URL.Query()
	fmt.Println("Query Params:", query)
	queryParams := map[string][]string{}
	for index, value := range query {
		queryParams[index] = value
	}
	gr.QueryParams = queryParams
	return nil
}

func (gr *Request) setPostForm(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	postForm := map[string][]string{}
	if contentType == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&postForm)
		if err != nil {
			return err
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form: %v", err)
		}
		for key, values := range r.PostForm {
			postForm[key] = values
		}
	}
	gr.PostParams = postForm

	return nil
}
