package pkg

import (
	"fmt"
	"golang.org/x/net/html"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

const (
	headerContentType = "Content-Type"
	attrHref          = "href"
	attrSrc           = "src"
	contentTypeHtml   = "text/html;"
	hrefMailto        = "mailto:"
	tagAnchor         = "a"
	tagLink           = "link"
	tagScript         = "script"
)

// Mirror writes the web page at u to data.outDir, and if the web page is HTML,
// writes all links of the same domain to data.outDir
//
// TODO: control access to data from one goroutine when we make this concurrent
func Mirror(u *url.URL, data *MirrorData) error {
	if data.Seen(u) {
		return nil
	}
	resp, err := fetchPage(u)
	if err != nil {
		return fmt.Errorf("error fetching web page %v: %v\n", u, err)
	}
	data.MarkSeen(u)
	contentType := resp.Header.Get(headerContentType)
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || exts == nil {
		return fmt.Errorf(
			"error determining file extension for type %q: %v\n",
			contentType,
			err,
		)
	}
	var writePage writeProtocol = httpResponseWriteProtocol(resp)
	var links []*url.URL
	if strings.HasPrefix(contentType, contentTypeHtml) {
		// this logic is specific for HTML pages — we need to parse the document
		// and comb for more links
		domBody, err := parseHtml(resp)
		if err != nil {
			return fmt.Errorf(
				"error parsing response body as html: %v\n",
				err,
			)
		}
		writePage = htmlNodeWriteProtocol(domBody)
		findLinks(domBody, u, &links, data)
		fmt.Printf("found %d links\n", len(links))
	}

	dir, fName := makePath(data.outDir, u.Path, exts[0])
	fPath := path.Join(dir, fName)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %q: %v\n", dir, err)
	}

	f, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("error opening file %q: %v\n", fPath, err)
	}
	err = writePage(f)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v\n", fPath, err)
	}

	for _, l := range links {
		// many links are relative.  Here we infer from previous link the
		// host and scheme so that we can successfully fetch it.
		if l.Host == "" {
			l.Host = u.Host
		}
		if l.Scheme == "" {
			l.Scheme = u.Scheme
		}
		err = Mirror(l, data)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"error mirroring page %v: %v\n",
				l,
				err,
			)
			continue
		}
	}

	return nil
}

// makePath builds a file path from the target directory outDir and a URL
// path uPath.  The function signature is intended to be similar to path.Split,
// which this function leans on.
// examples:
// 		makePath("./mirrored", "", ".html") ("/mirrored", "index.html")
// 		makePath("./mirrored", "/courses/ssba", .html") (
//			"/mirrored/courses",
//			"ssba.html",
//		)

// TODO: unit testing
func makePath(outDir string, uPath string, ext string) (
	dir string,
	fileName string,
) {
	dir, fileName = path.Split(uPath)
	// FIXME: only name the root-root "index" otherwise take the last of the
	//  directory and append ext thereon
	if fileName == "" {
		fileName = "index"
	}
	if path.Ext(fileName) == "" {
		fileName += ext
	}
	dir = path.Join(outDir, dir)
	return
}

// fetchPage fetches u and parses it as an HTML document, returning the root or
// any errors it encounters
func fetchPage(u *url.URL) (*http.Response, error) {
	// fix relative links
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("error fetching url %v: %v\n", u, err)
	}

	return resp, nil
}

func parseHtml(r *http.Response) (*html.Node, error) {
	domBody, err := html.Parse(r.Body)
	if err != nil {
		err = fmt.Errorf(
			"error parsing response body as HTML: %v\n",
			err,
		)
		return nil, err
	}
	return domBody, nil
}

// findLinks looks for links within n with the specified domain
func findLinks(
	n *html.Node,
	hostUrl *url.URL,
	urls *[]*url.URL,
	data *MirrorData,
) {
	if n.Type == html.ElementNode {
		var (
			u   *url.URL
			err error
		)
		// links can be links to other HTML documents (<a>), to JS (<script>),
		// or to other assets (<link>)
		switch n.Data {
		case tagAnchor:
			for _, a := range n.Attr {
				if a.Key != attrHref {
					continue
				}
				// we don't care about links within the same page or emails
				if strings.HasPrefix(a.Val, "#") ||
					strings.HasPrefix(a.Val, "/#") ||
					strings.HasPrefix(a.Val, hrefMailto) {
					continue
				}
				u, err = url.Parse(a.Val)
				if err != nil {
					fmt.Fprintf(
						os.Stderr,
						"error parsing anchor href: %v\n",
						err,
					)
					continue
				}
				// blank host is a relative link; fetchPage will fix it
				if data.Seen(u) {
					continue
				}
				if hasSameHost(u, hostUrl) {
					fmt.Printf("Found link: %q\n", u.String())
					*urls = append(*urls, u)
				}
			}
		case tagScript:
			for _, a := range n.Attr {
				if a.Key != attrSrc {
					continue
				}
				u, err = url.Parse(a.Val)
				if err != nil {
					fmt.Fprintf(
						os.Stderr,
						"error parsing script src: %v\n",
						err,
					)
					continue
				}
				if data.Seen(u) {
					continue
				}
				if hasSameHost(u, hostUrl) {
					fmt.Printf("Found script link: %q\n", u.String())
					*urls = append(*urls, u)
				}
			}
		case tagLink:
			for _, a := range n.Attr {
				if a.Key != attrHref {
					continue
				}
				u, err = url.Parse(a.Val)
				if err != nil {
					fmt.Fprintf(
						os.Stderr,
						"error parsing link href: %v\n",
						err,
					)
					continue
				}
				if data.Seen(u) {
					continue
				}
				if hasSameHost(u, hostUrl) {
					fmt.Printf("Found link: %q\n", u.String())
					*urls = append(*urls, u)
				}
			}
		default:
			break
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findLinks(c, hostUrl, urls, data)
	}
}

// hasSameHost returns true if u is a relative link or u is an absolute link
// matching original's Hostname
func hasSameHost(u *url.URL, original *url.URL) bool {
	return u.Hostname() == "" || u.Hostname() == original.Hostname()
}
