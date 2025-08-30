package g3

import "net/http"

type Response struct {
	Body       []byte
	statusCode int
	header     map[string]string
}

func (r *Response) SetStatusCode(code int) *Response {
	r.statusCode = code
	return r
}

func (r *Response) SetBody(body []byte) *Response {
	r.Body = body
	return r
}

func NewResponse() *Response {
	response := Response{}
	response.SetStatusCode(http.StatusOK)
	return &Response{}
}

func (r *Response) String(body string) Response {
	r.Body = []byte(body)
	r.SetHeader("Content-Type", "text/plain")
	return *r
}

func (r *Response) JSON(body []byte) Response {
	r.Body = body
	r.SetHeader("Content-Type", "application/json")
	return *r
}

func (r *Response) HTML(body string) Response {
	r.Body = []byte(body)
	r.SetHeader("Content-Type", "text/html")
	return *r
}

func (r *Response) XML(body string) Response {
	r.Body = []byte(body)
	r.SetHeader("Content-Type", "application/xml")
	return *r
}

func (r *Response) SetHeader(key, value string) *Response {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header[key] = value
	return r
}

func (r *Response) GetHeader(key string) string {
	if r.header == nil {
		return ""
	}
	return r.header[key]
}

func (r *Response) DelHeader(key string) {
	if r.header == nil {
		return
	}
	delete(r.header, key)
}

func (r *Response) ClearHeaders() *Response {
	if r.header == nil {
		return r
	}
	r.header = make(map[string]string)
	return r
}

func (r *Response) Headers() map[string]string {
	if r.header == nil {
		return make(map[string]string)
	}
	return r.header
}

func (r *Response) Status() int {
	if r.statusCode == 0 {
		return 200
	}
	return r.statusCode
}

func (r *Response) Redirect(url string, code int) *Response {
	if code < 300 || code > 399 {
		code = http.StatusFound
	}
	r.SetStatusCode(code)
	r.SetHeader("Location", url)
	return r
}
