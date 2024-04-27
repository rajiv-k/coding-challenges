package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	programName    = "cccurl"
	programVersion = "0.0.1"
	methodGet      = "GET"
	methodPost     = "POST"
	methodPut      = "PUT"
	methodDelete   = "DELETE"
)

var (
	version bool
	verbose bool
	method  string
	data    string
	header  string
)

func main() {
	flag.BoolVar(&version, "version", false, "show version")
	flag.BoolVar(&version, "V", false, "show version")
	flag.BoolVar(&verbose, "verbose", false, "verbose")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.StringVar(&method, "method", "GET", "HTTP method")
	flag.StringVar(&method, "X", "GET", "HTTP method")
	flag.StringVar(&data, "d", "", "request payload")
	flag.StringVar(&data, "data", "", "request payload")
	flag.StringVar(&header, "header", "", "request header")
	flag.StringVar(&header, "H", "", "request header")
	flag.Parse()

	if version {
		fmt.Printf("%v %v\n", programName, programVersion)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}
	uri := flag.Args()[0]
	u, err := url.Parse(uri)
	if err != nil {
		log.Fatalf("ERROR: invalid url: %v\n", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		log.Fatalf("ERROR: invalid protocol, only http and https are supported")
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	switch method {
	case methodGet:
		get(client, uri)

	case methodDelete:
		del(client, uri)

	case methodPost:
		post(client, uri, header, data)

	default:
		log.Fatalf("%v is not implemented\n", method)
	}
}

func usage() {
	fmt.Printf("Usage: %v [flags] <uri>\n", programName)
}

func dumpRequest(req *http.Request) *bytes.Buffer {
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("> %v %v %v\n", req.Method, req.URL.Path, req.Proto))
	for k, v := range req.Header {
		buf.WriteString(fmt.Sprintf("> %v: %v\n", k, v))
	}

	return buf
}

func dumpResponse(resp *http.Response) *bytes.Buffer {
	buf := &bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("< %v %v %v\n", resp.Proto, resp.StatusCode, resp.Status))
	buf.WriteString(fmt.Sprintf("< Date: %v\n", resp.Header.Get("Date")))
	for k, v := range resp.Header {
		if k == "Date" {
			continue
		}
		buf.WriteString(fmt.Sprintf("< %v: %v\n", k, v[0]))
	}

	return buf
}

func del(client *http.Client, uri string) {
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		log.Fatalf("ERROR: could not create HTTP request: %v", err)
	}

	do(client, req)
}

func get(client *http.Client, uri string) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Fatalf("ERROR: could not create HTTP request: %v", err)
	}

	do(client, req)
}

func post(client *http.Client, uri string, header string, data string) {
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(data))
	if err != nil {
		log.Fatalf("ERROR: could not create HTTP request: %v", err)
	}

	if header != "" {
		parts := strings.Split(header, ":")
		if len(parts) != 2 {
			log.Fatalf("ERROR: invalid value for header, must be of form -H 'foo: bar'")
		}

		headerKey, headerVal := parts[0], strings.TrimPrefix(parts[1], "")
		req.Header.Add(headerKey, headerVal)
	}
	do(client, req)
}

func do(client *http.Client, req *http.Request) {
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("ERROR: could not perform HTTP request: %v", err)
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ERROR: could not read HTTP response: %v", err)
	}

	if verbose {
		b := dumpRequest(req)
		fmt.Println(b.String())
		bb := dumpResponse(resp)
		fmt.Println(bb.String())
	}

	fmt.Fprintln(os.Stdout, string(data))
}
