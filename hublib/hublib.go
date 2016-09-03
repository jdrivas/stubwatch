package hublib

import (
  "encoding/xml"
  "fmt"
  "io/ioutil"
  "net/http"
  "github.com/dghubble/sling"
  "github.com/Sirupsen/logrus"
)

type StubHubCredentials struct {
  applicationToken string
}

func NewStubHubCredentials(applicationToken string) (StubHubCredentials) {
  return StubHubCredentials{
    applicationToken: applicationToken,
  }
}

const baseURI = "https://api.stubhub.com"
type StubhubService struct {
  Sling *sling.Sling
  authHeader []string
  creds StubHubCredentials
}

func NewStubhubService(creds StubHubCredentials) (*StubhubService) {
  return &StubhubService{
    Sling: baseRequest(creds),
    creds: creds,
  }
}

func baseRequest(creds StubHubCredentials) (*sling.Sling) {
  return sling.New().
    Base(baseURI).
    Set("Authorization", "Bearer " + creds.applicationToken).
    Set("accept", "JSON")
  // QueryStruct(&EventSearchQuery{Rows: 500, Parking: false})
}


const ResponseErrorGuard = 200
func checkForResponseError(resp *http.Response) (err error) {
  err = nil
  if resp.StatusCode > ResponseErrorGuard {
    var body []byte
    body, err = ioutil.ReadAll(resp.Body)
    log.Error(logrus.Fields{"responseCode": resp.StatusCode ,"responseBody": string(body[:]),}, "Bad reponse code.", nil)
    se := StubhubError{}
    err = xml.Unmarshal(body, &se)
    if err == nil {
      f := se.Fault()
      err = fmt.Errorf("Response %d Fault: %s. %s - %s", resp.StatusCode, f.Code, f.Message, f.Description)
    } else {
      // Let's try again .... this is some stupidity in the Stubhub machine.
      f := SHFault{}
      err = xml.Unmarshal(body, &f)
      if err == nil {
        err = fmt.Errorf("Response %d Fault: %s. %s - %s", resp.StatusCode, f.Code, f.Message, f.Description)
      } else {
        err = fmt.Errorf("Error Decoding XML return value: %s, %s", err, body)
      }
    }
  }
  return err  
}


func (sh *StubhubService) resetToBase() {
  sh.Sling = baseRequest(sh.creds)
}

// this sends the query along and fills the receiver from JSON returned by the request.
func (sh *StubhubService) send(receive interface{}) (error) {
  resp, err := sh.Sling.ReceiveSuccess(receive)

  if err == nil {
    log.Debug(logrus.Fields{"status": resp.StatusCode,},"Response results.")
    err = checkForResponseError(resp)
  }

  return err
}

func requestLogFields(req *http.Request) (f logrus.Fields) {
  f = logrus.Fields{
    "method": req.Method,
    "url": req.URL.String(),
    "host": req.Host,
  }
  return f
}

// TODO: This query string should probably be made clean.
const eventSearchPath = "/search/catalog/events/v3"
func (sh * StubhubService) SearchEvents(query string) (events Events, err error) {
  queryParam := &EventSearchQuery{Search: query}
  sh.Sling = sh.Sling.Get(eventSearchPath).QueryStruct(queryParam)

  req, err := sh.Sling.Request()
  if err != nil { return Events{}, err }
  f := requestLogFields(req)
  f["query"] = query
  log.Debug(f, "Request for the SearchCommand")

  err = sh.send(&events)
  return events, err
}

const listingSearchPath = "/search/inventory/v1"
func (sh *StubhubService) SearchListings(eventId int) (completeListings EventListings, err error) {
  queryParam  := defaultListingSearchQuery(eventId)
  sh.Sling = sh.Sling.Get(listingSearchPath).QueryStruct(queryParam)

  req, err := sh.Sling.Request()
  if err != nil { return completeListings, err }
  f := requestLogFields(req)
  f["eventID"] = eventId
  f["requestNo"] = 0
  log.Debug(f,"Request results.")

  firstListings := new(EventListings)
  err = sh.send(firstListings)
  if err == nil {
    totalListings :=  firstListings.TotalListings
    count := firstListings.Rows

    var allListings []*EventListings
    if count == 0 {
      allListings = make([]*EventListings, 0, 1) 
      log.Debug(nil, "Received no rows in last response.")
    } else {
      allListings = make([]*EventListings, 0, totalListings/count + 1) // count being the representative size of return.
    }
    allListings = append(allListings, firstListings)

    if totalListings == 0 {
      log.Debug(nil, "No listings available.")
    } else {

      f["totals-rows"] = totalListings
      reqCount := 1

      // TODO: This is a mess and needs cleaning up.
      // I don't think I like sling.
      for count < totalListings {
        newListings := new(EventListings)
        allListings = append(allListings, newListings)

        log.Debug(logrus.Fields{"requestNo": reqCount, "starting-from": count,}, "Need more rows.")
        
        sh.resetToBase()
        sh.Sling = sh.Sling.Get(listingSearchPath).QueryStruct(queryParam).QueryStruct(&ListingSearchQuery{Start: count})
        req, err = sh.Sling.Request()
        if err != nil { return completeListings, err}
        log.Debug(logrus.Fields{"uri": req.URL.String(), "method": req.Method}, "Next request.")

        err = sh.send(newListings)
        if err != nil { return completeListings, err}

        log.Debug(logrus.Fields{"received-rows": newListings.Rows}, "Successful response.")

        count += newListings.Rows
        reqCount++
      }
    }
    log.Debug(logrus.Fields{"total-no-listing-sests": len(allListings),}, "Finished with requests.")
    completeListings = combineListings(allListings)
  }

  return completeListings, err
}

func combineListings(allListings []*EventListings) (el EventListings) {

  fl := allListings[0]

  // These should remaing invariant.
  el.EventID = fl.EventID
  el.TotalListings = fl.TotalListings
  el.TotalTickets = fl.TotalTickets
  el.ZoneStats = fl.ZoneStats
  el.SectionStats = fl.SectionStats

  // These need some computation to get right
  // so I'm starting with just the first listings copy.
  el.MinQuantity = fl.MinQuantity
  el.MaxQuantity = fl.MaxQuantity
  el.PricingSummary = fl.PricingSummary

  // The rest are additive.
  for _, l := range allListings {
    el.Listings = append(el.Listings, l.Listings...)
    el.Errors = append(el.Errors, l.Errors...)
    el.ListingAttributeCategorySummary = append(el.ListingAttributeCategorySummary, l.ListingAttributeCategorySummary...)
    el.DeliveryTypeSummary = append(el.DeliveryTypeSummary, l.DeliveryTypeSummary...)
    el.Rows += l.Rows
  }

  return el
}

// Listings indexed by zone.
func (sh *StubhubService) DescribeZones(eventId int) (zoneStats []ZoneStats, err error) {

  // There may be a better way to do this with the APi (more efficient) ....
  listings, err := sh.SearchListings(eventId)
  if err != nil {return zoneStats, err}
  return listings.ZoneStats, err
}

// Listings indexed by Section
func (sh *StubhubService) DescribeSections(eventId int) (sectionStats []SectionStats, err error) {

  // There may be a better way to do this with the APi (more efficient) ....
  listings, err := sh.SearchListings(eventId)
  if err != nil {return sectionStats, err}
  return listings.SectionStats, err
}

