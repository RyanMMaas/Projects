package m4a

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
	{"..\\test_files\\m4a\\m4a - Copy (2).m4a", "M4A Album 1", "M4A Artist 2", "", nil},
	{"..\\test_files\\m4a\\m4a - Copy (3).m4a", "M4A Album 2", "M4A Artist 1", "", nil},
	{"..\\test_files\\m4a\\m4a - Copy (4).m4a", "M4A Album 1", "", "M4A Album Artist 1", nil},
	{"..\\test_files\\m4a\\m4a - Copy (5).m4a", "", "M4A Artist 3", "", nil},
	{"..\\test_files\\m4a\\m4a - Copy (6).m4a", "", "", "M4A Album Artist 2", nil},
	{"..\\test_files\\m4a\\m4a - Copy (7).m4a", "M4A Album 3", "", "", nil},
	{"..\\test_files\\m4a\\m4a - Copy.m4a", "M4A Album 1", "M4A Artist 1", "", nil},
	{"..\\test_files\\m4a\\m4a.m4a", "M4A Album 1", "M4A Artist 1", "M4A Album Artist 1", nil},
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
