package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

const (
	BUFSIZE int = 1024
	cd          = "cd"
	close       = "close"
	del         = "del"
	get         = "get"
	help        = "help"
	info        = "info"
	lcd         = "lcd"
	ldel        = "ldel"
	lls         = "lls"
	lpwd        = "lpwd"
	ls          = "ls"
	mdel        = "mdel"
	mget        = "mget"
	mkdir       = "mkdir"
	mput        = "mput"
	open        = "open"
	put         = "put"
	pwd         = "pwd"
	quit        = "quit"
)

func main() {
	fmt.Printf("Starting Client...\n")

	reader := bufio.NewReader(os.Stdin)
	var conn net.Conn
	if len(os.Args) > 1 {
		conn = openConn(reader, os.Args[1])
	} else {
		conn = openConn(reader)
	}
	err := os.Chdir("/")
	if err != nil {
		//handle error
	}

	for {
		fmt.Printf(">> ")
		line, err := reader.ReadBytes('\n')
		line = bytes.TrimRight(line, " \t\r\n")
		if err != nil {
			fmt.Printf("Error processing command: %s\n", err)
			continue
		}

		// Split the command into 2 byte slices.
		// comm[0] holds the command
		// comm[1] holds the arguments
		comm := bytes.SplitN(line, []byte(" "), 2)
		switch string(comm[0]) {
		case cd:
			changeDir(conn, comm[1])
		case close:
			conn.Close()
			fmt.Printf("Disconnected from server\n")
		case del:
			deleteFile(conn, comm[1])
		case get:
			if len(comm) == 2 {
				getFile(conn, comm[1])
			} else {
				fmt.Printf("Wrong usage of get.\n")
			}
		case help:
			commandHelp()
		case info:
			if len(comm) == 2 {
				getFileInfo(conn, comm[1])
			} else {
				fmt.Printf("Wrong usage of info.\n")
			}
		case lcd:
			if len(comm) == 2 {
				changeLocalDir(comm[1])
			} else {
				fmt.Printf("Wrong usage of lcd.\n")
			}
		case ldel:
			if len(comm) == 2 {
				localDeleteFile(comm[1])
			} else {
				fmt.Printf("Wrong usage of ldel.\n")
			}
		case lls:
			listLocalFiles()
		case lpwd:
			getLocalDir()
		case ls:
			listFiles(conn)
		case mdel:
			if len(comm) == 2 {
				multDeleteFile(conn, comm[1])
			} else {
				fmt.Printf("Wrong usage of mdel.\n")
			}
		case mget:
			if len(comm) == 2 {
				multGetFile(conn, comm[1])
			} else {
				fmt.Printf("Wrong usage of mget.\n")
			}
		case mkdir:
			if len(comm) == 2 {
				makeDir(conn, comm[1])
			} else {
				fmt.Printf("Wrong usage of mkdir.\n")
			}
		case mput:
			if len(comm) == 2 {
				multPutFile(conn, comm[1])
			} else {
				fmt.Printf("Wrong usage of mput.\n")
			}
		case open:
			if len(comm) == 2 {
				conn = openConn(reader, string(comm[1]))
			} else if len(comm) == 1 {
				conn = openConn(reader)
			}
		case put:
			if len(comm) == 2 {
				putFile(conn, comm[1])
			} else {
				fmt.Printf("Wrong usage of put.\n")
			}
		case pwd:
			getDir(conn)
		case quit:
			conn.Close()
			fmt.Printf("Disconnected from server\n")
			os.Exit(1)
		default:
			fmt.Printf("%s not a recognized command. \"help\" for command list\n", comm[0])
		}

	}

}

// Fills a byte slice b to be length s
func fillBuf(b []byte, s int) []byte {
	for len(b) < s {
		b = append(b, 0x00)
	}
	return b
}

