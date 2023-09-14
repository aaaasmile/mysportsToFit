package mytom

import (
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

type MytomState int

const (
	stLogin MytomState = iota
	stRequest
	stError
)

type Mytom struct {
	email        string
	password     string
	state        MytomState
	lasterr      error
	activitiesID []string
}

func NewMyTom(e, p string) *Mytom {
	r := Mytom{
		email:        e,
		password:     p,
		activitiesID: make([]string, 0),
	}
	return &r
}

func (mt *Mytom) DownloadFit(destDir string) error {
	log.Println("DownloadFit using target dir ", destDir)
	// create a new collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	// authenticate
	mt.state = stLogin
	c.OnRequest(func(r *colly.Request) {
		if mt.state != stLogin {
			return
		}
		log.Println("OnRequest CB for login... ")
		r.Headers.Set("Accept", "application/json, text/plain, */*")
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
	})
	//login is not a form but a json payload
	payload := fmt.Sprintf("{\"email\": \"%s\",\"password\": \"%s\"}", mt.email, mt.password)
	err := c.PostRaw("https://mysports.tomtom.com/service/webapi/v2/auth/user/login", []byte(payload))
	if err != nil {
		mt.state = stError
		return fmt.Errorf("login err: %v", err)
	} else {
		log.Println("Login ok")
	}
	mt.state = stRequest

	for _, actvno := range mt.activitiesID {
		//actvno := "205727472" //"531746638"
		log.Println("Processing activity ", actvno)
		c.OnResponse(func(r *colly.Response) {
			log.Println("response received", r.StatusCode) //, string(r.Body))
			fnn := fmt.Sprintf("%s/act_%s.fit", destDir, actvno)
			if err := os.WriteFile(fnn, r.Body, 0644); err != nil {
				log.Println("File write error ", err)
				mt.lasterr = err
				mt.state = stError
				return
			}
			log.Println("File written: ", fnn)
		})

		c.OnResponseHeaders(func(r *colly.Response) {
			log.Println("Response headers: ", r)
			if r.StatusCode == 403 {
				log.Println("Something is wrong with AUTH")
				mt.state = stError
			}
		})
		if mt.state == stError {
			return fmt.Errorf("scraper in wrong state %v", mt.lasterr)
		}
		// start scraping
		//c.Visit("https://mysports.tomtom.com/app/activities/")
		// You can see the activity in the browser: uri: fmt.Sprintf("https://mysports.tomtom.com/app/activity/%s/", actvno)
		// You download the activity using the web api
		uri := fmt.Sprintf("https://mysports.tomtom.com/service/webapi/v2/activity/%s?dv=1.3&format=fit", actvno)
		c.Visit(uri)
	}
	log.Println("Processed activities count ", len(mt.activitiesID))
	return nil
}
