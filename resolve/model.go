package resolve

type Country struct {
	// Name is the name of the country
	Name string `json:"name"`

	// Code is the code of the country in the website
	Code string `json:"code"`
}

type University struct {
	Country string `json:"country"`

	// Name is the name of the university
	Name string `json:"name"`

	// Domain is the website domain of the university
	Domain string `json:"domain"`
}
