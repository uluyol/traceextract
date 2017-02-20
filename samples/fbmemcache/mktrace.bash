#!/usr/bin/env bash
../../tracetracer \
	-color "#ff7d00" \
	-fuzzyThresh 0.5 \
	-xmin 0 \
	-xmax 86400 \
	-ymin 40000 \
	-ymax 77000 \
	-onlylongest=true \
	-longestgap 10 \
	24hr-facebook-trace.png \
	> 24hr-facebook-trace.csv
