package prompt

import (
	"fmt"
	"github.com/GreatGodApollo/ala/api"
	"github.com/GreatGodApollo/ala/models"
	"github.com/c-bata/go-prompt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var rCli *resty.Client
var baseUrl, username, password string

var suggestions = []prompt.Suggest{
	// Basics
	{"exit", "Quit ALC"},
	{"help", "List commands"},

	// API Stuff
	{"all", "Get valid licenses for a product"},
	{"new", "Generate a new license for a product"},
	{"invalidate", "Invalidate a license"},
	{"get", "Get a specific license"},
	{"check", "Check if a license is valid"},
}

func RunPrompt(client *resty.Client) {
	rCli = client
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">> "),
		prompt.OptionTitle("ALC"))

	username = viper.GetString("auth.username")
	password = viper.GetString("auth.password")
	baseUrl = viper.GetString("api.baseurl")

	p.Run()
}

func executor(in string) {
	in = strings.TrimSpace(in)

	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "exit":
		fmt.Println("Thanks for using ALC!")
		os.Exit(0)
	case "help":
		printCommands()
		break
	case "new":
		if len(blocks) > 2 {
			resp, err := api.CreateLicense(rCli, baseUrl, username, password, blocks[1], blocks[2])
			if err != nil {
				fmt.Println("An error occurred:")
				fmt.Println(err.Error())
			}
			if _, ok := resp.(models.LicenseResponse); ok {
				licenseObj := resp.(models.LicenseResponse)
				printLicenseResponse(licenseObj)
			}
			break
		} else {
			fmt.Println("new <email> <product>")
			break
		}
	case "all":
		if len(blocks) > 1 {
			resp, err := api.GetAll(rCli, baseUrl, username, password, blocks[1])
			if err != nil {
				fmt.Println("An error occurred:")
				fmt.Println(err.Error())
			}
			if _, ok := resp.(models.Licenses); ok {
				licensesObj := resp.(models.Licenses)
				if len(licensesObj.Licenses) != 0 {
					for _, license := range licensesObj.Licenses {
						printLicense(license)
					}
				} else {
					fmt.Println("No valid licenses found for that product!")
				}
				break
			} else if _, ok := resp.(models.BasicResponse); ok {
				respObj := resp.(models.BasicResponse)
				fmt.Println("An error occurred:")
				fmt.Println(respObj.Message)
			}
		} else {
			fmt.Println("all <product>")
			break
		}
	case "get":
		if len(blocks) > 1 {
			resp, err := api.GetSpecific(rCli, baseUrl, username, password, blocks[1])
			if err != nil {
				fmt.Println("An error occurred:")
				fmt.Println(err.Error())
			}
			if _, ok := resp.(models.License); ok {
				licenseObj := resp.(models.License)
				printLicense(licenseObj)
				break
			} else if _, ok := resp.(models.BasicResponse); ok {
				respObj := resp.(models.BasicResponse)
				if respObj.Message == "license nonexistent" {
					fmt.Println("That license doesn't exist!")
				} else {
					fmt.Println("An error occurred:")
					fmt.Println(respObj.Message)
				}
			}
		} else {
			fmt.Println("get <key>")
			break
		}
	case "check":
		if len(blocks) > 2 {
			product := blocks[2]
			key := blocks[1]
			valid := api.CheckValidity(rCli, baseUrl, key, product)
			fmt.Println("---")
			fmt.Printf("License Key: %s\n", key)
			fmt.Printf("Product: %s\n", product)
			fmt.Printf("Valid: %t\n", valid)
			break
		}
		fmt.Println("check <key> <product>")
	case "invalidate":
		if len(blocks) > 1 {
			resp, err := api.InvalidateLicense(rCli, baseUrl, username, password, blocks[1])
			if err != nil {
				fmt.Println("An error occurred:")
				fmt.Println(err.Error())
			}

			if respObj, ok := resp.(models.LicenseResponse); ok {
				printLicenseResponse(respObj)
			} else if respObj2, ok2 := resp.(models.BasicResponse); ok2 {
				fmt.Println("An error occurred:")
				fmt.Println(respObj2.Message)
			}
			break
		} else {
			fmt.Println("invalidate <license>")
			break
		}
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	w := d.GetWordBeforeCursor()
	if w == "" {
		return suggestions
	}

	return prompt.FilterHasPrefix(suggestions, w, true)
}

func printCommands() {
	for _, s := range suggestions {
		fmt.Printf("%s - %s\n", s.Text, s.Description)
	}
}

func printLicense(license models.License) {
	fmt.Println("---")
	fmt.Printf("License Key: %s\n", license.LicenseKey)
	fmt.Printf("Product: %s\n", license.Product)
	fmt.Printf("Valid: %t\n", license.Valid)
	fmt.Printf("Email: %s\n", license.Email)
}

func printLicenseResponse(response models.LicenseResponse) {
	fmt.Println("---")
	fmt.Printf("License Key: %s\n", response.LicenseKey)
	fmt.Printf("Status: %s\n", response.Status)
	fmt.Printf("Code: %d\n", response.Code)
	fmt.Printf("Message: %s\n", response.Message)
}
