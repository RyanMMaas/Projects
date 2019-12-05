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
	"path/filepath"
	"strings"

	"github.com/RyanMMaas/projects/ftp_go/user"
)

const (
	BUFSIZE int = 1024
	ls          = "001"
	cd          = "002"
	pwd         = "003"
	mkdir       = "004"
	get         = "005"
	put         = "006"
	del         = "007"
	info        = "008"
)

func main() {
	fmt.Printf("Starting Server\n")

	var listener net.Listener
	var err error
	var addr string

	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	//REMOVE
	addr = "127.0.0.1:200"

	for {
		reader := bufio.NewReader(os.Stdin)
		if addr == "" {
			fmt.Print("Enter address and port to start server (x.x.x.x:yyyy): ")
			addr, err = reader.ReadString('\n')
		}
		// Start server on address given as argument 1 on command line
		fmt.Printf("Server starting on %s\n", addr)
		addr = strings.TrimRight(addr, "\t\n\r")

		listener, err = net.Listen("tcp", addr)
		if err != nil {
			fmt.Printf("Error starting server on %s:%s\n", addr, err)
			fmt.Print("Enter address and port to start server (x.x.x.x:yyyy): ")
			addr, err = reader.ReadString('\n')
			continue
		} else {
			break
		}
	}

	fmt.Println("SERVER STARTED")

	err = os.Chdir("/")
	if err != nil {
		//handle error
	}
	// Wait for connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s", err)
			continue
		}

		fmt.Printf("CONNECTED: Client %v\n", conn.RemoteAddr())
		go handleClient(conn)
	}
}

// Fills a byte slice b to be length s
func fillBuf(b []byte, s int) []byte {
	for len(b) < s {
		b = append(b, 0x00)
	}
	return b
}

// After connecting, clients are sent here until they disconnect
func handleClient(conn net.Conn) {
	// pass userstate to functions that need it...
	var us user.State
	us.Init()

	defer conn.Close()

	// Commands are sent as 3 number string to reduce chances of errors
	var buf [3]byte
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("DISCONNECTED: Client %v\n", conn.RemoteAddr())
			return
		}
		comm := buf[:n]
		if bytes.Compare(comm, []byte(ls)) == 0 {
			listFiles(conn, &us)
		} else if bytes.Compare(comm, []byte(cd)) == 0 {
			changeDir(conn, &us)
		} else if bytes.Compare(comm, []byte(pwd)) == 0 {
			getDir(conn, &us)
		} else if bytes.Compare(comm, []byte(mkdir)) == 0 {
			makeDir(conn, &us)
		} else if bytes.Compare(comm, []byte(get)) == 0 {
			getFile(conn, &us)
		} else if bytes.Compare(comm, []byte(put)) == 0 {
			putFile(conn)
		} else if bytes.Compare(comm, []byte(del)) == 0 {
			deleteFile(conn, &us)
		} else if bytes.Compare(comm, []byte(info)) == 0 {
			getFileInfo(conn, &us)
		}
	}
}

// Send list of files in current directory to client
func listFiles(conn net.Conn, us *user.State) {
	d := us.CurrentDir()
	files, err := ioutil.ReadDir(d)
	if err != nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))

	// Get how many files/directories, convert to binary and send
	dir := make([]byte, 8)
	binary.LittleEndian.PutUint64(dir, uint64(len(files)))
	conn.Write(dir)

	// Loop through files, get the length of the file name
	// convert it to binary and send it to the client
	flen := make([]byte, 8)
	for _, fn := range files {
		conn.Write([]byte("OK"))

		binary.LittleEndian.PutUint64(flen, uint64(len(fn.Name())))
		conn.Write(flen)
		b := fillBuf([]byte(fn.Name()), BUFSIZE)
		conn.Write(b)
	}

	return
}

// Changes directory on server
func changeDir(conn net.Conn, us *user.State) {
	// Get directory name
	var dir [256]byte
	n, err := conn.Read(dir[:])
	if err != nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))

	err = us.ChangeDir(string(dir[:n]))
	if err != nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))
}

// Copy file to client
func getFile(conn net.Conn, us *user.State) {
	// Name of file to copy
	var name [256]byte
	n, err := conn.Read(name[:])
	if err != nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))

	// Open file to copy, send error message and break
	// if file does not exist or if there was an error opening the file
	file, err := os.Open(filepath.Join(us.CurrentDir(), string(name[:n])))
	defer file.Close()
	if err != nil {
		conn.Write([]byte("ER"))
		if err == os.ErrNotExist {
			conn.Write([]byte("File does not exist"))
		} else {
			conn.Write([]byte("Error opening file"))
		}
		return
	}
	conn.Write([]byte("OK"))

	// buf: bytes being sent to the client
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

