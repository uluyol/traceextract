#!/usr/bin/env Rscript

library(scales)
library(ggplot2)

data <- read.csv(file="72hr-facebook-trace.csv", header=FALSE)
colnames(data) <- c("Time", "QPS")
data$Time <- data$Time / 3600
pdf("72hr-facebook-trace.pdf", height=2, width=8)
ggplot(data, aes(x=Time, y=QPS)) +
	geom_line() +
	xlab("Time (hr)") + 
	ylab("Requests / sec") +
	scale_x_continuous(breaks=pretty_breaks(n=12)) +
	scale_y_continuous(limits=c(0, NA))
.junk <- dev.off()
