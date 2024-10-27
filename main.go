package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strings"
	"syscall"
)

const (
	RED     = "\033[0;31m"
	GREEN   = "\033[0;32m"
	YELLOW  = "\033[0;33m"
	BLUE    = "\033[0;34m"
	MAGENTA = "\033[0;35m"
	CYAN    = "\033[0;36m"
	WHITE   = "\033[0;37m"
	NC      = "\033[0m"
)

func cleanup() {
	fmt.Printf("\n%sProgram terminated by the user.%s\n", RED, NC)
	os.Exit(1)
}

func removeUnwantedHTML(filePath string) {
	// Read file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("%sError reading file: %v%s\n", RED, err, NC)
		return
	}

	// Define regex patterns to remove
	patterns := []string{
		`<div class="col-md-12 text-center">.*?</div>`,
		`<a class="twitter-follow-button">.*?</div>`,
		`Check your IP address [0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`,
		`<div class="padding-block">.*?</div>`,
		`<nav class="navbar navbar-default">.*?</nav>`,
		`<ul class="nav navbar-nav">.*?</ul>`,
		`<form class="navbar-form navbar-left".*?</form>`,
		`<img alt="Brand".*?>`,
	}

	htmlContent := string(content)

	// Apply all regex replacements
	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?s)` + pattern)
		htmlContent = re.ReplaceAllString(htmlContent, "")
	}

	// Write back to file
	err = ioutil.WriteFile(filePath, []byte(htmlContent), 0644)
	if err != nil {
		fmt.Printf("%sError writing file: %v%s\n", RED, err, NC)
		return
	}
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Function to validate IP address format
func isValidIP(ip string) bool {
	return regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$`).MatchString(ip)
}

func main() {
	// Define command line flags
	ipAddr := flag.String("ip", "", "IP address to search for")
	flag.Usage = func() {
		fmt.Printf("Usage: %s -ip <ip_address>\n", os.Args[0])
		fmt.Println("Example: ./torrentspyder -ip 8.8.8.8")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Set up interrupt signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
	}()

	clearScreen()

	// Print ASCII art banner
	fmt.Println(" _____                          _   __                 _           ")
	fmt.Println("/__   \\___  _ __ _ __ ___ _ __ | |_/ _\\_ __  _   _  __| | ___ _ __ ")
	fmt.Println("  / /\\/ _ \\| '__| '__/ _ \\ '_ \\| __\\ \\| '_ \\| | | |/ _` |/ _ \\ '__|")
	fmt.Println(" / / | (_) | |  | | |  __/ | | | |__\\ \\ |_) | |_| | (_| |  __/ |   ")
	fmt.Println(" \\/   \\___/|_|  |_|  \\___|_| |_|\\__\\__/ .__/ \\__, |\\__,_|\\___|_|   ")
	fmt.Println("                                      |_|    |___/                 ")
	fmt.Println("                                                                   ")
	fmt.Printf("%sUnveil the IP's Hidden Secrets with TorrentSpyder - AnonKryptiQuz\n%s\n", GREEN, NC)

	userIP := *ipAddr
	if userIP == "" {
		// If IP not provided via command line, prompt for input
		fmt.Print("Enter the IP address: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		userIP = strings.TrimSpace(input)
	}

	// Validate IP address format
	if !isValidIP(userIP) {
		fmt.Printf("%sInvalid IP address. Please enter a valid IP.%s\n", RED, NC)
		return
	}

	// Download webpage
	url := fmt.Sprintf("https://iknowwhatyoudownload.com/en/peer/?ip=%s", userIP)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%sUnable to download webpage.%s\n", RED, NC)
		return
	}
	defer resp.Body.Close()

	// Create and write to file
	file, err := os.Create("downloaded_page.html")
	if err != nil {
		fmt.Printf("%sUnable to create HTML file.%s\n", RED, NC)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("%sError writing to file.%s\n", RED, NC)
		return
	}

	removeUnwantedHTML("downloaded_page.html")

	fmt.Printf("%sWebpage downloaded successfully. Opening in default web browser...%s\n", GREEN, NC)

}
