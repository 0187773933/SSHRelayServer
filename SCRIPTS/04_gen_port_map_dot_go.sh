#!/bin/bash

# Hard Coded Ports
# Hardcoded Service Routes, People Can Port Forward To
# 0 - 1023 = Kernel
# 1024 - 65535 = Userspace
# Service Routes = 30090 - 30099

# Each User Gets 30 Ports To Start With
#          = 32000
# 30 * 254 = 07620
#          = 39619

# These are so the Go-Binary-Clients Can Have Easier Time
#echo "254 Users - 30 Ports Each"
printf "package portmap\n\n" > portmap.go
printf "var PORTS = [][2]uint32{ \n" >> portmap.go
for number in {1..254}; do
    Min=$(( 32000 + (30*(number-1)) ))
    Max=$(( Min + 29 ))
    printf "\t{ $Min , $Max } , \n" >> portmap.go
done
printf "}" >> portmap.go

cp portmap.go ../GoClientSource/Windows/v1/portmap/
mv portmap.go ../GoClientSource/Linux/v1/portmap/

# for number in {1..254}; do
#     Min=$(( 32000 + (30*(number-1)) ))
#     Max=$(( Min + 29 ))
#     if [ $number -eq 254 ]; then
#         printf "user$number === $Min === $Max" >> PortMap.txt
#     else
#         printf "user$number === $Min === $Max\n" >> PortMap.txt
#     fi
# done
# exit