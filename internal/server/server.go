package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	opts := pagerduty.ListIncidentsOptions{
		Includes: []string{"acknowledgers", "assignees", "first_trigger_log_entries", "services", "teams", "users"},
		// TeamIDs:   viper.GetStringSlice("pagerduty.team-ids"),
		Since:    time.Now().Add(-time.Hour * 24).Format(time.RFC3339),
		Statuses: []string{"triggered", "acknowledged"},
		//Statuses:  []string{"triggered"},
		Urgencies: []string{"high"},
		UserIDs:   viper.GetStringSlice("pagerduty.user-ids"),
	}
	resp, err := s.pagerdutyClient.ListIncidentsWithContext(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to list incidents: %v", err)
	}

	for _, incident := range resp.Incidents {
		log.Printf("Incident: %+v\n", incident)
		log.Println("Status: ", incident.Status)
		log.Println("Urgency: ", incident.Urgency)
	}
}
