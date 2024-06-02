package cmd

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"url/pkg/config"
	"url/pkg/mongodb"

	"github.com/spf13/cobra"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var addCmd = &cobra.Command{
	Use:   "add [url]",
	Short: "Generate short URL",
	Long: `Generate and add a shorted URL to the database. For example:
	urlcli add https://example.com`,
	Args: cobra.ExactArgs(1),
	Run:  addURL,
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().IntP("ml", "m", config.Config.DefaultLength, "Max length of short URL. Must be between 6 and 32.")
	addCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}

func addURL(cmd *cobra.Command, args []string) {
	url := args[0]
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	maxLength, err := cmd.Flags().GetInt("ml")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

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

	shortURL, err := getShortURL(mongoClient, url, verbose)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("URL:", shortURL)
}

func validateURLAndLength(url string, maxLength int) error {
	if err := validateURL(url); err != nil {
		return err
	}
	return validateMaxLength(maxLength)
}

func validateURL(url string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return errors.New("URL must start with http:// or https://")
	}

	pattern := `^(http(s)?:\/\/)?([\w-]+\.)+[\w-]+(\/[\w- ;,./?%&=]*)?$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(url) {
		return errors.New("invalid URL format")
	}

	return nil
}

func validateMaxLength(maxLengthValue int) error {
	if maxLengthValue < config.Config.MinLength || maxLengthValue > config.Config.MaxLength {
		return fmt.Errorf("max length of short URL should be between %d and %d", config.Config.MinLength, config.Config.MaxLength)
	}
	return nil
}

func getShortURL(client *mongo.Client, url string, verbose bool) (string, error) {
	value, err := mongodb.GetValue(client, config.Config.CollectionName, config.Config.DbName, url)
	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
	}

	if value == "" {
		shortURL := generateShortURL(config.Config.DefaultLength)
		err = mongodb.Set(client, config.Config.CollectionName, config.Config.DbName, url, shortURL)
		if err != nil {
			return "", err
		}
		if verbose {
			fmt.Println("URL", url, "added to database")
		}
		return "http://localhost:" + config.Config.ServerPort + "/" + shortURL, nil
	}

	if verbose {
		fmt.Println("URL", url, "already exists in the database")
	}
	return "http://localhost:" + config.Config.ServerPort + "/" + value, nil
}

func generateShortURL(maxLength int) string {
	b := make([]byte, maxLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
