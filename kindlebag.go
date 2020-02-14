package main

import (
	"flag"
	"fmt"
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
		log.Printf("handleFlags: non-flag args=%v", strings.Join(flag.Args(), " "))
	}
	// version first, because it directly exits here
	if *v {
		fmt.Printf("version %v\n", version)
		os.Exit(0)
	}
	// test verbose before debug because debug implies verbose
	if *verbose && !*debug && !*debugDebug {
		log.Printf("verbose mode")
	}
	if *debug && !*debugDebug {
		log.Printf("handleFlags: debug mode")
		// debug implies verbose
		*verbose = true
	}
	if *debugDebug {
		log.Printf("handleFlags: debug² mode")
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
		log.Println("GetAllEntries: first GetEntries call failed", err)
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
		fmt.Println(err.Error())
		os.Exit(1)
	}
	items, err := GetReaderEntries()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(items); i++ {
		exp, err := wallabago.ExportEntry(wallabago.APICall, items[i].ID, "mobi")
		if err != nil {
			panic(err)
		}
		fname := strings.ReplaceAll(path.Base(items[i].Title), ":", "")
		output := filepath.Join(*outfolder, fname+".mobi")
		out, err := os.Create(output)
		if err != nil {
			fmt.Errorf("failed to create output file: %v", err)
		}
		defer out.Close()
		n, err := out.Write(exp)
		if err != nil {
			fmt.Errorf("can't write file: %v", err)
		}
		if n >= 0 {
			log.Printf("wrote %d bytes (%s) in file %s", n, uint64(n), output)
		}
	}
}
