package g3

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	g := New(":5500")

	g.Post("/register", func(r *Request) (Response, error) {
		res := Response{}

		r.AddValidation("name", func(r *Request) (bool, string) {
			value, ok := r.PostParams["name"]
			if ok {
				if len(value[0]) < 5 {
					return false, "name must at least be 5 chars"
				}
			} else {
				return false, "name is required"
			}
			return true, ""
		}).AddValidation("age", func(r *Request) (bool, string) {
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

		res.Body = []byte(":)")
		return res, nil

	})

	// ✅
	form := strings.NewReader("name=Hossein&age=30")
	req := httptest.NewRequest("POST", "/register", form)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	g.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// ❌
	formErr := strings.NewReader("name=ali&age=30")
	reqErr := httptest.NewRequest("POST", "/register", formErr)
	reqErr.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	wErr := httptest.NewRecorder()
	g.ServeHTTP(wErr, reqErr)
	if wErr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", wErr.Code)
	}

}
