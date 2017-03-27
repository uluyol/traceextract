package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	start   = flag.Duration("start", 0, "start time")
	end     = flag.Duration("end", 0, "end time")
	inpath  = flag.String("in", "", "input csv path")
	outpath = flag.String("out", "", "output csv path")
)

func usage() {
	fmt.Fprintln(os.Stderr, "tracecutter -start S -end E -in in.csv -out out.csv")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("tracecutter: ")
	flag.Parse()
	if *inpath == "" || *outpath == "" {
		usage()
	}

	fin, err := os.Open(*inpath)
	if err != nil {
		log.Fatal(err)
	}

	fout, err := os.Create(*outpath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(fin)
	tmin := *end
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" || s[0] == '#' {
			continue
		}
		tend := strings.IndexByte(s, ',')
		if tend == -1 {
			log.Fatal("want time,val pairs: got missing ,")
		}
		if tend == len(s)-1 {
			log.Fatal("missing value")
		}
		t, err := strconv.ParseFloat(s[:tend], 64)
		if err != nil {
			log.Fatalf("unable to parse time: %v", err)
		}
		v, err := strconv.ParseFloat(s[tend+1:], 64)
		if err != nil {
			log.Fatalf("unable to parse val: %v", err)
		}
		tdur := time.Duration(t * float64(time.Second))
		if *start <= tdur && tdur <= *end {
			if tmin > tdur {
				tmin = tdur
			}
			fmt.Fprintf(fout, "%f,%f\n", float64(tdur-tmin)/float64(time.Second), v)
		}
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

	fout.Close()
	fin.Close()
}
