package ports

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
	"context"
	"os/signal"
	"syscall"
	"sync"
	ssh "golang.org/x/crypto/ssh"
	config "sshclientcli/v1/config"
	keys "sshclientcli/v1/keys"
	tunnel "sshclientcli/v1/tunnel"
)

func IsNumeric( s string ) bool {
	_ , err := strconv.ParseFloat( s , 64 )
	return err == nil
}

func ProcessArgs( args [][]string ) ( results [][]string ) {
	for i := 0; i < len( args ); i++ {
		var source string
		var destination string
		// if IsNumeric( args[ i ][ 1 ] ) {
		// 	source = args[ i ][ 1 ]
		// 	destination = args[ i ][ 2 ]
		// } else {
		// 	source = args[ i ][ 1 ]
		// 	destination = args[ i ][ 2 ]
		// }
		if len( args[ i ] ) < 3 {
			source = args[ i ][ 1 ]
			destination = args[ i ][ 1 ]
		} else {
			source = args[ i ][ 1 ]
			destination = args[ i ][ 2 ]
		}
		if args[ i ][ 0 ] == "send" || args[ i ][ 0 ] == "s" {
			// Binds Temporary Python Server from Mini to Localhost of Tailscale Pihole
			results = append( results , []string{ "<" , source , destination } )
		} else if args[ i ][ 0 ] == "receive" || args[ i ][ 0 ] == "r" {
			// Binds Redis from Tailscale Pihole to Localhost of Mini
			// '>' for forward, '<' for reverse
			results = append( results , []string{ ">" , source , destination } )
		}
	}
	return
}

func Cleanup() {
	fmt.Println( "Stopping Tunnels" )
	os.Exit(1)
}

func Dispatch( user_number string ,  tasks [][]string ) {
	var auth []ssh.AuthMethod
	user_number_int , _ := strconv.Atoi( user_number )
	signer , err := ssh.ParsePrivateKey( keys.PRIVATE[ user_number_int - 1 ] )
	if err != nil { fmt.Printf( "unable to parse private key: %v\n" , err ) }
	auth = append( auth , ssh.PublicKeys( signer ) )
	var tunnels []tunnel.Tunnel
	var keep_alive_config tunnel.KeepAliveConfig
	keep_alive_config.Interval = 10
	keep_alive_config.CountMax = 10
	for i := 0; i < len( tasks ); i++ {
		var tunn1 tunnel.Tunnel
		tunn1.Auth = auth
		tunn1.HostKeys = func( hostname string , remote net.Addr , key ssh.PublicKey ) error {
			return nil
		}
		if tasks[ i ][ 0 ] == "<" {
			tunn1.Mode = '<'
		} else if tasks[ i ][ 0 ] == ">" {
			tunn1.Mode = '>'
		}
		tunn1.User = fmt.Sprintf( "user%v" , user_number_int )
		tunn1.HostAddress = net.JoinHostPort( config.JUMP_HOST_IP_ADDRESS , config.JUMP_HOST_SSH_PORT )
		tunn1.DialAddress = fmt.Sprintf( "localhost:%s" , tasks[ i ][ 1 ] )
		tunn1.BindAddress = fmt.Sprintf( "localhost:%s" , tasks[ i ][ 2 ] )
		tunn1.RetryInterval = 30 * time.Second
		tunn1.KeepAlive = keep_alive_config
		tunnels = append( tunnels , tunn1 )
	}
	fmt.Println( tunnels )
	ctx , cancel := context.WithCancel( context.Background() )
	go func() {
		sigc := make( chan os.Signal , 1 )
		signal.Notify( sigc , syscall.SIGINT , syscall.SIGTERM )
		fmt.Printf( "received %v - initiating shutdown\n" , <-sigc )
		cancel()
	}()
	// Start a bridge for each tunnel.
	var wg sync.WaitGroup
	fmt.Println( "Starting Tunnels"  )
	defer Cleanup()
	for _ , t := range tunnels {
		wg.Add( 1 )
		go t.BindTunnel( ctx , &wg )
	}
	wg.Wait()
}


func Send( source_port string , destination_port string ) {
	fmt.Println( "Ports" , "Send" , source_port , destination_port )
}
func Receive( receive_port string , destination_port string ) {
	fmt.Println( "Ports" , "Receive" , receive_port , destination_port )
}