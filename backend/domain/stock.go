package domain

type Stock struct {
	Id         string
	Ticker     string
	TargetFrom string
	TargetTo   string
	Company    string
	Action     string
	Brokerage  string
	RatingFrom string
	RatingTo   string
	Time       string
}
