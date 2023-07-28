package pkg

import (
	"fmt"
	"golang.org/x/net/html"
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
	// in theory, we could check content type first, but I think you have to
	// read the body to do that.
	domBody, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf(
			"error parsing response body as HTML: %v\n",
			err,
		)
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
	// TODO: detect existing files to avoid overwriting / circular graph issues
	f, err := os.Create(fp)
	if err != nil {
		return fmt.Errorf("error opening file %q: %v\n", fp, err)
	}

	err = html.Render(f, domBody)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v\n", fp, err)
	}

	return nil
}
