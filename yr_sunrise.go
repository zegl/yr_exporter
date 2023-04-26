package main

type SunriseResponse struct {
	Location Location `json:"location"`
	Meta     Meta     `json:"meta"`
}

type HighMoon struct {
	Desc      string `json:"desc"`
	Elevation string `json:"elevation"`
	Time      string `json:"time"`
}

type LowMoon struct {
	Desc      string `json:"desc"`
	Elevation string `json:"elevation"`
	Time      string `json:"time"`
}

type Moonphase struct {
	Desc  string `json:"desc"`
	Time  string `json:"time"`
	Value string `json:"value"`
}

type Moonposition struct {
	Azimuth   string `json:"azimuth"`
	Desc      string `json:"desc"`
	Elevation string `json:"elevation"`
	Phase     string `json:"phase"`
	Range     string `json:"range"`
	Time      string `json:"time"`
}

type Moonrise struct {
	Desc string `json:"desc"`
	Time string `json:"time"`
}

type Moonset struct {
	Desc string `json:"desc"`
	Time string `json:"time"`
}

type Moonshadow struct {
	Azimuth   string `json:"azimuth"`
	Desc      string `json:"desc"`
	Elevation string `json:"elevation"`
	Time      string `json:"time"`
}

type Solarmidnight struct {
	Desc      string  `json:"desc"`
	Elevation float64 `json:"elevation,string"`
	Time      string  `json:"time"`
}

type Solarnoon struct {
	Desc      string  `json:"desc"`
	Elevation float64 `json:"elevation,string"`
	Time      string  `json:"time"`
}

type Sunrise struct {
	Desc string `json:"desc"`
	Time string `json:"time"`
}

type Sunset struct {
	Desc string `json:"desc"`
	Time string `json:"time"`
}

type Time struct {
	Date          string        `json:"date"`
	HighMoon      HighMoon      `json:"high_moon,omitempty"`
	LowMoon       LowMoon       `json:"low_moon,omitempty"`
	Moonphase     Moonphase     `json:"moonphase,omitempty"`
	Moonposition  Moonposition  `json:"moonposition"`
	Moonrise      Moonrise      `json:"moonrise,omitempty"`
	Moonset       Moonset       `json:"moonset,omitempty"`
	Moonshadow    Moonshadow    `json:"moonshadow,omitempty"`
	Solarmidnight Solarmidnight `json:"solarmidnight,omitempty"`
	Solarnoon     Solarnoon     `json:"solarnoon,omitempty"`
	Sunrise       Sunrise       `json:"sunrise,omitempty"`
	Sunset        Sunset        `json:"sunset,omitempty"`
}

type Location struct {
	Height    string `json:"height"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Time      []Time `json:"time"`
}
