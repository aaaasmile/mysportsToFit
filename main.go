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
	Buildnr = "00.02.20230914-00"
)

func main() {
	var ver = flag.Bool("ver", false, "Prints the current version")
	var configfile = flag.String("config", "config.toml", "Configuration file path")
	var outdir = flag.String("outdir", "./dest", "output directory for downloaded fit files")
	var hardcoded = flag.Bool("hardcoded", false, "use hard coded ids for the download (you need to rebuild if you want to use your ids)")
	var idonly = flag.String("idonly", "", "for example 205727472 to download only one fit activity id")
	var chunksize = flag.Int("chunksize", 5, "number of parallel download processes")
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

	if *hardcoded {
		mt.UseHardCodedIds()
	} else if *idonly != "" {
		mt.UseThisIdOnly(*idonly)
	}

	if err := mt.DownloadFit(*outdir); err != nil {
		log.Fatal(err)
	}
	log.Println("That's all folks!")
}
