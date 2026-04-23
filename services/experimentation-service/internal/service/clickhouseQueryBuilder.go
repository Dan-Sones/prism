package service

import (
	"errors"
	eventModel "experimentation-service/internal/model/event"
	"fmt"
	"time"

	"github.com/Dan-Sones/prismdbmodels/model/event"
	"github.com/Dan-Sones/prismdbmodels/model/metric"
)

type QueryBuilder interface {
	BuildQueryFor(experimentKey string, m metric.EnrichedMetric) (eventModel.QueryString, error)
}

type ClickhouseQueryBuilder struct {
}

// Kind of unnessary struct?
func NewClickhouseQueryBuilder() *ClickhouseQueryBuilder {
	return &ClickhouseQueryBuilder{}
}

func (c *ClickhouseQueryBuilder) BuildQueryFor(experimentKey string, m metric.EnrichedMetric, startTime, endTime time.Time) (eventModel.QueryString, error) {
	if len(m.MetricComponents) == 0 {
		return "", errors.New("metric must have at least one component")
	}

	if m.MetricType == metric.MetricTypeRatio {
		return c.buildForRatioMetric(experimentKey, m, startTime, endTime)
	}

	return "", nil

}

func (c *ClickhouseQueryBuilder) buildForRatioMetric(experimentKey string, m metric.EnrichedMetric, startTime, endTime time.Time) (eventModel.QueryString, error) {
	var query eventModel.ClickhouseQuery

	query.WHERE = append(query.WHERE, "experiment_key = '"+experimentKey+"'")
	query.WHERE = append(query.WHERE, c.BuildInEventKeyWhere(m))
	query.WHERE = append(query.WHERE, c.BuildTimeRangeWhere(startTime, endTime))

	query.SELECT = append(query.SELECT, "variant_key")
	for _, component := range m.MetricComponents {
		switch component.Role {
		case metric.ComponentRoleNumerator:
			numeratorSelect, err := c.BuildSelectForNumeratorComponent(component)
			if err != nil {
				return "", errors.New("error building select for numerator component: " + err.Error())
			}
			query.SELECT = append(query.SELECT, numeratorSelect)
		case metric.ComponentRoleDenominator:
			denominatorSelect, err := c.BuildSelectForDenominatorComponent(component)
			if err != nil {
				return "", errors.New("error building select for denominator component: " + err.Error())
			}
			query.SELECT = append(query.SELECT, denominatorSelect)
		}
	}

	query.FROM = "events"

	query.GroupBy = append(query.GroupBy, "variant_key")

	return query.BuildQueryString(), nil
}

func (c *ClickhouseQueryBuilder) BuildSelectForNumeratorComponent(component metric.EnrichedMetricComponent) (string, error) {
	switch component.AggregationOperation {
	case metric.AggregationOperationCountDistinct:
		return c.BuildSelectItemForCountDistinct(component)
	default:
		return "", fmt.Errorf("unsupported aggregation operation: %s", component.AggregationOperation)
	}
}

func (c *ClickhouseQueryBuilder) BuildSelectForDenominatorComponent(component metric.EnrichedMetricComponent) (string, error) {
	switch component.AggregationOperation {
	case metric.AggregationOperationCountDistinct:
		return c.BuildSelectItemForCountDistinct(component)
	default:
		return "", fmt.Errorf("unsupported aggregation operation: %s", component.AggregationOperation)
	}
}

func (c *ClickhouseQueryBuilder) BuildSelectItemForCountDistinct(component metric.EnrichedMetricComponent) (string, error) {
	// See if the count distinct is on a system column or an event field
	if component.SystemColumnName != nil {
		return fmt.Sprintf("uniqExactIf(%s, event_key = '%s') AS %s", *component.SystemColumnName, component.EventType.EventKey, component.Role), nil
	} else if component.AggregationField != nil {
		// then identify the type of the field so we know which map to look within the event table
		switch component.AggregationField.DataType {
		case event.DataTypeString:
			return fmt.Sprintf("uniqExactIf(string_properties['%s'], event_key = '%s') AS %s", component.AggregationField.FieldKey, component.EventType.EventKey, component.Role), nil
		case event.DataTypeFloat:
			return fmt.Sprintf("uniqExactIf(float_properties['%s'], event_key = '%s') AS %s", component.AggregationField.FieldKey, component.EventType.EventKey, component.Role), nil
		case event.DataTypeInt:
			return fmt.Sprintf("uniqExactIf(int_properties['%s'], event_key = '%s') AS %s", component.AggregationField.FieldKey, component.EventType.EventKey, component.Role), nil
		}
	}

	return "", fmt.Errorf("invalid component configuration, component does not have SystemColumnName OR AggregationField: %v", component)
}

func (c *ClickhouseQueryBuilder) BuildInEventKeyWhere(m metric.EnrichedMetric) string {
	var eventKeys []string
	for _, component := range m.MetricComponents {
		eventKeys = append(eventKeys, component.EventType.EventKey)
	}

	var inClause string
	for i, key := range eventKeys {
		inClause += "'" + key + "'"
		if i < len(eventKeys)-1 {
			inClause += ", "
		}
	}
	return fmt.Sprintf("event_key IN (%s)", inClause)
}

func (c *ClickhouseQueryBuilder) BuildTimeRangeWhere(startTime, endTime time.Time) string {
	const layout = "2006-01-02 15:04:05"
	return fmt.Sprintf(
		"sent_at >= '%s' AND sent_at <= '%s'", startTime.UTC().Format(layout), endTime.UTC().Format(layout),
	)
}
