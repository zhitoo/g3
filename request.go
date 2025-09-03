package g3

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
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

func setField(field reflect.Value, value string) error {
	if !field.CanSet() {
		return fmt.Errorf("cannot set field")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int64:
		if val, err := strconv.ParseInt(value, 10, 64); err == nil {
			field.SetInt(val)
		} else {
			return err
		}
	default:
		return fmt.Errorf("unsupported field type: %v", field.Kind())
	}
	return nil
}

func (rg *Request) bindFormParams(obj any) error {
	fmt.Println("bindFormParams", reflect.ValueOf(obj))
	fmt.Println("PostParams:", rg.PostParams)
	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("form")
		if tag == "" {
			tag = field.Name
		}
		if value, exists := rg.PostParams[tag]; exists && len(value) > 0 {
			fv := val.Field(i)
			fmt.Println("fv", fv)

			if fv.Kind() == reflect.Slice {
				sliceType := fv.Type().Elem()
				slice := reflect.MakeSlice(fv.Type(), 0, len(value))

				for _, v := range value {
					elem := reflect.New(sliceType).Elem()
					if err := setField(elem, v); err != nil {
						return fmt.Errorf("Bind: failed to set slice param %s: %v", tag, err)
					}
					slice = reflect.Append(slice, elem)
				}

				fv.Set(slice)
			} else {
				if err := setField(fv, value[0]); err != nil {
					return fmt.Errorf("Bind: failed to set form param %s: %v", tag, err)
				}
			}
		}
	}
	return nil
}

func (rg *Request) bindQueryParams(obj any) error {
	val := reflect.ValueOf(obj).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("form")
		if tag == "" {
			tag = field.Name
		}
		if value, exists := rg.QueryParams[tag]; exists && len(value) > 0 {
			if err := setField(val.Field(i), value[0]); err != nil {
				return fmt.Errorf("Bind: failed to set query param %s: %v", tag, err)
			}
		}
	}
	return nil
}

func (r *Request) Bind(obj any) error {
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("Bind: input must be a non-nil pointer")
	}
	if val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Bind: input must be a pointer to a struct")
	}

	if err := r.bindQueryParams(obj); err != nil {
		return err
	}

	if r.PostParams != nil {
		if err := r.bindFormParams(obj); err != nil {
			return err
		}
	}

	return nil
}
