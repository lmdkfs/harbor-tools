package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"harbor-tools/harbor-tools/controllers"
	"harbor-tools/harbor-tools/db"
	"harbor-tools/harbor-tools/models"
	"log"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch tags from harbor",
	Run: func(cmd *cobra.Command, args []string) {
		MyDB, _, err := db.DB()
		defer MyDB.Close()
		if err != nil {
			log.Panicf("db init fail error:", err)
		}
		MyDB.AutoMigrate(&models.ImageTag{}, &models.DiffTag{})

		fmt.Println("Start fetch tags")
		//controllers.AddTagDemo()
		controllers.FetchTags()

	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().String("db-host", "", "Postgres address")
	fetchCmd.Flags().String("db-user", "", "Postgres db user")
	fetchCmd.Flags().String("db-pwd", "", "Postgres db password")
	fetchCmd.Flags().String("db-name", "", "Postgres db name")
	fetchCmd.Flags().String("db-port", "5432", "Postgres db port")
	viper.BindPFlag("db.host", fetchCmd.Flags().Lookup("db-host"))
	viper.BindPFlag("db.user", fetchCmd.Flags().Lookup("db-user"))
	viper.BindPFlag("db.password", fetchCmd.Flags().Lookup("db-pwd"))
	viper.BindPFlag("db.name", fetchCmd.Flags().Lookup("db-name"))
	viper.BindPFlag("db.port", fetchCmd.Flags().Lookup("db-port"))
}
