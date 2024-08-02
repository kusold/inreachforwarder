package server

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/kusold/inreachforwarder/internal/inreachparser"
	"github.com/spf13/viper"
)

func Start() {
	fmt.Println("Start command executed")

	log.Println(viper.AllSettings())
	srv := NewServer(
		viper.GetString("pagerduty.user-api-token"),
		viper.GetStringSlice("pagerduty.team-ids"),
	)
	srv.Start()
}

func NewServer(apiToken string, teamIds []string) *Server {
	pagerdutyClient := pagerduty.NewClient(apiToken)
	return &Server{
		pagerdutyClient: pagerdutyClient,
	}
}

type Server struct {
	pagerdutyClient *pagerduty.Client
}

func (s *Server) Start() {
	fmt.Println("Server started")
	s.pollForIncidents()
}

func (s *Server) pollForIncidents() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	go func() {
		err := s.handleIncidents()
		if err != nil {
			log.Printf("Error handling incidents: %j\n", err)
		}
		for range ticker.C {
			log.Println("Polling for incidents")
			err := s.handleIncidents()
			if err != nil {
				log.Printf("Error handling incidents: %j\n", err)
			}
		}
	}()
	select {}
}

func (s *Server) handleIncidents() error {
	previousInreachUrl := viper.GetString("storage.inreach-url")

	log.Println("Fetching inreach message url")
	inreachUrl, err := WatchForInReachMessages()
	if err != nil {
		return err
	}

	activeIncidentIds, err := s.notifyIfActiveIncident(inreachUrl)
	if err != nil {
		return err
	}

	// Only read messages if the inreachUrl has changed
	if previousInreachUrl != inreachUrl {
		log.Println("New Inreach Message Received")
		msg, err := inreachparser.ReadMessageFromInReach(inreachUrl)
		if err != nil {
			return err
		}
		log.Printf("Message: %s\n", msg)
		if strings.Contains(strings.ToLower(msg), "ack") {
			log.Printf("Acknowledging incident\n")
			s.AcknowledgeIncidents(activeIncidentIds)
		}
		viper.Set("storage.inreach-url", inreachUrl)
	}
	return viper.WriteConfig()
}

func (s *Server) notifyIfActiveIncident(inreachUrl string) ([]string, error) {
	activeIncidents := s.GetActiveIncidents()
	var activeIncidentIds []string

	if len(activeIncidents) > 0 {
		log.Println("Active incidents found", activeIncidents)
		// Get all active incident IDs
		for _, incident := range activeIncidents {
			activeIncidentIds = append(activeIncidentIds, incident.ID)
		}
		// Get all incident IDs that have been notified
		notifiedIncidents := viper.GetStringSlice("storage.notified-incidents")
		// Find the intersection of the two lists
		var newIncidents []string
		for _, incidentId := range activeIncidentIds {
			if !contains(notifiedIncidents, incidentId) {
				newIncidents = append(newIncidents, incidentId)
			}
		}
		if len(newIncidents) > 0 {
			err := inreachparser.SendMessageToInReach(inreachUrl, fmt.Sprintf("There are active incidents: %v", newIncidents))
			if err != nil {
				return nil, err
			}
			log.Printf("Successfully notified about active incidents: %v\n", newIncidents)
		}

		// Persist the incidents we've notified
		viper.Set("storage.notified-incidents", append(notifiedIncidents, activeIncidentIds...))
	}
	return activeIncidentIds, nil
}

// contains checks if a slice contains a specific element
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
