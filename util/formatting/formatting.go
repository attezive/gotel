package formatting

import (
	"bufio"
	"fmt"
	"net/http"
)

func PrintResponse(resp *http.Response) {
	fmt.Println("response Status:", resp.Status)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