// Opens connection to the server
func openConn(r *bufio.Reader, a ...string) net.Conn {
	// ac is used to determine if an address was given or not
	ac := len(a)
	var conn net.Conn
	var err error
	var addr string
	for {
		// ac == 1: address was given when the funcion was called
		// if ac == anything else: ask the user to enter address
		if ac == 1 {
			addr = a[0]
			fmt.Printf("Connecting to %s\n", addr)
		} else {
			fmt.Printf("IP and port to connect (x.x.x.x:y): ")
			addr, err = r.ReadString('\n')
		}
		addr = strings.TrimRight(addr, "\t\n\r")

		conn, err = net.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("Error connecting to %s:%s\n", addr, err)
			// Set ac == 0 so that if the first try failed because the funciton was
			// called with a wrong address, the user can enter a new address the next loop
			ac = 0
			continue
		} else {
			break
		}
	}
	return conn
}

// Change directory on server
func changeDir(conn net.Conn, name []byte) {
	conn.Write([]byte("002"))
	conn.Write(name)

	// Error checking on server
	var ok [2]byte
	n, err := conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error changing directory: %s\n", err)
		return
	}

	// Error checking on server
	n, err = conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error changing directory: %s\n", err)
		return
	}

}

// Delete multiple files
func multDeleteFile(conn net.Conn, files []byte) {
	f := bytes.Split(files, []byte(" \""))
	for _, a := range f {
		fn := bytes.Trim(a, "\"")
		deleteFile(conn, fn)
	}
}

// Delete file on server
func deleteFile(conn net.Conn, name []byte) {
	conn.Write([]byte("007"))
	conn.Write(name)

	// Error checking on server
	var ok [2]byte
	n, err := conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error deleting file: %s\n", err)
		return
	}

	// Error checking on server
	n, err = conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error deleting file: %s\n", err)
		return
	}

}

// Copy multiple files from server
func multGetFile(conn net.Conn, files []byte) {
	f := bytes.Split(files, []byte(" \""))
	for _, a := range f {
		fn := bytes.Trim(a, "\"")
		getFile(conn, fn)
	}
}

// Copy file from server
func getFile(conn net.Conn, name []byte) {
	// Check if a file with given name exists on the client so it doesn't overwrite that file
	_, err := os.Stat(string(name))
	if err == nil {
		fmt.Printf("\"%s\" already exists in the current directory\n", string(name))
		return
	}

	// Create the file on the client
	file, err := os.Create(string(name[:]))
	defer file.Close()
	if err != nil {
		fmt.Printf("Error creating file %s: %s\n", string(name[:]), err)
		return
	}

	conn.Write([]byte("005"))
	conn.Write(name)

	// Error checking on server
	var ok [2]byte
	n, err := conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error copying file: %s\n", err)
		return
	}

	// Error checking on server, errB is the error from the server
	var errB [256]byte
	n, err = conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		n, err = conn.Read(errB[:])
		if err != nil {
			fmt.Printf("Error copying file: %s\n", err)
		} else {
			fmt.Printf("%s", err)
		}
		return
	}

	// buf: bytes sent from the server are read into buf
	// toRead: the amount of bytes being read. This is used so
	//		that when a file sent from the server was run through fillBuf(),
	//		the null characters are not written to the file. By sending this
	// 		as 4 bytes rather than a string there is less chance of error
	// 		ex: writing to much to the file, getting stuck in the loop
	var buf [BUFSIZE]byte
	var toRead [4]byte
	for {
		n, err = conn.Read(ok[:])
		// Break if there is an error on server otherwise continue reading into buf
		if string(ok[:n]) == "ER" || err != nil {
			break
		} else if string(ok[:n]) == "OK" {
			// Get how many bytes are being read going to be read
			n, err = conn.Read(toRead[:])
			if err != nil {
				fmt.Printf("Error copying file %s\n", err)
				return
			}
			tr := int32(binary.LittleEndian.Uint32(toRead[:n]))

			// Read BUFSIZE into buf
			n, err = conn.Read(buf[:])
			if err != nil {
				fmt.Printf("Error copying file %s\n", err)
				return
			}

			// Write to the file tr (amount bytes that aren't null characters from fillBuf) from buf
			n, err = file.Write(buf[:tr])
			if err != nil {
				fmt.Printf("Error copying file %s\n", err)
				return
			}
		} else {
			break
		}

	}
}

