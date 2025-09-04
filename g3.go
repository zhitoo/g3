package g3

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type group struct {
	Prefix      string
	Middlewares []Middleware
	Regex       *regexp.Regexp
	ParamNames  []string
}

type G3 struct {
	Server      http.Server
	Addr        string
	routes      map[string]func(*Request) (Response, error)
	path_prefix string
	groups      []group
	middlewares []Middleware
}

func New(Addr string) *G3 {
	fmt.Println("Creating...")
	g3 := G3{}
	g3.Addr = Addr
	g3.Server.Handler = &g3
	g3.routes = map[string]func(*Request) (Response, error){}

	return &g3
}

func (g3 *G3) Serve() error {
	g3.Server.Addr = g3.Addr
	return g3.Server.ListenAndServe()
}

func (g3 *G3) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	response, err := g3.runController(r)
	if err != nil {
		var validationError ValidationError
		if errors.As(err, &validationError) {
			jsonResp, er := json.Marshal(validationError)
			if er != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("%v", er)))
			}
			response.Body = jsonResp
			w.WriteHeader(response.statusCode)
			w.Write(response.Body)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	statusCode := response.statusCode

	//add response headers to w
	for key, value := range response.header {
		w.Header().Add(key, value)
	}

	if statusCode > 300 && statusCode < 399 {
		if location, ok := response.header["Location"]; ok {
			println(location)
			http.Redirect(w, r, location, statusCode) // 302
			return
		}
	}

	w.Header().Add("accept", r.Header.Get("accept"))

	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	w.WriteHeader(statusCode)
	w.Write(response.Body)
}
