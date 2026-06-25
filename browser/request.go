package browser

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

func NewClientRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	ApplyCookiesToReq(req)

	req.Header.Set("User-Agent",
		"Mozilla/5.0 (X11; Linux x86_64; rv:152.0) Gecko/20100101 Firefox/152.0")

	return req, nil
}

func PerformPost(url string, jsonData []byte) error {
	req, err := NewClientRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(dump))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println("Status:", resp.Status)
	fmt.Println("Body:", string(body))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("сервер вернул %s", resp.Status)
	}

	return nil
}
