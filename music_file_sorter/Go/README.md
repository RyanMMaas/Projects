# Music Sorter - Go

# Why
This project was made after creating the python file sorter. It was a way to better understand the go language as well as learning more about the way the files were structured.

# Keywords
>Bitwise math, file manipulation, Go, golang, mp3, m4a, regex

# Features
* [x] Moving/sorting mp3's
* [x] Moving/sorting m4a's
* [x] Creating directories
* [x] Ignoring other types of files

# Usage
Set root (Line 15) in music_sorter.go as the path of unsorted files.
Set dest (Line 15) in music_sorter.go as the destination path.

```bash
go run music_sorter.go
```
# Issues
Some mp3 files of certain versions *may* read correctly.

# After Thoughts
Learning the different structures and how to read the meta-data of each file was difficult at first. I hadn't expected it to require reading through the bits of the file individually but that is part of why I did this project a second time. Doing it in python was really simple because there were modules that found the metadata of the files for me.

This project I found much more fulfilling since I had to do research on the structure of mp3/m4a files on my own.