package main

import (
	"context"
	"os"
	"os/signal"

	"asset-measurements-assignment/internal/entrypoints/simulator"
	"github.com/GLCharge/otelzap"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	devxCfg "github.com/xBlaz3kx/DevX/configuration"
	"go.uber.org/zap"
)

const serviceName = "simulator"

var (
	// Configuration
	configurationFilePathFlag string

	rootCmd = &cobra.Command{
		Use:   "simulator",
		Short: "simulator",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, args []string) {
			devxCfg.InitConfig(configurationFilePathFlag, "./config", ".")
			devxCfg.SetupEnv(serviceName)
			devxCfg.SetDefaults(serviceName)
			viper.SetDefault("environment", "development")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get configuration from Viper
			otelzap.L().Info("Parsing configuration")

			cfg := &simulator.Config{}
			devxCfg.GetConfiguration(viper.GetViper(), cfg)

			// Run the simulator
			return simulator.Run(cmd.Context(), *cfg)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configurationFilePathFlag, "config", "c", "", "config file path")

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		otelzap.L().Fatal("Unable to run", zap.Error(err))
	}
}
