package main

import (
	"flag"
	"fmt"
	"log"
	"mytom/conf"
	"mytom/mytom"
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
	var withexample = flag.Bool("withexample", false, "use example ids to check if it works")
	var idonly = flag.String("idonly", "", "for example 205727472 to download only one fit activity id")
	flag.Parse()

	if *ver {
		fmt.Printf("%s  version %s", Appname, Buildnr)
		os.Exit(0)
	}
	current, err := conf.ReadConfig(*configfile)
	if err != nil {
		log.Fatal("Error on read config file", err)
	}
	mt := mytom.NewMyTom(current.Email, current.Password)

	if *withexample {
		mt.UseExampleIds()
	} else if *idonly != "" {
		mt.UseThisIdOnly(*idonly)
	}

	if err := mt.DownloadFit(*outdir); err != nil {
		log.Fatal(err)
	}
	log.Println("That's all folks!")
}
