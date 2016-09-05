package server

import (
  // "fmt"
  // "html/template"
  "net/http"
  // "strconv"
  // "sort"
  // "strings"
  "time"
  "github.com/Sirupsen/logrus"
)

func LogHandler(handler http.Handler) http.Handler {
  return http.HandlerFunc( func (w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    f := logrus.Fields{
      "host": r.Host, 
      "hostAddr": r.RemoteAddr, 
      "method": r.Method, 
      "url":  r.URL.String(),
      "requestTime": 0,
    }
    log.Info(f, ServerName + ": Started request." )

    handler.ServeHTTP(w,r)

    f["requestTime"] = time.Since(start)
    // TODO: Need to determine what the response will be.
    log.Info(f, ServerName + ": Completed request")
  })
}