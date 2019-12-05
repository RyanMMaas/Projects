package user

import (
	"errors"
	"os"
	"path/filepath"
)

/*
State stores data about the user connected to the server
*/
type State struct {
	currentDir string
}

// Init initialize the State of the user
func (us *State) Init() {
	dir, err := os.Getwd()
	if err != nil {
		//handle error
	}
	us.currentDir = dir
}

//ChangeDir changes the directory on the server
func (us *State) ChangeDir(dir string) error {
	chdirError := errors.New("Error changing directory")
	if dir == ".." {
		us.currentDir = filepath.Dir(us.currentDir)
	}
	dname, de, err := dirExists(us.currentDir, dir)
	if de == false {
		//directory doesnt exist
	}
	if err != nil {
		//handle error
	}
	if de == true {
		us.currentDir = filepath.Join(us.currentDir, dname)
		return nil
	}
	return chdirError
}

//CurrentDir returns the current directory that the user is in
func (us *State) CurrentDir() string {
	return us.currentDir
}

func dirExists(current, dir string) (string, bool, error) {
	f, err := os.Stat(filepath.Join(current, dir))
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	if err == nil && f.IsDir() {
		return f.Name(), true, nil
	}
	//something went worng
	return "", false, nil
}
