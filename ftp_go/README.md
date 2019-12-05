# FTP Client

# Why
This project was a way to learn something new while improving my programming ability in go. I had not worked on many projects using sockets before so this was a fun project with a lot of learning opportunities.

# Keywords
>Go, golang, ftp, sockets, tcp

# Features
* [x] Multiple clients connect at once
* [x] Change directory of client/server
* [x] Get current working directory
* [x] Close connection
* [x] Open new connection
* [x] Delete file(s) on server/client
* [x] Copy file(s) from/to server
* [x] Help command to list all commands
* [x] Get info of file from server
* [x] List files in current directory on client/server
* [x] Create directory on server

# Usage
Start the *server*:
```bash
go run server.go xxx.xxx.xxx.xxx:yyyy
```
> `xxx.xxx.xxx.xxx` is the ip address and `yyyy` is the port number

OR
```bash
go run server.go
```
and enter the ip and port from there.
___
Start the *client*:
```bash
go run client.go xxx.xxx.xxx.xxx:yyyy
```
> `xxx.xxx.xxx.xxx` is the ip address and `yyyy` is the port number

OR
```bash
go run client.go
```
and enter the ip and port from there.

# Issues
* The get/put commands currently copy wrong when going from linux/windows to windows/linux

# After Thoughts
This project took more time than any other project I had worked on alone before. I read a lot on sockets/tcp and used a lot of functions in Go that I had not used before. I really enjoyed working on it and implementing new functions. Seeing new functions work correctly was very rewarding.

I had trouble early on understanding how the sockets were working. I had errors where entering too many commands would eventually cause an error because some bytes were being read on the server to late. Once I realized how it was being sent in a stream it became much easier and I was able to resolve these issues.

Testing worked completely perfectly when done on a local machine. Then I tried testing between a virtual machine running Linux and my machine running Windows. This is when the errors with copying files popped up. Files would Show up with and extra '0' byte in the destination. This is something I am currently working to fix though.