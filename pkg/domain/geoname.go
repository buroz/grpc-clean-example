package domain

import "time"

type Geoname struct {
	Id                    int       `csv:"geonameid" json:"geoname_id"`
	Name                  string    `csv:"name" json:"name"`
	AsciiName             string    `csv:"asciiname" json:"ascii_name"`
	AlternateNames        string    `csv:"alternatenames"`
	Latitude              string    `csv:"latitude" json:"latitude"`
	Longitude             string    `csv:"longitude" json:"longitude"`
	FeatureClass          string    `csv:"feature class" json:"feature_class"`
	FeatureCode           string    `csv:"feature code" json:"feature_code"`
	CountryCode           string    `csv:"country code" json:"country_code"`
	AlternateCountryCodes string    `csv:"cc2" json:"cc2,omitempty"`
	Admin1Code            string    `csv:"admin1 code" json:"admin1_code,omitempty"`
	Admin2Code            string    `csv:"admin2 code" json:"admin2_code,omitempty"`
	Admin3Code            string    `csv:"admin3 code" json:"admin3_code,omitempty"`
	Admin4Code            string    `csv:"admin4 code" json:"admin4_code,omitempty"`
	Population            int       `csv:"population" json:"population,omitempty"`
	Elevation             int       `csv:"elevation,omitempty" json:"elevation,omitempty"`
	DigitalElevationModel int       `csv:"dem,omitempty" json:"dem,omitempty"`
	Timezone              string    `csv:"timezone" json:"timezone"`
	ModificationDate      time.Time `csv:"modification date" json:"modification_date,omitempty"`
}

type GeoAlternateName struct {
	Id              int    `csv:"alternateNameId" json:"alternate_name_id"`
	GeonameId       int    `csv:"geonameid" json:"geoname_id"`
	IsoLanguage     string `csv:"isolanguage" json:"iso_language"`
	AlternateName   string `csv:"alternate name" json:"alternate_name"`
	IsPreferredName bool   `csv:"isPreferredName" json:"is_preferredName"`
	IsShortName     bool   `csv:"IsShortName" json:"is_short_name"`
	IsColloquial    bool   `csv:"isColloquial" json:"is_colloquial"`
	IsHistoric      bool   `csv:"isHistoric" json:"is_historic"`
	From            string `csv:"from" json:"from"`
	To              string `csv:"to" json:"to"`
}

type GeoCountry struct {
	Iso                string  `csv:"ISO" json:"iso"`                                 // ISO
	Iso3Code           string  `csv:"ISO3" json:"iso3"`                               // ISO3
	IsoNumeric         string  `csv:"ISO-Numeric" json:"iso_numeric"`                 // ISO-Numeric
	Fips               string  `csv:"fips" json:"fips"`                               // fips
	Name               string  `csv:"Country" json:"country"`                         // Country
	Capital            string  `csv:"Country Capital" json:"country_capital"`         // Capital
	Area               float64 `csv:"Area(in sq km)" json:"area"`                     // Area(in sq km)
	Population         uint64  `csv:"Population" json:"population"`                   // Population
	Continent          string  `csv:"Continent" json:"continent"`                     // Continent
	Tld                string  `csv:"tld" json:"tld"`                                 // tld
	CurrencyCode       string  `csv:"CurrencyCode json:"currency_code"`               // CurrencyCode
	CurrencyName       string  `csv:"CurrencyName" json:"currency_name"`              // CurrencyName
	Phone              string  `csv:"Phone" json:"phone"`                             // Phone
	PostalCodeFormat   string  `csv:"Postal Code Format" json:"postal_code_format"`   // Postal Code Format
	PostalCodeRegex    string  `csv:"PostalCode Regex" json:"postal_code_regex"`      // Postal Code Regex
	Languages          string  `csv:"Languages"`                                      // Languages
	GeonameID          int     `csv:"geonameid" json:"geoname_id"`                    // geonameid
	Neighbours         string  `csv:"neighbours"`                                     // neighbours
	EquivalentFipsCode string  `csv:"EquivalentFipsCode" json:"equivalent_fips_code"` // EquivalentFipsCode

	// LanguagesJson  []string `json:"languages"`
	// NeighboursJson []string `json:"neighbours"`
}

type GeoIsoLanguage struct {
	Iso3 string `csv:"ISO 639-3" json:"iso3"`     // ISO 639-3
	Iso2 string `csv:"ISO 639-2" json:"iso2"`     // ISO 639-2
	Iso1 string `csv:"ISO 639-1" json:"iso1"`     // ISO 639-1
	Name string `csv:"Language Name" json:"name"` // ISO 639-1
}

type GeoTimeZone struct {
	CountryCode string `csv:"CountryCode" json:"country_code"`
	TimeZoneId  string `csv:"TimeZoneId" json:"time_zone_id"`
	FmtOffset   string `csv:"GMT offset" json:"gmt_offset"`
	DstOffset   string `csv:"DST offset" json:"dst_offset"`
	RawOffset   string `csv:"rawOffset" json:"raw_offset"`
}

type GeoIsoHierarchy struct {
	ParentId  int    `csv:"ParentId" json:"parent_id"`
	ChildId   int    `csv:"ChildId" json:"child_id"`
	AdminCode string `csv:"AdminCode" json:"admin_code"`
}
