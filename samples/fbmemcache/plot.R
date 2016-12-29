#!/usr/bin/env Rscript

library(ggplot2)

data <- read.csv(file="24hr-facebook-trace.csv", header=FALSE)
colnames(data) <- c("Time", "QPS")
pdf("24hr-facebook-trace.pdf", height=2, width=8)
ggplot(data, aes(x=Time, y=QPS)) +
	geom_line() +
	xlab("Time (s)") + 
	ylab("Requests / sec")
.junk <- dev.off()
