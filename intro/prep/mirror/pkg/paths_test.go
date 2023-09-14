package pkg

import "testing"

func TestMakePathFromUrl(t *testing.T) {
	type wanted struct {
		dir      string
		filename string
	}

	data := []struct {
		outDir string
		path   string
		ext    string
		want   wanted
	}{
		{
			outDir: "./mirrored",
			ext:    ".html",
			want: wanted{
				dir:      "mirrored",
				filename: "index.html",
			},
		},
		{
			outDir: "./mirrored",
			path:   "/courses/ssba",
			ext:    ".html",
			want: wanted{
				dir:      "mirrored/courses",
				filename: "ssba.html",
			},
		},
		{
			outDir: "./mirrored",
			path:   "/courses/ssba/",
			ext:    ".html",
			want: wanted{
				dir:      "mirrored/courses/ssba",
				filename: "index.html",
			},
		},
	}
	for _, d := range data {
		dir, filename := makePathFromUrl(d.outDir, d.path, d.ext)
		if dir != d.want.dir || filename != d.want.filename {
			t.Errorf(
				"makePathFromUrl(%q, %q, %q) = (%q, %q), want (%q, %q)",
				// args
				d.outDir,
				d.path,
				d.ext,
				// results
				dir,
				filename,
				// wants
				d.want.dir,
				d.want.filename,
			)
		}
	}
}
