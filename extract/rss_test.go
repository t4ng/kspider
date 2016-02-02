package extract

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestParseRss(t *testing.T) {
	url := "http://36kr.com/feed"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36")
	req.Header.Set("Cache-Control", "max-age=0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http get url error: %s", err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	pageInfo, err := ExtractRss(url, body)
	if err != nil {
		t.Fatalf("extract error: %s", err)
	}

	jsonText, _ := json.MarshalIndent(pageInfo, "", "\t")
	fmt.Printf("content: %s", jsonText)
}
