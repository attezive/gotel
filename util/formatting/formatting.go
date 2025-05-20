package formatting

import (
	"bufio"
	"fmt"
	"github.com/attezive/gotel/data"
	"net/http"
)

func PrintResponse(resp *http.Response) {
	fmt.Println("response Status:", resp.Status)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func UpdateF(handleFunc func()) func(*data.Update) {
	return func(update *data.Update) {
		handleFunc()
	}
}

func MessageF(handleFunc func()) func(*data.Message) {
	return func(update *data.Message) {
		handleFunc()
	}
}
