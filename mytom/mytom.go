package mytom

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
)

type MytomState int

const (
	stLogin MytomState = iota
	stRequest
	stError
	stOK
)

type InfoDownFit struct {
	state   MytomState
	lasterr error
	actvno  string
}

type Mytom struct {
	email        string
	password     string
	state        MytomState
	activitiesID []string
	chuncksize   int
}

func NewMyTom(e, p string, size int) *Mytom {
	r := Mytom{
		email:        e,
		password:     p,
		chuncksize:   size,
		activitiesID: make([]string, 0),
	}
	return &r
}

func (mt *Mytom) UseHardCodedIds() {
	mt.activitiesID = append(mt.activitiesID, "205727472", "531746638", "108044668", "108044654")
	log.Println("using hard coded array activity ids", len(mt.activitiesID))
}

func (mt *Mytom) UseThisIdOnly(id string) {
	mt.activitiesID = make([]string, 1)
	mt.activitiesID[0] = id
	log.Println("using only one id ", id)
}

func (mt *Mytom) DownloadFit(destDir string) error {
	log.Println("DownloadFit using target dir ", destDir)
	start := time.Now()
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
	chRes := make(chan *InfoDownFit)

	c_para := 0
	chunkNr := 0
	chunkSize := mt.chuncksize
	var res *InfoDownFit
	log.Printf("Request for %d activities in download chunk size %d\n", len(mt.activitiesID), chunkSize)

	for ix, actvno := range mt.activitiesID {
		go processActivity(actvno, c, destDir, chRes)

		c_para += 1
		if c_para == chunkSize || (ix == len(mt.activitiesID)-1) {
			// blocking and wait for the chunk download termination
			chunkNr += 1
			log.Printf("[ix => %d] blocking in chunk %d with %d download processes \n", ix, chunkNr, c_para)
			chTimeout := make(chan struct{})
			timeout := 60 * time.Second
			time.AfterFunc(timeout, func() {
				chTimeout <- struct{}{}
			})
			num_in_chunk := c_para
		loop:
			for {
				select {
				case res = <-chRes:
					if res.state == stError {
						return res.lasterr
					}
					c_para -= 1
					if c_para == 0 {
						log.Printf("[ix => %d] chunk %d OK (size %d)\n", ix, chunkNr, num_in_chunk)
						break loop // continue with the next chunk
					}
				case <-chTimeout:
					// something wrong
					return fmt.Errorf("download timeout error")
				}
			}
		}
	}
	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("Processed %d Fit activities, time duration %v", len(mt.activitiesID), elapsed)
	return nil
}

func processActivity(actvno string, corig *colly.Collector, destDir string, chRes chan *InfoDownFit) {
	log.Printf("[%s] start processing", actvno)
	infoFit := InfoDownFit{actvno: actvno, state: stOK}
	sent := false
	c := corig.Clone() // reset callbacks but not auth cookies so that we don't need a new login
	c.OnResponse(func(r *colly.Response) {
		log.Printf("[%s] response received status %d size %d", infoFit.actvno, r.StatusCode, len(r.Body))
		fnn := fmt.Sprintf("%s/act_%s.fit", destDir, actvno)
		if err := os.WriteFile(fnn, r.Body, 0644); err != nil {
			log.Println("File write error ", err)
			infoFit.lasterr = err
			infoFit.state = stError
		}
		chRes <- &infoFit
		sent = true
		log.Printf("[%s] file written: %s", infoFit.actvno, fnn)
	})

	c.OnResponseHeaders(func(r *colly.Response) {
		if r.StatusCode != 200 {
			log.Printf("[%s] response headers %d is suspect", infoFit.actvno, r.StatusCode)
		}
		if r.StatusCode == 403 {
			log.Printf("[%s] something is wrong with AUTH", actvno)
			infoFit.state = stError
			infoFit.lasterr = fmt.Errorf("error with AUTH inside the dowload (login was successfully?)")
			chRes <- &infoFit
			sent = true
		}
	})

	c.OnError(func(e *colly.Response, err error) {
		log.Println("Error on scrap", err)
		if !sent {
			infoFit.lasterr = err
			infoFit.state = stError
			chRes <- &infoFit
			sent = true
		}
	})

	// start scraping
	//c.Visit("https://mysports.tomtom.com/app/activities/")
	// You can see the activity in the browser: uri: fmt.Sprintf("https://mysports.tomtom.com/app/activity/%s/", actvno)
	// You download the activity using the web api
	uri := fmt.Sprintf("https://mysports.tomtom.com/service/webapi/v2/activity/%s?dv=1.3&format=fit", actvno)
	c.Visit(uri)

	if !sent {
		infoFit.lasterr = fmt.Errorf("[%s] download was somehow not working", infoFit.actvno)
		infoFit.state = stError
		chRes <- &infoFit
		sent = true
	}
}
