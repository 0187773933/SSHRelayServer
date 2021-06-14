#!/bin/bash
for number in {1..254}; do
	username=$(echo "user$number")
	echo $username
	useradd --home /home/$username -s '/bin/bash' -m $username
	mkdir -p /home/$username
	usermod -a -G sshrelay $username
	printf "#!/bin/bash\n" > "/home/$username/startup.sh"
	printf "mkdir -p /home/$username/.ssh\n" >> "/home/$username/startup.sh"
	printf "chmod 700 /home/$username/.ssh\n" >> "/home/$username/startup.sh"
	# printf "ssh-keygen -t ed25519 -b 521 -a 100 -f '/home/$username/.ssh/$username' -q -N ''\n" >> "/home/$username/startup.sh"
	printf "cp /KEYS/$username.pub /home/$username/.ssh/authorized_keys\n" >> "/home/$username/startup.sh"
	# printf "mv '/home/$username/.ssh/$username.pub' '/home/$username/.ssh/authorized_keys'\n" >> "/home/$username/startup.sh"
	printf "chmod 600 /home/$username/.ssh/authorized_keys\n" >> "/home/$username/startup.sh"
	chmod +x /home/$username/startup.sh
	su $username -c "/home/$username/startup.sh"
done
exit