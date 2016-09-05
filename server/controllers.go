package server

import (
 "fmt"
  // "html/template"
  "net/http"
  "strconv"
  "sort"
  // "time"
  // "stubwatch/config"
  "stubwatch/hublib"
  // "github.com/gorilla/context"
  "github.com/gorilla/mux"
  "github.com/Sirupsen/logrus"
)

func MainAppHandler(w http.ResponseWriter, req *http.Request) {
  Render(w, req, nil)
}

func ListingsHandler(w http.ResponseWriter, req *http.Request) {
  vars := mux.Vars(req)
  eidS := vars["eventId"]
  eventId, err := strconv.Atoi(eidS)
  if err != nil {
    log.Error(nil,"Can't convert event ID in request.", err)
    http.Error(w, fmt.Sprintf("evenId must be an integer. Got: %s",eidS), http.StatusBadRequest)
  }

  f := logrus.Fields{"eventId": eventId}

  s, err := hublib.NewDefaultService()
  if err != nil {
    log.Error(nil, "Can't get stubhub service.", err)
    http.Error(w, fmt.Sprintf("Don't have valide stubhub credentials."), http.StatusUnauthorized)
  }

  log.Debug(f,"Getting Listings.")
  listings, err := s.SearchListings(eventId) 

  if err == nil {
    f["numberOfListings"] = len(listings.Listings)
    log.Info(f, "Received listings.")

    sort.Sort(hublib.ByZoneSectionRowSeat(listings.Listings))
    renderValues := struct {
      Listings hublib.EventListings
    } {Listings: listings}
    Render(w, req, renderValues)
 
  } else { // TODO: Need hublib to give us better data to e able to say something here.
    log.Error(f, "Failed to get listings.", err)
    http.Error(w, fmt.Sprintf("Can get lists: %s", err), http.StatusServiceUnavailable)
  }
}

