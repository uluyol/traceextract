#!/usr/bin/env bash
../../tracetracer \
	-color "#ff4d00" \
	-fuzzyThresh 0.5 \
	-xmin 0 \
	-xmax 979200 \
	-ymin 35000 \
	-ymax 88000 \
	-onlylongest=true \
	-longestgap 20 \
	long-facebook-trace.png \
	> long-facebook-trace.csv