// List commands/usage of commands
func commandHelp(command ...string) {
	fmt.Printf("Commands\n\n")

	fmt.Printf("%-8s- %s\n", cd, "Changes directory on server")
	fmt.Printf("%-8s- %s\n", close, "Closes current connection")
	fmt.Printf("%-8s- %s\n", del, "Deletes a file from the server")
	fmt.Printf("%-8s- %s\n", get, "Copys a file from the server")
	fmt.Printf("%-8s- %s\n", help, "Lists commands and gets usage of commands")
	fmt.Printf("%-8s- %s\n", info, "Gets info of a file from the server")
	fmt.Printf("%-8s- %s\n", lcd, "Changes the directory on the client")
	fmt.Printf("%-8s- %s\n", ldel, "Deletes a file on the client")
	fmt.Printf("%-8s- %s\n", lls, "Lists files in current local directory")
	fmt.Printf("%-8s- %s\n", lpwd, "Lists current working directory on the cleint")
	fmt.Printf("%-8s- %s\n", ls, "Lists files on the current working diretory of the server")
	fmt.Printf("%-8s- %s\n", mdel, "Deletes multiple files")
	fmt.Printf("%-8s- %s\n", mget, "Copys multiple files")
	fmt.Printf("%-8s- %s\n", mkdir, "Creates a directory on the server")
	fmt.Printf("%-8s- %s\n", mput, "Copys multiple files to the server")
	fmt.Printf("%-8s- %s\n", open, "Opens a connection")
	fmt.Printf("%-8s- %s\n", put, "Copys a file to the server")
	fmt.Printf("%-8s- %s\n", pwd, "Get the current working directory on the server")
	fmt.Printf("%-8s- %s\n", quit, "Closes the connection and exits")

	fmt.Printf("\n\n")
}

// Gets file info (size, directory/file, permissions) of a file on the server
func getFileInfo(conn net.Conn, name []byte) {
	conn.Write([]byte("008"))
	conn.Write(name)

	// Error checking on server
	var ok [2]byte
	n, err := conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error getting file info\n")
		return
	}

	// Error checking on server
	var errB [256]byte
	n, err = conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		n, err = conn.Read(errB[:])
		if err != nil {
			fmt.Printf("Error getting file info: %s\n", err)
		} else {
			fmt.Printf("%s", err)
		}
		return
	}

	// Error checking on server
	n, err = conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error getting file info\n")
		return
	}

	// Getting file size as bytes rather than string so that 8 bytes are always read
	var fs [8]byte
	n, err = conn.Read(fs[:])
	if err != nil {
		fmt.Printf("Error getting file size\n")
		return
	}
	size := int64(binary.LittleEndian.Uint64(fs[:n]))

	// Get whether file is a directory or file
	var d [1]byte
	n, err = conn.Read(d[:])
	if err != nil {
		fmt.Printf("Error getting if directory\n")
		return
	}
	var dir string
	if string(d[:n]) == "1" {
		dir = "Dir"
	} else {
		dir = "File"
	}

	// Get file permissions
	var m [10]byte
	n, err = conn.Read(m[:])
	if err != nil {
		fmt.Printf("Error getting file mode\n")
		return
	}
	mode := string(m[:n])

	fmt.Printf("Size: %d bytes\n", size)
	fmt.Printf("File/Dir: %s\n", dir)
	fmt.Printf("File mode: %s\n", mode)
}

// Change directory on client
func changeLocalDir(dir []byte) {
	err := os.Chdir(string(dir[:len(dir)]))
	if err != nil {
		fmt.Printf("Error changing client directory: %s\n", err)
	}
}

// Delete file on client
func localDeleteFile(name []byte) {
	err := os.Remove(string(name))
	if err != nil && !os.IsNotExist(err) {
		fmt.Printf("Error deleting file: %s\n", err)
	}
}

