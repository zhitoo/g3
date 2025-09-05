package g3

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRoute(t *testing.T) {
	g := New(":5500")

	g.Get("/hello", func(r *Request) (Response, error) {
		response := NewResponse()
		return response.String("world")
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if body := w.Body.String(); body != "world" {
		t.Errorf("expected body 'world', got '%s'", body)
	}
}

func TestRouteWithParam(t *testing.T) {
	g := New(":5500")

	g.Get("/user/{id:[0-9]+}", func(r *Request) (Response, error) {
		return NewResponse().String("id=" + r.PathParams["id"])
	})

	req := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Body.String() != "id=123" {
		t.Errorf("expected 'id=123', got '%s'", w.Body.String())
	}

	g.Get("/user/{id:[0-9]+}", func(r *Request) (Response, error) {
		return NewResponse().String("id=" + r.PathParams["id"])
	})

	req = httptest.NewRequest("GET", "/user/test", nil)
	w = httptest.NewRecorder()

	g.ServeHTTP(w, req)

	if w.Body.String() != "not found :)" {
		t.Errorf("expected 'id=123', got '%s'", w.Body.String())
	}
}

func TestMiddleware(t *testing.T) {
	g := New(":8080")

	called := false
	g.Use(func(next Controller) Controller {
		return func(r *Request) (Response, error) {
			called = true
			return next(r)
		}
	})

	g.Get("/ping", func(r *Request) (Response, error) {
		return NewResponse().String("world")
	})

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)

	if !called {
		t.Errorf("expected middleware to be called")
	}
}
