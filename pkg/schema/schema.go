package schema

// Analytics schema that forwards the json data payload to our backend analytics system (Couchbase)
type Analytics struct {
	TrackingId  string         `json:"trackingid"`
	To          PageDetail     `json:"to"`
	From        PageDetail     `json:"from"`
	Location    LocationDetail `json:"location"`
	Currency    CurrencyDetail `json:"currency"`
	Event       EventDetail    `json:"event"`
	CampaignId  string         `json:"campaignid"`
	AffiliateId string         `json:"affiliateid"`
	Timestamp   int64          `json:"timestamp"`
	Platform    PlatformDetail `json:"platform"`
	ProductName string         `json:"productname"`
}

// PageDetail schema
type PageDetail struct {
	Url      string `json:"url"`
	PageName string `json:"pagename"`
	PageType string `json:"pagetype"`
}

// LocationDetail schema - used in the Analytics schema
type LocationDetail struct {
	Ip      string        `json:"ip"`
	Carrier string        `json:"carrier"`
	Country CountryDetail `json:"country"`
}

// CountryDetail schema
type CountryDetail struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Capitial string `json:"capital"`
}

// CurrencyDetail schema - used in the Analytics schema
type CurrencyDetail struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// EventDetail schema
type EventDetail struct {
	Type       string `json:"type"`
	TimeonPage int    `json:"timeonpage"`
}

// PlatformDetail schema
type PlatformDetail struct {
	AppcodeName string `json:"appCodeName"`
	AppName     string `json:"appName"`
	AppVersion  string `json:"appVersion"`
	Language    string `json:"language"`
	Os          string `json:"os"`
	Product     string `json:"product"`
	ProductSub  string `json:"productSub"`
	UserAgent   string `json:"userAgent"`
	Vendor      string `json:"vendor"`
}
