package mp3

import (
	"os"
	"testing"
)

type GetTagsTest struct {
	file     string
	alb      string
	art      string
	contrArt string
	err      error
}

var CreatePassGetTagsTests = []GetTagsTest{
	{"..\\test_files\\mp3\\mp3 - Copy (2).mp3", "MP3 Album 1", "MP3 Artist 2", "", nil},
	{"..\\test_files\\mp3\\mp3 - Copy (3).mp3", "MP3 Album 2", "MP3 Artist 1", "", nil},
	{"..\\test_files\\mp3\\mp3 - Copy (4).mp3", "MP3 Album 1", "", "MP3 Album Artist 1", nil},
	{"..\\test_files\\mp3\\mp3 - Copy (5).mp3", "", "MP3 Artist 3", "", nil},
	{"..\\test_files\\mp3\\mp3 - Copy (6).mp3", "", "", "MP3 Album Artist 2", nil},
	{"..\\test_files\\mp3\\mp3 - Copy (7).mp3", "MP3 Album 3", "", "", nil},
	{"..\\test_files\\mp3\\mp3 - Copy.mp3", "MP3 Album 1", "MP3 Artist 1", "", nil},
	{"..\\test_files\\mp3\\mp3.mp3", "MP3 Album 1", "MP3 Artist 1", "MP3 Album Artist 1", nil},
}

func TestGetTags(t *testing.T) {
	for _, test := range CreatePassGetTagsTests {
		f, err := os.Open(test.file)
		if err != nil {
			t.Errorf("%s, Error opening file %s", err, test.file)
		}
		if al, ar, cAr, e := GetTags(f); al != test.alb && ar != test.art && cAr != test.contrArt && e != test.err {
			t.Errorf("Expected/got: album %s/%s, artist: %s/%s, contributing artist: %s/%s, error: %s/%s",
				test.alb, al, test.art, ar, test.contrArt, ar, test.err, e)
		}
	}
}
