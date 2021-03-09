package hwapro

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
)

type UserSession struct {
	identifier  string
	vid         string
	fidentifier string
}

// perform login
func (sess *UserSession) Login(username string, password string) string {
	sess.getLoginFormToken()

	fmt.Println("[+] Trying to log in as", username)
	urlPath := "https://hardverapro.hu/muvelet/hozzaferes/belepes.php?url=%2Findex.html"

	reqBody := map[string]string{
		"fidentifier":  sess.fidentifier,
		"email":        username,
		"pass":         password,
		"all":          "1",
		"stay":         "1",
		"no_ip_check":  "1",
		"leave_others": "1",
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for key, val := range reqBody {
		_ = writer.WriteField(key, val)
	}
	_ = writer.Close()

	req, _ := http.NewRequest("POST", urlPath, body)
	req.Header.Set("Content-Type", "multipart/form-data;  boundary="+writer.Boundary())

	_, headers := performHTTPRequest(req, sess)
	sess.identifier = getCookieValue("identifier", headers)

	if sess.identifier != "" {
		fmt.Println("[+] Login passed, session cookie parsed as: " + sess.identifier)
	} else {
		panic("login failed")
	}

	return "asdf"
}

// get a form token for logging in
func (sess *UserSession) getLoginFormToken() {
	fmt.Println("[+] Getting form token")

	urlPath := "https://hardverapro.hu/muvelet/hozzaferes/belepes.php?url=%2Findex.html"
	req, _ := http.NewRequest("GET", urlPath, nil)
	respBody, headers := performHTTPRequest(req, sess)

	sess.fidentifier = trimBetween(string(respBody), `name=\"fidentifier\" value=\"`, `\"`)
	fmt.Println("[+] Login form token parsed as: ", sess.fidentifier)
	sess.vid = getCookieValue("vid", headers)
	fmt.Println("[+] vid cookie parsed as: ", sess.vid)
}

// get a form token for posting ad
func (sess *UserSession) getAdPostFormToken() {
	fmt.Println("[+] Getting form token")

	urlPath := "https://hardverapro.hu/hirdetesfeladas/uj.php"
	req, _ := http.NewRequest("GET", urlPath, nil)
	respBody, headers := performHTTPRequest(req, sess)

	sess.fidentifier = trimBetween(string(respBody), `<body class="ha" data-token="`, `">`)
	fmt.Println("[+] Ad posting form token parsed as: ", sess.fidentifier)
	sess.vid = getCookieValue("vid", headers)
	fmt.Println("[+] vid cookie parsed as: ", sess.vid)
}
