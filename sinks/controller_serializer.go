package sinks

import (
	"fmt"

	"github.com/cloudfoundry/sonde-go/events"
)

type DataPoint struct {
	Metric  string
	Value   int64
	Allowed bool
	MetricType string
}

type ControllerEventSerializer struct {
	tier string
}

func NewControllerEventSerializer(tier_name string) *ControllerEventSerializer {
	return &ControllerEventSerializer{tier: tier_name}
}

func (w *ControllerEventSerializer) BuildHttpStartStopEvent(event *events.Envelope) interface{} {
	allowed := false
	return &DataPoint{Metric: "", Value: int64(0), Allowed: allowed, MetricType: HttpStartStopEvent}
}

func (w *ControllerEventSerializer) BuildLogMessageEvent(event *events.Envelope) interface{} {
	allowed := false
	return &DataPoint{Metric: "", Value: int64(0), Allowed: allowed, MetricType: LogMessageEvent}
}

func (w *ControllerEventSerializer) BuildValueMetricEvent(event *events.Envelope) interface{} {
	valueMetric := event.GetValueMetric()
	return w.makeDataPoint(valueMetric.GetName(), int64(valueMetric.GetValue()), event, ValueMetricEvent)
}

func (w *ControllerEventSerializer) BuildCounterEvent(event *events.Envelope) interface{} {
	counterMetric := event.GetCounterEvent()
	return w.makeDataPoint(counterMetric.GetName(), int64(counterMetric.GetDelta()), event, CounterEvent)
}

func (w *ControllerEventSerializer) makeDataPoint(name string, value int64, event *events.Envelope, metric_type string) *DataPoint {
	origin := event.GetOrigin()
	alias, present := FilterMetrics(origin, name)
	if present {
		deployment, index, job := event.GetDeployment(), event.GetIndex(), event.GetJob()
		prefix := fmt.Sprintf("Server|Component:%v|Custom Metrics|PCF Firehose Monitor", w.tier)
		metric_name := fmt.Sprintf("%v.%v", origin, name)
		return &DataPoint{
			Metric:  fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", prefix, alias, origin, deployment, job, index, metric_name),
			Value:   value,
			Allowed: present,
			MetricType: metric_type}
	} else {
		return &DataPoint{Metric: "", Value: int64(0), Allowed: present, MetricType: metric_type}
	}
}

func (w *ControllerEventSerializer) BuildErrorEvent(event *events.Envelope) interface{} {
	allowed := false
	return &DataPoint{Metric: "", Value: int64(0), Allowed: allowed, MetricType: ErrorEvent}
}

func (w *ControllerEventSerializer) BuildContainerEvent(event *events.Envelope) interface{} {
	allowed := false
	return &DataPoint{Metric: "", Value: int64(0), Allowed: allowed, MetricType: ContainerEvent}
}

func FilterMetrics(eventOrigin, eventName string) (string, bool) {
	filters, present := MetricFilter[eventOrigin]
	if present {
		alias := MetricAlias[eventOrigin]
		for _, allowedMetric := range filters {
			if allowedMetric == eventName {
				return alias, true
			}
		}
	}
	return "", false
}
