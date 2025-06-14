package handlers

import (
	"context"
	"fmt"

	"github.com/mixpanel/mixpanel-go"
	"github.com/sirupsen/logrus"
)

// TODO: move to dependency injection
var globalMixPanelSender *MixpanelSender

func init() {
	globalMixPanelSender = NewMixpanelSender("your-project-token")
}

type MixpanelSender struct {
	cl    *mixpanel.ApiClient
	queue chan *mixpanel.Event
}

func NewMixpanelSender(token string) *MixpanelSender {
	s := &MixpanelSender{
		cl:    mixpanel.NewApiClient(token),
		queue: make(chan *mixpanel.Event, 100000),
	}
	go s.run()
	return s
}

func (m *MixpanelSender) run() {
	for event := range m.queue {
		err := m.sendSingle(context.Background(), event)
		if err != nil {
			logrus.Errorf("Mixpanel event tracking failed: %v", err)
		}
	}
}

func (m *MixpanelSender) sendSingle(ctx context.Context, event *mixpanel.Event) error {
	err := m.cl.Track(ctx, []*mixpanel.Event{event})
	if err != nil {
		return fmt.Errorf("mixpanel event tracking failed: %v", err)
	}
	return nil
}

func (m *MixpanelSender) Send(name string, distinctID string, properties map[string]any) {
	logrus.WithFields(logrus.Fields{
		"name":       name,
		"distinctID": distinctID,
		"properties": properties,
	}).Info("Pseudo sending mixpanel event")
	// TODO: disabled until proper token was integrated
	//m.queue <- m.cl.NewEvent(name, distinctID, properties)
}
