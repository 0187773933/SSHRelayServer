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
	SealedPrivateKeyData=$(../SecretBoxBinaries/linux/arm/secretbox --key "$SecretBoxKey" seal message "$PrivateKeyData")
	printf "\t[]byte(\`$SealedPrivateKeyData\`) ,\n" >> keys.go
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