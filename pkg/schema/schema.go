package schema

import "time"

type NewFormat struct {
	Type        string `json:"type"`
	UtmVariable struct {
		CustomerNumber string `json:"customerNumber,omitempty"`
		Pagename       string `json:"pagename,omitempty"`
		Affiliate      string `json:"affiliate"`
		Pagetype       string `json:"pagetype"`
		PromoCode      string `json:"promoCode"`
		UtmCampaign    string `json:"utm_campaign"`
		UtmMedium      string `json:"utm_medium"`
		UtmContent     string `json:"utm_content"`
		UtmSource      string `json:"utm_source"`
	} `json:"utm_variable"`
	Value     int    `json:"value"`
	Spec      string `json:"spec"`
	UserAgent string `json:"userAgent"`
}

type SegmentIO struct {
	Id          string `json:"id,omitempy"`
	AnonymousID string `json:"anonymousId,omitempty"`
	Context     struct {
		Campaign struct {
			Content string `json:"content"`
			Name    string `json:"name"`
			Source  string `json:"source"`
		} `json:"campaign"`
		IP      string `json:"ip"`
		Library struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"library"`
		Locale string `json:"locale"`
		Page   struct {
			Path     string `json:"path"`
			Referrer string `json:"referrer"`
			Search   string `json:"search"`
			Title    string `json:"title"`
			URL      string `json:"url"`
		} `json:"page,omitempty"`
		UserAgent string `json:"userAgent"`
	} `json:"context"`
	Event        string `json:"event,omitempty"`
	Integrations struct {
	} `json:"integrations,omitempty"`
	MessageID         string    `json:"messageId,omitempty"`
	OriginalTimestamp time.Time `json:"originalTimestampi,omitempty"`
	Properties        struct {
		IrisPlusData struct {
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
		} `json:"irisPlusData"`
		Type        string `json:"type"`
		Spec        string `json:"spec"`
		UtmVariable struct {
			CustomerNumber string `json:"customerNumber,omitempty"`
			Pagename       string `json:"pagename,omitempty"`
			Affiliate      string `json:"affiliate"`
			Pagetype       string `json:"pagetype"`
			PromoCode      string `json:"promoCode"`
			UtmCampaign    string `json:"utm_campaign"`
			UtmMedium      string `json:"utm_medium"`
			UtmContent     string `json:"utm_content"`
			UtmSource      string `json:"utm_source"`
		} `json:"utm_variable"`
		Value int `json:"value"`
	} `json:"properties"`
	ReceivedAt time.Time `json:"receivedAt"`
	SentAt     time.Time `json:"sentAt"`
	Timestamp  time.Time `json:"timestamp"`
	Type       string    `json:"type"`
	UserID     string    `json:"userId"`
	Version    int       `json:"version"`
}
