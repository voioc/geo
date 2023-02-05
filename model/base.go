package model

type Geo interface {
	Cmd()
	_version()
	_update()
	Analyze()
}

type Location struct {
	IP       string  `json:"ip"`
	Country  string  `json:"country"`
	Province string  `json:"province"`
	City     string  `json:"city"`
	District string  `json:"district"`
	Area     string  `json:"area"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}
