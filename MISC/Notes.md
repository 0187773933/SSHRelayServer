
# Hard Coded Ports
> Hardcoded Service Routes, People Can Port Forward To
-  0 - 1023 = Kernel
- 1024 - 65535 = Userspace
- Service Routes = 30090 - 30099

# Each User Gets 30 Ports To Start With
```
          = 32000
 30 * 254 = 07620
          = 39620
```


> Mac OSX --> Parallels VM --> XUbuntu --> python3 -m http.server 8000 --> SSH Reverse Port Forward -R 30090:localhost:8000 , Vultr Cloud VPS , Docker Container : 10092
#30090 on Docker Container is Mapped to 30090 on VPS host, which is reverse_proxied from a caddy server to a subdomain
