#!/bin/bash
# https://stackoverflow.com/questions/44492576/looking-for-large-text-files-for-testing-compression-in-all-sizes
tr -dc "A-Za-z 0-9" < /dev/urandom | fold -w100|head -n 100000 > bigfile.txt
