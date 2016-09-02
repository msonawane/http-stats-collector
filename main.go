package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var (
	http_host  string
	image_path string
)

func ProcessAndSend(params map[string]string, query url.Values) {
	fmt.Println(params, query)
}

func PageViewsHandler(image_path string) (handle func(http.ResponseWriter, *http.Request)) {

	//cache file once
	image, err := ioutil.ReadFile(image_path)
	if err != nil {
		fmt.Errorf("error", err)
	}
	//content length
	length := strconv.Itoa(len(image))

	// return closure with access to image and length
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "image/gif")
		res.Header().Add("Cache-Control", "private, no-cache, no-cache=Set-Cookie, proxy-revalidate no-store, proxy-revalidate")
		res.Header().Add("Pragma", "no-cache")
		res.Header().Add("Expires", "Fri, 11 Jan 2000 01:00:00 GMT")
		res.Header().Add("Content-Length", length)
		res.Write(image)

		// Get request data
		params := make(map[string]string)
		params["userAgent"] = req.UserAgent()
		query, _ := url.ParseQuery(req.URL.RawQuery)

		go ProcessAndSend(params, query)
	}
}

func main() {
	flag.StringVar(&http_host, "http_host", "localhost:8888", "host/port to serve tracking image")
	flag.StringVar(&image_path, "image_path", "1.png", "absolute path to image")
	flag.Parse()

	http.HandleFunc("/", PageViewsHandler(image_path))
	http.ListenAndServe(http_host, nil)
}
