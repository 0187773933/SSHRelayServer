#!/bin/bash
printf "package config\n\n" > config.go
printf "var USER_ID = 99\n" >> config.go
printf "var JUMP_HOST_IP_ADDRESS = \"45.77.156.253\"\n" >> config.go
printf "var JUMP_HOST_SSH_PORT = "10092"\n" >> config.go

cp config.go ../GoClientSource/Windows/v1/config/
mv config.go ../GoClientSource/Linux/v1/config/

# for number in {1..254}; do
# 	PublicKeyData=$(<./BUILD/KEYS/user${number}.pub)
# 	printf "\t[]byte(\`$PublicKeyData\`) ,\n" >> keys.go
# done
#printf "}\n\n" >> keys.go