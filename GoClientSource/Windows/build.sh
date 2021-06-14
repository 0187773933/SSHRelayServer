#!/bin/bash
# BinaryName="$1"
BinaryName="l3r"
rm -rf ./bin/

# https://stackoverflow.com/questions/25051623/golang-compile-for-all-platforms-in-windows-7-32-bit
# https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04

declare -a windows_architectures=(
	"amd64"
)

for architecture in "${windows_architectures[@]}"
do
	echo "Building Windows: $architecture"
	GOOS=windows GOARCH=$architecture go build -o bin/windows/$architecture/$BinaryName.exe
done