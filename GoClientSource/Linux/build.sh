#!/bin/bash
# BinaryName="$1"
BinaryName="l3r"
rm -rf ./bin/

# https://stackoverflow.com/questions/25051623/golang-compile-for-all-platforms-in-windows-7-32-bit
# https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04

declare -a linux_architectures=(
	"amd64"
	"arm"
	"arm64"
)
declare -a darwin_architectures=(
	"amd64"
)

# mkdir -p ./bin/plugins
# GOOS=windows go build -buildmode=plugin -o ./v1/server/plugins/windows.go ./bin/plugins/windows.so
# GOOS=linux go build -buildmode=plugin -o ./v1/server/plugins/linux.go ./bin/plugins/linux.so
# GOOS=darwin go build -buildmode=plugin -o ./v1/server/plugins/darwin.go ./bin/plugins/darwin.so

for architecture in "${linux_architectures[@]}"
do
	echo "Building Linux: $architecture"
	GOOS=linux GOARCH=$architecture go build -o bin/linux/$architecture/$BinaryName
done

for architecture in "${darwin_architectures[@]}"
do
	echo "Building Darwin: $architecture"
	GOOS=darwin GOARCH=$architecture go build -o bin/darwin/$architecture/$BinaryName
done