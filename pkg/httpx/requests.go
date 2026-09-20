package httpx

import (
	"net/http"
	"net/textproto"
	"strings"

	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

const (
	formKey           = "form"
	pathKey           = "path"
	headerKey         = "header"
	maxMemory         = 32 << 20 // 32MB
	separator         = ";"
	tokensInAttribute = 2
)

var (
	formUnmarshaler   = mapping.NewUnmarshaler(formKey, mapping.WithStringValues())
	pathUnmarshaler   = mapping.NewUnmarshaler(pathKey, mapping.WithStringValues())
	headerUnmarshaler = mapping.NewUnmarshaler(headerKey, mapping.WithStringValues(),
		mapping.WithCanonicalKeyFunc(textproto.CanonicalMIMEHeaderKey))
)

// Parse parses the request form, path, header, and json body.
func Parse(r *http.Request, v interface{}) error {
	if err := ParsePath(r, v); err != nil {
		return err
	}

	if err := ParseForm(r, v); err != nil {
		return err
	}

	if err := ParseHeaders(r, v); err != nil {
		return err
	}

	return httpx.ParseJsonBody(r, v)
}

// ParsePath parses the path parameters.
func ParsePath(r *http.Request, v interface{}) error {
	vars := pathvar.Vars(r)
	m := make(map[string]interface{}, len(vars))
	for k, val := range vars {
		m[k] = val
	}
	return pathUnmarshaler.Unmarshal(m, v)
}

// ParseHeaders parses the request headers.
func ParseHeaders(r *http.Request, v interface{}) error {
	m := map[string]interface{}{}
	for k, val := range r.Header {
		if len(val) == 1 {
			m[k] = val[0]
		} else {
			m[k] = val
		}
	}
	return headerUnmarshaler.Unmarshal(m, v)
}

// ParseForm parses the query and form parameters.
func ParseForm(r *http.Request, v interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}

	params := make(map[string]interface{}, len(r.Form))
	for name := range r.Form {
		formValue := r.Form.Get(name)
		if len(formValue) > 0 {
			params[name] = formValue
		}
	}

	return formUnmarshaler.Unmarshal(params, v)
}
