package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"netbox_go/internal/worker"
)

func main() {
	syncAll := flag.Bool("all", false, "Synchronize all data sources")
	flag.Parse()

	args := flag.Args()

	if !*syncAll && len(args) == 0 {
		fmt.Println("Error: Must specify at least one data source, or set --all.")
		os.Exit(1)
	}

	// In a real implementation, we would initialize the database and service layers here
	fmt.Println("Initializing synchronization...")

	if *syncAll {
		fmt.Println("Syncing all data sources...")
		// Logic to fetch all data sources and call .Sync() on each
	} else {
		for _, name := range args {
			fmt.Printf("Syncing data source: %s...\n", name)
			// Logic to fetch data source by name and call .Sync()
		}
	}

	fmt.Println("Finished.")
}
