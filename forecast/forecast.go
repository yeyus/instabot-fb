package forecast

import (
	"github.com/doublerebel/bellows"
	"github.com/nickpankow/yql"
	"log"
)

func GetForecast(location string) (c map[string]interface{}, err error) {
	y := yql.YQL{"https://query.yahooapis.com/v1/public/yql", "http://datatables.org/alltables.env", "json"}
	var r *yql.Response
	r, err = y.Query("select item.condition from weather.forecast where woeid in (select woeid from geo.places(1) where text=\"" + location + "\")")

	if err != nil {
		log.Printf("error retrieving forecast for %s: %s", location, err)
		return nil, err
	}

	log.Printf("forecast response: %v", r)

	return bellows.Flatten(r.Results), err
}
