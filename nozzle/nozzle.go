package nozzle

import (
	"errors"
	"log"
	"time"

	"github.com/cloudfoundry/sonde-go/events"
)

type Nozzle interface {
	Run(flushWindow time.Duration) error
}

type ForwardingNozzle struct {
	client             Client
	eventSerializer    EventSerializer
	includedEventTypes map[events.Envelope_EventType]bool
	eventsChannel      <-chan *events.Envelope
	errorsChannel      <-chan error
	batch              []interface{}
	logger             *log.Logger
}

type Client interface {
	PostBatch([]interface{}) error
}

type EventSerializer interface {
	BuildHttpStartStopEvent(event *events.Envelope) interface{}
	BuildLogMessageEvent(event *events.Envelope) interface{}
	BuildValueMetricEvent(event *events.Envelope) interface{}
	BuildCounterEvent(event *events.Envelope) interface{}
	BuildErrorEvent(event *events.Envelope) interface{}
	BuildContainerEvent(event *events.Envelope) interface{}
}

func NewForwarder(clientlient Client, eventSerializer EventSerializer, selectedEventTypes []events.Envelope_EventType, eventsChannel <-chan *events.Envelope, errors <-chan error, logger *log.Logger) Nozzle {
	nozzle := &ForwardingNozzle{
		client:          clientlient,
		eventSerializer: eventSerializer,
		eventsChannel:   eventsChannel,
		errorsChannel:   errors,
		batch:           make([]interface{}, 0),
		logger:          logger,
	}

	nozzle.includedEventTypes = map[events.Envelope_EventType]bool{
		events.Envelope_HttpStartStop:   false,
		events.Envelope_LogMessage:      false,
		events.Envelope_ValueMetric:     false,
		events.Envelope_CounterEvent:    false,
		events.Envelope_Error:           false,
		events.Envelope_ContainerMetric: false,
	}
	for _, selectedEventType := range selectedEventTypes {
		nozzle.includedEventTypes[selectedEventType] = true
	}

	return nozzle
}

func (s *ForwardingNozzle) Run(flushWindow time.Duration) error {
	ticker := time.Tick(flushWindow)
	for {
		select {
		case event, ok := <-s.eventsChannel:
			if !ok {
				return errors.New("eventsChannel channel closed")
			}
			s.handleEvent(event)
		case err, ok := <-s.errorsChannel:
			if !ok {
				return errors.New("errorsChannel closed")
			}
			s.handleError(err)
		case <-ticker:
			if len(s.batch) > 0 {
				s.logger.Printf("Posting %d events", len(s.batch))
				err := s.client.PostBatch(s.batch)
				if err != nil {
					return err
				}
				s.batch = make([]interface{}, 0)
			} else {
				s.logger.Print("No events to post")
			}
		}
	}
}

func (s *ForwardingNozzle) handleEvent(envelope *events.Envelope) {
	var event interface{} = nil

	eventType := envelope.GetEventType()
	if !s.includedEventTypes[eventType] {
		return
	}

	switch eventType {
	case events.Envelope_HttpStartStop:
		event = s.eventSerializer.BuildHttpStartStopEvent(envelope)
	case events.Envelope_LogMessage:
		event = s.eventSerializer.BuildLogMessageEvent(envelope)
	case events.Envelope_ValueMetric:
		event = s.eventSerializer.BuildValueMetricEvent(envelope)
	case events.Envelope_CounterEvent:
		event = s.eventSerializer.BuildCounterEvent(envelope)
	case events.Envelope_Error:
		event = s.eventSerializer.BuildErrorEvent(envelope)
	case events.Envelope_ContainerMetric:
		event = s.eventSerializer.BuildContainerEvent(envelope)
	}

	if event != nil {
		s.batch = append(s.batch, event)
	}
}

func (s *ForwardingNozzle) handleError(err error) {
	s.logger.Printf("Error from firehose", err)
}
