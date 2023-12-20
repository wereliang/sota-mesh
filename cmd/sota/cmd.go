package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wereliang/sota-mesh/pkg/config"
	"github.com/wereliang/sota-mesh/pkg/server"
)

var RootCmd = &cobra.Command{
	Use:   "./sota",
	Short: "Sota agent",
	Long:  `Sota is a smart open traffic agent`,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start sota agent",
	Long:  `Start sota agent`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := cmd.Flags().GetString("config")
		if cfg == "" {
			cmd.Help()
			return
		}
		ctype, _ := cmd.Flags().GetString("type")
		sota, err := server.NewSota(cfg, config.CfgType(ctype))
		if err != nil {
			panic(err)
		}
		sota.Start()
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("config", "c", "", "config file")
	startCmd.Flags().StringP("type", "t", "sota", "config type[envoy|sota]")
}
