package g3

import (
	"fmt"
	"net/http"
	"regexp"
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

	key := method + ":" + path

	g3.routes[key] = controller

	parts := strings.Split(path, "/")

	for i, part := range parts {

		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			param := part[1 : len(part)-1]
			paramParts := strings.SplitN(param, ":", 2)
			paramNamePart := paramParts[0]

			if strings.HasSuffix(paramNamePart, "?") {

				optionalParts := append([]string{}, parts[:i]...)
				optionalParts = append(optionalParts, parts[i+1:]...)
				optionalPath := strings.Join(optionalParts, "/")

				if optionalPath == "" {
					optionalPath = "/"
				}

				optionalKey := method + ":" + optionalPath
				
				if _, exists := g3.routes[optionalKey]; !exists {
					g3.routes[optionalKey] = controller
				}
			}
		}
	}
}

func (g3 *G3) getController(r *http.Request) (Controller, error) {
	method := r.Method
	path := r.URL.Path
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	var matchedRoute string
	var pathParams map[string]string

	for route := range g3.routes {
		routeParts := strings.SplitN(route, ":", 2)
		if len(routeParts) != 2 {
			continue
		}

		routeMethod := routeParts[0]
		routePattern := routeParts[1]

		if routeMethod != method {
			continue
		}

		routePatternParts := strings.Split(strings.Trim(routePattern, "/"), "/")
		if len(pathParts) != len(routePatternParts) {
			continue
		}

		tmpParams := map[string]string{}
		matched := true

		for i, rp := range routePatternParts {
			pp := pathParts[i]

			if strings.HasPrefix(rp, "{") && strings.HasSuffix(rp, "}") {
				paramName := strings.Trim(rp, "{}")

				// regex {id:[0-9]+}
				if strings.Contains(paramName, ":") {
					parts := strings.SplitN(paramName, ":", 2)
					key := parts[0]
					regex := parts[1]

					if ok, _ := regexp.MatchString("^"+regex+"$", pp); !ok {
						matched = false
						break
					}
					tmpParams[key] = pp
				} else {
					tmpParams[paramName] = pp
				}
			} else {

				if rp != pp {
					matched = false
					break
				}
			}
		}

		if matched {
			matchedRoute = route
			pathParams = tmpParams
			break
		}
	}

	if matchedRoute == "" {
		return nil, fmt.Errorf("not found :)")
	}

	controller := g3.routes[matchedRoute]

	rq := &Request{
		Method:     r.Method,
		Path:       r.URL.Path,
		PathParams: pathParams,
	}
	return func(*Request) (Response, error) {
		return controller(rq)
	}, nil
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
