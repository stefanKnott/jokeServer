//TODO: finish writing unit tests to get coverage 80%+
//buildNameJoke and makeHTTPRequest tests would be good starts
package jokester

import "testing"

func getJokester() *Jokester {
	return &Jokester{
		responseCh: make(chan string, responseBufferSize),
		nameCh:     make(chan fullName, nameBufferSize),
		jokeCh:     make(chan string, jokeBufferSize),
	}
}

func TestInit(t *testing.T) {
	j := getJokester()
	err := j.Init()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeinit(t *testing.T) {
	j := getJokester()
	j.Deinit()

	if nj := <-j.responseCh; nj != "" {
		t.Fatalf("Deinit did not complete correctly")
	}

	if n := <-j.nameCh; n.first != "" || n.last != "" {
		t.Fatalf("Deinit did not complete correctly")
	}

	if j := <-j.jokeCh; j != "" {
		t.Fatalf("Deinit did not complete correctly")
	}
}
