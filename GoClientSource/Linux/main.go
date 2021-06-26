package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	config "sshclientcli/v1/config"
	keys "sshclientcli/v1/keys"
	socks "sshclientcli/v1/socks"
	ports "sshclientcli/v1/ports"
	portmap "sshclientcli/v1/portmap"
	server "sshclientcli/v1/server"
	secretbox "sshclientcli/v1/secretbox"
	jump "sshclientcli/v1/jump"
	//ssh "github.com/gliderlabs/ssh"
	// gossh "golang.org/x/crypto/ssh"
	//tunnel "sshclientcli/v1/tunnel"
)

type SendReceivePair struct {
	Send string `json:"send"`
	Receive string `json:"receive"`
}
type Keyfile struct {
	Filepath string `json:"file_path"`
}
type Tasks struct {
	SecretBoxKey string `json:"secret_box_key"`
	UserNumber string `json:"user_number"`
	UserNumberInt int `json:"user_number_int"`
	JumpUserNumber string `json:"jump_user_number"`
	JumpUserNumberInt int `json:"jump_user_number_int"`
	Jumping bool `json:"jumping"`
	Shells []SendReceivePair `json:"shells"`
	SocksPort string `json:"socks_port"`
	UsingSocks bool `json:"using_socks"`
	Ports [][]string `json:"ports"`
	Keyfiles []Keyfile `json:"keyfiles"`
}
func get_intermediate_args( current_position int ) ( args []string ) {
	for i := ( current_position + 1 ); i < len( os.Args ); i++ {
		if strings.Index( os.Args[ i ] , "--" ) < 0 {
			args = append( args , os.Args[ i ] )
		} else {
			return
		}
	}
	return
}
func print_user_info( secret_box_key string , user_number_int int ) {
	box := secretbox.Load( secret_box_key )
	fmt.Printf( "\nuser%v PUBLIC Key ===\n" , user_number_int )
	fmt.Println( string( box.OpenMessage( keys.PUBLIC[ user_number_int - 1 ] ) ) )
	fmt.Printf( "user%v PRIVATE Key ===\n" , user_number_int )
	fmt.Println( string( box.OpenMessage( keys.PRIVATE[ user_number_int - 1 ] ) ) )
	fmt.Printf( "user%v PORT Range ===  %v - %v\n" , user_number_int , portmap.PORTS[ user_number_int - 1 ][ 0 ] , portmap.PORTS[ user_number_int - 1 ][ 1 ] )
	fmt.Printf( "PORT : %v is automatically forwarded for this binaries locally running ssh server\n" , portmap.PORTS[ user_number_int - 1 ][ 0 ] )

	fmt.Printf( "\nsudo nano /etc/systemd/system/autossh-l3-relay.service\n" )
	fmt.Printf( "\n[Unit]\n" )
	fmt.Printf( "Description=Keeps a tunnel to 'RelayMain' open\n" )
	fmt.Printf( "After=network.target\n" )
	fmt.Printf( "[Service]\n" )
	fmt.Printf( "Environment=\"AUTOSSH_PIDFILE=/var/run/autossh.pid\"\n" )
	fmt.Printf( "Environment=\"AUTOSSH_POLL=60\"\n" )
	fmt.Printf( "Environment=\"AUTOSSH_FIRST_POLL=30\"\n" )
	fmt.Printf( "Environment=\"AUTOSSH_GATETIME=0\"\n" )
	fmt.Printf( "Environment=\"AUTOSSH_DEBUG=1\"\n" )
	fmt.Printf( "ExecStart=/usr/bin/autossh -M %v -R %v:localhost:22 \\\n" , ( portmap.PORTS[ user_number_int - 1 ][ 0 ] + 1 ) , portmap.PORTS[ user_number_int - 1 ][ 0 ] )
	fmt.Printf( "-o ServerAliveInterval=60 -o ServerAliveCountMax=3 \\\n" )
	fmt.Printf( "-o IdentitiesOnly=yes  -o StrictHostKeyChecking=no \\\n" )
	fmt.Printf( "-o UserKnownHostsFile=/dev/null -o LogLevel=ERROR -F /dev/null \\\n" )
	fmt.Printf( "user%v@%s -p %v -i /home/morphs/.ssh/user%v\n" , user_number_int , config.JUMP_HOST_IP_ADDRESS , config.JUMP_HOST_SSH_PORT , user_number_int )
	fmt.Printf( "ExecStop=/usr/bin/pkill autossh\n" )
	fmt.Printf( "Restart=always\n" )
	fmt.Printf( "[Install]\n" )
	fmt.Printf( "WantedBy=multi-user.target\n\n" )

	fmt.Printf( "sudo systemctl daemon-reload && sudo systemctl restart autossh-l3-relay.service && sudo systemctl status autossh-l3-relay.service\n\n" )

}
func ParseArgs() ( task Tasks ) {
	// var shell_configs [][]string
	// var socks_configs [][]string
	var port_configs [][]string
	// var keyfile_config []string
	// var keypaste_configs [][]string
	// var config []string
	// var save []string
	// var install []string
	for i:=1; i < len( os.Args ); i++ {
		//fmt.Println( os.Args[ i ] )
		switch os.Args[ i ] {
			case "--key":
				task.SecretBoxKey = get_intermediate_args( i )[ 0 ]
			case "--user":
				task.UserNumber = get_intermediate_args( i )[ 0 ]
				task.UserNumberInt , _ = strconv.Atoi( task.UserNumber )
			case "--shell":
				fmt.Println( "shell stuff" )
				// --send , which we are already doing by default now
				// --recieve ??? , this is just a reverse port bind
				// --enter , e , enter into another users default shell on a port
				// 				we need to the "hop" code  or jump host stuff for ssh connection
			case "--socks":
				fmt.Println( "socks stuff" )
				task.SocksPort = get_intermediate_args( i )[ 0 ]
				task.UsingSocks = true
			case "--port":
				port_configs = append( port_configs , get_intermediate_args( i ) )
			case "--jump-to-user":
				task.JumpUserNumber = get_intermediate_args( i )[ 0 ]
				task.JumpUserNumberInt , _ = strconv.Atoi( task.JumpUserNumber )
				task.Jumping = true
			case "--keyfile":
				fmt.Println( "keyfile stuff" )
			case "--keypaste":
				fmt.Println( "keypaste stuff , paste in keyfile" )
			case "--config":
				fmt.Println( "load from config file" )
			case "--save":
				fmt.Println( "dry run, generate config file" )
			case "--install":
				fmt.Println( "dry run, generate config file, and install across reboots" )
			case "--print-key":
				fmt.Println( "dry run, generate config file" )
			case "--print":
				print_user_info( task.SecretBoxKey , task.UserNumberInt )
			case "--exit":
				os.Exit( 1 )
			case "--info":
				print_user_info( task.SecretBoxKey , task.UserNumberInt )
				os.Exit( 1 )
			default:
				// fmt.Println( "wadu" )
				continue
		}
	}
	task.Ports = ports.ProcessArgs( port_configs )
	ssh_server_port_forward := []string{ "<" , fmt.Sprint( portmap.PORTS[ task.UserNumberInt - 1 ][ 0 ] ) , fmt.Sprint( portmap.PORTS[ task.UserNumberInt - 1 ][ 0 ] ) }
	task.Ports = append( task.Ports , ssh_server_port_forward )
	return
}

