package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Strubbl/wallabago"
)

const version = "0.1"
const defaultConfigJSON = "config.json"

var debug = flag.Bool("d", false, "get debug output (implies verbose mode)")
var debugDebug = flag.Bool("dd", false, "get even more debug output like data (implies debug mode)")
var v = flag.Bool("v", false, "print version")
var verbose = flag.Bool("verbose", false, "verbose mode")
var configJSON = flag.String("config", defaultConfigJSON, "file name of config JSON file")
var outfolder = flag.String("outfolder", "out", "file name of config JSON file")

func handleFlags() {
	flag.Parse()
	if *debug && len(flag.Args()) > 0 {
		log.Println("handleFlags: non-flag args=%v", strings.Join(flag.Args(), " "))
	}
	// version first, because it directly exits here
	if *v {
		log.Println("version %v\n", version)
		os.Exit(0)
	}
	// test verbose before debug because debug implies verbose
	if *verbose && !*debug && !*debugDebug {
		log.Println("verbose mode")
	}
	if *debug && !*debugDebug {
		log.Println("handleFlags: debug mode")
		// debug implies verbose
		*verbose = true
	}
	if *debugDebug {
		log.Println("handleFlags: debug mode")
		// debugDebug implies debug
		*debug = true
		// and debug implies verbose
		*verbose = true
	}
}

func GetReaderEntries() ([]wallabago.Item, error) {
	page := -1
	perPage := -1
	e, err := wallabago.GetEntries(wallabago.APICall, -1, -1, "", "", page, perPage, "toreader")
	if err != nil {
		log.Println("GetAllEntries: call failed", err)
		return nil, err
	}
	allEntries := e.Embedded.Items
	if e.Total > len(allEntries) {
		secondPage := e.Page + 1
		perPage = e.Limit
		pages := e.Pages
		for i := secondPage; i <= pages; i++ {
			e, err := wallabago.GetEntries(wallabago.APICall, -1, -1, "", "", i, perPage, "")
			if err != nil {
				log.Printf("GetAllEntries: GetEntries for page %d failed: %v", i, err)
				return nil, err
			}
			tmpAllEntries := e.Embedded.Items
			allEntries = append(allEntries, tmpAllEntries...)
		}
	}
	return allEntries, err
}

func main() {
	log.SetOutput(os.Stdout)
	handleFlags()
	// check for config
	if *verbose {
		log.Println("reading config", *configJSON)
	}
	err := wallabago.ReadConfig(*configJSON)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	items, err := GetReaderEntries()
	if err != nil {
		log.Fatal(err)
	}
	if *verbose {
		log.Println("Got", len(items), "entries")
	}
	for i := 0; i < len(items); i++ {
		fname := strings.ReplaceAll(path.Base(items[i].Title), ":", "")
		if *verbose {
			log.Println("Getting", i+1, fname)
		}
		export, err := wallabago.ExportEntry(wallabago.APICall, items[i].ID, "mobi")
		if err != nil {
			log.Fatal(err)
		}
		output := filepath.Join(*outfolder, fname+".mobi")
		if *verbose {
			log.Println("Creating", i+1, fname)
		}
		out, err := os.Create(output)
		if err != nil {
			log.Println("failed to create output file:", err)
		}
		defer out.Close()
		if *verbose {
			log.Println("Saving", i+1, fname)
		}
		n, err := out.Write(export)
		if err != nil {
			log.Println("can't write file:", err)
			if n == 0 {
				log.Println("can't write file, probably output folder is missing: ")
			}
		}
		if n > 0 {
			log.Println("wrote", n, "bytes in file:", output)
		}
	}
}
