package lion

import (
	"fmt"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/voioc/geo/model"
)

type IP struct {
}

func (t *IP) Analyze(ip string) (res *model.Location, err error) {
	var dbPath = "./db/ip2region.xdb"
	searcher, err := xdb.NewWithFileOnly(dbPath)
	if err != nil {
		return nil, err
	}

	defer searcher.Close()

	// do the search
	// var ip = "1.2.3.4"
	// var tStart = time.Now()
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		fmt.Printf("failed to SearchIP(%s): %s\n", ip, err)
		return nil, err
	}

	// 国家|区域|省份|城市|ISP
	data := strings.Split(region, "|")
	if data[0] == "0" && data[1] == "0" && data[2] == "0" && data[3] == "0" {
		return nil, fmt.Errorf("info empty")
	}

	location := model.Location{
		IP:       ip,
		Country:  data[0],
		Province: data[2],
		City:     data[3],
		Area:     data[4],
	}

	// fmt.Printf("{region: %s, took: %s}\n", region, time.Since(tStart))
	return &location, nil
}
