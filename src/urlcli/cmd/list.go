package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"url/pkg/config"
	"url/pkg/mongodb"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entries in the database",
	Long:  `This command lists all entries in the database. No arguments are required.`,
	Run:   listEntries,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listEntries(cmd *cobra.Command, args []string) {
	mongoClient, err := mongodb.MongoConnect(config.Config.DbHost, config.Config.DbPort, config.Config.DbUser, config.Config.DbPass)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer mongoClient.Disconnect(context.TODO())

	if err := printAllEntries(mongoClient); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func printAllEntries(client *mongo.Client) error {
	allUrls, err := mongodb.GetAll(client, config.Config.CollectionName, config.Config.DbName)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "URL\tShort URL")

	for shortURL, url := range allUrls {
		fmt.Fprintf(w, "%s\t%s\n", shortURL, url)
	}

	w.Flush()
	return nil
}
