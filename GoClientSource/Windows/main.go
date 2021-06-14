package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"time"
	user "os/user"
	"io/ioutil"
	"encoding/json"
	keys "sshclientcli/v1/keys"
	//socks "sshclientcli/v1/socks"
	ports "sshclientcli/v1/ports"
	portmap "sshclientcli/v1/portmap"
	server "sshclientcli/v1/server"
	//ssh "github.com/gliderlabs/ssh"
	// gossh "golang.org/x/crypto/ssh"
	//tunnel "sshclientcli/v1/tunnel"
	// "log"
	// service "github.com/kardianos/service"
	robustly "github.com/VividCortex/robustly"
	// try "github.com/manucorporat/try"
)

var tasks Tasks
var config_file_base_path string
var config_file_path string

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

func ReadConfig() ( read_tasks Tasks ) {
	data , err := ioutil.ReadFile( config_file_path )
	if err != nil { fmt.Print( err ) }
	unmarshal_error := json.Unmarshal( data , &read_tasks )
	if unmarshal_error != nil { fmt.Println( unmarshal_error ) }
	return
}

func WriteConfig() {
	json_data , json_data_error := json.Marshal( tasks )
	if json_data_error != nil { fmt.Println( json_data_error ) }
	ioutil.WriteFile( config_file_path , json_data , os.ModePerm )
	// file , file_error := os.OpenFile( config_file_path , os.O_CREATE , os.ModePerm )
	// if file_error != nil { fmt.Println( file_error ) }
	// defer file.Close()
	// encoder := json.NewEncoder( file )
	// encoder.Encode( tasks )
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
func print_user_info( user_number_int int ) {
	fmt.Printf( "user%v PUBLIC Key ===\n" , user_number_int )
	fmt.Println( string( keys.PUBLIC[ user_number_int - 1 ] ) )
	fmt.Printf( "user%v PRIVATE Key ===\n" , user_number_int )
	fmt.Println( string( keys.PRIVATE[ user_number_int - 1 ] ) )
	fmt.Printf( "user%v PORT Range ===  %v - %v\n" , user_number_int , portmap.PORTS[ user_number_int - 1 ][ 0 ] , portmap.PORTS[ user_number_int - 1 ][ 1 ] )
	fmt.Printf( "PORT : %v is automatically forwarded for this binaries locally running ssh server\n" , portmap.PORTS[ user_number_int - 1 ][ 0 ] )
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
			case "--print":
				print_user_info( task.UserNumberInt )
			case "--exit":
				os.Exit( 1 )
			case "--info":
				print_user_info( task.UserNumberInt )
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
	ports.Dispatch( tasks.UserNumber , tasks.Ports )
}

func RobustlyStartProgram() {
	WriteConfig()
	go DispatchTasks( tasks )
	server.Serve( tasks.UserNumberInt )
}

func StartProgram() {
	robustly.Run( RobustlyStartProgram , &robustly.RunOptions{
		RateLimit:  1.0,
		Timeout:    time.Second ,
		PrintStack: false ,
		RetryDelay: 3 * time.Second ,
	})
}

// var logger service.Logger
// type program struct{}
// func ( p *program ) Start( s service.Service ) error {
// 	fmt.Println( "l3r --> Start()" )
// 	go p.run()
// 	return nil
// }
// func ( p *program ) run() {
// 	fmt.Println( "l3r --> run()" )
// 	StartProgram()
// }
// func ( p *program ) Stop( s service.Service ) error {
// 	// Stop should not block. Return with a few seconds.
// 	fmt.Println( "l3r --> Stop()" )
// 	return nil
// }

func main() {
	user_info , _ := user.Current()
	config_file_base_path = fmt.Sprintf( "%s\\AppData\\Local\\hasdfs\\" , user_info.HomeDir )
	_ = os.Mkdir( config_file_base_path , 0700 )
	config_file_path = config_file_base_path + "config.json"

	if len( os.Args ) < 2 {
		fmt.Println( "interactive repl mode" )
		fmt.Println( "or just we are reading from config file" )
		tasks = ReadConfig()
	} else {
		tasks = ParseArgs()
	}
	StartProgram()
	// https://pkg.go.dev/github.com/kardianos/service#Config


	// svcConfig := &service.Config{
	// 	Name: "l4333r" ,
	// 	DisplayName: "l4333r" ,
	// 	Description: "l4333r" ,
	// 	// Arguments: []string{ config_file_path } ,
	// }
	// prg := &program{}
	// var s service.Service
	// var s_error error
	// try.This( func() {
	// 	s , s_error = service.New( prg , svcConfig )
	// 	if s_error != nil { fmt.Println( s_error ) }
	// }).Catch( func ( e try.E ) {
	// 	fmt.Println( e )
	// })
	// logger , _ = s.Logger( nil )
	// try.This( func() {
	// 	status , status_error := s.Status()
	// 	if status_error != nil { logger.Error( status_error ) }
	// 	fmt.Printf( "Status === %v\n", string( status ) )
	// }).Catch( func ( e try.E ) {
	// 	fmt.Println( e )
	// })
	// try.This( func() {
	// 	fmt.Println( "Stopping Any Existing Previous Service" )
	// 	stop_error := s.Stop()
	// 	if stop_error != nil { logger.Error( stop_error ) }
	// }).Catch( func ( e try.E ) {
	// 	fmt.Println( e )
	// })
	// try.This( func() {
	// 	fmt.Println( "Uninstalling Any Existing Service" )
	// 	uninstall_error := s.Uninstall()
	// 	if uninstall_error != nil { logger.Error( uninstall_error ) }
	// 	time.Sleep( 1 * time.Second )
	// }).Catch( func ( e try.E ) {
	// 	fmt.Println( e )
	// })
	// try.This( func() {
	// 	fmt.Println( "Installing" )
	// 	install_error := s.Install()
	// 	if install_error != nil { logger.Error( install_error ) }
	// }).Catch( func ( e try.E ) {
	// 	fmt.Println( e )
	// })
	// try.This( func() {
	// 	fmt.Println( "Running" )
	// 	running_error := s.Run()
	// 	if running_error != nil { logger.Error( running_error ) }
	// }).Catch( func ( e try.E ) {
	// 	fmt.Println( e )
	// })
}