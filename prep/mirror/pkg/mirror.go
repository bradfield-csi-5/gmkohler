package pkg

import (
	"fmt"
	"io"
	"net/http"
)

func Mirror(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching url %s: %v\n", url, err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v\n", err)
	}
	fmt.Printf("%s\n", body)
	return nil
}
