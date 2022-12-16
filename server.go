package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/localyyz/go-shopify"

	// Import the godotenv package to load the environment variables from the .env file
	"github.com/joho/godotenv"
)

// Initialize a new Shopify client
func newShopifyClient(shopName string, accessToken string) *shopify.Client {
	return shopify.NewClient(shopName, accessToken, nil)
}

func main() {
	r := mux.NewRouter()

	// Handle the homepage
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to your Shopify app!")
	})

	// Handle the OAuth callback
	r.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		// Get the API key and secret from the environment or .env file
		apiKey := os.Getenv("SHOPIFY_API_KEY")
		if apiKey == "" {
			// Load the API key and secret from the .env file
			err := godotenv.Load(".env")
			if err != nil {
				fmt.Fprintf(w, "Error loading .env file: %s", err)
				return
			}
			apiKey = os.Getenv("SHOPIFY_API_KEY")
			apiSecret := os.Getenv("SHOPIFY_API_SECRET")
		}

		// Get the shop name from the query string
		shopName := r.URL.Query().Get("shop")

		// Use the shopify.New() function to create a new Shopify client
		client := shopify.New(apiKey, apiSecret)

		// Use the GetAccessToken() method of the Shopify client to get a SHOPIFY_ACCESS_TOKEN
		accessToken, err := client.GetAccessToken()
		if err != nil {
			fmt.Fprintf(w, "Error getting access token: %s", err)
			return
		}

		// Fetch a list of products for the authenticated store
		products, err := client.Product.List(context.Background(), nil)
		if err != nil {
			fmt.Fprintf(w, "Error fetching products: %s", err)
			return
		}

		// Print the titles of the first 10 products in the list
		for i := 0; i < 10; i++ {
			product := products[i]
			fmt.Fprintln(w, product.Title)
		}
	})

}
