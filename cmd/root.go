package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stalwartgiraffe/cmr/internal/config"
	//"go.opentelemetry.io/otel"
	//"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	//"go.opentelemetry.io/otel/metric"
	//sdk "go.opentelemetry.io/otel/sdk/metric"
)

type CmdConfig struct {
	Config *config.Config
}

func NewRootCmd(cfg *CmdConfig) *cobra.Command {

	var cfgFilepath string
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "cmr",
		Short: "cmr is a dev workflow automation tool.",
		Long: `cmr is a dev workflow automation tool.
For example:

FIXME - write examples here.`,

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			cfg.Config, err = config.LoadConfigFile(cfgFilepath)
			if err != nil {
				log.Fatalf("Could not load config %s: %s", cfgFilepath, err)
			}
		},

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("No argument specified.")
			fmt.Println("")
			err := cmd.Help()
			if err != nil {
				log.Fatalf("Could not load config %s: %s", cfgFilepath, err)
			}
		},
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFilepath, "config", "", "config file (default is $HOME/.cmr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return rootCmd
}

func AddRootCommand(cancel context.CancelFunc) *cobra.Command {
	cfg := &CmdConfig{}

	rootCmd := NewRootCmd(cfg)

	addInitCommand(cfg, rootCmd)
	rootCmd.AddCommand(NewLabCommand(cfg))
	rootCmd.AddCommand(NewViewProjectsCommand(cancel, cfg))
	rootCmd.AddCommand(NewEventsCommand(cancel, cfg))
	rootCmd.AddCommand(NewMergeRequestCommand(cancel, cfg))
	rootCmd.AddCommand(NewPrjEventsCommand(cancel, cfg))
	rootCmd.AddCommand(NewCloneCommand(cfg))
	rootCmd.AddCommand(NewPullCommand(cfg))
	rootCmd.AddCommand(NewLintCommand(cfg))
	rootCmd.AddCommand(NewGacCommand(cfg, nil))
	rootCmd.AddCommand(NewPushCommand(cfg, nil))
	rootCmd.AddCommand(NewSecretToolCommand(cfg))
	return rootCmd
}
