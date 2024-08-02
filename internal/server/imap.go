package server

import (
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/spf13/viper"
)

func WatchForInReachMessages() (string, error) {
	return getMostRecentInReachUrl()
}

func getMostRecentInReachUrl() (string, error) {
	c, err := imapclient.DialTLS(viper.GetString("imap.server"), nil)
	if err != nil {
		return "", err
	}
	defer c.Close()

	if err := c.Login(viper.GetString("imap.user"), viper.GetString("imap.password")).Wait(); err != nil {
		return "", err
	}

	selectedMbox, err := c.Select("INBOX", nil).Wait()
	if err != nil {
		return "", err
	}
	log.Printf("INBOX contains %v messages", selectedMbox.NumMessages)

	inreachUrl := ""
	if selectedMbox.NumMessages > 0 {
		messages, err := searchForInReachMessages(c)
		if err != nil {
			return "", err
		}
		inreachUrl = getInReachUrlFromMessages(messages)
	}

	if err := c.Logout().Wait(); err != nil {
		return "", err
	}
	return inreachUrl, nil
}

func searchForInReachMessages(c *imapclient.Client) ([]*imapclient.FetchMessageBuffer, error) {
	searchCriteria := &imap.SearchCriteria{
		Header: []imap.SearchCriteriaHeaderField{
			{Key: "Subject", Value: "inReach message from Katie"}, // TODO: Make configurable
		},
	}
	res, err := c.UIDSearch(searchCriteria, &imap.SearchOptions{}).Wait()
	if err != nil {
		log.Fatalf("failed to search INBOX: %v", err)
	}
	// for seq := range res.AllSeqNums() {
	// 	fetchOptions := &imap.FetchOptions{Envelope: true}
	// 	message, err := c.Fetch(seq, fetchOptions).Collect()
	// 	log.Println(message)
	// }
	// seqSet := imap.SeqSetNum(2)
	seqSet := res.All
	fetchOptions := &imap.FetchOptions{
		Envelope:    true,
		BodySection: []*imap.FetchItemBodySection{{}},
	}
	return c.Fetch(seqSet, fetchOptions).Collect()
}

func getInReachUrlFromMessages(messages []*imapclient.FetchMessageBuffer) string {
	inreachUrl := ""

	// Sort messages by most recent
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Envelope.Date.After(messages[j].Envelope.Date)
	})
	for _, msg := range messages {
		// log.Println(msg)
		log.Println(msg.Envelope.Subject)
		log.Println(msg.Envelope.From)

		for _, section := range msg.BodySection {
			body := string(section)

			// TODO: Requires the email address to end in .com
			re := regexp.MustCompile(`(?misU)(https:\/\/.*\/textmessage\/txtmsg.*com)`)
			match := re.FindStringSubmatch(body)
			if len(match) > 0 {
				// Strip newlines from the the URL

				// log.Println(match[0])
				cleaned := strings.ReplaceAll(match[0], "=\r\n", "")
				inreachUrl = strings.ReplaceAll(cleaned, "3D", "")
				// log.Println(inreachUrl)
				return inreachUrl
			}
		}
	}
	return inreachUrl
}
