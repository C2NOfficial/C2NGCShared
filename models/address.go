package models

type Address struct {
	Line    string `json:"line"`
	City    string `json:"city"`
	State   string `json:"state"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}

func (a *Address) SetCity(city string) {
	a.City = city
}

func (a *Address) SetState(state string) {
	a.State = state
}
