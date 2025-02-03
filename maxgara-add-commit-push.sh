#!/bin/zsh

while getopts ":m:" opt; do
	msg=$OPTARG
	git add . && git commit -m $msg && git push origin main;
done;

