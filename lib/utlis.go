package hwapro

import (
	// for making http requests
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	/*
	   //for debug proxy
	   "crypto/tls"
	   "net/url"
	*/)

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// Uploads a picture for a new ad, only 5 pictures can be in a new ad, any more and this will fail.
func uploadFile(path string) *http.Request {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}
	req, _ := http.NewRequest("POST", "https://hardverapro.hu/muvelet/apro/feltolt.php", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

// Generic request manager for easy debugging
func performHTTPRequest(req *http.Request, sess *UserSession) ([]byte, []string) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	// form token is bound to vid
	req.Header.Set(`Cookie`, `vid=`+sess.vid+`; identifier=`+sess.identifier+`; login-options={"stay":true,"no_ip_check":true,"leave_others":true}; prf_ls_uad=price.a.200.normal; rtif-legacy=1; login-options={"stay":true,"no_ip_check":true,"leave_others":true}`)

	/*
	  // this is for debug proxying
	  proxy, _ :=url.Parse("http://127.0.0.1:8080")
	  tr := &http.Transport{
	  	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	  	Proxy: http.ProxyURL(proxy),
	  }
	*/

	tr := &http.Transport{}
	// for avoiding infinite redirect loops
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[!] HTTP request failed to" + req.URL.Host + req.URL.Path)
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[!] HTTP request failed to" + req.URL.Host + req.URL.Path)
		panic(err)
	}
	// fmt.Println(string(resp.Header.Values("Set-Cookie")[0]))

	return respBody, resp.Header.Values("Set-Cookie")
}

// gets the value of a cookie that was set in a response header
func getCookieValue(name string, headers []string) string {
	value := ""
	for _, val := range headers {
		if strings.Contains(val, name) {
			value = val
		}
	}

	return trimBetween(value, "=", ";")
}

// extracts the string between start and end (exclusive)
func trimBetween(s string, start string, end string) string {
	if idx := strings.Index(s, start); idx != -1 {
		s = s[idx+len(start):]
	}

	if idx := strings.Index(s, end); idx != -1 {
		s = s[:idx]
	}

	return s
}
