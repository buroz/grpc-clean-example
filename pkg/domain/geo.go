package domain

type Continent struct {
	Iso  string `json:"iso"`
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Country struct {
	Iso              string   `json:"iso"`
	Name             string   `json:"name"`
	Continent        string   `json:"continent"`
	CurrencyCode     string   `json:"currency_code"`
	CurrencyName     string   `json:"currency_name"`
	Phone            string   `json:"phone"`
	PostalCodeFormat string   `json:"postal_code_format"`
	PostalCodeRegex  string   `json:"postal_code_regex"`
	Languages        []string `json:"languages"`
	Id               int      `json:"id"`
}
