package jokester

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Jokester struct {
	responseCh chan string
	jokeCh     chan string
	nameCh     chan fullName
}

type jokeResponse struct {
	Typ   string `json:"type"`
	Value []struct {
		ID         int      `json:"id"`
		Joke       string   `json:"joke"`
		Categories []string `json:"categories"`
	} `json:"value"`
}

type nameResponse []struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Gender  string `json:"gender"`
	Region  string `json:"region"`
}

type fullName struct {
	first string
	last  string
}

const (
	nameURL            = "http://uinames.com/api/?amount=100"
	jokeURL            = "http://api.icndb.com/jokes/random/100?firstName=John&lastName=Doe&limitT=[nerdy]"
	responseBufferSize = 500
	nameBufferSize     = 100
	jokeBufferSize     = 100
)

func makeHTTPRequest(url string) (*http.Response, error) {
	if url == "" {
		return nil, fmt.Errorf("Empty URL in HTTP request")
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non 200")
	}

	return resp, nil
}

func (j *Jokester) HandleNameJoke(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(<-j.responseCh))
}

func buildNameJoke(firstName, lastName, joke string) string {
	joke = strings.Replace(joke, "John", firstName, -1)
	joke = strings.Replace(joke, "Doe", lastName, -1)
	return joke
}

func (j *Jokester) buildJokes() {
	for {
		fullName := <-j.nameCh
		joke := <-j.jokeCh

		if fullName.first != "" && fullName.last != "" && joke != "" {
			joke = strings.Replace(joke, "John", fullName.first, -1)
			j.responseCh <- strings.Replace(joke, "Doe", fullName.last, -1)
		}
	}
}

func (j *Jokester) checkResponseBuffer() {
	t := time.NewTicker(100 * time.Millisecond)
	for {
		if len(j.responseCh) < 401 {
			err := j.makeJokeReqs()
			if err != nil {
				log.Printf("Err from makeJokeReqs: %v", err)
			}
		}
		<-t.C
	}
}

func (j *Jokester) makeJokeReqs() error {
	var nr nameResponse
	var jr jokeResponse

	nameResp, err := makeHTTPRequest(nameURL)
	if err != nil {
		return err
	}

	jokeResp, err := makeHTTPRequest(jokeURL)
	if err != nil {
		return err
	}

	if nameResp != nil && jokeResp != nil {
		json.NewDecoder(nameResp.Body).Decode(&nr)
		json.NewDecoder(jokeResp.Body).Decode(&jr)

		for _, resp := range nr {
			j.nameCh <- fullName{first: resp.Name, last: resp.Surname}
		}

		for _, resp := range jr.Value {
			j.jokeCh <- resp.Joke
		}
	}
	return nil
}

func (j *Jokester) Deinit() {
	log.Println("deinit jokester")
	if j.jokeCh == nil || j.nameCh == nil || j.responseCh == nil {
		return
	}

	close(j.jokeCh)
	close(j.nameCh)
	close(j.responseCh)
}

func (j *Jokester) Init() error {
	var err error
	log.Println("init jokester")
	j.responseCh = make(chan string, responseBufferSize)
	j.jokeCh = make(chan string, jokeBufferSize)
	j.nameCh = make(chan fullName, nameBufferSize)

	go j.buildJokes()

	//fill buffered response channel
	for i := 0; i < 5; i++ {
		err = j.makeJokeReqs()
		if err != nil {
			return err
		}
	}

	go j.checkResponseBuffer()

	return err
}
