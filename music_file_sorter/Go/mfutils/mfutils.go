package mfutils

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GetFiles returns a slice of mp3 and m4a files in
// the directory root.
func GetFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if s := string(path[len(path)-3:]); s == "mp3" || s == "m4a" {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

//CreateDestination creates the destination folder using
// the artist, album, and contributing artist; if the folder
// doesn't already exist. It replaces characters that cannot but
// used in a filepath on windows with a '_'
func CreateDestination(art, alb, conArt string) string {
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

// MoveFile moves the file from the source (srcPath)
// to the destination (dstPath) removing the original file from the source
func MoveFile(srcPath, dstPath string) error {

	_, err := os.Stat(dstPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dstPath, os.ModePerm)
		} else {
			return err
		}
	}

	inputFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}

	s := strings.Split(srcPath, "\\")
	outputFile, err := os.Create(dstPath + "\\" + s[len(s)-1])
	if err != nil {
		inputFile.Close()
		return err
	}

	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return err
	}

	err = os.Remove(srcPath)
	if err != nil {
		return err
	}
	return nil
}

//ParseTextField searches 'buf' starting at 'pos' for 'size'
//bytes and return the string of what was found.
func ParseTextField(buf []byte, pos int, size int) string {
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
