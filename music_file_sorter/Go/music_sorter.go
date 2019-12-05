package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"strings"
)

func main() {
	root := ""

	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if string(path[len(path)-3:]) == "mp3" || string(path[len(path)-3:]) == "m4a" {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, mfPath := range files {
		END := ""
		if filetype := filepath.Ext(mfPath); filetype == ".mp3" {
			var m mp3ID3v2file

			file, err := os.Open(mfPath)
			if err != nil {
				log.Fatal(err)
			}
			getMP3Tags(file, &m)
			file.Close()
			ep := createDestinationPath(m.artist, m.album, m.contrArtist)
			END += ep
			moveFile(mfPath, END)
		} else if filetype == ".m4a" {
			var m m4afile
			file, err := os.Open(mfPath)
			if err != nil {
				log.Fatal(err)
			}
			getM4ATags(file, &m)
			file.Close()
			ep := createDestinationPath(m.artist, m.album, m.contrArtist)
			END += ep
			moveFile(mfPath, END)
		} else {
			ep := createDestinationPath("", "", "")
			END += ep
			moveFile(mfPath, END)
		}
	}
}

func createDestinationPath(art, alb, conArt string) string {
	var destPath bytes.Buffer
	re := regexp.MustCompile("[\\/:*?\"<>|]")
	if art != "" && alb != "" {
		destPath.WriteString(re.ReplaceAllLiteralString(art, "_"))
		destPath.WriteString("\\")
		destPath.WriteString(re.ReplaceAllLiteralString(alb, "_"))
	} else if conArt != "" && alb != "" {
		destPath.WriteString(re.ReplaceAllLiteralString(conArt, "_"))
		destPath.WriteString("\\")
		destPath.WriteString(re.ReplaceAllLiteralString(alb, "_"))
	} else if art != "" {
		destPath.WriteString(re.ReplaceAllLiteralString(art, "_"))
	} else if conArt != "" {
		destPath.WriteString(re.ReplaceAllLiteralString(conArt, "_"))
	} else {
		destPath.WriteString("NO_TAGS_FOLDER")
	}
	return destPath.String()
}

func moveFile(srcPath, dstPath string) error {

	_, err := os.Stat(dstPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dstPath, os.ModePerm)
		}
	}

	inputFile, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(err)
	}

	s := strings.Split(srcPath, "\\")
	outputFile, err := os.Create(dstPath + "\\" + s[len(s)-1])
	if err != nil {
		inputFile.Close()
		log.Fatal(err)
	}

	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Remove(srcPath)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func getMP3Tags(file *os.File, m *mp3ID3v2file) {

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
			m.album = parseTextField(buf, (pos + id3FrameSize), framesize)
		}
		if string(framename) == "TPE2" {
			m.artist = parseTextField(buf, (pos + id3FrameSize), framesize)
		}
		if string(framename) == "TPE1" {
			m.contrArtist = parseTextField(buf, (pos + id3FrameSize), framesize)
		}
		pos += int(framesize) + id3FrameSize
		continue
	}
}

// search for 4 character codes in bytes
// http://atomicparsley.sourceforge.net/mpeg-4files.html
func getM4ATags(file *os.File, m *m4afile) {
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	//search for aART (artist)
	if i := bytes.Index(buf, []byte{97, 65, 82, 84}); i != -1 {
		size := (buf[i-1]) | ((buf[i-2]) << 8) | ((buf[i-3]) << 16) | ((buf[i-4]) << 24)
		m.artist = parseTextField(buf, (i + 19), (int(size) - 20))
	}
	//search for 0xa9alb (album)
	if i := bytes.Index(buf, []byte{169, 97, 108, 98}); i != -1 {
		size := (buf[i-1]) | ((buf[i-2]) << 8) | ((buf[i-3]) << 16) | ((buf[i-4]) << 24)
		m.album = parseTextField(buf, (i + 19), (int(size) - 20))
	}
	//search for 0xa9art (contributing artists)
	if i := bytes.Index(buf, []byte{169, 65, 82, 84}); i != -1 {
		size := (buf[i-1]) | ((buf[i-2]) << 8) | ((buf[i-3]) << 16) | ((buf[i-4]) << 24)
		m.contrArtist = parseTextField(buf, (i + 19), (int(size) - 20))
	}

}

func parseTextField(buf []byte, pos int, size int) string {
	// still need to deal with decoding types so
	// special characters show up
	// pos+3 to get rid of [1 255 254]
	var s string
	// 255 254 is byte order mark denoting utf-16 little endian
	if buf[pos] == 1 {
		s = string(bytes.Replace(buf[pos+3:pos+size], []byte("\x00"), []byte{}, -1))
	} else {
		s = string(bytes.Replace(buf[pos+1:pos+size], []byte("\x00"), []byte{}, -1))
	}
	return s
}

type mp3ID3v2file struct {
	fileID  []uint8
	version []uint8
	flags   uint8
	size    []uint32
	// headerData
	album       string
	artist      string
	contrArtist string
}
type m4afile struct {
	album       string
	artist      string
	contrArtist string
}
