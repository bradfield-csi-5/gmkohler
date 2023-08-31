package pkg

import (
	"fmt"
	"net/url"
	"path"
)

// TODO move to new package to hide private members
type MirrorData struct {
	outDir        string
	hostName      string
	seenUrls      map[string]string // consider custom types
	memoryChannel chan *memory      // used to write to seenUrls in a goroutine
}
type memory struct {
	Url       *url.URL
	Extension string
}

// BelongsToSite returns true if u is a relative link or u is an absolute link
// whose hostname matches m.hostName
func (m *MirrorData) BelongsToSite(u *url.URL) bool {
	return u.Hostname() == "" || u.Hostname() == m.hostName
}

// Seen returns whether we've seen u.
func (m *MirrorData) Seen(u *url.URL) bool {
	_, ok := m.seenUrls[urlKey(u)]
	return ok
}

// MarkSeen updates our data to indicate we've seen u and its corresponding
// file path fp has been computed (extension relies on response headers being
// made, so we want to avoid fetching the page again)
func (m *MirrorData) MarkSeen(u *url.URL, fp string) {
	m.memoryChannel <- &memory{Url: u, Extension: path.Ext(fp)}
}

func (m *MirrorData) FilePathFor(u *url.URL) (string, error) {
	fp, exists := m.seenUrls[urlKey(u)]
	if !exists {
		return "", fmt.Errorf(
			"MirrorData.FilePathFor(u): url %q has not been seen",
			u.String(),
		)
	}
	return fp, nil
}

// urlKey defines how we consider each URL to be unique.  This will be easier to
// change than raw map access, should we decide to change the key of seenUrls.
// Right now, we are using Path because Host/Scheme get modified, and we have
// decided not to care about query parameters for now (should be easy to add).
func urlKey(u *url.URL) string {
	return u.Path
}

// NewMirrorData builds a new *MirrorData to start a recursive mirroring,
// and spawns a goroutine to update the map based on data sent to the channel.
func NewMirrorData(outDir string) *MirrorData {
	data := &MirrorData{
		outDir:        outDir,
		seenUrls:      make(map[string]string),
		memoryChannel: make(chan *memory),
	}

	go func() {
		for m := range data.memoryChannel {
			d, fp := makePathFromUrl(data.outDir, m.Url.Path, m.Extension)
			data.seenUrls[urlKey(m.Url)] = path.Join(d, fp)
		}
	}()

	return data
}
