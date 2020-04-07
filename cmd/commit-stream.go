package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/x1sec/commit-stream/pkg"
	"os"
)

func printAscii() {
	h := `
 ██████╗ ██████╗ ███╗   ███╗███╗   ███╗██╗████████╗   ███████╗████████╗██████╗ ███████╗ █████╗ ███╗   ███╗
██╔════╝██╔═══██╗████╗ ████║████╗ ████║██║╚══██╔══╝   ██╔════╝╚══██╔══╝██╔══██╗██╔════╝██╔══██╗████╗ ████║
██║     ██║   ██║██╔████╔██║██╔████╔██║██║   ██║█████╗███████╗   ██║   ██████╔╝█████╗  ███████║██╔████╔██║
██║     ██║   ██║██║╚██╔╝██║██║╚██╔╝██║██║   ██║╚════╝╚════██║   ██║   ██╔══██╗██╔══╝  ██╔══██║██║╚██╔╝██║
╚██████╗╚██████╔╝██║ ╚═╝ ██║██║ ╚═╝ ██║██║   ██║      ███████║   ██║   ██║  ██║███████╗██║  ██║██║ ╚═╝ ██║
 ╚═════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝     ╚═╝╚═╝   ╚═╝      ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝ 
https://github.com/x1sec/commit-stream        

`

	fmt.Fprintf(os.Stderr, h)
}

func init() {
	flag.Usage = func() {
		printAscii()

		h := "Stream Github commit authors in realtime\n\n"

		h += "Usage:\n"
		h += "  commit-stream [OPTIONS]\n\n"

		h += "Options:\n"
		h += "  -e, --email       Match email addresses field (specify multiple with comma). Omit to match all.\n"
		h += "  -n, --name        Match author name field (specify multiple with comma). Omit to match all.\n"
		h += "  -t, --token       Github token (if not specified, will use environment variable 'CSTREAM_TOKEN')\n"
		h += "  -a  --all-commits Search through previous commit history (default: false)\n"
		h += "\n\n"
		fmt.Fprintf(os.Stderr, h)
	}
}

func main() {

	var (
		authToken        string
		rate             int
		filter           commitstream.FilterOptions
		searchAllCommits bool
	)

	flag.StringVar(&filter.Email, "email", "", "")
	flag.StringVar(&filter.Email, "e", "", "")

	flag.StringVar(&filter.Name, "name", "", "")
	flag.StringVar(&filter.Name, "n", "", "")

	flag.StringVar(&authToken, "token", "", "")
	flag.StringVar(&authToken, "t", "", "")
	flag.IntVar(&rate, "r", 0, "")
	flag.IntVar(&rate, "rate", 0, "")

	flag.BoolVar(&searchAllCommits, "a", false, "")
	flag.BoolVar(&searchAllCommits, "all-commits", false, "")

	flag.Parse()

	if filter.Email == "" && filter.Name == "" {
		filter.Enabled = false
	} else {
		filter.Enabled = true
	}

	if authToken == "" {
		authToken = os.Getenv("CSTREAM_TOKEN")
		if authToken == "" {
			fmt.Fprintf(os.Stderr, "Please specify Github authentication token with '-t' or by setting the environment variable CSTREAM_TOKEN\n")
			os.Exit(1)
		}

	}

	streamOpt := commitstream.StreamOptions{AuthToken: authToken, SearchAllCommits: searchAllCommits}
	commitstream.DoIngest(streamOpt, filter, handleResult)
}

func handleResult(s []string) {
	w := csv.NewWriter(os.Stdout)
	w.Write(s)
	w.Flush()
}