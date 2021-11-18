package eskomlol

type Municipalities []Municipality

type Municipality struct {
	ID       string `json:"Value,omitempty"`
	Name     string `json:"Text,omitempty"`
	Disabled bool   `json:"Disabled,omitempty"`
	Selected bool   `json:"Selected,omitempty"`
	Group    string `json:"Group,omitempty"`
}
