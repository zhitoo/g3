package g3

import (
	"encoding/json"
	"net/http"
)

type Request struct {
	Method      string
	Path        string
	PathParams  map[string]string
	QueryParams map[string]string
	PostParams  map[string]any
}

func (r *Request) Get(name string) string {
	if value, ok := r.QueryParams[name]; ok {
		return value
	}
	return ""
}

func (r *Request) Post(name string) any {
	if value, ok := r.PostParams[name]; ok {
		return value
	}
	return nil
}

func (r *Request) Input(name string) any {
	value := r.Post(name)
	if value == nil {
		value = r.Get(name)
	}
	return value
}

func (gr *Request) setQueryParams(r *http.Request) error {
	query := r.URL.Query()
	queryParams := map[string]string{}
	for index, value := range query {
		queryParams[index] = value[0]
	}
	gr.QueryParams = queryParams
	return nil
}

func (gr *Request) setPostForm(r *http.Request) error {
	contentType := r.Header.Get("Content-Type")
	postForm := map[string]any{}
	if contentType == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&postForm)
		if err != nil {
			return err
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			return err
		}
		for index, value := range r.PostForm {
			postForm[index] = value[0]
		}
	}
	gr.PostParams = postForm

	return nil
}
