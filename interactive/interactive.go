package interactive 

import (
  "fmt"
  "io"
  "os"
  "os/user"
  "sort"
  "strings"
  "stubwatch/hublib"
  "text/tabwriter"
  "github.com/alecthomas/kingpin"
  "github.com/bobappleyard/readline"
  "github.com/spf13/viper"
)
var (

  app *kingpin.Application

  exit *kingpin.CmdClause
  quit *kingpin.CmdClause
  verboseCmd *kingpin.CmdClause
  verbose bool
  testString []string

  searchCmd *kingpin.CmdClause
  searchEventsCmd *kingpin.CmdClause
  searchStringArg string

  describeCmd *kingpin.CmdClause
  describeListingsCmd *kingpin.CmdClause
  describeZonesCmd *kingpin.CmdClause
  describeSectionsCmd *kingpin.CmdClause
  eventIdArg int
)

func init() {
  // keyValues = make(map[string]string)
  app = kingpin.New("", "Interactive mode.").Terminate(doTerminate)

  // state
  verboseCmd = app.Command("verbose", "toggle verbose mode.")
  exit = app.Command("exit", "exit the program. <ctrl-D> works too.")
  quit = app.Command("quit", "exit the program.")

  searchCmd = app.Command("search","Search stubhub.")
  searchEventsCmd = searchCmd.Command("events", "Search for events that match the search-string")
  searchEventsCmd.Arg("search-string", "What to sarch for.").Required().StringVar(&searchStringArg)


  describeCmd = app.Command("describe", "More detail on events")
  describeListingsCmd = describeCmd.Command("listings", "Search for listings for the event-id")
  describeListingsCmd.Arg("event-id", "Get listings for the event descrived by event-id").Required().IntVar(&eventIdArg)
  describeZonesCmd = describeCmd.Command("zones", "describe zones stats for the event.")
  describeZonesCmd.Arg("event-id","Get zones stats for this event.").Required().IntVar(&eventIdArg)
  describeSectionsCmd = describeCmd.Command("sections", "describe section stats for the event.")
  describeSectionsCmd.Arg("event-id","Get section stats for this event.").Required().IntVar(&eventIdArg)

  // keyValueCmd = app.Command("kv", "Test key value args.")
  // keyValueCmd.Arg("key value", "Test for key values pairs.").StringMapVar(&keyValues)
  initializeConfig()
}


func DoICommand(line string, creds hublib.StubHubCredentials) (err error) {

  // This is due to a 'peculiarity' of kingpin: it collects strings as arguments across parses.
  testString = []string{}

  // Prepare a line for parsing
  line = strings.TrimRight(line, "\n")
  fields := []string{}
  fields = append(fields, strings.Fields(line)...)
  if len(fields) <= 0 {
    return nil
  }

  command, err := app.Parse(fields)
  if err != nil {
    fmt.Printf("Command error: %s.\nType help for a list of commands.\n", err)
    return nil
  } else {
    switch command {
      case verboseCmd.FullCommand(): err = doVerbose()
      case exit.FullCommand(): err = doQuit()
      case quit.FullCommand(): err = doQuit()
      case searchEventsCmd.FullCommand(): err = doEventSearch(creds)
      case describeListingsCmd.FullCommand(): err = doDescribeListings(creds)
      case describeZonesCmd.FullCommand(): err = doDescribeZones(creds)
      case describeSectionsCmd.FullCommand(): err = doDescribeSections(creds)
    }
  }
  return err
}

func doEventSearch(creds hublib.StubHubCredentials) (err error) {
  s := hublib.NewStubhubService(creds)
  events, err := s.SearchEvents(searchStringArg)
  if err == nil {
    fmt.Printf("There are %d events\n", events.Count)
    w := tabwriter.NewWriter(os.Stdout, 2, 5, 2, ' ', 0)
    for i, event := range events.Events {
      fmt.Fprintf(w, "%d. ********************\n", i+1)
      fmt.Fprintf(w, "%d %s\t\t%s %.0f\n", event.ID, event.Name, event.DateLocal, event.Distance)
      performers := event.Ancestors.Performers
      fmt.Fprintf(w, "Performers: ")
      for _, p := range performers {
        fmt.Fprintf(w, "%s ", p.Name)
      }
      fmt.Fprintln(w)
      // ti := event.TicketInfo
      // fmt.Fprintf(w, "%d\t%g\t%g\n", ti.TotalTickets, ti.MinPrice, ti.MaxPrice)
      v := event.Venue
      fmt.Fprintf(w, "%s\t%s, %s\n", v.Name, v.City, v.State)
      fmt.Fprintf(w, "Description: ")
      fmt.Fprintf(w, "%s\n", event.Description)
    }
    w.Flush()
  }
  return err
}

