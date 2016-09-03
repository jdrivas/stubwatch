package hublib

import (
  "fmt"
  "regexp"
  "strconv"
  "strings"
  // "github.com/Sirupsen/logrus"
)

type ListingSearchQuery struct {
  EventID int                         `url:"eventId,omitempty"`
  ZoneIDist []int                     `url:"zoneidlist,omitempty"`                    // List of ZoneIDS restrict listings to.
  SecionIDList []int                  `url:"sectionidlist,omitempty"`                 // Section filter.
  Quantity int                        `url:"quantity,omitempty"`                      // how many tickets do you want.
  PriceMin float32                    `url:"pricemin,omitempty"`                      // minimum price for listings.
  PriceMax float32                    `url:"pricemax,omitempty"`                      // Maximum price for listings.
  ListingAttributeList []int          `url:listingatributelist,omitempty""`           // Ticket traits (see ListingAttribute list below.)
  ListingAttributeCategoryList []int  `url:"listingattributecategorylist,omitempty"`  // 1: ObstructedView, 
                                                      // 2: Wheelchar acesseible
                                                      // 3: Alcohol-free
                                                      // 4: Parking pass included.
                                                      // 5: Piggyback Seats
                                                      // 6: Aisle
  DeliveryTypeList []int               `url:"deliverytypelist,omitempty"`               // 1. Electronic, 2. Instance, 3. LMS, 4 UPS.
  Sort int                             `url:"sort,omitempty"`
  Start int                            `url:"start,omitempty"`
  Rows int                             `url:"rows,omitempty"`
  ZoneStats bool                       `url:"zonestats,omitempty"`
  SectionStats bool                    `url:"sectionstats,omitempty"`
  PricingSummary bool                  `url:"pricingsummary,omitempty"`
}

// Fills out some defaults
func defaultListingSearchQuery(eventId int) (lsq *ListingSearchQuery) {
    return &ListingSearchQuery{
      EventID: eventId,
      Rows: 500,
      ZoneStats: true,
      SectionStats: true,
      PricingSummary: true,
    }
}

type EventListings struct {
  EventID int                     `json:"eventId"`          // Event ID
  TotalListings int               `json:"totalListings"`    // number of listings available
  Start int                       `json:"start"`            // the start index of this batch.
  TotalTickets int                `json:"totalTickets"`     // Total number of tickets available.
  MinQuantity int                  `json:"minQuantity"`      // minimum number of tickets for purchase in one listing.
  MaxQuantity int                 `json:"maxQuantity"`      // maximum number of tickets avialable.
  Listings []Listing              `json:"listing"`          // collection of listings
  PricingSummary PricingSummary    `jsons:"pricingSummary"` // summary of prices for these listings.
  SectionStats  []SectionStats     `json:"section_stats"`   // venue section descriptions.
  ZoneStats []ZoneStats            `json:"zone_stats"`      // zones for searching groups of sections.
  Errors []Errors                  `json:"errors"`
  ListingAttributeCategorySummary []string  `json:"listingAttributeCategorySummary"`
  DeliveryTypeSummary []string     `json:"deliveryTypeSummary"`
  Rows int                         `json:"rows"`             // the number of listings here.
}

type Listing struct {
  ListingID int                      `json:"listingId"`      // The ID for this listing.
  CurrentPrice Money                 `json:"currentPrice"`   // price of a ticker for this listing.
  SectionId int                      `json:"sectionId"`      // ID for this section.
  Quantity int                       `json:"quantity"`       // number of available tickets.
  SectionName string                 `json:"sectionName"`    // Where these tickets are located.
  Row string                         `json:"row"`            // 
  SeatNumbers string                 `json:"SeatNumbers"`    //
  ZoneID int                         `json:"zoneId"`         // StubHub zone where the ticekts are located.
  ZoneName string                    `json:"zoneName"`       // Display name for the zone.
  DeliveryTypeList []int             `json:"deliveryTypeList"` // 1: electonic, 2: instance downloadl, 4:last minute service scenter, 5: UPS
  ListingAttributeList []int         `json:"listingAttributeList"` // StubHub IDs of attributes that apply.
  ListingAttributeCategoryList []int `json:"listingAttributeCategoryList"` // Id's of categories that apply.
  DirtyTicketInd bool                 `json:"dirtyTicketInd"` // BOOL indicated wether or not the ticket locations can be mapped to the venue.
  SplitVector []int                  `json:"splitVector"`    // Numbers of ticket you can buy: e.g. 1, 2, 4.
  FaceValue Money                    `json:"faceValue"`      // Issue price of ticket if known.
  ServiceFee Money                   `json:"serviceFee"`     // currently null.
  DeliveryFee Money                  `json:"deliveryFee"`    // currently null.
  TotalCost Money                    `json:"totalCost"`      // currently null.
  SellerOwnInd int                   `json:"sellerOwnInd"`      // BOOL does the seller actually poses these.
}

// Return a string for SplitVector.
func (l Listing) SplitsString() (s string) {
  for _, sp := range l.SplitVector {
    s += fmt.Sprintf("%d,", sp)
  }
  return s
}

type PricingSummary struct {
  Name string             `json:"name"`
  MinTicketPrice float64  `json:"minTicketPrice"`
  MaxTicketPrice float64   `json:"maxTicketPrice"`
  AverageTicketPrice float64 `json:"averageTicketPrice"`
  TotalListings int         `json:"totalListings"`
  MedianTicketPrice float64 `json:"medianTicketPrice"`
}

