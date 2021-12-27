package domain

import "encoding/xml"

type TcmbTodayResponse struct {
	XMLName  xml.Name `xml:"Tarih_Date"`
	Text     string   `xml:",chardata"`
	Tarih    string   `xml:"Tarih,attr"`
	Date     string   `xml:"Date,attr"`
	BultenNo string   `xml:"Bulten_No,attr"`
	Currency []struct {
		Text            string `xml:",chardata"`
		CrossOrder      string `xml:"CrossOrder,attr"`
		Kod             string `xml:"Kod,attr" json:"code"`
		CurrencyCode    string `xml:"CurrencyCode,attr" json:"currency_code"`
		Unit            string `xml:"Unit" json:"unit"`
		Isim            string `xml:"Isim" json:"name"`
		CurrencyName    string `xml:"CurrencyName" json:"currency_name"`
		ForexBuying     string `xml:"ForexBuying" json:"forex_buying"`
		ForexSelling    string `xml:"ForexSelling" json:"forex_selling"`
		BanknoteBuying  string `xml:"BanknoteBuying" json:"banknote_buying"`
		BanknoteSelling string `xml:"BanknoteSelling" json:"banknote_selling"`
		CrossRateUSD    string `xml:"CrossRateUSD" json:"cross_rate_usd"`
		CrossRateOther  string `xml:"CrossRateOther" json:"cross_rate_other"`
	} `xml:"Currency"`
}
