package m4a

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/RyanMMaas/Projects/music_file_sorter/Go/mfutil"
)

type M4afile struct {
	album       string
	artist      string
	contrArtist string
}

// GetTags returns the album, artist, contributing artist, and
// an error for m4a files
func GetTags(file *os.File) (album, artist, contrArtist string, e error) {
	// search for 4 character codes in bytes
	var m M4afile

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return "", "", "", err
	}
	//search for aART (artist)
	if i := bytes.Index(buf, []byte{97, 65, 82, 84}); i != -1 {
		size := (buf[i-1]) | ((buf[i-2]) << 8) | ((buf[i-3]) << 16) | ((buf[i-4]) << 24)
		m.artist = mfutil.ParseTextField(buf, (i + 19), (int(size) - 20))
	}
	//search for 0xa9alb (album)
	if i := bytes.Index(buf, []byte{169, 97, 108, 98}); i != -1 {
		size := (buf[i-1]) | ((buf[i-2]) << 8) | ((buf[i-3]) << 16) | ((buf[i-4]) << 24)
		m.album = mfutil.ParseTextField(buf, (i + 19), (int(size) - 20))
	}
	//search for 0xa9art (contributing artists)
	if i := bytes.Index(buf, []byte{169, 65, 82, 84}); i != -1 {
		size := (buf[i-1]) | ((buf[i-2]) << 8) | ((buf[i-3]) << 16) | ((buf[i-4]) << 24)
		m.contrArtist = mfutil.ParseTextField(buf, (i + 19), (int(size) - 20))
	}
	return m.album, m.artist, m.contrArtist, nil
}

// http://atomicparsley.sourceforge.net/mpeg-4files.html
// character codes of metadata in m4a files
