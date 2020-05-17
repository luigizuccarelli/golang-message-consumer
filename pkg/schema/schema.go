package schema

type Trackmate struct {
	JourneyId string `json:"journey_id"`
	MessageId string `json:"message_id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Page      struct {
		ReferrerName string `json:"referrerName"`
		Referrer     string `json:"referrer"`
		URL          string `json:"url"`
		URLName      string `json:"urlName"`
	} `json:"page,omitempty"`
	UtmVariable struct {
		CustomerNumber string `json:"customerNumber,omitempty"`
		Affiliate      string `json:"affiliate"`
		Pagename       string `json:"pagename,omitempty"`
		Pagetype       string `json:"pagetype"`
		FromPage       string `json:"frompage"`
		PromoCode      string `json:"promoCode"`
		UtmCampaign    string `json:"utm_campaign"`
		UtmMedium      string `json:"utm_medium"`
		UtmContent     string `json:"utm_content"`
		UtmSource      string `json:"utm_source"`
	} `json:"utm_variable"`
	Value     interface{} `json:"value"`
	Spec      string      `json:"spec"`
	UserAgent string      `json:"userAgent"`
	Timestamp uint64      `json:"timestamp"`
}

type IrisPlusData struct {
	CREATIVEName                   string `json:"CREATIVE.name"`
	CREATIVEStatus                 string `json:"CREATIVE.status"`
	EFFORTAdvertisementName        string `json:"EFFORT.advertisement_name"`
	EFFORTCampaign                 string `json:"EFFORT.campaign"`
	EFFORTDate                     string `json:"EFFORT.date"`
	EFFORTDomain                   string `json:"EFFORT.domain"`
	EFFORTEffortDestination        string `json:"EFFORT.effort_destination"`
	EFFORTID                       string `json:"EFFORT.id"`
	EFFORTJourneyName              string `json:"EFFORT.journey_name"`
	EFFORTPromocode                string `json:"EFFORT.promocode"`
	EFFORTType                     string `json:"EFFORT.type"`
	EFFORTWhatAreYouPromoting      string `json:"EFFORT.what_are_you_promoting"`
	EFFORTWhereIsTheMarketingGoing string `json:"EFFORT.where_is_the_marketing_going"`
	FORMINFOAddress                string `json:"FORM_INFO.address"`
	FORMINFOAddress2               string `json:"FORM_INFO.address2"`
	FORMINFOAddress3               string `json:"FORM_INFO.address3"`
	FORMINFOCity                   string `json:"FORM_INFO.city"`
	FORMINFOCompanyName            string `json:"FORM_INFO.companyName"`
	FORMINFOCountryCode            string `json:"FORM_INFO.countryCode"`
	FORMINFOCountryName            string `json:"FORM_INFO.countryName"`
	FORMINFOEmail                  string `json:"FORM_INFO.email"`
	FORMINFOFaxNumber              string `json:"FORM_INFO.faxNumber"`
	FORMINFOFirstName              string `json:"FORM_INFO.firstName"`
	FORMINFOLastName               string `json:"FORM_INFO.lastName"`
	FORMINFOPhoneNumber            string `json:"FORM_INFO.phoneNumber"`
	FORMINFOPhoneNumber2           string `json:"FORM_INFO.phoneNumber2"`
	FORMINFOPhoneNumber3           string `json:"FORM_INFO.phoneNumber3"`
	FORMINFOPostalCode             string `json:"FORM_INFO.postalCode"`
	FORMINFOStateCode              string `json:"FORM_INFO.stateCode"`
	FORMINFOStateName              string `json:"FORM_INFO.stateName"`
	FORMINFOSuffix                 string `json:"FORM_INFO.suffix"`
	FORMINFOTitle                  string `json:"FORM_INFO.title"`
	JOURNEYCreativeSequence        string `json:"JOURNEY.creative_sequence"`
	JOURNEYName                    string `json:"JOURNEY.name"`
	JOURNEYStatus                  string `json:"JOURNEY.status"`
}
