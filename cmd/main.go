package main

import (
	"fmt"
	"net/url"

	rqp "github.com/nfidel/rest-query-parser/v2"
)

func main() {

	// Filter is parameter provided in the Query part of the URL
	//   The lib handles system filters:
	//     * fields - list of fields separated by comma (",") for SELECT statement. Should be validated.
	//     * sort   - list of fields separated by comma (",") for ORDER BY statement. Should be validated. Could includes prefix +/- which means ASC/DESC sorting. Eg. &sort=-id will be ORDER BY id DESC.
	//     * limit  - number for LIMIT statement. Should be greater then 0 by default.
	//     * offset - number for OFFSET statement. Should be greater then or equal to 0 by default.
	//   and user defined filters.
	//
	// Validation is a function for validate some Filter
	//
	// Field is enumerated in the Filter "fields" field which lib must put into SELECT statement.

	url, _ := url.Parse("http://localhost/?fields=frequency_cap&pace.pace[is]=null&page_size=10&page=0") //&flights[is]=null&frequency_cap[is]=null&frequency_cap.impressions.test[is]=NULL")
	q := rqp.NewQV(url.Query(),
		rqp.Validations{
			"fields": rqp.In("pace", "frequency_cap", "flights", "targeting"),
		},
		rqp.QueryDatabaseMap{
			"pace.pace":                 {Name: "global_bid_rate", Table: "campaign_pace", Type: rqp.FieldTypeFloat},
			"pace.pacing_strategy":      {Name: "pacing_strategy", Table: "campaign", Type: rqp.FieldTypeString},
			"frequency_cap":             {Name: "frequency_cap", Table: "campaign", Type: rqp.FieldTypeCustom},
			"frequency_cap.impressions": {Name: "frequency_cap.impressions", Table: "campaign", Type: rqp.FieldTypeInt, IsNested: true}, //true},
			"frequency_cap.duration":    {Name: "frequency_cap.duration", Table: "campaign", Type: rqp.FieldTypeInt, IsNested: true},    //true},
			//"frequency_cap.impressions.test": {Name: "frequency_cap.impressions.test", Table: "campaign", Type: rqp.FieldTypeFloat},
			"targeting.geo": {Table: "targeting", Name: "geo_targeting", Type: rqp.FieldTypeCustom},
			"targeting.ip":  {Table: "targeting", Name: "ip_targeting", Type: rqp.FieldTypeCustom},
		})
	q.IgnoreUnknownFilters(false)
	err := q.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Println(q.SQL("campaign")) // SELECT * FROM table WHERE id = ? AND i = ? AND s = ? AND (email LIKE ? OR name LIKE ?) ORDER BY name, id DESC LIMIT 10
}
