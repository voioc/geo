package lite

import (
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/voioc/geo/model"
)

type Geo struct{}

func (geo *Geo) Cmd() {

}

// func (geo *Geo) _version() {}
// func (geo *Geo) _update()  {}

func (geo *Geo) Analyze(ip string) (res *model.Location, err error) {
	db, err := geoip2.Open("./db/GeoLite2-City.mmdb")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ipParse := net.ParseIP(ip)
	record, err := db.City(ipParse)
	if err != nil {
		return nil, err
	}

	province := ""
	if record.City.Names["zh-CN"] != "" {
		if len(record.Subdivisions) > 0 {
			province = record.Subdivisions[0].Names["zh-CN"]
		}

		res = &model.Location{
			Country:  record.Country.Names["zh-CN"],
			Province: province, // record.Subdivisions[0].Names["zh-CN"],
			City:     record.City.Names["zh-CN"],
			District: "",
			Lat:      record.Location.Latitude,
			Lon:      record.Location.Longitude,
		}
	}

	return res, nil
}
