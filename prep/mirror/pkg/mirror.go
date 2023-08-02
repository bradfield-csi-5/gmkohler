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
	"sync"
)

const (
	headerContentType = "Content-Type"
	attrHref          = "href"
	attrSrc           = "src"
	contentTypeHtml   = "text/html;"
	hrefMailto        = "mailto:"
	tagAnchor         = "a"
	tagImage          = "img"
	tagLink           = "link"
	tagScript         = "script"
)

// link holds a URL to be fetched and an attr whose Val should be modified
// with a filepath corresponding to the URL name
type link struct {
	url  *url.URL
	attr *html.Attribute
}

// Mirror writes the web page at u to data.outDir, and if the web page is HTML,
// writes all the document's links within the same domain to data.outDir
func Mirror(
	u *url.URL,
	data *MirrorData,
	wg *sync.WaitGroup,
) error {
	defer wg.Done()
	if data.Seen(u) {
		return nil
	}
	fmt.Printf("mirroring %s", u.String())
	resp, err := fetchPage(u)
	if err != nil {
		return fmt.Errorf("error fetching web page %v: %v\n", u, err)
	}
	var writeToFile writeProtocol = httpResponseWriteProtocol(resp)
	contentType := resp.Header.Get(headerContentType)
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil || exts == nil {
		return fmt.Errorf(
			"error determining file extension for type %q: %v\n",
			contentType,
			err,
		)
	}
	dir, fName := makePathFromUrl(data.outDir, u.Path, exts[0])
	fPath := path.Join(dir, fName)
	data.MarkSeen(u, fPath)

	var links []*link
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
		writeToFile = htmlNodeWriteProtocol(domBody)
		findLinks(domBody, u, &links, data)

		fmt.Printf("found %d links in %s\n", len(links), u.String())
	}

	var childWg sync.WaitGroup
	for _, l := range links {
		// many links are relative.  Here we infer from previous link the
		// host and scheme so that we can successfully fetch it.
		if l.url.Host == "" {
			l.url.Host = u.Host
		}
		if l.url.Scheme == "" {
			l.url.Scheme = u.Scheme
		}
		childWg.Add(1)
		go func(u *url.URL) {
			err = Mirror(u, data, &childWg)
			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"error mirroring page %v: %v\n",
					u,
					err,
				)
			}
		}(l.url)
	}

	childWg.Wait()
	for j := range links {
		l := links[j]
		localPath, err := data.FilePathFor(l.url)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"error getting file path for %q\n",
				l.url.String(),
			)
		} else {
			l.attr.Val = localPath
		}
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %q: %v\n", dir, err)
	}

	f, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("error opening file %q: %v\n", fPath, err)
	}
	err = writeToFile(f)
	if err != nil {
		return fmt.Errorf("error writing file %s: %v\n", fPath, err)
	}

	return nil
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

// findLinks collects links within n with the same Host as hostUrl
// right now, it is serial because it's all in-memory and the network requests
// should be more of a bottleneck, but we could probably make it concurrent.
// This can also be DRYed out with a helper that takes an *[]*html.Attribute
// because most of/the/logic is/the/same, save for some special filtering in the
// anchor tags and a different choice of attribute name for where we know the
// link will be encoded.
func findLinks(
	n *html.Node,
	hostUrl *url.URL,
	links *[]*link,
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
			for j, a := range n.Attr {
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
				if data.BelongsToSite(u) {
					fmt.Printf("Found anchor link: %q\n", u.String())
					*links = append(*links, &link{
						url:  u,
						attr: &n.Attr[j], // need element access later
					})
				}
			}
		case tagScript:
			for j, a := range n.Attr {
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
				if data.BelongsToSite(u) {
					fmt.Printf("Found script link: %q\n", u.String())
					*links = append(*links, &link{
						url:  u,
						attr: &n.Attr[j], // need element access
					})
				}
			}
		case tagLink:
			for j := range n.Attr {
				a := n.Attr[j]
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
				if data.BelongsToSite(u) {
					fmt.Printf("Found link: %q\n", u.String())
					*links = append(*links, &link{
						url:  u,
						attr: &n.Attr[j], // need element access
					})
				}
			}
		default:
			break
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findLinks(c, hostUrl, links, data)
	}
}