type SectionStats struct {
  SectionId int           `json:"sectionId"`
  SectionName string      `json:"sectionName"`
  TotalTickets int        `json:"totalTickets"`
  TotalListings int    `json:"totalListings"`
  MinTicketPrice float64  `json:"minTicketPrice"`
  MedianTicketPrice float64 `json:"medianTicketPrice"`
  AverageTicketPrice float64 `json:"averageTicketPrice"`
  MaxTicketPrice float64  `json:"maxTicketPrice"`
  MinTicketQuantity int `json:"minTicketQuantity"`
  MaxTicketQuanity int `json:"maxTicketQuantity"`
  Percentiles []Percentile `json:"percentiles"`
}


// Bet there is a cheaper way to do this.
const startsAlphaRE = "^[A-Za-z]"
func startsAlpha(s string) bool {
  re := regexp.MustCompile(startsAlphaRE)
  return re.MatchString(s)
}

type ByZoneSectionRowSeat []Listing
func (a ByZoneSectionRowSeat) Len() int { return len(a) }
func (a ByZoneSectionRowSeat) Swap(i, j int ) { a[i], a[j] = a[j], a[i] }
func (a ByZoneSectionRowSeat) Less(i , j int) bool {
  if a[i].ZoneName != a[j].ZoneName {
    return a[i].ZoneName < a[j].ZoneName 
  }
  if a[i].SectionName != a[j].SectionName { 
    return a[i].SectionName < a[j].SectionName 
  }

  // Rows with alpha are lower than rows with numbers.
  // This is venue by venue of course, but AT&T park
  // is like this so ......
  if a[i].Row != a[j].Row {
    return alphaNumLess(a[i].Row, a[j].Row)
  }

  // Just look at the first seat.
  // iS := strings.Split(a[i].SeatNumbers, ",")
  // jS := strings.Split(a[j].SeatNumbers, ",")
  // iAlpha := startsAlpha(iS)
  // jAlpha := startsAlpha(js)
  // iS, err  := strconv.Atoi(iSeats[0])
  // if err != nil { log.Panic(nil, "Bad seats.", err)}
  // jS, err := strconv.Atoi(jSeats[0])
  // if err != nil { log.Panic(nil, "Bad seats.", err) }
  iS := strings.Split(a[i].SeatNumbers, ",")[0]
  jS := strings.Split(a[j].SeatNumbers, ",")[0]
  return alphaNumLess(iS, jS)

}

const isNumRE = "^[0-9]+$"
func alphaNumLess(i, j string) (bool)  {
  re := regexp.MustCompile(isNumRE)
  iIsNum := re.MatchString(i)
  jIsNum:= re.MatchString(j)

  if !iIsNum && !jIsNum { return i < j} // both strings
  if iIsNum && jIsNum {  // both nums
    iN, err := strconv.Atoi(i)
    if err != nil { log.Panic(nil, "hublib: Trying to compare strings as ints.", err) }
    jN, err := strconv.Atoi(j)
    if err != nil { log.Panic(nil, "hublib: Trying to compare strings as ints.", err) }

    return iN < jN
  }

  // one is alpha and one is num ... string is less than num
  if iIsNum { return false}
  return true
}



type BySection []SectionStats
func (a BySection) Len() int { return len(a) }
func (a BySection) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySection) Less(i, j int) bool { return a[i].SectionName < a[j].SectionName }

type ZoneStats struct {
  ZoneID int              `json:"zoneId"`
  ZoneName string         `json:"zoneName"`
  TotalListings int       `json:"totalListings"`
  TotalTickets int        `json:"totalTickets"`
  MinTicketPrice float64  `json:"minTicketPrice"`
  MaxTicketPrice float64  `json:"maxTicketPrice"`
  MinTicketQuantity int   `json:"minTicketQuanity"`
  MaxTicketQuanity int    `json:"maxTicketQuantity"`
  AverageTicketPrice float64 `json:"averageTicketPrice"`
  Percentiles []Percentile `json:"percentiles"`
}

type ByZone []ZoneStats
func (a ByZone) Len() int { return len(a) }
func (a ByZone) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByZone) Less(i, j int) bool { return a[i].ZoneName < a[j].ZoneName }


type Money struct {
  Amount float64        `json:"amount"`
  Currency string       `json:"currency"`
}

func (m Money) String() (s string) {
  switch m.Currency {
    case "USD": s = "$"
    default: s = m.Currency + " "
  }
  s += fmt.Sprintf("%.02f", m.Amount)
  return s  
}

type Percentile struct {
  Name float64         `json:"name"`
  Value float64         `json:"value"`
}

type Errors struct {
  ErrorDescription string `json:"errorDescription"`    // This could be 'type'. The spec is slightly ambiguous.
  ErrorType string        `json:"errorType"`
  ErrorTypeId string      `json:"errorTypeId"`
  ErrorMessage string     `json:"errorMessage"`
  ErrorParamater string   `json:"errorParamater"`
  Message string          `json:"message"`
  Type string             `json:"type"`
  Paramater string        `json:"paramater"`
}

type ListingsMap map[string]Listing

// Listings indexed by ZoneNames
func (EventListings) ZoneListingsMap() (ListingsMap) {
  zlm := make(ListingsMap)
  return zlm
}

// Listings index by SectionNames
func (EventListings) SectionListingsMaps() (ListingsMap) {
  slm := make(ListingsMap)
  return slm
}






