#!/bin/bash

SecretBoxKey=$(../SecretBoxBinaries/linux/arm/secretbox new)
echo $SecretBoxKey > ../SecretBoxKey.txt

printf "package keys\n\n" > keys.go

printf "var PUBLIC = []string{\n" >> keys.go
for number in {1..254}; do
	PublicKeyData=$(<../BUILD/KEYS/user${number}.pub)
	SealedPublicKeyData=$(../SecretBoxBinaries/linux/arm/secretbox --key "$SecretBoxKey" seal message "$PublicKeyData")
	printf '\t"%s" ,\n' $SealedPublicKeyData >> keys.go
done
printf "}\n\n" >> keys.go

printf "var PRIVATE = []string{\n" >> keys.go
for number in {1..254}; do
	PrivateKeyData=$(<../BUILD/KEYS/user${number})
	../SecretBoxBinaries/linux/arm/secretbox --key "$SecretBoxKey" seal file "../BUILD/KEYS/user${number}" "../BUILD/KEYS/user${number}.sealed"
	SealedPrivateKeyData=$(<../BUILD/KEYS/user${number}.sealed)
	printf '\t"%s" ,\n' SealedPrivateKeyData >> keys.go
done
printf "}\n" >> keys.go

cp keys.go ../GoClientSource/Windows/v1/keys/
mv keys.go ../GoClientSource/Linux/v1/keys/