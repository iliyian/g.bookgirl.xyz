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

var googleReg = regexp.MustCompile("www.google.com")
var cookieReg = regexp.MustCompile(".google.com")

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

const site string = "g.bookgirl.xyz"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		agent := r.UserAgent()
		addr := r.RemoteAddr
		unescUri, err := url.PathUnescape(uri)
		checkErr(err)
		log.Println(addr, r.Method, "UserAgent:", agent, "URI:", unescUri)
		// 获取请求的基础信息

		if r.Method != "GET" {
			log.Fatal("Find none GET method.")
		}
		// 以防万一，因为目前只写了GET

		resp := getGoogle("https://www.google.com" + uri, r.Header)
		content, err := io.ReadAll(resp.Body)
		checkErr(err)
		defer resp.Body.Close()
		// 服务器获取指定内容

		html := googleReg.ReplaceAll(content, []byte(site))
		w.Write(html)
		// 更换掉文本中的www.google.com

		for k, v := range resp.Header {
			value := strings.Join(v, "")
			if k == "Set-Cookie" {
				value = cookieReg.ReplaceAllString(value, ".bookgirl.xyz")
				// log.Println("GOT SET-COOKIE, value:", value)
				// 更换掉 cookie 中的.google.com域名，实现cookie存储
			}
			w.Header().Set(k, value)
		}
		log.Println("Response headers", resp.Header)
	})
	log.Println("Start to listen on :1984")
	err := http.ListenAndServe(":1984", nil)
	checkErr(err)
}