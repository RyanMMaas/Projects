package mfutil

import (
	"bufio"
	"os"
	"path/filepath"
	"testing"
)

func TestGetFiles(t *testing.T) {
	testData := "..\\test_files\\files_list.txt"
	testRoot := "..\\test_files"
	files, err := GetFiles(testRoot)
	if err != nil {
		t.Error("Error reading files")
	}
	data, err := os.Open(testData)
	if err != nil {
		t.Error("Error opening files_list.txt")
	}
	defer data.Close()

	var list []string
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}

	for i, file := range files {
		if filepath.Base(file) != list[i] {
			t.Errorf("Expected %s, got: %s", list[i], filepath.Base(file))
		}
	}
}

type DestinationTest struct {
	art    string
	alb    string
	conArt string
	out    string
}

var createDestTests = []DestinationTest{
	{"foo", "bar", "baz", "foo\\bar"},
	{"fo\\o", "ba?r", "baz", "fo_o\\ba_r"},
	{"foo", "bar", "", "foo\\bar"},
	{"f:oo", "ba*r", "", "f_oo\\ba_r"},
	{"", "bar", "baz", "baz\\bar"},
	{"", "b\"ar", "b<az", "b_az\\b_ar"},
	{"foo", "", "", "foo"},
	{">foo", "", "", "_foo"},
	{"", "", "baz", "baz"},
	{"", "", "baz|", "baz_"},
	{"", "bar", "", "NO_TAGS"},
	{"", "b?ar", "", "NO_TAGS"},
	{"", "", "", "NO_TAGS"},
}

func TestCreateDestination(t *testing.T) {
	for _, test := range createDestTests {
		actual := CreateDestination(test.art, test.alb, test.conArt)
		if actual != test.out {
			t.Errorf("(%s, %s, %s) Expected: %s, Got: %s", test.art, test.alb, test.conArt, test.out, actual)
		}
	}
}

type MoveFileTest struct {
	src string
	dst string
}

var tempFiles []string = []string{"a", "b"}

/*
-source exists and dest doesnt
-source that doesnt exist
-dest has file with same name as source
*/
var createPassMoveFileTests = []MoveFileTest{
	{"..\\test_files\\temp\\src\\a", "..\\test_files\\temp\\dst"},
	{"..\\test_files\\temp\\src\\b", "..\\test_files\\temp\\dst"},
}
var createFailMoveFileTests = []MoveFileTest{
	{"..\\test_files\\temp\\src\\ba?dpa*th", "..\\test_files\\temp\\dst"},
	{"..\\test_files\\temp\\src\\a", "..\\test_files\\temp\\dst\\ba?dpa*th"},
	{"..\\test_files\\temp\\src\\a2\\a", "..\\test_files\\temp\\dst"},
	{"..\\test_files\\temp\\src\\b", "..\\test_files\\temp\\dst\\does_not_exist"},
	{"..\\test_files\\temp\\src\\c", "..\\test_files\\temp\\dst"},
}

func TestMoveFile(t *testing.T) {
	if err := os.MkdirAll("..\\test_files\\temp\\src\\a2", os.ModePerm); err != nil {
		t.Errorf("Error creating tempory source directory")
		return
	}
	if err := os.MkdirAll("..\\test_files\\temp\\dst", os.ModePerm); err != nil {
		t.Errorf("Error creating tempory destination directory")
		return
	}

	//Create temporary files
	for _, file := range tempFiles {
		tf, err := os.Create(filepath.Join("..\\test_files\\temp\\src", file))
		if err != nil {
			t.Errorf("%s (Error creating temporary file %s)", err, file)
		}
		tf.Close()
	}
	tf, err := os.Create("..\\test_files\\temp\\src\\a2\\a")
	if err != nil {
		t.Errorf("Error creating temporary file a2\\a")
	}
	tf.Close()
	defer os.RemoveAll("..\\test_files\\temp")

	for _, test := range createPassMoveFileTests {
		err := MoveFile(test.src, test.dst)
		if err != nil {
			t.Errorf("Error (%s) moving file: %s, to: %s", err, test.src, test.dst)
		}
	}
	for _, test := range createFailMoveFileTests {
		err := MoveFile(test.src, test.dst)
		if err == nil {
			t.Errorf("Error: no error (source: %s, dest: %s)", test.src, test.dst)
		}
	}
}
