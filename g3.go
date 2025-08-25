package g3

import (
	"fmt"
	"net/http"
	"strings"
)

type G3 struct {
	Server http.Server
	Addr   string
	routes map[string]func(*Request) (Response, error)
}

type Controller func(*Request) (Response, error)

func (g3 *G3) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	controller, err := g3.getController(r)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(fmt.Sprintf("%v", err)))
		return
	}

	fmt.Printf("controller %v", controller)

	response, err := controller(&Request{
		Method: r.Method,
		Path:   r.URL.Path,
	})

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

func (g3 *G3) addRoute(method, path string, controller Controller) {
	g3.routes[method+":"+path] = controller
}

func (g3 *G3) getController(r *http.Request) (Controller, error) {
	fmt.Printf("%v\n", r.Form)
	method := r.Method
	path := r.URL.Path
	fmt.Printf("%v\n", path)
	fmt.Printf("%v\n", g3.routes)

	//find route with path
	pathParts := strings.Split(path, "/")
	fmt.Println(pathParts)

	findRoute := ""

	for route := range g3.routes {
		routeParts := strings.Split(route, ":")
		//check method
		fmt.Println("route parts 0:", routeParts[0])
		if routeParts[0] != method {
			continue
		}
		routeWithoutMethod := routeParts[1]
		routeWithoutMethodParts := strings.Split(routeWithoutMethod, "/")
		fmt.Println("other route parts :", routeParts)
		fmt.Println("pathParts :", pathParts)

		if len(pathParts) == len(routeWithoutMethodParts) {
			findRoute = route
		}
	}

	//todo: we should set PathParams if exists

	controller, ok := g3.routes[findRoute]
	fmt.Printf("%v\n", g3.routes)
	if ok {
		return controller, nil
	}
	return nil, fmt.Errorf("not found :)")

}

type Request struct {
	Method     string
	Path       string
	PathParams map[string]string
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

func (g3 *G3) Get(path string, controller Controller) *G3 {
	g3.addRoute("GET", path, controller)
	return g3
}

func (g3 *G3) Post(path string, controller Controller) *G3 {
	g3.addRoute("POST", path, controller)
	return g3
}

func (g3 *G3) Put(path string, controller Controller) *G3 {
	g3.addRoute("PUT", path, controller)
	return g3
}

func (g3 *G3) Patch(path string, controller Controller) *G3 {
	g3.addRoute("PATCH", path, controller)
	return g3
}

func (g3 *G3) Delete(path string, controller Controller) *G3 {
	g3.addRoute("Delete", path, controller)
	return g3
}
