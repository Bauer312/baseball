# Go data pipelines

A data pipeline is a series of pipeline stages that are passing data from one to the next.  The output of one stage is the input to the next.  In this Go-oriented design, each pipeline stage is running in its own Goroutine(s).  Each pipeline stage receives input via one or more channel and sends output to one or more channel.

## General usage

### Starting
* For the first time
* Restarting

### Stopping
* Cancel all processing and stop immediately (abort)
* Stop accepting new work and stop when finished (quiesce)

### Finalize
* Done