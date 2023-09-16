package main

import (
	"flag"
	"fmt"
	"log"
	"mysportsToFit/conf"
	"mysportsToFit/mytom"
	"os"
)

const (
	Appname = "Mysports to Fit"
	Buildnr = "00.02.20230916-00"
)

func main() {
	var ver = flag.Bool("ver", false, "Prints the current version")
	var configfile = flag.String("config", "config.toml", "Configuration file path")
	var outdir = flag.String("outdir", "./dest", "output directory for downloaded fit files")
	var hardcoded = flag.Bool("hardcoded", false, "use hard coded ids for the download (you need to rebuild if you want to use your ids)")
	var idonly = flag.String("idonly", "", "for example 205727472 to download only one fit activity id")
	var chunksize = flag.Int("chunksize", 5, "number of parallel download processes")
	var source = flag.String("source", "", "Hal file path to get activity ids")
	var nodownload = flag.Bool("nodownload", false, "prepare the activity list but do not exec the download")
	flag.Parse()

	if *ver {
		fmt.Printf("%s  version %s", Appname, Buildnr)
		os.Exit(0)
	}
	current, err := conf.ReadConfig(*configfile)
	if err != nil {
		log.Fatal("Error on read config file", err)
	}
	mt := mytom.NewMyTom(current.Email, current.Password, *chunksize)

	log.Println("settings: ", *outdir, *hardcoded, *idonly, *chunksize, *source)

	if *hardcoded {
		mt.UseHardCodedIds()
	} else if *idonly != "" {
		mt.UseThisIdOnly(*idonly)
	} else if *source != "" {
		if err := mt.UseHalSource(*source); err != nil {
			log.Fatal("Hal processing error: ", err)
		}
	} else {
		log.Fatal("No activity ids to downlad")
	}

	if err := mt.DownloadFit(*outdir, *nodownload); err != nil {
		log.Fatal(err)
	}
	log.Println("That's all folks!")
}
