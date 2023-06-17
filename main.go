package main

import (
	// "compress/zlib"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getGoogle(url string, header http.Header) *http.Response {
	// log.Println("url:", url)
	req, err := http.NewRequest("GET", url, nil)
	req.Header = header
	req.Header.Set("Accept-Encoding", "null")
	checkErr(err)
	// log.Println("header", header)
	client := http.Client{}
	resp, err := client.Do(req)
	checkErr(err)
	return resp
}

const site string = "127.0.0.1"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		agent := r.UserAgent()
		addr := r.RemoteAddr
		unescUri, err := url.PathUnescape(uri)
		checkErr(err)
		log.Println(addr, r.Method, "UserAgent:", agent, "URI:", unescUri)

		if r.Method != "GET" {
			log.Fatal("Find none GET method.")
		}

		resp := getGoogle("https://www.google.com" + uri, r.Header)
		content, err := io.ReadAll(resp.Body)
		checkErr(err)
		// log.Println("Response body", string(content))
		defer resp.Body.Close()

		// w.Write(content)
		reg, err := regexp.Compile("www.google.com")
		checkErr(err)
		html := reg.ReplaceAll(content, []byte(site))
		w.Write(html)
		// log.Println("html", string(html))

		for k, v := range resp.Header {
			w.Header().Set(k, strings.Join(v, ""))
		}
	})
	log.Println("Start to listen on :1984")
	err := http.ListenAndServe(":1984", nil)
	checkErr(err)
}