package hwapro

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type advertisment struct {
	Category    string   // category id of the ad
	Title       string   // item title
	Price       string   // item price
	Description string   // ad description
	imageUrls   []string // urls for advert's images
	LocalImages []string // picture locations on disk
	Location    string   // geographic location
	fdrid       string   // image folder id
	mainpic     string   // thumbnail id for the ad
}

type MapEntry struct {
	Id    int
	Title string
}

// Parse an ARCHIVED advertisment based on url
func ParseAdvertisment(url string) advertisment {
	// Get the advert
	req, _ := http.NewRequest("GET", url, nil)
	respBody, _ := performHTTPRequest(req, &UserSession{})
	ad := advertisment{}

	// Cut around the meat of the response for easier handling
	div := trimBetween(string(respBody), `<div id="middle">`, `Tetszik</button>\n\t\t</div>`)

	// Parse the title (onyl works on archived ads)
	ad.Title = trimBetween(div, `Archív –
					<s>`, `</s>`)
	ad.Title = html.UnescapeString(ad.Title)

	// Parse the geographic location
	location := trimBetween(div, `<span class="fas fa-map-marker"></span>`, `</span>`)
	ad.Location = mapToMap(location)

	// Parse the price
	price := trimBetween(div, `<h2 class="text-center text-md-left">`, `Ft</h2>`)
	ad.Price = strings.Join(strings.Fields(price), "")

	// Parse the description
	ad.Description = trimBetween(div, `<div class="mb-3 rtif-content">`, `</div>`)
	// for double newlines
	ad.Description = strings.ReplaceAll(ad.Description, ` class="mgt2">`, `></p><p></p><p>`)
	ad.Description = strings.ReplaceAll(ad.Description, ` class="mgt1">`, `></p><p>`)
	ad.Description = strings.TrimSpace(ad.Description)

	// Parse the category
	cat := strings.Split(trimBetween(string(respBody), `"breadcrumb"`, `</div>`), "\n")
	for _, val := range cat {
		if strings.Contains(val, "<a href=") {
			ad.Category = trimBetween(val, `html">`, `</a`)
		}
	}
	ad.Category = mapToCategory(ad.Category)

	// Parse the image urls
	images := strings.Split(trimBetween(div, `<div class="carousel-inner">`, `fa-map-marker">`), "\n")
	for _, val := range images {
		if strings.Contains(val, "<img") {
			ad.imageUrls = append(ad.imageUrls, "https://hardverapro.hu"+trimBetween(val, `="`, `"`))
		}
	}

	fmt.Println("[+] Advert parsed")
	ad.Print()
	return ad
}

// Saves an advertisment to the current directory
func (ad *advertisment) SaveAdvertisment() {
	if _, err := os.Stat("ads"); os.IsNotExist(err) {
		os.Mkdir("ads", 0o755)
	}
	os.Mkdir("ads/"+ad.Title, 0o755)

	// download images to local folder
	for key, val := range ad.imageUrls {
		downloadFile("ads/"+ad.Title+"/"+strconv.Itoa(key)+filepath.Ext(val), val)
		ad.LocalImages = append(ad.LocalImages, "ads/"+ad.Title+"/"+strconv.Itoa(key)+filepath.Ext(val))
	}

	// serialize the info to a text file
	json, _ := json.Marshal(ad)
	f, _ := os.Create("ads/" + ad.Title + "/info")
	defer f.Close()
	_, _ = f.Write(json)
	fmt.Println(`[+] Advert saved to: ads/` + ad.Title)
}

// Post the advert to the site
// requires saving beforend so that we have the images to upload
func (ad *advertisment) RepostSaved(sess *UserSession) {
	// tell hardverapro you want to post a new ad
	req, _ := http.NewRequest("GET", "https://hardverapro.hu/muvelet/apro/uj.php?url=%2Fhirdetesfeladas%2Fuj.php", nil)
	resp, _ := performHTTPRequest(req, sess) // I had to create a resp variable anyway
	// upload the imges
	for _, val := range ad.LocalImages {
		req := uploadFile(val)
		resp, _ = performHTTPRequest(req, sess)
	}
	fmt.Println(`[+] Images uploaded`)

	// last one uploaded will be used as the cover picture
	var bare map[string]interface{}
	json.Unmarshal(resp, &bare)
	ad.mainpic = strconv.Itoa(int(bare["uploaded"].(map[string]interface{})["id"].(float64)))
	ad.fdrid = strconv.Itoa(int(bare["fdrid"].(float64)))

	// update form token
	sess.getAdPostFormToken()

	// prepare the form
	reqBody := map[string]string{
		"fidentifier": sess.fidentifier,
		"title":       ad.Title,
		"mce_0":       "",
		"content":     ad.Description,
		"price":       ad.Price,
		"cities":      ad.Location,
		"shipping":    "",
		"buying":      "0",
		"udrid":       ad.Category,
		"bundle":      "",
		"cmpid":       "",
		"picid":       ad.mainpic,
		"fdrid":       ad.fdrid,
		"cmp_text":    "",
		"infourl":     "",
		"brandnew":    "",
	}
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for key, val := range reqBody {
		_ = writer.WriteField(key, val)
	}
	_ = writer.Close()

	// post the request
	req, _ = http.NewRequest("POST", "https://hardverapro.hu/muvelet/apro/uj.php?url=%2Fhirdetesfeladas%2Fuj.php", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, _ = performHTTPRequest(req, sess)
	fmt.Println(`[+] Advert Posted Succesfully!`)
}

// Loads a previously saved advertisment
func LoadAdvertisment(path string) advertisment {
	dat, _ := ioutil.ReadFile(path + "/info")
	ad := advertisment{}
	json.Unmarshal(dat, &ad)
	return ad
}

func (ad advertisment) Print() {
	fmt.Println(`-- Category: `+ad.Category+`
-- Title: `+ad.Title+`
-- Price: `+ad.Price+`
-- Description: `+ad.Description+`
-- location: `+ad.Location+`
-- Images: `, ad.imageUrls)
}

// Maps city names to the ID used by hardverapro
func mapToMap(location string) string {
	location = strings.TrimSpace(location)
	mapUrl := "https://hardverapro.hu/muvelet/telepules/listaz.php"
	req, _ := http.NewRequest("GET", mapUrl, nil)
	resp, _ := performHTTPRequest(req, &UserSession{})
	var mapArray []MapEntry
	json.Unmarshal(resp, &mapArray)

	// Get the ID where the title (city name) matches
	for _, val := range mapArray {
		if strings.Contains(val.Title, location) {
			return strconv.Itoa(val.Id)
		}
	}
	return "Not Found"
}

// Maps the category names to the ID used by hardverapro
func mapToCategory(category string) string {
	mapUrl := "https://hardverapro.hu/muvelet/apro_konyvtar/listaz.php"
	req, _ := http.NewRequest("GET", mapUrl, nil)
	resp, _ := performHTTPRequest(req, &UserSession{})

	var bare map[string]interface{}
	json.Unmarshal(resp, &bare)
	list := bare["list"].(map[string]interface{})

	for i := 1; i < 1000; i++ { // iterate on all entries and check the titles for a match
		idx := strconv.Itoa(i)
		if list[idx] != nil {
			if list[idx].(map[string]interface{})["title"] != nil {
				title := list[idx].(map[string]interface{})["title"].(string)
				if strings.Contains(title, category) {
					return idx
				}
			}
		}
	}
	return "Not Found"
}
