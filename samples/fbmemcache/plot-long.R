#!/usr/bin/env Rscript

library(scales)
library(ggplot2)

data <- read.csv(file="long-facebook-trace.csv", header=FALSE)
colnames(data) <- c("Time", "QPS")
data$Time <- data$Time / (24 * 3600)
pdf("long-facebook-trace.pdf", height=2, width=14)
ggplot(data, aes(x=Time, y=QPS)) +
	geom_line() +
	xlab("Time (days)") + 
	ylab("Requests / sec") +
	scale_x_continuous(breaks=pretty_breaks(n=12)) +
	scale_y_continuous(limits=c(0, NA))
.junk <- dev.off()
