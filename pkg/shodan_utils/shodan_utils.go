package shodan_utils

import (
	"io/ioutil"

	"github.com/Jeffail/gabs"
	"github.com/jayateertha043/goshod/pkg/httpclient"
)

func ShodanSearchHostCount(params map[string]string) (statuscode int, hostcount int, responsebody string) {
	sres, err := httpclient.GetRequestP("https://api.shodan.io/shodan/host/count", params, nil, 25)

	if err != nil {
		return 500, -100, "Error: Unable to make request to shodan"
	}

	defer sres.Body.Close()
	status_code := sres.StatusCode
	sresbytes, _ := ioutil.ReadAll(sres.Body)
	response_body := string(sresbytes)
	if status_code == 200 {
		res, err := gabs.ParseJSON([]byte(sresbytes))
		if err != nil {
			response_body = "Error: Unable to parse result from shodan"
			return 500, -101, response_body
		}
		total_exist := res.ExistsP("total")
		if total_exist {
			total := res.Path("total").Data().(float64)
			pages := int(((total - 1) / 100)) + 1
			return 200, pages, response_body
		}
	} else {
		return 500, -100, "Error: Unknown"
	}

	return 500, -100, "Error: Unknown"
}

func GetShodanResultForPage(params map[string]string) (status_code int, response_body string) {

	sres, err := httpclient.GetRequestP("https://api.shodan.io/shodan/host/search", params, nil, 25)
	if err != nil {
		return 500, `{"error": "Internal connection error with shodan"}`
	}

	defer sres.Body.Close()
	statuscode := sres.StatusCode
	sresbytes, _ := ioutil.ReadAll(sres.Body)
	responsebody := string(sresbytes)
	return statuscode, responsebody
}
