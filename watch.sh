#!/bin/bash

#specify commands that should be executed automatically when the shell receives certain signals
trap "echo 'Received SIGINT or SIGTERM. Exiting...'; exit" SIGINT SIGTERM

#entr => run arbitrary commands when files change
while true; do
	find srcs/ -name "*.go" | entr -rd make run;
done