// List files in working directory of server
func listFiles(conn net.Conn) {
	conn.Write([]byte("001"))

	// Error checking on server
	var ok [2]byte
	n, err := conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error getting files/directories: %s\n", err)
		return
	}

	// Read how many files/directories there are and convert it from binary to uint64
	var f [8]byte
	n, err = conn.Read(f[:])
	dirLen := int64(binary.LittleEndian.Uint64(f[:n]))

	// buf: bytes sent from the server are read into buf
	// toRead: the amount of bytes being read. This is used so
	//		that when a file sent from the server was run through fillBuf(),
	//		the null characters are not read. By sending this
	// 		as 8 bytes rather than a string there is less chance of error
	// 		ex: writing to much to the file, getting stuck in the loop
	toRead := make([]byte, 8)
	var buf [BUFSIZE]byte
	for i := 0; i < int(dirLen); i++ {
		// Error checking on server
		n, err = conn.Read(ok[:])
		if string(ok[:n]) != "OK" || err != nil {
			if string(ok[:n]) == "ER" {
				break
			}
			fmt.Printf("Error getting files/directories: %s\n", err)
			break
		}

		// Get how many bytes to read into buf
		n, err = conn.Read(toRead[:])
		if err != nil {
			fmt.Printf("Error getting files/directories: %s\n", err)
			break
		}
		fnlen := int64(binary.LittleEndian.Uint64(toRead[:n]))

		// Read from conn into buf
		n, err = conn.Read(buf[:])
		if err != nil {
			fmt.Printf("Error getting files/directories: %s\n", err)
			break
		}
		// Output the file name
		fmt.Printf("%s\n", string(buf[:int(fnlen)]))

	}

	return
}

// Get the current directory on the client
func getLocalDir() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting local directory %s\n", err)
	} else {
		fmt.Println(wd)
	}
}

// List the files in the current local directory
func listLocalFiles() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Printf("Error getting local files: %s\n", err)
	} else {
		for _, fn := range files {
			fmt.Printf("%s\n", fn.Name())
		}
	}

}

// Create a directory on the server
func makeDir(conn net.Conn, dirName []byte) {
	conn.Write([]byte("004"))
	conn.Write(dirName)

	// Error checking on server
	var ok [2]byte
	n, err := conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		return
	}

	// Error checking on server
	n, err = conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		return
	}
}

// Copy multiple files to server
func multPutFile(conn net.Conn, files []byte) {
	f := bytes.Split(files, []byte(" \""))
	for _, a := range f {
		fn := bytes.Trim(a, "\"")
		putFile(conn, fn)
	}
}

// Copy a file to server
func putFile(conn net.Conn, fileName []byte) {
	conn.Write([]byte("006"))
	conn.Write(fileName)

	// Error checking on server
	var ok [2]byte
	n, err := conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error copying file to server: %s\n", err)
		return
	}

	// Error checking on server
	n, err = conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		if err != nil {
			fmt.Printf("Error copying file to server: %s\n", err)
		} else {
			fmt.Printf("File named %s already exists on server\n", string(fileName))
			return
		}
	}

	// Error checking on server
	n, err = conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error creating file on server: %s\n", err)
		return
	}

	// Open file to copy to server
	file, err := os.Open(string(fileName[:]))
	defer file.Close()
	if err != nil {
		if err == os.ErrNotExist {
			fmt.Printf("Error: File %s does not exist on client\n", string(fileName[:]))
		} else {
			fmt.Printf("Error opening file %s on client\n", string(fileName[:]))
		}
		return
	}

	// buf: bytes being sent to the server
	// toSend: the amount of bytes being being sent that aren't
	//		null characters added by fillBuf().
	buf := make([]byte, BUFSIZE)
	toSend := make([]byte, 4)
	for {
		n, err = file.Read(buf[:])
		if err == io.EOF {
			conn.Write([]byte("FE"))
			break
		} else if err != nil {
			conn.Write([]byte("ER"))
			break
		} else {
			conn.Write([]byte("OK"))

			binary.LittleEndian.PutUint32(toSend, uint32(n))
			conn.Write(toSend)

			ts := fillBuf(buf, BUFSIZE)
			conn.Write(ts)
		}
	}

}

// Get current working directory of server
func getDir(conn net.Conn) {
	conn.Write([]byte("003"))

	// Error checking on server
	var ok [2]byte
	n, err := conn.Read(ok[:])
	if string(ok[:n]) != "OK" || err != nil {
		fmt.Printf("Error getting directory: %s\n", err)
		return
	}

	var dir [256]byte
	n, err = conn.Read(dir[:])
	if err != nil {
		fmt.Printf("Error getting directory: %s\n", err)
	}
	fmt.Println(string(dir[:n]))
}
