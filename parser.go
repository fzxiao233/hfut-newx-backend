package main

import "strings"

type Parser struct {
}

//type TimePlace struct {
//	Lesson string `json:"lesson"`
//	Place string `json:"place"`
//	Week string `json:"week"`
//}

func (p *Parser) ParseTimePlace(selectResults []map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for i, aResult := range selectResults {
		classes := strings.Split(strings.Replace(aResult["courseTimePlace"].(string), "\n", "", -1), ";")
		var timePlace []map[string]interface{}
		for _, class := range classes {
			splitClass := strings.Split(class, " ")
			if len(splitClass) > 5 {
				splitClass = splitClass[1:]
			}
			if len(splitClass) == 5 {
				timePlace = append(timePlace, map[string]interface{}{
					"lesson": splitClass[1] + "," + splitClass[2],
					"place":  splitClass[4],
					"week":   splitClass[0],
				})
			}
		}
		selectResults[i]["timePlace"] = timePlace
		result = append(result, selectResults[i])
	}
	return result
}
