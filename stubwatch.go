package main

import (
  "fmt"
  "os"
  "strings"
  "stubwatch/config"
  "stubwatch/interactive"
  "stubwatch/server"
  "text/template"
  "github.com/alecthomas/kingpin"
  // "github.com/spf13/viper"
)

var (
  app                               *kingpin.Application
  verbose                           bool
  debug                             bool
  region                            string

  interactiveCmd *kingpin.CmdClause
  serveCmd *kingpin.CmdClause
  serverAddressArg string

  command1 *kingpin.CmdClause
  sub1_command1 *kingpin.CmdClause

  templateCmd *kingpin.CmdClause
  testTemplate *template.Template
  whoArg string


)

func init() {
  app = kingpin.New("craft-config.go", "Command line to to manage minecraft configs.")
  app.Flag("verbose", "Describe what is happening, as it happens.").Short('v').BoolVar(&verbose)
  app.Flag("Debug", "Tell us in detail what's going on.").Short('d').BoolVar(&debug)

  interactiveCmd = app.Command("interactive", "Prompt for commands.")
  serveCmd = app.Command("serve", "Become an app server.")
  serveCmd.Arg("server-address", "IP address of server").Default("127.0.0.1:3100").StringVar(&serverAddressArg)

  command1 = app.Command("command1", "Do stuff in a command-1 context.")
  sub1_command1= command1.Command("sub1", "Sub1 command for command-1")

  templateCmd = app.Command("template", "Do stuff with templates.")
  templateCmd.Arg("who", "Who you might say hello to.").Default("World").StringVar(&whoArg)

  kingpin.CommandLine.Help = `A command-line minecraft config tool.`

  config.InitializeConfig()
}

func main() {

  // Parse the command line to fool with flags and get the command we'll execeute.
  // command := kingpin.MustParse(app.Parse(os.Args[1:]))
  app.Terminate(doTerminate)
  command, err := app.Parse((os.Args[1:]))
  if err != nil {
    fmt.Printf("Error on parse: %s\n", err)
  }

  testTemplate = template.Must(template.New("example").Parse("It's {{ .When.String }}, hello {{ .HelloString }}.\n"))

   if verbose {
    fmt.Printf("Starting up.")
   }

   serverAddress := getServerAddress(serverAddressArg)

  // List of commands as parsed matched against functions to execute the commands.
  commandMap := map[string]func() {
    sub1_command1.FullCommand(): func() { fmt.Println("Hello Sub Command 1.")},
  }

  // Execute the command.
  switch command {
    case interactiveCmd.FullCommand(): interactive.DoInteractive(config.Viper)
    case serveCmd.FullCommand(): server.DoServe(config.Viper, serverAddress)
    default: commandMap[command]()
  }
}

func doTerminate(i int) {
  server.DoServe(config.Viper, getServerAddress(serverAddressArg))
}

// Using the PORT environment variable to exist with gin the app reloader.
func getServerAddress(address string) (string) {
  a := strings.Split(address, ":")
  ip := "127.0.0.1"
  port := "3100"
  if len(a) == 2 {
    ip = a[0]
    port = a[1]
  }
  if p, ok := os.LookupEnv("PORT"); ok {
    port = p
  }
  return fmt.Sprintf("%s:%s", ip, port)
}

