#!/bin/bash
APP_NAME="public-ssh-server"
RELAY_NAME="l3-relay"

# 1.) Create Shared Docker Network
sudo docker network rm $RELAY_NAME || echo "failed to delete network"
sudo docker network create \
--driver=bridge \
--subnet=10.4.4.0/24 \
--gateway=10.4.4.1 \
$RELAY_NAME || echo "failed to create network"
#--ip-range=10.4.4.0/24 \

# 2.) Build Debian SSH Server
sudo docker rm $APP_NAME -f || echo "failed to remove existing ssh server"
sudo docker build -t $APP_NAME .

# 5.) Start Debian SSH Server
# Note: we are mapping ports 30090-30099 to be used as "system/service" ports
# This is how we get smartphones, and everything else inside.
# you can try to add it , but they just always fail somehow
# you just bind to one of these service ports in the future, and then they will be accessable via
# s1.34353.org { reverse_proxy localhost:30090 } in the /etc/caddy/Caddyfile
# we are binding/mapping them here with this docker container, because why not. there isn't a better place yet to put them
# add more as necessary
id=$(sudo docker run -dit --restart='always' \
--name  \
--net $RELAY_NAME \
-p 10092:10092 \
-p 30090-30099:30090-30099 \
$APP_NAME)
echo "ID = $id"
# sudo docker exec -it $id /bin/bash