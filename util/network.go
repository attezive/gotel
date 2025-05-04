package util

import (
	"bytes"
	"mime/multipart"
	"net/http"
)

const AddressUrl string = "https://api.telegram.org/bot"

func GetRequest(token string, command string, params map[string]string) (*http.Response, error) {
	sendParams := ""
	if params != nil {
		sendParams += "?"
		for k, v := range params {
			sendParams += k + "=" + v + "&"
		}
	}
	resp, err := http.Get(AddressUrl + token + "/" + command + sendParams)
	return resp, err
}

func PostRequest(token string, command string, params map[string]string, body *bytes.Buffer, writer *multipart.Writer) (*http.Response, error) {
	sendParams := ""
	if params != nil {
		sendParams += "?"
		for k, v := range params {
			sendParams += k + "=" + v + "&"
		}
	}
	req, _ := http.NewRequest("POST", AddressUrl+token+"/"+command+sendParams, body)
	if writer != nil {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	client := &http.Client{}
	return client.Do(req)
}
