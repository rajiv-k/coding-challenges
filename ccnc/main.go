package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

const (
	programName    = "ccnc"
	programVersion = "0.0.1"
)

var (
	version        bool
	listenMode     bool
	listenPort     int
	executablePath string
)

func main() {
	flag.BoolVar(&version, "version", false, "show version")
	flag.BoolVar(&listenMode, "l", false, "start in listening mode")
	flag.IntVar(&listenPort, "p", 0, "port to listen on")
	flag.StringVar(&executablePath, "e", "/bin/cat", "path to the executable")
	flag.Parse()
	if version {
		log.Printf("%v %v\n", programName, programVersion)
		os.Exit(0)
	}

	if listenMode {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%v", listenPort))
		if err != nil {
			log.Fatalf("ERROR: could not start listener: %v", err)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("WARN: could not accept(): %v\n", err)
				continue
			}
			handleConnection(conn, executablePath)
		}
	}

}

func handleConnection(conn net.Conn, executablePath string) {
	remoteAddr := conn.RemoteAddr()
	log.Printf("%v> connected!\n", remoteAddr)

	for {
		cmd := exec.Command(executablePath)
		cmd.Stderr = conn
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("%v> bye!\n", remoteAddr)
				return // NOTE: the client has disconnected.
			}
			log.Printf("%> ERROR: could not read from client: %v\n", remoteAddr, err)
			conn.Close()
			return
		}
		log.Printf("%v> read %v bytes\n", remoteAddr, n)
		if string(buf[:n]) == "quit\n" {
			conn.Close()
			break
		}

		cmd.Stdin = bytes.NewReader(buf)
		out, err := cmd.Output()
		if err != nil {
			log.Printf("%v> ERROR: err during command execution: %v\n", remoteAddr, err)
		}
		// log.Printf("%v> out: %v\n", remoteAddr, string(out))
		conn.Write(out)
	}

	conn.Close()
	log.Printf("%v> closed!\n", remoteAddr)
}
