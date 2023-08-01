package pkg

import "net/url"

type MirrorData struct {
	outDir   string
	seenUrls map[string]bool
}

// Seen returns whether we've seen u.  This will be easier to change than raw
// map accesses should we decide to change the key of seen.  Right now, we are
// using Path because Host/Scheme get modified, but this is not perfect because
// you could have same path at different subdomains.
func (m *MirrorData) Seen(u *url.URL) bool {
	return m.seenUrls[u.Path]
}

// MarkSeen updates our data to indicate we've seen u
func (m *MirrorData) MarkSeen(u *url.URL) {
	m.seenUrls[u.Path] = true
}

// NewMirrorData builds a new *MirrorData to start a recursive mirroring
//
// TODO: consider not exporting this type and giving cmd an interface with
// string
func NewMirrorData(outDir string) *MirrorData {
	return &MirrorData{
		outDir:   outDir,
		seenUrls: make(map[string]bool),
	}
}
