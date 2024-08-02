package inreachparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type MapSharePayload struct {
	ReplyAddress string
	ReplyMessage string
	MessageID    string
	Guid         string
}

func SendMessageToInReach(replyUrl, text string) error {
	payload, err := scrapeMapShareForPayload(replyUrl)
	if err != nil {
		return err
	}
	payload.ReplyMessage = text
	log.Printf("Payload: %+v\n", payload)

	// So far I've seen the host be either "us0.explore.garmin.com" or "eur.explore.garmin.com"
	parsedUrl, err := url.Parse(replyUrl)
	if err != nil {
		return err
	}

	return sendPayloadToInReach(parsedUrl.Host, payload)
}

func ReadMessageFromInReach(replyUrl string) (string, error) {
	return scrapeMapShareForMessage(replyUrl)
}

func scrapeMapShareForMessage(replyUrl string) (string, error) {
	doc, err := scrapeMapShare(replyUrl)
	if err != nil {
		return "", err
	}

	var message string
	// Get a div with class "message-text"
	doc.Find("div.message-text").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
		message = s.Text()
	})
	return message, nil
}

func sendPayloadToInReach(host string, payload *MapSharePayload) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://"+host+"/TextMessage/TxtMsg", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("response Status: %s\n", resp.Status)
	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}

func scrapeMapShareForPayload(url string) (*MapSharePayload, error) {
	doc, err := scrapeMapShare(url)
	if err != nil {
		log.Fatal(err)
	}

	var messageID, guid, replyAddress string
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("id")
		value, _ := s.Attr("value")

		if name == "MessageId" {
			messageID = value
		} else if name == "Guid" {
			guid = value
		} else if name == "ReplyAddress" {
			replyAddress = value
		}
	})

	payload := &MapSharePayload{
		ReplyAddress: replyAddress,
		// ReplyMessage: text,
		MessageID: messageID,
		Guid:      guid,
	}
	return payload, nil
}

func scrapeMapShare(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