func doDescribeListings(creds hublib.StubHubCredentials) (err error) {
  s := hublib.NewStubhubService(creds)
  listings, err := s.SearchListings(eventIdArg)
  if err == nil {

    // Summary
    fmt.Printf("Summary of listings for event: %d\n", eventIdArg)

    pSum := listings.PricingSummary
    w := tabwriter.NewWriter(os.Stdout, 2, 5, 2, ' ', 0)
    fmt.Fprintf(w, "Listings\tCount\tTickets\tMax\tAvg\tMin\tSections\tZones\tErrors\tCat Sum\tDelivery Sum\n")
    fmt.Fprintf(w, "%d\t%d\t%d\t%.2f\t%.2f\t%.2f\t%d\t%d\t%d\t%s\t%s\n",
      listings.TotalListings, listings.Rows, listings.TotalTickets, 
      pSum.MaxTicketPrice, pSum.AverageTicketPrice, pSum.MinTicketPrice,
      len(listings.SectionStats), len(listings.ZoneStats), len(listings.Errors),
      listings.ListingAttributeCategorySummary, listings.DeliveryTypeSummary)
    w.Flush()

    // Listings (reuse the same tabwriter)
    fmt.Fprintf(w,"Price\tQuantity\tZone\tSection\tRow\tSeats\tSplit\n")
    offers := listings.Listings
    sort.Sort(hublib.ByZoneSectionRowSeat(offers))
    for _, l := range offers {
      sName := l.SectionName
      if l.DirtyTicketInd { sName = "*" + sName}
      fmt.Fprintf(w,"%s\t%d\t%s\t%s\t%s\t%s\t%s\n",
        l.CurrentPrice, l.Quantity, l.ZoneName, l.SectionName, 
        l.Row, l.SeatNumbers, l.SplitsString())
    }
    w.Flush()
  }
  return err
}

func doDescribeZones(creds hublib.StubHubCredentials) (error) {
  s := hublib.NewStubhubService(creds)
  zoneStats, err := s.DescribeZones(eventIdArg)
  if err == nil {
    sort.Sort(hublib.ByZone(zoneStats))
    w := tabwriter.NewWriter(os.Stdout, 2, 5, 2, ' ', 0)
    fmt.Fprintf(w,"Zone\tListings\tTickets\tMin\tAvg\tMax\n")
    for _, zone := range zoneStats {
      fmt.Fprintf(w, "%s\t%d\t%d\t%.2f\t%.2f\t%.2f\n", 
        zone.ZoneName, zone.TotalListings, zone.TotalTickets,
        zone.MinTicketPrice, zone.AverageTicketPrice, zone.MaxTicketPrice)
    }
    w.Flush()
  }
  return err
}

func doDescribeSections(creds hublib.StubHubCredentials) (error) {
  s := hublib.NewStubhubService(creds)
  sectionStats, err := s.DescribeSections(eventIdArg)
  if err == nil {
    sort.Sort(hublib.BySection(sectionStats))
    w := tabwriter.NewWriter(os.Stdout, 2, 5, 2, ' ', 0)
    fmt.Fprintf(w,"Section\tListings\tTickets\tMin\tAvg\tMax\n")
    for _, s := range sectionStats {
      fmt.Fprintf(w, "%s\t%d\t%d\t%.2f\t%.2f\t%.2f\n", 
        s.SectionName, s.TotalListings, s.TotalTickets,
        s.MinTicketPrice, s.AverageTicketPrice, s.MaxTicketPrice)
    }
    w.Flush()
  }
  return err
}
func priceString(min, avg, med, max float64) (string) {
  return fmt.Sprintf("max $%.2f, avg $%.2f, med $%.2f, min $%.2f", max, avg, med, min)
}



// func doKeyValue() (error) {
//   fmt.Printf("There were %d key values pairs.\n", len(keyValues))
//   for key, value := range keyValues {
//     fmt.Printf("%s = %s.\n", key, value)
//   }
//   return nil
// }


//
// Support.
//
const (
  ConfigFileName = ".stubwatch"
  ConfigFileExt = "yml"
  StubHubApplicationTokenKey = "StubhubApplicationToken"
)
func initializeConfig() {

  viper.SetConfigName(ConfigFileName)
  viper.AddConfigPath(".")
  u, err := user.Current()
  if err == nil {
    homePath := u.HomeDir
    viper.AddConfigPath(homePath)
  }
  err = viper.ReadInConfig()
  if err != nil {
    fmt.Printf("Error in reading in configuration: %s\n", err)
  }
}

func toggleVerbose() bool {
  verbose = verbose
  return verbose
}

func doVerbose() (error) {
  if toggleVerbose() {
    fmt.Println("Verbose is on.")
  } else {
    fmt.Println("Verbose is off.")
  }
  return nil
}

func doQuit() (error) {
  return io.EOF
}

func doTerminate(i int) {}

func promptLoop(prompt string, process func(string) (error)) (err error) {
  errStr := "Error - %s.\n"
  for moreCommands := true; moreCommands; {
    line, err := readline.String(prompt)
    if err == io.EOF {
      moreCommands = false
    } else if err != nil {
      fmt.Printf(errStr, err)
    } else {
      readline.AddHistory(line)
      err = process(line)
      if err == io.EOF {
        moreCommands = false
      } else if err != nil {
        fmt.Printf(errStr, err)
      }
    }
  }
  return nil
}

// This gets called from the main program, presumably from the 'interactive' command on main's command line.
func DoInteractive() {
  if !viper.IsSet(StubHubApplicationTokenKey) {
    fmt.Printf("ApplicationTokenString not set in credentials file: %s.%s\n", ConfigFileName, ConfigFileExt)
  }
  creds := hublib.NewStubHubCredentials(viper.GetString(StubHubApplicationTokenKey))
  xICommand := func(line string) (err error) {return DoICommand(line, creds)}
  prompt := "> "
  err := promptLoop(prompt, xICommand)
  if err != nil {fmt.Printf("Error - %s.\n", err)}
}




