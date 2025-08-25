package g3

import (
	"fmt"
	"io"
	"net/http"
)

type G3 struct {
	Server http.Server
	Addr   string
	routes map[string]func(*Request) (Response, error)
}

func (g *G3) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := fmt.Sprintf("%v", r.URL)
	if controller, ok := g.routes[method+path]; ok {
		response, err := controller(&Request{
			Method: r.Method,
			Path:   path,
		})
		if err != nil {
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
		return

	}

	w.WriteHeader(404)
	io.WriteString(w, "Not Found :)")
}

type Request struct {
	Method string
	Path   string
}

type Response struct {
	Body       []byte
	StatusCode int
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

func (g3 *G3) Get(path string, controller func(*Request) (Response, error)) *G3 {
	g3.routes["GET"+path] = controller
	return g3
}

func (g3 *G3) Post(path string, controller func(*Request) (Response, error)) *G3 {
	g3.routes["POST"+path] = controller
	return g3
}

func (g3 *G3) Put(path string, controller func(*Request) (Response, error)) *G3 {
	g3.routes["PUT"+path] = controller
	return g3
}

func (g3 *G3) Patch(path string, controller func(*Request) (Response, error)) *G3 {
	g3.routes["PATCH"+path] = controller
	return g3
}

func (g3 *G3) Delete(path string, controller func(*Request) (Response, error)) *G3 {
	g3.routes["Delete"+path] = controller
	return g3
}
