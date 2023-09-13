package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

func main() {
	// create a new collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	// authenticate
	state := "login"
	c.OnRequest(func(r *colly.Request) {
		if state != "login" {
			return
		}
		log.Println("OnRequest CB for login... ")
		r.Headers.Set("Accept", "application/json, text/plain, */*")
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
	})
	//login is not a form but a json payload
	payload := "<get this from settings.json>"
	err := c.PostRaw("https://mysports.tomtom.com/service/webapi/v2/auth/user/login", []byte(payload))
	if err != nil {
		log.Fatalln("Login err:", err)
	}
	state = "request"
	actvno := "205727472" //"531746638"

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode) //, string(r.Body))
		fnn := fmt.Sprintf("./dest/act_%s.fit", actvno)
		os.WriteFile(fnn, r.Body, 0644)
		log.Println("File written: ", fnn)
	})

	c.OnResponseHeaders(func(r *colly.Response) {
		log.Println("Response headers: ", r)
		if r.StatusCode == 403 {
			log.Println("Something is wrong with AUTH")
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		log.Println("link is ", link)
		// // If link start with browse or includes either signup or login return from callback
		// if !strings.HasPrefix(link, "/browse") || strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 {
		// 	return
		// }
		// // start scaping the page under the link found
		// e.Request.Visit(link)
	})

	// start scraping
	//c.Visit("https://mysports.tomtom.com/app/activities/")
	// You can see the activity in the browser: uri: fmt.Sprintf("https://mysports.tomtom.com/app/activity/%s/", actvno)
	// You download the activity using the web api
	uri := fmt.Sprintf("https://mysports.tomtom.com/service/webapi/v2/activity/%s?dv=1.3&format=fit", actvno)
	c.Visit(uri)
}
