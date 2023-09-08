#!/bin/bash
#NOTE: This script requires the "parallel" command line utility
# The script was taken from this very useful website: https://www.simonholywell.com/post/2015/06/parallel-benchmark-many-urls-with-apachebench/
# It can be installed with Homebrew:
# brew install parallel

# The expectation is that the count should be 0 after this call
(echo "127.0.0.1:9095/increment"; echo "127.0.0.1:9095/decrement") | parallel 'ab -n 200 -c 10 {}'