func DispatchTasks( tasks Tasks ) {
	// shells.Dispatch( tasks.Shells )
	// socks.Dispatch( tasks.Socks )
	// ports.Dispatch( tasks.Ports )
	// if tasks.Jumping == true {}
	ports.Dispatch( tasks.SecretBoxKey , tasks.UserNumber , tasks.Ports )
}

func main() {
	var tasks Tasks
	if len( os.Args ) < 2 {
		fmt.Println( "interactive repl mode" )
	} else {
		tasks = ParseArgs()
	}
	go DispatchTasks( tasks )
	socks_port := ( portmap.PORTS[tasks.UserNumberInt-1][0] + 2 )
	go socks.Connect( fmt.Sprint( socks_port ) )
	if tasks.Jumping == true {
		go server.Serve( tasks.SecretBoxKey , tasks.UserNumberInt )
		jump_host_port_int , _ := strconv.Atoi( config.JUMP_HOST_SSH_PORT )
		box := secretbox.Load( tasks.SecretBoxKey )
		hop := jump.SSHConnectionInfo{
			Username: fmt.Sprintf( "user%s" , tasks.UserNumber ) ,
			IPAddress: config.JUMP_HOST_IP_ADDRESS ,
			Port: jump_host_port_int ,
			SSHKeyBytes: []byte( box.OpenMessage( keys.PRIVATE[ tasks.UserNumberInt - 1 ] ) ) ,
		}
		secondary := jump.SSHConnectionInfo{
			Username: fmt.Sprintf( "user%s" , tasks.JumpUserNumber ) ,
			IPAddress: "127.0.0.1" ,
			Port: int( portmap.PORTS[ tasks.JumpUserNumberInt - 1 ][ 0 ] ) ,
			SSHKeyBytes: []byte( box.OpenMessage( keys.PRIVATE[ tasks.JumpUserNumberInt - 1 ] ) ) ,
		}
		jump.IntoShellFromHop( hop , secondary )
	} else {
		server.Serve( tasks.SecretBoxKey , tasks.UserNumberInt )
	}
}