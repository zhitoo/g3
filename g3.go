package g3

import (
	"fmt"
	"net/http"
)

type G3 struct {
	Server http.Server
	Addr   string
	routes map[string]func(*Request) (Response, error)
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
		//todo: check for validation error or any other type of error
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	w.Header().Add("accept", r.Header.Get("accept"))
	statusCode := response.StatusCode
	if statusCode == 0 {
		statusCode = 200
	}
	w.WriteHeader(statusCode)
	w.Write(response.Body)
}
