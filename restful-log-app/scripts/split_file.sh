#!/bin/bash

# Create chunks directory if it does not exist
mkdir -p ./chunks

# Split the log file into smaller chunks
split -l 10000 $1 ./chunks/chunk_
