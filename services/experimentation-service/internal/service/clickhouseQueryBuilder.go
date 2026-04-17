package service

import (
	"fmt"

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

func (c *ClickhouseQueryBuilder) buildForRatioMetric(experimentKey string, metric metric.Metric) (string, error) {

	var query Query

	query.WHERE = append(query.WHERE, "experiment_key = '"+experimentKey+"'")
	query.WHERE = append(query.WHERE, c.BuildInEventKeyWhere(metric))

	return "", nil
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
