package event

import "fmt"

type QueryString string

type Query interface {
	BuildQueryString() QueryString
}

type ClickhouseQuery struct {
	SELECT  []string
	FROM    string
	WHERE   []string
	GroupBy []string
}

func (q *ClickhouseQuery) BuildQueryString() QueryString {
	commaSeperatedSelect := ""
	for i, selectItem := range q.SELECT {
		commaSeperatedSelect += selectItem
		if i < len(q.SELECT)-1 {
			commaSeperatedSelect += ", "
		}
	}

	query := fmt.Sprintf("SELECT %s FROM %s", commaSeperatedSelect, q.FROM)

	andSeperatedWhere := ""
	for i, whereClause := range q.WHERE {
		andSeperatedWhere += whereClause
		if i < len(q.WHERE)-1 {
			andSeperatedWhere += " AND "
		}
	}

	query = fmt.Sprintf("%s WHERE %s", query, andSeperatedWhere)

	if len(q.GroupBy) > 0 {
		commaSeperatedGroupBy := ""
		for i, groupByItem := range q.GroupBy {
			commaSeperatedGroupBy += groupByItem
			if i < len(q.GroupBy)-1 {
				commaSeperatedGroupBy += ", "
			}
		}
		query = fmt.Sprintf("%s GROUP BY %s", query, commaSeperatedGroupBy)
	}

	query = fmt.Sprintf("%s;", query)

	return QueryString(query)
}
