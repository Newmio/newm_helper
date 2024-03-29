package newm_helper

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/Newmio/newm_logger"
)

func RenderHtml(directory string, data interface{}) (string, error) {
	buffer := new(strings.Builder)

	tmpl, err := template.ParseFiles(directory)
	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(buffer, data); err != nil {
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

	log := newm_logger.Log{
		Url:       param.Url,
		Method:    param.Method,
		RequestId: param.RequestId,
	}

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

	log.BodyReq = string(body)
	log.DateStart = time.Now().Format("2006-01-02 15:04:05")

	req, err := http.NewRequest(param.Method, param.Url, bytes.NewBuffer(body))
	if err != nil {
		return 500, nil, err
	}

	for key, value := range param.Headers {
		log.HeadersReq += key + ": " + value.(string)
		req.Header.Set(key, value.(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		if resp != nil {
			return 504, nil, err
		}

		log.Status = 404
		log.DateStop = time.Now().Format("2006-01-02 15:04:05")
		return 404, nil, err
	}
	defer resp.Body.Close()

	log.DateStop = time.Now().Format("2006-01-02 15:04:05")
	log.Status = resp.StatusCode

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return 500, nil, err
	}
	log.BodyResp = string(bodyResp)

	for key, values := range resp.Header {
		log.HeadersResp = fmt.Sprintf("%s: %v", key, values)
	}

	log.Success = true

	if param.CreateLog {
		newm_logger.CreateLog(log)
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
