package server

import(
  // "fmt"
  // "html/template"
  "net/http"
  // "strconv"
  // "sort"
  // "strings"
  // "time"
  "stubwatch/config"
  "stubwatch/hublib"
  "github.com/gorilla/context"
  "github.com/gorilla/mux"
  "github.com/jdrivas/sl"
  "github.com/spf13/viper"
  "github.com/Sirupsen/logrus"
)

var(
  debug = true
  verbose = false
  log = sl.New()
)

func init() {
  configureLogs()
}

// Server/router.
const ServerName = "Stubwatch"

// Context keys
type contextKey int
const (
  ViewsKey contextKey = iota
  ViewsListKey
)

// Entry point from main()
func DoServe(v *viper.Viper, address string) {
  shCreds := hublib.NewStubHubCredentials(viper.GetString(config.StubHubApplicationTokenKey))
  hublib.SetDefaultCredentials(shCreds)
  serve(address)
}

func serve(address string) error {
  log.Info(logrus.Fields{"serverAddress": address,},ServerName + " serving.")

  r := mux.NewRouter()

  // Routes
  r.HandleFunc("/", MainAppHandler)
  r.HandleFunc("/listings/{eventId}/", ListingsHandler)

  // TODO: move this to app/assets and be done with it in one go.
  // Asset files.
  r.PathPrefix("/app/styles/").
    Handler(http.StripPrefix("/app/styles/", http.FileServer(http.Dir("./app/styles/"))))
  r.PathPrefix("/app/js/").
    Handler(http.StripPrefix("/app/js/", http.FileServer(http.Dir("./app/js/"))))
  r.PathPrefix("/app/images/").
    Handler(http.StripPrefix("/app/images/", http.FileServer(http.Dir("./app/images/"))))


  // TODO: This could do with some thought and configuraiton.
  // Set up handler chain.
  hChain := context.ClearHandler(SetViewsHandler(LogHandler(r)))
  // Off you go.
  err := http.ListenAndServe(address, hChain)
  log.Error(logrus.Fields{"serverAddress": address,},ServerName + " shutting down.", err)

  return nil
}


func configureLogs() {
  setFormatter()
  updateLogLevel()
}

const (
  jsonLog = "json"
  textLog = "text"
)

func setFormatter() {
  switch textLog {
  // case jsonLog:
  //   f := new(logrus.JSONFormatter)
  //   log.SetFormatter(f)
  //   // mclib.SetLogFormatter(f)
  case textLog:
    f := new(sl.TextFormatter)
    f.FullTimestamp = true
    log.SetFormatter(f)
    // mclib.SetLogFormatter(f)
  }
}

func updateLogLevel() {
  l := logrus.InfoLevel
  if debug || verbose {
    l = logrus.DebugLevel
  }
  log.SetLevel(l)
  // mclib.SetLogLevel(l)
}

func SetLogLevel(l logrus.Level) {
  log.SetLevel(l)
}