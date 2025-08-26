package g3

type Request struct {
	Method      string
	Path        string
	PathParams  map[string]string
	QueryParams map[string]string
}

func (r *Request) Get(name string) string {
	if value, ok := r.QueryParams[name]; ok {
		return value
	}
	return ""
}
