package pkg

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// Mirror writes the web page at url to outDir
func Mirror(webPage string, outDir string) error {
	resp, err := http.Get(webPage)
	if err != nil {
		return fmt.Errorf("error fetching url %s: %v\n", webPage, err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v\n", err)
	}
	// TODO: decide how to handle "directory exists".  For now,
	//  we fail (simple), but maybe we could just remove its contents
	//  (we can check os.IsExist(err)).
	err = os.Mkdir(outDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf(
			"error creating directory %q: %v\n",
			outDir,
			err,
		)
	}

	u, err := url.Parse(webPage)
	if err != nil {
		return fmt.Errorf("error parsing url %s: %v\n", webPage, err)
	}
	// we will put the index at "index.html"
	var fileName string
	if u.Path == "" {
		fileName = "index.html"
	}
	fp := filepath.Join(outDir, u.Path, fileName)
	err = os.WriteFile(fp, body, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v\n", fp, err)
	}

	return nil
}
