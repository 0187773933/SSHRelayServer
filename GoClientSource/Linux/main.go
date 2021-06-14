package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	// keys "sshclientcli/v1/keys"
	//socks "sshclientcli/v1/socks"
	ports "sshclientcli/v1/ports"
	portmap "sshclientcli/v1/portmap"
	server "sshclientcli/v1/server"
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
	UserNumber string `json:"user_number"`
	UserNumberInt int `json:"user_number_int"`
	Shells []SendReceivePair `json:"shells"`
	Socks []SendReceivePair `json:"socks"`
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
			case "--port":
				port_configs = append( port_configs , get_intermediate_args( i ) )
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
	ports.Dispatch( tasks.UserNumber , tasks.Ports )
}

func main() {
	var tasks Tasks
	if len( os.Args ) < 2 {
		fmt.Println( "interactive repl mode" )
	} else {
		tasks = ParseArgs()
	}
	go DispatchTasks( tasks )
	server.Serve( tasks.UserNumberInt )
}package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	// keys "sshclientcli/v1/keys"
	//socks "sshclientcli/v1/socks"
	ports "sshclientcli/v1/ports"
	portmap "sshclientcli/v1/portmap"
	server "sshclientcli/v1/server"
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
	UserNumber string `json:"user_number"`
	UserNumberInt int `json:"user_number_int"`
	Shells []SendReceivePair `json:"shells"`
	Socks []SendReceivePair `json:"socks"`
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
			case "--port":
				port_configs = append( port_configs , get_intermediate_args( i ) )
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
	ports.Dispatch( tasks.UserNumber , tasks.Ports )
}

func main() {
	var tasks Tasks
	if len( os.Args ) < 2 {
		fmt.Println( "interactive repl mode" )
	} else {
		tasks = ParseArgs()
	}
	go DispatchTasks( tasks )
	server.Serve( tasks.UserNumberInt )
}