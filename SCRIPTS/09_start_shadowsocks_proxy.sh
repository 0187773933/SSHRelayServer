#!/bin/bash
# Start Shadowsocks Proxy
# https://github.com/hangim/kcp-shadowsocks-docker
# https://github.com/hangim/shadowsocks-docker/blob/master/Dockerfile
# https://github.com/shadowsocks/shadowsocks-libev
# sudo ufw allow 9444/udp || echo "failed to allow 9444 in UFW firewall" && \
# sudo docker rm public-shadowsocks -f || echo "failed to remove existing shadowsocks server" && \
# sudo docker run -dit --restart=always \
# --name=public-shadowsocks \
# --net l3-relay \
# -e "SS_PORT=9444" -e "SS_PASSWORD=3a1871f34ac8b1f7aabd112c66b4005a" \
# -e "SS_METHOD=chacha20-ietf-poly1305" -e "SS_TIMEOUT=600" \
# -e "KCP_PORT=9443" -e "KCP_MODE=fast" -e "MTU=1400" \
# -e "SNDWND=1024" -e "RCVWND=1024" \
# -p 9444:9444 -p 9444:9444/udp \
# -p 9443:9443/udp \
# imhang/kcp-shadowsocks-docker