// Send file info to client
func getFileInfo(conn net.Conn, us *user.State) {
	// Name of file
	var name [256]byte
	n, err := conn.Read(name[:])
	if err != nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))

	// Attempt to open file
	file, err := os.Open(filepath.Join(us.CurrentDir(), string(name[:n])))
	defer file.Close()
	if err != nil {
		conn.Write([]byte("ER"))
		if err == os.ErrNotExist {
			conn.Write([]byte("File does not exist"))
		} else {
			conn.Write([]byte("Error opening file"))
		}
		return
	}
	conn.Write([]byte("OK"))

	// Get file info into fi
	fi, err := file.Stat()
	if err != nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))

	// Send size as 8 bytes
	size := make([]byte, 8)
	binary.LittleEndian.PutUint64(size, uint64(fi.Size()))
	conn.Write(size)

	// send 1 byte for if directory or file
	var dir [1]byte
	if fi.IsDir() {
		dir[0] = 0x01
	} else {
		dir[0] = 0x00
	}
	conn.Write(dir[:])

	// Send file permissions
	conn.Write([]byte(fi.Mode().String()))
}

// Delete file
func deleteFile(conn net.Conn, us *user.State) {
	// Name of file
	var name [256]byte
	n, err := conn.Read(name[:])
	if err != nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))

	err = os.Remove(filepath.Join(us.CurrentDir(), string(name[:n])))
	if err != nil && !os.IsNotExist(err) {
		conn.Write([]byte("ER"))
	} else {
		conn.Write([]byte("OK"))
	}
}

// Create directory
func makeDir(conn net.Conn, us *user.State) {
	// Name of directory
	var dir [256]byte
	n, err := conn.Read(dir[:])
	if err != nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))

	err = os.Mkdir(filepath.Join(us.CurrentDir(), string(dir[:n])), os.ModePerm)
	if err != nil {
		conn.Write([]byte("ER"))
	} else {
		conn.Write([]byte("OK"))
	}

}

// Copy file to server
func putFile(conn net.Conn) {
	// Name of file to copy
	var name [256]byte
	n, err := conn.Read(name[:])
	if err != nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))

	// Check if a file with given name exists on the server so it doesn't overwrite that file
	_, err = os.Stat(string(name[:n]))
	if err == nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))

	// Create the file on the server
	file, err := os.Create(string(name[:n]))
	defer file.Close()
	if err != nil {
		conn.Write([]byte("ER"))
		return
	}
	conn.Write([]byte("OK"))

	// buf: bytes sent from the client are read into buf
	// toRead: the amount of bytes being read. This is used so
	//		that when a file sent from the server was run through fillBuf(),
	//		the null characters are not written to the file. By sending this
	// 		as 4 bytes rather than a string there is less chance of error
	// 		ex: writing to much to the file, getting stuck in the loop
	var buf [BUFSIZE]byte
	var toRead [4]byte
	var ok [2]byte
	for {
		n, err = conn.Read(ok[:])
		// Break if there is an error on client otherwise continue reading into buf
		if string(ok[:n]) == "ER" || err != nil {
			fmt.Printf("Error copying file to server: %s\n", err)
		} else if string(ok[:n]) == "FE" {
			break
		} else if string(ok[:n]) == "OK" {
			// Get how many bytes are being read going to be read
			n, err = conn.Read(toRead[:])
			if err != nil {
				fmt.Printf("Error copying file to server %s\n", err)
				return
			}
			tr := int32(binary.LittleEndian.Uint32(toRead[:n]))

			// Read BUFSIZE into buf
			n, err = conn.Read(buf[:])
			if err != nil {
				fmt.Printf("Error copying file to server %s\n", err)
				return
			}

			// Write to the file tr (amount bytes that aren't null characters from fillBuf) from buf
			n, err = file.Write(buf[:tr])
			if err != nil {
				fmt.Printf("Error copying file to server %s\n", err)
				return
			}
		} else {
			break
		}

	}
}

// Get current directory
func getDir(conn net.Conn, us *user.State) {
	conn.Write([]byte("OK"))

	wd := us.CurrentDir()
	conn.Write([]byte(wd))
}
