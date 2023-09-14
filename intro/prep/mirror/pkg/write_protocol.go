package pkg

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
)

// writeProtocol lets us write either response body or parsed HTML tree without
// duplicating the filename-making-code and without coming up with some way
// to read a body twice
type writeProtocol func(f *os.File) error

func httpResponseWriteProtocol(resp *http.Response) writeProtocol {
	return func(f *os.File) error {
		bytesWritten, err := io.Copy(f, resp.Body)
		if err != nil {
			return err
		}
		if bytesWritten < resp.ContentLength {
			return fmt.Errorf("error copying all response bytes to file")
		}
		return nil
	}
}
func htmlNodeWriteProtocol(root *html.Node) writeProtocol {
	return func(f *os.File) error {
		return html.Render(f, root)
	}
}
