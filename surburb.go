package eskomlol

type Suburbs []Suburb

type Suburb struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"text,omitempty"`
	Total int    `json:"Tot,omitempty"`
}

// OmitEmpty filters for Suburbs with a non-zero Total.
func (s Suburbs) OmitEmpty() Suburbs {
	filtered := make(Suburbs, 0)
	for _, suburb := range s {
		if suburb.Total > 0 {
			filtered = append(filtered, suburb)
		}
	}
	return filtered
}

type SuburbResult struct {
	Results Suburbs `json:"Results,omitempty"`
	Total   int     `json:"Total,omitempty"`
}

type SearchSuburbs []SearchSuburb

// OmitEmpty filters for SearchSuburbs with a non-zero Total.
func (s SearchSuburbs) OmitEmpty() SearchSuburbs {
	filtered := make(SearchSuburbs, 0)
	for _, suburb := range s {
		if suburb.Total > 0 {
			filtered = append(filtered, suburb)
		}
	}
	return filtered
}

type SearchSuburb struct {
	MunicipalityName string `json:"MunicipalityName,omitempty"`
	ProvinceName     string `json:"ProvinceName,omitempty"`
	Name             string `json:"Name,omitempty"`
	ID               int    `json:"Id,omitempty"`
	Total            int    `json:"Total,omitempty"`
}
