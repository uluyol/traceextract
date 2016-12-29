package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"strings"

	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	searchColor = flag.String("color", "#000", "hex color to search for, will fuzzy match")
	minVal      = flag.Float64("ymin", 0, "y value at minimum (used to scale)")
	maxVal      = flag.Float64("ymax", 1, "y value at maximum (used to scale)")
	minIndex    = flag.Float64("xmin", 0, "x value at minimum (used to scale)")
	maxIndex    = flag.Float64("xmax", 1, "x value at maximum (used to scale)")
	onlyLongest = flag.Bool("onlylongest", true, "only get the longest continuous run of the color")
	longestGap  = flag.Int("longestgap", 0, "number of missing pixels allowed to still be considered continuous")
	fuzzyThresh = flag.Float64("fuzzyThresh", 0.85, "match if rgb values are this similar")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [flags] image\n", os.Args[0])
	flag.PrintDefaults()
}

type point struct {
	x int
	y float64
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("tracetracer: ")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}

	color, err := parseHexColor(*searchColor)
	if err != nil {
		log.Fatalf("invalid color %q: %v", *searchColor, err)
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("failed to open %s: %v", flag.Arg(0), err)
	}

	plot, _, err := image.Decode(f)
	f.Close()
	if err != nil {
		log.Fatalf("failed to read image format: %v", err)
	}

	var vals []point

	{
		xmin := plot.Bounds().Min.X
		xmax := plot.Bounds().Max.X

		for x := xmin; x < xmax; x++ {
			if y, ok := meanOfFirstContinuousRun(plot, color, x, *fuzzyThresh); ok {
				vals = append(vals, point{x: x, y: y})
			}
		}

		if len(vals) == 0 {
			log.Fatal("did not find any values")
		}

		if *onlyLongest {
			vals = longestRun(vals, *longestGap)
		}
	}

	ymin := math.Inf(1)
	ymax := math.Inf(-1)
	xmin := math.Inf(1)
	xmax := math.Inf(-1)

	for _, p := range vals {
		if p.y < ymin {
			ymin = p.y
		}
		if p.y > ymax {
			ymax = p.y
		}
		if float64(p.x) < xmin {
			xmin = float64(p.x)
		}
		if float64(p.x) > xmax {
			xmax = float64(p.x)
		}
	}

	xnorm := (*maxIndex - *minIndex) / (xmax - xmin)
	ynorm := (*maxVal - *minVal) / (ymax - ymin)

	for _, v := range [...]float64{ymin, ymax, xmin, xmax, xnorm, ynorm} {
		if (math.IsInf(v, 0) || math.IsNaN(v)) && len(vals) > 0 {
			log.Fatal("logic error: got inf bounds or norm terms")
		}
	}

	for _, p := range vals {
		// scale
		sx := (float64(p.x)-xmin)*xnorm + *minIndex
		sy := (p.y-ymin)*ynorm + *minVal
		fmt.Printf("%f,%f\n", sx, sy)
	}
}

func longestRun(vs []point, maxGap int) []point {
	var (
		start int = -1
		last  int
		runs  [][2]int
	)
	for i, p := range vs {
		if start == -1 {
			start = i
			last = p.x
		} else if p.x > last+(1+maxGap) {
			runs = append(runs, [2]int{start, i})
			start = i
			last = p.x
		} else {
			last = p.x
		}
	}
	if start != -1 {
		runs = append(runs, [2]int{start, len(vs)})
	}
	if len(runs) == 0 {
		return nil
	}
	var (
		bigLen int
		bigIdx int
	)
	for i, r := range runs {
		if r[1]-r[0] > bigLen {
			bigLen = r[1] - r[0]
			bigIdx = i
		}
	}
	return vs[runs[bigIdx][0]:runs[bigIdx][1]]
}

func meanOfFirstContinuousRun(p image.Image, c color.Color, x int, matchThresh float64) (float64, bool) {
	ymin := p.Bounds().Min.Y
	ymax := p.Bounds().Max.Y

	var runVals []int64
	foundRun := false
	for y := ymin; y < ymax; y++ {
		if fuzzySame(c, p.At(x, y), matchThresh) {
			foundRun = true
			runVals = append(runVals, int64(ymax-y))
		} else if foundRun {
			break
		}
	}

	var sum int64
	for _, v := range runVals {
		sum += v
	}
	if len(runVals) == 0 {
		return 0, false
	}
	return float64(sum) / float64(len(runVals)), true
}

func fuzzySame(c1, c2 color.Color, thresh float64) bool {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()
	fuzzyDiff := (1 - thresh) * 0xffff
	for _, p := range [3][2]uint32{{r1, r2}, {g1, g2}, {b1, b2}} {
		p1 := float64(p[0])
		p2 := float64(p[1])
		if math.Abs(p1-p2) > fuzzyDiff {
			return false
		}
	}
	return true
}

func parseHexColor(hex string) (color.Color, error) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) != 6 {
		return nil, errors.New("must contain 3 or 6 hex digits")
	}
	var rgb [3]byte
	for i := range rgb {
		vbig, ok := fromHexChar(hex[2*i])
		if !ok {
			return nil, errors.New("found invalid digits")
		}
		vsmall, ok := fromHexChar(hex[2*i+1])
		if !ok {
			return nil, errors.New("found invalid digits")
		}
		rgb[i] = vbig*16 + vsmall
	}
	c := color.NRGBA{
		R: uint8(rgb[0]),
		G: uint8(rgb[1]),
		B: uint8(rgb[2]),
		A: 255,
	}
	return c, nil
}

func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}
