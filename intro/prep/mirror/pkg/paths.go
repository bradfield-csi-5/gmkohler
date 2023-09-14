package pkg

import "path"

// makePathFromUrl builds a file path from the target directory outDir and a URL
// path uPath.  The function signature is intended to be similar to path.Split,
// which this function leans on.
// examples:
//
//	makePathFromUrl("./mirrored", "", ".html") ("/mirrored", "index.html")
//	makePathFromUrl("./mirrored", "/courses/ssba", .html") (
//		"/mirrored/courses",
//		"ssba.html",
//	)
func makePathFromUrl(outDir string, uPath string, ext string) (
	dir string,
	fileName string,
) {
	dir, fileName = path.Split(uPath)
	if fileName == "" {
		fileName = "index"
	}
	if path.Ext(fileName) == "" {
		fileName += ext
	}
	dir = path.Join(outDir, dir)
	return
}
