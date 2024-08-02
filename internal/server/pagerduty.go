package server

import (
	"context"
	"log"
	"time"

	"github.com/PagerDuty/go-pagerduty"
	"github.com/spf13/viper"
)

func (s *Server) GetActiveIncidents() []pagerduty.Incident {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	opts := pagerduty.ListIncidentsOptions{
		Includes: []string{"acknowledgers", "assignees", "first_trigger_log_entries", "services", "teams", "users"},
		Since:    time.Now().Add(-time.Hour * 24).Format(time.RFC3339),
		// Statuses: []string{"triggered", "acknowledged"},
		Statuses:  []string{"triggered"},
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
	return resp.Incidents
}

func (s *Server) AcknowledgeIncidents(incidentIds []string) error {
	// ackOpts := pagerduty.AcknowledgeIncidentOptions{
	// 	IncidentID: incidentId,
	// 	Payload:    inreachUrl,
	// }
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	var opts []pagerduty.ManageIncidentsOptions
	for _, incidentId := range incidentIds {
		inciOpt := pagerduty.ManageIncidentsOptions{
			ID:     incidentId,
			Status: "acknowledged",
			Type:   "incident_reference",
		}
		opts = append(opts, inciOpt)
		log.Println(inciOpt)

	}
	_, err := s.pagerdutyClient.ManageIncidentsWithContext(ctx, viper.GetString("pagerduty.user-email"), opts)
	if err != nil {
		return err
	}
	return nil
}
