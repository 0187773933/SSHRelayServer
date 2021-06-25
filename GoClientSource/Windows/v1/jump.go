package jump

import (
	"fmt"
	"log"
	"net"
	"time"
	"io"
	"os"
	"context"
	ssh "golang.org/x/crypto/ssh"
)

type TerminalConfig struct {
	Term   string
	Height int
	Weight int
	Modes  ssh.TerminalModes
}
type RemoteShell struct {
	client *ssh.Client
	requestPty bool
	terminalConfig *TerminalConfig
	stdin io.Reader
	stdout io.Writer
	stderr io.Writer
}
type SSHConnectionInfo struct {
	Username string
	IPAddress string
	Port int
	SSHKeyBytes []byte
}

// newClientConn is a wrapper around ssh.NewClientConn
// https://github.com/gravitational/teleport/blob/5ad1a9025cdd9e5d5fad46a0d316128615efc29b/lib/client/client.go#L807
func new_client_connection( conn net.Conn , nodeAddress string , config *ssh.ClientConfig ) ( ssh.Conn , <-chan ssh.NewChannel , <-chan *ssh.Request , error ) {
	var ctx = context.Background()
	type response struct {
		conn   ssh.Conn
		chanCh <-chan ssh.NewChannel
		reqCh  <-chan *ssh.Request
		err    error
	}
	respCh := make( chan response , 1 )
	go func() {
		conn, chans, reqs, err := ssh.NewClientConn(conn, nodeAddress, config)
		respCh <- response{conn, chans, reqs, err}
	}()
	select {
		case resp := <-respCh:
			if resp.err != nil {
				return nil, nil, nil , nil
			}
			return resp.conn, resp.chanCh, resp.reqCh, nil
		case <-ctx.Done():
			errClose := conn.Close()
			if errClose != nil {
				fmt.Println( errClose )
			}
			// drain the channel
			<-respCh
			return nil, nil, nil , nil
	}
}



func ConnectToHop( hop_connection_info SSHConnectionInfo ) ( ssh_client *ssh.Client ) {
	HOP_SSH_KEY_FILE_PASSWORD := ""
	var auths []ssh.AuthMethod
	if HOP_SSH_KEY_FILE_PASSWORD != "" {
		auths = append( auths , ssh.Password( HOP_SSH_KEY_FILE_PASSWORD ) )
	}
	hop_signer , hop_signer_error := ssh.ParsePrivateKey( hop_connection_info.SSHKeyBytes )
	if hop_signer_error != nil {
		fmt.Printf( "unable to parse hop private key: %v\n" , hop_signer_error )
		panic( "unable to parse hop private key" )
	}
	auths = append( auths , ssh.PublicKeys( hop_signer ) )
	ssh_config := &ssh.ClientConfig{
		User: hop_connection_info.Username ,
		Auth: auths ,
		HostKeyCallback: func( hostname string , remote net.Addr , key ssh.PublicKey ) error {
			return nil
		} ,
		Timeout: 10 * time.Second ,
	}
	address_string := fmt.Sprintf( "%s:%d" , hop_connection_info.IPAddress , hop_connection_info.Port )
	ssh_client , ssh_connection_error := ssh.Dial( "tcp" , address_string , ssh_config )
	if ssh_connection_error != nil {
		log.Fatalf( "unable to connect to [%s]: %v" , address_string , ssh_connection_error )
	}
	//defer ssh_client.Close()
	return
}

func ConnectToSecondary( jump_host_ssh_client *ssh.Client , secondary_user_info SSHConnectionInfo ) ( ssh_client *ssh.Client ) {
	SECONDARY_SSH_KEY_FILE_PASSWORD := ""
	var auths []ssh.AuthMethod
	if SECONDARY_SSH_KEY_FILE_PASSWORD != "" {
		auths = append( auths , ssh.Password( SECONDARY_SSH_KEY_FILE_PASSWORD ) )
	}
	secondary_ssh_key_signer , secondary_ssh_key_signer_error := ssh.ParsePrivateKey( secondary_user_info.SSHKeyBytes )
	if secondary_ssh_key_signer_error != nil {
		fmt.Printf( "unable to parse private key: %v\n" , secondary_ssh_key_signer_error )
		panic( "unable to parse private key" )
	}
	auths = append( auths , ssh.PublicKeys( secondary_ssh_key_signer ) )
	ssh_config := &ssh.ClientConfig{
		User: secondary_user_info.Username ,
		Auth: auths ,
		HostKeyCallback: func( hostname string , remote net.Addr , key ssh.PublicKey ) error {
			return nil
		} ,
		Timeout: 10 * time.Second ,
	}
	address_string := fmt.Sprintf( "%s:%d" , secondary_user_info.IPAddress , secondary_user_info.Port )
	ssh_proxy_connection , ssh_proxy_connection_error := jump_host_ssh_client.Dial( "tcp" , address_string )
	//ssh_client , ssh_connection_error := jump_host_ssh_client.Dial( "tcp" , address_string , ssh_config )
	// ssh_client , ssh_connection_error := ssh.Dial( "tcp" , address_string , ssh_config )
	if ssh_proxy_connection_error != nil {
		log.Fatalf( "unable to connect to ssh proxy [%s]: %v" , address_string , ssh_proxy_connection_error )
	}
	//defer ssh_client.Close()
	conn , chans , _ , err := new_client_connection( ssh_proxy_connection , address_string , ssh_config )
	if err != nil {
		// if strings.Contains( trace.Unwrap( err ).Error() , "ssh: handshake failed" ) {
			// ssh_proxy_connection.Close()
			// return nil
		// }
		ssh_proxy_connection.Close()
		return nil
	}
	// We pass an empty channel which we close right away to ssh.NewClient
	// because the client need to handle requests itself.
	emptyCh := make( chan *ssh.Request )
	close( emptyCh )
	ssh_client = ssh.NewClient( conn , chans , emptyCh )
	return
}

// Single Hop
func IntoShellFromHop( hop_connection_info SSHConnectionInfo , secondary_user_info SSHConnectionInfo ) {
	fmt.Println( "Jumping into Shell of Other User" )
	fmt.Printf( "Step [%d] of 4\n" , 1 )
	jump_host_connection := ConnectToHop( hop_connection_info )
	fmt.Printf( "Step [%d] of 4\n" , 2 )
	secondary_host_connection := ConnectToSecondary( jump_host_connection , secondary_user_info )
	fmt.Printf( "Step [%d] of 4\n" , 3 )
	session , session_error := secondary_host_connection.NewSession()
	if session_error != nil { panic( session_error ) }
	defer session.Close()
	fmt.Printf( "Step [%d] of 4\n" , 4 )
	terminal_config := &TerminalConfig {
		Term: "xterm" ,
		Height: 40 ,
		Weight: 80 ,
		Modes: ssh.TerminalModes {
			ssh.TTY_OP_ISPEED: 14400 , // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400 , // output speed = 14.4kbaud
		} ,
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	session_pty_err := session.RequestPty( terminal_config.Term , terminal_config.Height , terminal_config.Weight , terminal_config.Modes );
	if session_pty_err != nil { panic( session_pty_err ) }

	session_shell_err := session.Shell()
	if session_shell_err != nil { panic( session_shell_err ) }

	session_shell_wait_err := session.Wait()
	if session_shell_wait_err != nil { panic( session_shell_wait_err ) }
}

func IntoPortFromHop() {

}

// // Two Hoppers
// func IntoShellFromHopFromHop() {

// }

// func IntoPortFromHopFromHop() {

// }