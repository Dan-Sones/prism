package service

import (
	"errors"
	"fmt"

	"github.com/Dan-Sones/prismdbmodels/model/event"
	"github.com/Dan-Sones/prismdbmodels/model/metric"
)

type QueryBuilder interface {
	BuildQueryFor(variantKey, experimentKey string, m metric.Metric) (string, error)
}

type Query struct {
	SELECT  []string
	FROM    string
	WHERE   []string
	GroupBy []string
}

type ClickhouseQueryBuilder struct {
}

func (c *ClickhouseQueryBuilder) BuildQueryForExperimentMetric(experimentKey string, m metric.Metric) (string, error) {

	if m.MetricType == metric.MetricTypeRatio {
		return c.buildForRatioMetric(experimentKey, m)
	}

	return "", nil

}

func (c *ClickhouseQueryBuilder) buildForRatioMetric(experimentKey string, m metric.Metric) (string, error) {
	var query Query

	query.WHERE = append(query.WHERE, "experiment_key = '"+experimentKey+"'")
	query.WHERE = append(query.WHERE, c.BuildInEventKeyWhere(m))

	for _, component := range m.MetricComponents {
		switch component.Role {
		case metric.ComponentRoleNumerator:
			numeratorSelect, err := c.BuildSelectForNumeratorComponent(component)
			if err != nil {
				return "", errors.New("error building select for numerator component: " + err.Error())
			}
			query.SELECT = append(query.SELECT, numeratorSelect)
		case metric.ComponentRoleDenominator:
			query.SELECT = append(query.SELECT, c.BuildSelectForDenominatorComponent(component))
		}
	}

	return "", nil
}

func (c *ClickhouseQueryBuilder) BuildSelectForNumeratorComponent(component metric.MetricComponent) (string, error) {
	switch component.AggregationOperation {
	case metric.AggregationOperationCountDistinct:
		return c.BuildSelectItemForCountDistinct(component)
	default:
		return "", fmt.Errorf("unsupported aggregation operation: %s", component.AggregationOperation)
	}
}

func (c *ClickhouseQueryBuilder) BuildSelectForDenominatorComponent(component metric.MetricComponent) string {

	// TODO: IMPLEMENT
	return ""
}

func (c *ClickhouseQueryBuilder) BuildSelectItemForCountDistinct(component metric.MetricComponent) (string, error) {
	// See if the count distinct is on a system column or an event field
	if component.SystemColumnName != nil {
		return fmt.Sprintf("uniqExactIf(%s, event_key = '%s') AS %s", *component.SystemColumnName, component.EventType.EventKey, component.Role), nil
	} else if component.AggregationField != nil {
		// then identify the type of the field so we know which map to look within the event table
		switch component.AggregationField.DataType {
		case event.DataTypeString:
			return fmt.Sprintf("uniqExactIf(string_properties[%s], event_key = '%s') AS %s", component.AggregationField.FieldKey, component.EventType.EventKey, component.Role), nil
		case event.DataTypeFloat:
			return fmt.Sprintf("uniqExactIf(float_properties[%s], event_key = '%s') AS %s", component.AggregationField.FieldKey, component.EventType.EventKey, component.Role), nil
		case event.DataTypeInt:
			return fmt.Sprintf("uniqExactIf(int_properties[%s], event_key = '%s') AS %s", component.AggregationField.FieldKey, component.EventType.EventKey, component.Role), nil
		}
	}

	return "", fmt.Errorf("invalid component configuration, component does not have SystemColumnName OR AggregationField: %v", component)
}

func (c *ClickhouseQueryBuilder) BuildInEventKeyWhere(m metric.Metric) string {
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
	return fmt.Sprintf("event_key in (%s)", inClause)
}

//SELECT
//variant_key,
//uniqExactIf(user_id, event_key = 'experiment_exposure') AS exposed_users,
//uniqExactIf(user_id, event_key = 'purchase_complete')   AS converted_users
//FROM events
//WHERE experiment_key = 'exp_key'
//AND event_key IN ('experiment_exposure', 'purchase_complete')
//GROUP BY variant_key;
