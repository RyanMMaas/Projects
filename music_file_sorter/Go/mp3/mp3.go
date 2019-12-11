package mp3

import (
	"io"
	"log"
	"os"

	"github.com/RyanMMaas/projects/music_file_sorter/Go/mfutils"
)

//Mp3File holds the data of the mp3 file
type Mp3File struct {
	fileID  []uint8
	version []uint8
	flags   uint8
	size    []uint32
	// headerData
	album       string
	artist      string
	contrArtist string
}

// GetTags returns the album, artist, contributing artist, and
// an error for mp3 files
func GetTags(file *os.File) (album, artist, contArtist string, e error) {
	var m Mp3File

	buf := make([]byte, 10)
	if _, err := io.ReadFull(file, buf); err != nil {
		log.Fatal(err)
	}

	m.fileID = append(m.fileID, uint8(buf[0]), uint8(buf[1]), uint8(buf[2]))
	m.version = append(m.version, uint8(buf[3]), uint8(buf[4]))
	m.flags = buf[5]
	m.size = append(m.size, uint32(buf[6]), uint32(buf[7]), uint32(buf[8]), uint32(buf[9]))

	tagVersion := int(m.version[0])
	tagSize := (m.size[3] & 0xFF) | ((m.size[2] & 0xFF) << 7) | ((m.size[1] & 0xFF) << 14) | ((m.size[0] & 0xFF) << 21) + 10

	var usesSync bool
	if (m.flags & 0x80) != 0 {
		usesSync = true
	} else {
		usesSync = false
	}

	buf = make([]byte, tagSize)
	if _, err := io.ReadFull(file, buf); err != nil {
		log.Fatal(err)
	}

	length := len(buf)
	if usesSync {
		newPos := 0
		newBuffer := make([]byte, tagSize)

		for i := 0; i < len(buf); i++ {
			if i < len(buf)-1 && (buf[i]&0xFF) == 0xFF && buf[i+1] == 0 {
				newBuffer[newPos] = byte(0xFF)
				newPos++
				i++
				continue
			}
			newBuffer[newPos] = buf[i]
			newPos++
		}
		length = newPos
		buf = newBuffer
	}

	pos := 0
	var id3FrameSize int
	if tagVersion < 3 {
		id3FrameSize = 6
	} else {
		id3FrameSize = 10
	}

	for i := 0; i < 30; i++ {
		rembytes := length - pos
		if rembytes < id3FrameSize {
			break
		}
		if buf[pos] < 'A' || buf[pos] > 'Z' {
			break
		}
		var framename []byte
		var framesize int

		if tagVersion < 3 {
			framename = append(framename, buf[pos], buf[pos+1], buf[pos+2])
			framesize = (int((buf[pos+5] & 0xFF)) << 8) | (int((buf[pos+4] & 0xFF)) << 16) | (int((buf[pos+3] & 0xFF)) << 24)
		} else {
			framename = append(framename, buf[pos], buf[pos+1], buf[pos+2], buf[pos+3])
			framesize = int((buf[pos+7] & 0xFF)) | (int((buf[pos+6] & 0xFF)) << 8) | (int((buf[pos+5] & 0xFF)) << 16) | (int((buf[pos+4] & 0xFF)) << 24)
		}
		if pos+framesize > length {
			break
		}
		if string(framename) == "TALB" {
			m.album = mfutils.ParseTextField(buf, (pos + id3FrameSize), framesize)
		}
		if string(framename) == "TPE2" {
			m.artist = mfutils.ParseTextField(buf, (pos + id3FrameSize), framesize)
		}
		if string(framename) == "TPE1" {
			m.contrArtist = mfutils.ParseTextField(buf, (pos + id3FrameSize), framesize)
		}
		pos += int(framesize) + id3FrameSize
		continue
	}

	return m.album, m.artist, m.contrArtist, nil
}
