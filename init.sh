#!/bin/bash
cd ./SCRIPTS
./01_gen_keypairs.sh
./02_gen_config_dot_go.sh
./03_gen_keys_dot_go.sh
./04_gen_port_map_dot_go.sh

cd ..
./dockerRun.sh