package cmd

import (
	"context"
	"fmt"
	"os"
	"url/pkg/config"
	"url/pkg/mongodb"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a URL from the database",
	Long:  `This command deletes a URL from the database. It takes one argument which is the URL to delete.`,
	Args:  cobra.ExactArgs(1),
	Run:   deleteURL,
}

func init() {

	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}

func deleteURL(cmd *cobra.Command, args []string) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	url := args[0]

	mongoClient, err := mongodb.MongoConnect(config.Config.DbHost, config.Config.DbPort, config.Config.DbUser, config.Config.DbPass)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer mongoClient.Disconnect(context.TODO())

	if err := performDelete(mongoClient, url, verbose); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func performDelete(client *mongo.Client, url string, verbose bool) error {
	_, err := mongodb.GetValue(client, config.Config.CollectionName, config.Config.DbName, url)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("URL %s not found in the database", url)
	}

	if err := mongodb.Delete(client, config.Config.CollectionName, config.Config.DbName, url); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("URL %s deleted from the database\n", url)
	}
	return nil
}
