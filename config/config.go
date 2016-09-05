package config

import(
  "fmt"
  "os/user"
  "github.com/spf13/viper"
)

var (
  FileName = ".stubwatch"
  FileExt = "yml"
  StubHubApplicationTokenKey = "StubhubApplicationToken"

  Viper = viper.GetViper()
)

func InitializeConfig() {

  Viper.SetConfigName(FileName)
  Viper.AddConfigPath(".")
  u, err := user.Current()
  if err == nil {
    homePath := u.HomeDir
    Viper.AddConfigPath(homePath)
  }
  err = Viper.ReadInConfig()
  if err != nil {
    fmt.Printf("Error in reading in configuration: %s\n", err)
  }
}