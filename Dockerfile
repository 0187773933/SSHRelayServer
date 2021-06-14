FROM debian:stable-slim
RUN apt-get update -y
RUN apt-get install gcc -y
RUN apt-get install nano -y
RUN apt-get install tar -y
RUN apt-get install bash -y
RUN apt-get install sudo -y
RUN apt-get install openssl -y
RUN apt-get install git -y
RUN apt-get install make -y
RUN apt-get install wget -y
RUN apt-get install curl -y
RUN apt-get install openssh-server -y
RUN apt-get install openssh-client -y
RUN apt-get install python3 -y
RUN apt-get install net-tools -y
RUN apt-get install iproute2 -y
RUN apt-get install iputils-ping -y

RUN printf "Port 10092\nAddressFamily inet\nPermitRootLogin no\nPasswordAuthentication no\nChallengeResponseAuthentication no\nUsePAM yes\nX11Forwarding no\nPrintMotd no\nBanner none\nAcceptEnv LANG LC_*\nSubsystem   sftp    /usr/lib/openssh/sftp-server\nClientAliveInterval 3\nGatewayPorts yes\nPubkeyAuthentication yes\n" > /etc/ssh/sshd_config
RUN /etc/init.d/ssh restart
RUN groupadd -g 1274 sshrelay
ADD ./BUILD/KEYS /KEYS
RUN chgrp -R sshrelay /KEYS
RUN chmod -R g+r  /KEYS
RUN ls -lah /KEYS
RUN mkdir -p /etc/ssh/keys
WORKDIR /etc/ssh/keys
RUN ssh-keygen -A
WORKDIR /
COPY ./SCRIPTS/05_gen_users_debian.sh /
RUN chmod +x 05_gen_users_debian.sh.sh
RUN ./05_gen_users_debian.sh.sh
ENTRYPOINT service ssh restart && bash