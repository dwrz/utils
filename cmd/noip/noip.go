package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	domain   = flag.String("domain", "", "domain to update (required)")
	user     = flag.String("user", "", "noip username (required)")
	password = flag.String("password", "", "noip password (required)")
)

func main() {
	flag.Parse()

	if domain == nil || *domain == "" {
		fmt.Fprintf(os.Stderr, "missing domain flag\n")
		return
	}
	if user == nil || *user == "" {
		fmt.Fprintf(os.Stderr, "missing user flag\n")
		return
	}
	if password == nil || *password == "" {
		fmt.Fprintf(os.Stderr, "missing password flag\n")
		return
	}

	url := "https://dynupdate.no-ip.com/nic/update?hostname=" + *domain

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create request: %v\n", err)
		return
	}
	request.SetBasicAuth(*user, *password)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "update request failed: %v\n", err)
		return
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Fprintf(
			os.Stderr, "failed to read response body: %v\n", err,
		)
		return
	}

	fmt.Println(string(content))
}
