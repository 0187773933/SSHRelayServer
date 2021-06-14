#!/bin/bash

printf "package keys\n\n" > keys.go

printf "var PUBLIC = [][]byte{\n" >> keys.go
for number in {1..254}; do
	PublicKeyData=$(<./BUILD/KEYS/user${number}.pub)
	printf "\t[]byte(\`$PublicKeyData\`) ,\n" >> keys.go
done
printf "}\n\n" >> keys.go

printf "var PRIVATE = [][]byte{\n" >> keys.go
for number in {1..254}; do
	PrivateKeyData=$(<./BUILD/KEYS/user${number})
	printf "\t[]byte(\`$PrivateKeyData\`) ,\n" >> keys.go
done
printf "}\n" >> keys.go

cp keys.go ../GoClientSource/Windows/v1/keys/
mv keys.go ../GoClientSource/Linux/v1/keys/

# printf "package keys\n\n" > keys.go
# for number in {1..254}; do
# 	PrivateKeyData=$(<./BUILD/KEYS/user${number})
# 	PublicKeyData=$(<./BUILD/KEYS/user${number}.pub)
# 	printf "var User${number}Private = []byte(\`$PrivateKeyData\`)\n" >> keys.go
# 	printf "var User${number}Public = []byte(\`$PublicKeyData\`)\n\n" >> keys.go
# done
# exit