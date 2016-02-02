package extract

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestParse(t *testing.T) {
	return
	//url := "http://36kr.com/p/5041782.html"
	url := "http://www.huxiu.com/article/135849/1.html"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36")
	req.Header.Set("Cache-Control", "max-age=0")
	//req.Header.Set("Cookie", "__utma=47937701.1057136139.1428327029.1428327029.1428327029.1; kr_stat_uuid=dRW2X24156714; gr_user_id=37dcdb17-7173-4a39-a36c-c80ba4a62b8f; Hm_lvt_e8ec47088ed7458ec32cde3617b23ee3=1449409798; _kr_p_rm=UwTQa8WVDDXu+VLbLscUJFPRMg9ga7/1hAbflkxKNcdmr1b36CM8PHNOd/LkTCySKRUJ+VZ5niS5PkFAcbXU/WDwckSpJPypoIEonvoHx2P4D8k2DPs0hIfyF7967aztp/w+hmmnh4neRE/Ackh94EHrANH+mBnz1Yxx0zqYik/M9sj1x5ePpM6V2IVMSztyliykHW0JzxQC4LMrzJey7LvoG5zMyg8xPxWizErz/TpfJ/7UVsmo81OTWUdZDgioxtqEmqoWXjoJ3NRv5XzzO619x1GeCq9JLwak5Zp7G5dshqTNTCtQ2fYORZbCvcMwVlmmYz0NC9fC2t3m8WpCwzBgdSZooQlPL9Npyq8IQ2TqtbEMe/9p/R3iMpi7Y8FDRbTs/eqB1BQPUIhKUB5bIT0YmcS1+orouQ9zgMZNIjpIDftpt26YpHaKMkSBIm1lp9GixtxQ/KtNgKCEHv4ypmCY5lj6Ka8sjT6RjgnZ0gxXNXci8gsWO2LEhETN+hD4jA8r2XhuBuWyofIlVvLKRoshh9o0iStzqwM/I3bwrQZpOQF4pdyjrIvkoHsxU6YK8JvN84wnJKTy3MHVWyavn7rXEqCqyWO74dKREWvcNqDHrojVn/mREBJh3k7IJ4vbR3CrxXIHd+IIzM6daPBagTeBeemW1TiXt2A+sKdh1TS1c1yu1iYbVOSAQ4aElkVmCEi9+XopFSClWEeLM6D3UXTVeqRdY12ZgLBPB9PZcN84UMXJmDZYXdpDVc7NQHqLOJg/5PP+oyw5ycZHwgMD+/ev4goVBn0+nYbD6SHLCfXcR4M0cVvYsBYW5W0l8mRgINszGCYPuASWcdagj8KDmxbt/u2NlPUSgyMEi16j/FPQ5uDSQavWcRsgb17Bwnmm; aliyungf_tc=AQAAAHYiAjdLNA0ApdAcPD83s0MOm6Ay; krid_user_id=550137; krid_user_version=1001; kr_plus_id=396318; kr_plus_token=0c8962db4dfed2fd95c1cdbed673686d17b29424; _gat=1; _kr_p_se=f14906c3-9a54-41ab-a5cc-8268914158b9; _alicdn_sec=568690cad3ffdeeeb8e248208cf67dbfc2d7d622; Hm_lvt_713123c60a0e86982326bae1a51083e1=1449402863,1449656385,1451630952; Hm_lpvt_713123c60a0e86982326bae1a51083e1=1451659467; _ga=GA1.2.1057136139.1428327029; _krypton_session=ffc2f2bc6a44884c3e9f0a2b586e9632")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http get url error: %s", err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	pageInfo, err := Extract(url, body)
	if err != nil {
		t.Fatalf("extract error: %s", err)
	}

	jsonText, _ := json.MarshalIndent(pageInfo, "", "\t")
	fmt.Printf("content: %s", jsonText)
}
