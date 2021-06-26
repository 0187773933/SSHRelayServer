# Todo

1. Generate GoClient Binaries
    - CLI
        - Shell
            - Send (AutoSSH , etc)
            - Receive (jump into shell)
        - SOCKS Proxy
            - Send Local Connection
            - Receive Forwarded Connection
        - Port(s)
            - Send to Cloud MITM VPS
            - Receive from Cloud MITM VPS
<br>
<br>

2. Rewrite Dockerfile to PreGenerate ssh key-pairs , copy them in , and use for each user , deterministic ( if restart )

3. Auto Generate AutoSSH Script / Systemd Service ForEach User (linux)

4. Auto Generate Utility Scripts and Copy/Pastable Things Like Shell Hop Example Where Everything Is Self-Contained

5. Add Option for Importing/Reading From File another Private Key to Use after getting inside jump host for some other ephemeral user/port/key

<br>

# Shell Hop Example
>  - Ubuntu VM is autossh remote port forwarding openssh server port 22 to user98 inside MITM docker vm
>  - asdf
>  - asdf
>  - asdf

```bash
echo "-----BEGIN OPENSSH PRIVATE KEY-----
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==
-----END OPENSSH PRIVATE KEY-----" > /tmp/linuxMisc3RelayMainUser99 && \
chmod 600 /tmp/linuxMisc3RelayMainUser99 && \
echo -e "Host relay-l3
\t  HostName 11.22.33.44
\t  Port 10092
\t  User user99
\t  ServerAliveInterval=15
\t  ServerAliveCountMax=3
\t  IdentityFile /tmp/linuxMisc3RelayMainUser99
Host user98
\t  HostName localhost
\t  Port 10098
\t  User morphs
\t  ServerAliveInterval=15
\t  ServerAliveCountMax=3
\t  IdentityFile /Users/morpheous/Documents/Misc/SSH2/KEYS/p1
\t  ProxyJump relay-l3" > /tmp/relay-l3-ssh-config && \
ssh -F /tmp/relay-l3-ssh-config -J relay-l3 user98 \
-o ServerAliveInterval=15 -o ServerAliveCountMax=3 -o IdentitiesOnly=yes \
-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
-o ExitOnForwardFailure=yes \
-o LogLevel=ERROR && \
rm /tmp/relay-l3-ssh-config

```


```bash
./l3r --key dgxHgUUGgudqlPsGC1tqN0TI9va1hXcYPsWjjb1ydllAQEA9PT1AQEC6V5AueIjr8WyJXTfFf8eaiNGddQUsjh8= \
--user 44 --port send 9003 9003
```


```bash
./l3r --key dgxHgUUGgudqlPsGC1tqN0TI9va1hXcYPsWjjb1ydllAQEA9PT1AQEC6V5AueIjr8WyJXTfFf8eaiNGddQUsjh8= \
--user 45 --port receive 9003 9003
```


```bash
./l3r --key dgxHgUUGgudqlPsGC1tqN0TI9va1hXcYPsWjjb1ydllAQEA9PT1AQEC6V5AueIjr8WyJXTfFf8eaiNGddQUsjh8= \
--user 45 --jump-to-user 44
```


/usr/local/bin/l3r is a script / bash_rc alias alternative that calls the actuall go binary
so that we might just call this from systemd
```bash
#!/bin/bash
/usr/local/bin/l3r_binary --key dgxHgUUGgudqlPsGC1tqN0TI9va1hXcYPsWjjb1ydllAQEA9PT1AQEC6V5AueIjr8WyJXTfFf8eaiNGddQUsjh8= --user 11 "$@"
```