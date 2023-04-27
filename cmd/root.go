package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.eng.vmware.com/vdp/kafka-metadata-cleaner/config"
)

var version = "0.0.1"
var rootCmd = &cobra.Command{
	Use:     "kafka-metadata-cleaner",
	Version: version,
	Short:   "kafka-metadata-cleaner - a simple CLI to clean kafka old topics in kafka clusters",
	Long:    `kafka-metadata-cleaner is a CLI to clean your metadata in a kafka cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
var cfgFile string
var logLevel string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file (default is config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "Log level for the program")
	rootCmd.PersistentFlags().Bool("dry-run", true, "Run program in dry-run")
	viper.BindPFlag("dryRun", rootCmd.PersistentFlags().Lookup("dry-run"))
	viper.SetDefault("author", "jcriadomarco@vmware.com")
	viper.SetDefault("license", "vmware")

	cobra.OnInitialize(initCommand)
}

func initCommand() {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	loadConfig()

	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(logLevel)
	if err != nil {
		ll = logrus.DebugLevel
	}
	// set global log level
	logrus.SetLevel(ll)
	//metrics.InitMetrics(config.Instance.Kafka)
}
func loadConfig() {

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Info("Using config file: " + viper.ConfigFileUsed())
	}

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("fatal error config file: " + err.Error())
	}
	config.Instance = new(config.Config)
	err = viper.Unmarshal(&config.Instance)
	if err != nil {
		logrus.Fatal("fatal error unmashal config file: " + err.Error())
	}
	if viper.GetBool("dryRun") {
		config.Instance.DryRun = true
	}
}
