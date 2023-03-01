package grabber

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Fetcher interface {
	Get(req *Request) ([]byte, error)
}

type baseFetch struct{}

func (f *baseFetch) Get(req *Request) ([]byte, error) {
	resp, err := http.Get(req.URL)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code:%d", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := determinEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	return ioutil.ReadAll(utf8Reader)
}

type browserFetch struct{}

func (b *browserFetch) Get(req *Request) ([]byte, error) {
	client := &http.Client{}

	newReq, err := http.NewRequest("GET", req.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("request URL error:%v\n", err)
	}

	newReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.57")

	resp, err := client.Do(newReq)
	if err != nil {
		return nil, err
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := determinEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	return ioutil.ReadAll(utf8Reader)
}

func determinEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)

	if err != nil {
		fmt.Printf("fetch failed:%v\n", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
