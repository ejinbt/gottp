package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Domain struct {
	url    string
	status string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkUrl(url_value string) string {
	if !(strings.Contains(url_value, "https://")) {
		url_value = "https://" + url_value
	}
	return url_value
}

func client(c http.Client, url string) http.Response {
	parsed_url := checkUrl(url)
	req, err := http.NewRequest("GET", parsed_url, nil)
	if err != nil {
		fmt.Printf("Error %s \n", err)
	}
	req.Header.Add("Accept", `application/json`)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("error %s \n", err)
	}

	defer resp.Body.Close()
	return *resp
}

func read_file(c http.Client, filepath string) []Domain {
	var domains []Domain
	f, err := os.Open(filepath)
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		resp := client(c, scanner.Text())
		temp_domain := Domain{url: scanner.Text(), status: resp.Status}
		domains = append(domains, temp_domain)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return domains
}

func file_output(domains []Domain) string {
	output_file, err := os.Create("output.txt")
	check(err)
	defer output_file.Close()

	for _, domain := range domains {
		_, err := output_file.WriteString(domain.url + ":" + domain.status + "\n")
		check(err)
	}
	return output_file.Name()
}

func main() {
	c := http.Client{Timeout: time.Duration(5) * time.Second}
	var value string
	url := flag.String("url", "", "url to send")
	url_file := flag.String("file", "", "file to read url (need to have http:// or https:// prefix)")
	flag.Parse()
	if len(*url) != 0 {
		value = *url
		fmt.Println(value)
		client(c, value)
	} else if len(*url_file) != 0 {
		value = *url_file
		domains := read_file(c, value)
		filename := file_output(domains)
		fmt.Println("data saved to ", filename)
	}

}
