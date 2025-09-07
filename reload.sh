#!/bin/bash
#trap "exit;" SIGINT;
while true; do
  make && time build/cmr mvc reload
done

