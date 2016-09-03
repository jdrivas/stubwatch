package hublib


type EventSearchQuery struct {
  Search string             `url:"q,omitempty"`
  Name string               `url:"name,omitempty"`
  Rows int                  `url:"rows,omitempty"`
  Parking bool              `url:"parking,omitempty"`
}

type Events struct {
  Count  int         `json:"numFound"`
  geoExpansion bool   `json:"geoExpansion"`
  Events []Event      `json:"events"`
}

type Event struct {
  ID int                  `json:"id"`
  Status string           `json:"status"`
  Locale string           `json:"local"`
  Name string             `json:"name"`
  Description string      `json:"description"`
  WebURI string           `json:"webURI"`
  DateLocal string   `json:"eventDateLocal"`
  DateUTC  string    `json:"eventDateUTC"`
  Distance float64        `json:"distance"`
  Venue Venue             `json:"venue"`
  TicketInfo TicketInfo   `json:"ticketInfo"`
  Ancestors Ancestors     `json:"ancestors"`
  Images  []Image         `json:"images"`
  DisplayAttributes       `json:"displayAttributes"`
}



type Venue struct {
  ID int                  `json:"id"`
  Name string             `json:"name"`
  WebURI string           `json:"webURI"`
  Latitude float32        `json:"latitude"`
  Longitude float32       `json:"longitude"`
  Timzone string          `json:"timezone"`
  Address1 string         `json:"address1"`
  Address2 string         `json:address2"`
  City string             `json:"city"`
  State string            `json:"state"`
  PostCode string         `json:"postalCode"`
  Country string          `json:"country"`
}
type TicketInfo struct {
  MinPrice float64         `json:"minPrice"`
  MaxPrice float64         `json:"maxPrice"`
  TotalTickets int         `json:"totalTickets"`
  CurrencyCode string       `json:"currencyCode"`
}
type Ancestors struct {
  Categories []AncestorRef  `json:"categories"`
  Groupings  []AncestorRef  `json:"groupings"`
  Performers []AncestorRef  `json:"performers"`
}

// These are objects that are brought forward
// from other objects in the overall 
// e.g. subfields from Category, Group, Performer
type AncestorRef struct {
  Id int                    `json:"id"`
  Name string               `json:"name"`
  webURI string             `json:"webURI"`
}

type Category struct {}
type Grouping struct {}
type Performer struct {}
type Image struct {}

type DisplayAttributes struct {}

