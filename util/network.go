package util

import "net/http"

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
