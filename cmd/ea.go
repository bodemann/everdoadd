package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const AutoVersion = "0.3.1"

const ApiKeyEnvName = "EVERDO_API_KEY"
const IpAddressEnvName = "EVERDO_IP_ADDRESS"

const DefaultApiKey = ""
const DefaultIpAddress = "localhost:11111"

func main() {
	todoTitle, todoDescription := handleParameter()
	postBody := createPostBody(todoTitle, todoDescription)
	responseBody := bytes.NewBuffer(postBody)
	apiKey, ipAddress := getEnvironmentVariables()
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // skip certificate check for everdo
	postUrl := fmt.Sprintf("https://%s/api/items?key=%s", ipAddress, apiKey)
	resp, err := http.Post(postUrl, "application/json", responseBody)
	if err != nil {
		fmt.Printf("error in POST request to everdo occured: %v\n", err)
		os.Exit(1)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
}

func usage() {
	fmt.Println("everdoadd (ea), a program to add tasks to a Everdo (everdo.net) inbox.")
	fmt.Println("Version: ", AutoVersion)
	fmt.Println("everdoadd uses two environment variables:")
	fmt.Printf("\t%s to tell everdoadd on which IP/port Everdo is listening. Default is %s\n", IpAddressEnvName, DefaultIpAddress)
	fmt.Printf("\t%s to tell everdoadd which API-key Everdo is expecting.\n", ApiKeyEnvName)
	fmt.Println("Usage:")
	fmt.Println("\tea TEXT             Adds a new todo item to the everdo inbox ")
	fmt.Println("\tea TEXT DESCRIPTION Adds a new todo item and a description to the everdo inbox")
	fmt.Println("\t                    The first word is the title, all following words are the description")
	fmt.Println("\t                    Example:")
	fmt.Println("\t                    ea Vacation Plan vacation with all family members")
	fmt.Println("\t                    Vacation will be the todo name, Plan... the description")
	fmt.Println("\t                    Use \" \" for multi word titles")
	fmt.Println("\tea --debug          CAUTION: prints IP address and API-key to stdout")
	fmt.Println("\tea --help           Prints this text")
	fmt.Println("\tea --version        Prints version information")
	fmt.Println()
}

// getEnvironmentVariables reads the environment variables
// or returns default values
func getEnvironmentVariables() (apiKey, ipAddress string) {
	apiKey = os.Getenv(ApiKeyEnvName)
	if apiKey == "" {
		apiKey = DefaultApiKey
	}
	if apiKey == "" {
		fmt.Printf("ERROR: missing API-key\n")
		fmt.Printf("everdoadd (ea) needs an API key to communicate with Everdo. You can find/set one in Everdo->Settings->API. See the README of everdoadd for more information.\n")
		fmt.Printf("everdoadd uses the environment variable %s to find the API key.\n", ApiKeyEnvName)
		os.Exit(1)
	}
	ipAddress = os.Getenv(IpAddressEnvName)
	if ipAddress == "" {
		ipAddress = DefaultIpAddress
	}
	if !strings.Contains(ipAddress, ":") {
		fmt.Printf("ERROR: missing port number in Everdo IP address, set by %s\n", IpAddressEnvName)
		fmt.Printf("Everdo IP address also needs a port, e.g. 192.168.10.10:11223 The Everdo default port is :11111\n")
		os.Exit(1)
	}
	return apiKey, ipAddress
}

// handleParameter handles all program options like --VersionAutoCounter, --debug
// returns title and description no options were used
func handleParameter() (string, string) {
	if len(os.Args[1:]) == 0 {
		fmt.Println("ERROR: missing argument")
		usage()
		os.Exit(1)
	}
	todoTitleOrParameter := os.Args[1]
	todoDescription := strings.Join(os.Args[2:], " ")
	switch todoTitleOrParameter {
	case "--help":
		usage()
		os.Exit(0)
	case "--version":
		fmt.Println(AutoVersion)
		os.Exit(0)
	case "--debug":
		apiKey, ipAddress := getEnvironmentVariables()
		fmt.Printf("using everdo at IP address:%s with API key:%s\n", ipAddress, apiKey)
		os.Exit(0)
	}
	return todoTitleOrParameter, todoDescription
}

// createPostBody returns a body for http.Post from a todo_title and its description
func createPostBody(title, description string) []byte {
	parameter := make(map[string]string)
	parameter["title"] = title
	parameter["note"] = description
	postBody, err := json.Marshal(parameter)
	if err != nil {
		fmt.Printf("ERROR: could not convert parameter to JSON structure %s\n", err)
		os.Exit(1)
	}
	return postBody
}
