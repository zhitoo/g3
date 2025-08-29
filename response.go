package g3

type Response struct {
	Body       []byte
	StatusCode int
	Header     map[string]string
}

func (r *Response) SetStatusCode(code int) {
	r.StatusCode = code
}

func (r *Response) SetBody(body []byte) {
	r.Body = body
}

func NewResponse() Response {
	return Response{
		StatusCode: 200,
	}
}

func (r *Response) String(body string) {
	r.Body = []byte(body)
	r.SetHeader("Content-Type", "text/plain")
}

func (r *Response) JSON(body []byte) {
	r.Body = body
	r.SetHeader("Content-Type", "application/json")
}

func (r *Response) HTML(body string) {
	r.Body = []byte(body)
	r.SetHeader("Content-Type", "text/html")
}

func (r *Response) XML(body string) {
	r.Body = []byte(body)
	r.SetHeader("Content-Type", "application/xml")
}

func (r *Response) SetHeader(key, value string) {
	if r.Header == nil {
		r.Header = make(map[string]string)
	}
	r.Header[key] = value
}

func (r *Response) GetHeader(key string) string {
	if r.Header == nil {
		return ""
	}
	return r.Header[key]
}

func (r *Response) DelHeader(key string) {
	if r.Header == nil {
		return
	}
	delete(r.Header, key)
}

func (r *Response) ClearHeaders() {
	if r.Header == nil {
		return
	}
	r.Header = make(map[string]string)
}

func (r *Response) Headers() map[string]string {
	if r.Header == nil {
		return make(map[string]string)
	}
	return r.Header
}

func (r *Response) Status() int {
	if r.StatusCode == 0 {
		return 200
	}
	return r.StatusCode
}
