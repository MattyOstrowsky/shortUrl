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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a URL in the database",
	Long:  `This command updates a URL in the database. It takes one argument which is the URL to update.`,
	Args:  cobra.ExactArgs(1),
	Run:   updateURL,
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	updateCmd.Flags().IntP("ml", "m", config.Config.DefaultLength, "Max length of short URL. Must be between 6 and 32.")
}

func updateURL(cmd *cobra.Command, args []string) {
	url := args[0]
	verbose, _ := cmd.Flags().GetBool("verbose")
	maxLength, _ := cmd.Flags().GetInt("ml")

	if err := validateURLAndLength(url, maxLength); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	mongoClient, err := mongodb.MongoConnect(config.Config.DbHost, config.Config.DbPort, config.Config.DbUser, config.Config.DbPass)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer mongoClient.Disconnect(context.TODO())

	if err := performUpdate(mongoClient, url, maxLength, verbose); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func performUpdate(client *mongo.Client, url string, maxLength int, verbose bool) error {
	shortURL := generateShortURL(maxLength)

	_, err := mongodb.GetValue(client, config.Config.CollectionName, config.Config.DbName, url)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("URL %s not found in the database", url)
	}

	if err := mongodb.Update(client, config.Config.CollectionName, config.Config.DbName, url, shortURL); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("URL %s updated in the database\n", url)
	}
	fmt.Println("New URL:", "http://localhost:"+config.Config.ServerPort+"/"+shortURL)
	return nil
}
