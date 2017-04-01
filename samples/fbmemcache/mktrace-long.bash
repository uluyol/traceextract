#!/usr/bin/env bash
../../tracetracer \
	-color "#ff7d00" \
	-fuzzyThresh 0.5 \
	-xmin 0 \
	-xmax 950400 \
	-ymin 35000 \
	-ymax 88000 \
	-onlylongest=true \
	-longestgap 10 \
	long-facebook-trace.png \
	> long-facebook-trace.csv
