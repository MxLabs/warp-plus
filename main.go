package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}
	return string(b)
}

func main() {
	var ids []string
	flag.Func("ids", "List of ids", func(flagValue string) error {
		ids = strings.Fields(flagValue)
		return nil
	})
	flag.Parse()
	if len(ids) == 0 {
		log.Fatal("No ids specified")
	}

	for {
		for _, id := range ids {
			go func(accountId string) {
				if err := register(accountId); err == nil {
					log.Printf("1 GB has been successfully added to %s account.\n", accountId)
				}
			}(id)
		}
		time.Sleep(time.Second * 15)
	}
}

func register(referrer string) error {

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	installId := generateCode(22)

	data := Payload{
		Key:         fmt.Sprintf("%s=", generateCode(43)),
		InstallID:   installId,
		Tos:         fmt.Sprintf("%s+02:00", time.Now().Format("2006-01-02T15:04:05.000")),
		FcmToken:    fmt.Sprintf("%s:APA91b%s", installId, generateCode(134)),
		Referrer:    referrer,
		WarpEnabled: false,
		Type:        "Android",
		Locale:      "es_ES",
	}

	payload, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	url := fmt.Sprintf("https://api.cloudflareclient.com/v0a%003d/reg", rnd.Intn(999))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("User-Agent", "okhttp/3.12.1")
	req.Header.Set("Host", "api.cloudflareclient.com")
	req.Header.Set("Connection", "Keep-Alive")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	result := resp.StatusCode == 200
	if !result {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}

type Payload struct {
	Key         string `json:"key"`
	InstallID   string `json:"install_id"`
	FcmToken    string `json:"fcm_token"`
	Referrer    string `json:"referrer"`
	WarpEnabled bool   `json:"warp_enabled"`
	Tos         string `json:"tos"`
	Type        string `json:"type"`
	Locale      string `json:"locale"`
}
