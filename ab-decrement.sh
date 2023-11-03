#!/bin/bash
ab -n 1000 -c 100 "127.0.0.1:9095/decrement"