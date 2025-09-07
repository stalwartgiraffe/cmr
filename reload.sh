#!/bin/bash
#trap "exit;" SIGINT;
while true; do
  make || break
  time build/cmr mvc reload || break
done

