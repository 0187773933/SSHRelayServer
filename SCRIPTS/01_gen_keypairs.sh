#!/bin/bash
cd ../
sudo rm -rf ./BUILD
mkdir -p ./BUILD
mkdir -p ./BUILD/KEYS
mkdir -p ./BUILD/BINARIES
cd ./BUILD/KEYS
for number in {1..254}; do
	username=$(echo "user$number")
	echo $username
	ssh-keygen -t ed25519 -b 521 -a 100 -f $username -q -N ''
	chmod 600 ./$username.pub
	chmod 600 ./$username
done
cd ../../
chmod 700 ./BUILD/KEYS
exit