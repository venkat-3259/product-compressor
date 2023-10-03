package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func PrepareHTTPRequest(v any, reqMethod, reqURL string) (*http.Request, error) {

	byteData, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// log.Println("Request", string(byteData))

	req, err := http.NewRequest(reqMethod, reqURL, bytes.NewBuffer(byteData))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	return req, nil
}

func SendHTTPRequest(req *http.Request) (int, []byte, error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, body, err

}
