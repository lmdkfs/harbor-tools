package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"harbor-tools/harbor-tools/server"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "start http server",
	Run: func(cmd *cobra.Command, args []string) {
		server := server.NewServer()
		server()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	cobra.OnInitialize(initServerConfig)

	serverCmd.Flags().String("port", "8080", "listen port")
	serverCmd.Flags().String("host", "0.0.0.0", "listen ip")
	serverCmd.Flags().String("mode", "info", "log level")
	serverCmd.Flags().String("logname", "harbor-tools", "log name")

	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
	fmt.Println("server init get server.port:", serverCmd.Flags().Lookup("port").Value.String()  )
	viper.BindPFlag("server.host", serverCmd.Flags().Lookup("host"))
	viper.BindPFlag("server.logname", serverCmd.Flags().Lookup("logname"))

	viper.SetDefault("server.readtimeout", "1000")
	viper.SetDefault("server.writetimeout", "10000")
}

func initServerConfig(){

}

