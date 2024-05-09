package newm_helper

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"
)

func СontainsSQLInjection(query string) bool {
	sqlInjectionPattern := regexp.MustCompile(`(?i)\b(SELECT|UPDATE|DELETE|FROM|WHERE|DROP|UNION|AND|OR)\b`)
	return sqlInjectionPattern.MatchString(query)
}
func RenderHtml(directory string, data interface{}) (string, error) {
	buffer := new(strings.Builder)

	funcMap := template.FuncMap{
		"add": func(x, y int) int {
			return x + y
		},

		"idForStr": func(str string) string {
			hash := md5.New()
			hash.Write([]byte(str))
			return "x" + hex.EncodeToString(hash.Sum(nil))
		},
	}

	name := filepath.Base(directory)

	tmpl, err := template.New(name).Funcs(funcMap).ParseFiles(directory)
	if err != nil {
		return "", err
	}

	if err := tmpl.ExecuteTemplate(buffer, name, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func ErrorResponse(err string) map[string]string {
	return map[string]string{
		"status":      "error",
		"description": err,
	}
}

func RequestHTTP(param Param) (int, []byte, error) {
	var body []byte

	param.Url = strings.Replace(param.Url, " ", "%20", -1)
	param.Url = strings.Replace(param.Url, "+", "%2B", -1)

	client := &http.Client{}

	if param.BodyType == "JSON" {

		b, err := json.Marshal(param.Body)
		if err != nil {
			return 500, nil, err
		}
		body = b

	} else if param.BodyType == "XML" {

		b, err := xml.Marshal(param.Body)
		if err != nil {
			return 500, nil, err
		}
		body = b

	} else {
		body = nil
	}

	req, err := http.NewRequest(param.Method, param.Url, bytes.NewBuffer(body))
	if err != nil {
		return 500, nil, err
	}

	for key, value := range param.Headers {
		req.Header.Set(key, value.(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		if resp != nil {
			return 504, nil, err
		}

		return 404, nil, err
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return 500, nil, err
	}

	return resp.StatusCode, bodyResp, nil
}

func Trace(err error, any ...interface{}) error {
	if err == nil {
		return nil
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return err
	}

	var str string

	for _, value := range any {
		str += fmt.Sprint(value)
	}

	return fmt.Errorf("%s%s%s%s(*_*) %s:%d (*_*)", err.Error(), "\n", str, "\n", file, line)
}
