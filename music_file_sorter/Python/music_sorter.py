import os
import time
import sys
import re
from subprocess import call
import mutagen


# resize window
os.system("mode con cols=50 lines=20")

# sf: startFolder/startDestination - folder that holds everything to be copied
# e: endfolder/enddestinatino - where the files will be copied
# name: name of file in folder to be copied
def robocopy(sf, e, name):
	call([r"robocopy", sf, e, name, "/mov", "/NFL", "/NDL", "/NJH", "/NJS", "/nc", "/ns", "/np"])

def get_path_from_metadata(artist, album, contArtist):
	if artist and album:
		end = moveTo+"\\"+re.sub(r'[\/:*?"<>|]', '_', artist)+"\\"+re.sub(r'[\/:*?"<>|]', '_', album)
	elif contArtist and album:
		end = moveTo+"\\"+re.sub(r'[\/:*?"<>|]', '_', contArtist)+"\\"+re.sub(r'[\/:*?"<>|]', '_', album)
	elif artist:
		end = moveTo+"\\"+re.sub(r'[\/:*?"<>|]', '_', artist)
	elif contArtist:
		end = moveTo+"\\"+re.sub(r'[\/:*?"<>|]', '_', contArtist)
	else:
		end = moveTo+"\\"+'NO_TAGS_FOLDER'

	return end

def move_file(path, sFolder):
	# songFile stores metadata about file (eg. artist, album, songname)
	songFile = mutagen.File(path)
	
	if type(songFile) == mutagen.mp4.MP4:
		artist, album, contArtist = None, None, None
		if 'aART' in songFile:
			artist = songFile['aART'][0]
		if '\xa9alb' in songFile:
			album = songFile['\xa9alb'][0]
		if '\xa9ART' in songFile:
			contArtist = songFile['\xa9ART'][0]

		end = get_path_from_metadata(artist, album, contArtist)

	elif type(songFile) == mutagen.mp3.MP3:
		artist, album, contArtist = None, None, None
		if 'TPE2' in songFile:
			artist = songFile['TPE2'][0]
		if 'TALB' in songFile:
			album = songFile['TALB'][0]
		if 'TPE1' in songFile:
			contArtist = songFile['TPE1'][0]
		
		end = get_path_from_metadata(artist, album, contArtist)

	else:
		end = moveTo+"\\"+'NOT_MP3_MP4'
		
	robocopy(sFolder, end, name)
	
def main():
	start_time = time.time()

	# root is folder where all unsorted files are
	root = r''


	# where the sorted files move to
	moveTo = r''

	if moveTo:
		for dirpath, subdirs, files in os.walk(root):
			for name in files:
				move_file(os.path.abspath(os.path.join(dirpath, name)), os.path.abspath(os.path.join(dirpath)))

		print("--- %s seconds ---" % (time.time() - start_time))
	else:
		print("No end location set (moveTo variable)")
	input("")

if __name__=="__main__":
	main()