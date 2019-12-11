package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RyanMMaas/projects/music_file_sorter/Go/m4a"
	"github.com/RyanMMaas/projects/music_file_sorter/Go/mfutils"
	"github.com/RyanMMaas/projects/music_file_sorter/Go/mp3"
)

func main() {
	root, dest := "", ""
	if root == "" || dest == "" {
		reader := bufio.NewReader(os.Stdin)
		if root == "" {
			fmt.Print("Please set the root path: ")
			root, _ = reader.ReadString('\n')
		}
		if dest == "" {
			fmt.Print("Please set the destinatoin path: ")
			dest, _ = reader.ReadString('\n')
		}
	}

	files, err := mfutils.GetFiles(root)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, path := range files {
		moveDest, alb, art, cArt := "", "", "", ""
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening: %s\n", path)
			break
		}

		switch fileType := filepath.Ext(path); fileType {
		case ".mp3":
			alb, art, cArt, err = mp3.GetTags(file)
			break
		case ".m4a":
			alb, art, cArt, err = m4a.GetTags(file)
			break
		default:
			break
		}
		file.Close()
		if err != nil {
			fmt.Printf("%s\nError getting file info for: %s\n", err, path)
			break
		}

		dp := mfutils.CreateDestination(art, alb, cArt)

		moveDest = filepath.Join(dest, dp)
		mfutils.MoveFile(path, moveDest)
	}
}
