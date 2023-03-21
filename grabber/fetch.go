package grabber

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/unnull0/crawler/extensions"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type FetchType int

const (
	BaseFetchType FetchType = iota
	BrowserFetchType
)

type Fetcher interface {
	Get(req *Request) ([]byte, error)
}

func NewFetcher(tp FetchType) Fetcher {
	switch tp {
	case BaseFetchType:
		return &baseFetch{}
	case BrowserFetchType:
		return &BrowserFetch{}
	default:
		return &BrowserFetch{}
	}
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

type BrowserFetch struct {
	Timeout time.Duration
}

func (b *BrowserFetch) Get(req *Request) ([]byte, error) {
	client := &http.Client{
		Timeout: b.Timeout,
	}

	newReq, err := http.NewRequest("GET", req.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("request URL error:%v", err)
	}

	if len(req.Task.Cookie) > 0 {
		newReq.Header.Set("Cookie", req.Task.Cookie)
	}
	newReq.Header.Set("User-Agent", extensions.GenerateRandomUA())

	resp, err := client.Do(newReq)
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

func determinEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)

	if err != nil {
		fmt.Printf("fetch failed:%v\n", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
