package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)


type firstLine struct {
	Method string
	Path string
	Query string 
	HttpVersion string
}


type request struct {
	FirstLine firstLine
	Headers  map[string]string
	Body string
}


func writeAFile(){
	f, err := os.Create("req.txt")
	if err != nil {
		panic(err.Error())
	}

	defer f.Close()

	l, err := f.WriteString("GET /api/search?q=query_string&page=1&limit=10 HTTP/1.1\r\nHost: api.example.com\r\nUser-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36\r\nAccept: application/json\r\nAuthorization: Bearer your_access_token\r\n\r\n")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println( l ,"bytes written successfully")

}


func readFile() (firstLine string, headers []string, body string, err error) {
	
	f, err := os.Open("req.txt")
	if err != nil {
		return "", nil, "", err
	}
	defer f.Close()

	var c strings.Builder
	_, err = io.Copy(&c, f)

	if err != nil {
		return "", nil, "", err
	}
	
	cs := c.String()

	lines := strings.Split(cs, "\r\n")

	firstLine = lines[0]

	var headerEnd int 

	for i, line := range lines {
		if line == "" {
			headerEnd = i
			break
		}
	}

	headers = lines [1:headerEnd]

	body = strings.Join(lines[headerEnd+1:], "\r\n")

	return firstLine, headers, body, nil
}

//first line

func firstLineConstractor(f string) firstLine {

	var fl firstLine
	flt := strings.Split(f, " ")
	fl.Method = flt[0]
	
	if strings.Contains(flt[1],"?") {
		slt := strings.Split(flt[1],"?")
		fl.Path = slt[0]
		fl.Query = slt[1]
	}else {
		fl.Path = flt[1]
	}


	fl.HttpVersion = flt[2]
	
	return fl
}

//headers 
func headersConstractor(headers []string) map[string]string{
	ht := make(map[string]string)
	for i := 0; i < len(headers); i++ {
		parts := strings.SplitN(headers[i], ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			ht[key] = value
		}else {
			fmt.Println("can't constract heasers")
		}
	}

	return ht

}

//body 

func bodyConstractor(body string) string{
	return body
}

//request 

func requestBuilder (f firstLine, headers map[string]string, body string) request{
	var r request
	r.FirstLine = f
	r.Headers  = headers
	r.Body = body 

	return r
}
//requist handler 


func handleRequest (r request) {
	if r.FirstLine.Method == "POST" {
		fmt.Println("this is a post request")
	} else if r.FirstLine.Method == "GET" {
		fmt.Println("this is a get request")
	}else {
		fmt.Println("request method not supported ")
	}
}


func main() {
	firstLine, headers, body, err := readFile()
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}

	r := requestBuilder(firstLineConstractor(firstLine), headersConstractor(headers), bodyConstractor(body) )


	fmt.Println(r.FirstLine.Method)
	fmt.Println(r.FirstLine.Path)
	fmt.Println(r.FirstLine.Query)
	fmt.Println(r.FirstLine.HttpVersion)

	for key, value := range r.Headers {
		fmt.Printf("%s: %s\n", key, value)
	}

	fmt.Println(r.Body)

	handleRequest(r)
}