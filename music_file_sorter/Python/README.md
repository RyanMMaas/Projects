# Music Sorter - Python

# Why
This project is a music file sorter that reads the metadata of music files (mp3 & m4a) and sorts them in directories based on artist/album name. It was created to help automate a task that a friend was doing.

# Keywords
>File manipulation, mp3, m4a, mutagen, python, regex

# Features
* [x] Moving/sorting mp3's
* [x] Moving/sorting m4a's
* [x] Creating directories
* [x] Moving all other types of files
* [x] Displays total time taken.

# Usage
Set root (Line 69) in music_sorter.py as the path of unsorted files.

Set moveTo (Line 73) in music_sorter.py as the destination path.

```bash
python music_sorter.py
```

# Issues
No known issues.

# After Thoughts
Doing this project got me more interested in the structure of the files. I had expected to be able to just get the metadata from the files without needing to use a module but it was more complicated than I thought. After reading the documentation and looking at a few examples of how the mutagen module worked the rest of the project was very simple.