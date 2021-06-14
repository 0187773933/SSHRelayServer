package server
// https://github.com/awgh/sshell/blob/master/sshell.go
import (
	"fmt"
	// "strings"
	"io"
	"os"
	"encoding/binary"
	"net"
	"errors"
	"syscall"
	"unsafe"
	"sync"
	"os/exec"
	"runtime"
	ssh "golang.org/x/crypto/ssh"
	pty "github.com/kr/pty"
	keys "sshclientcli/v1/keys"
	portmap "sshclientcli/v1/portmap"
	// config "sshclientcli/v1/config"
	// terminal "golang.org/x/crypto/ssh/terminal"
	// config "sshclientcli/v1/config"
	// "github.com/awgh/sshell/commands"
	// commands "github.com/awgh/sshell/commands"
)

var DEFAULT_SHELL = "/bin/bash"
var USER_NUMBER_INT = 0
var SSH_SERVER_PORT string
// var AUTHORIZED_KEYS = map[string]string{
// 	"user": "AAAAC3NzaC1lZDI1NTE5AAAAIADi9ZoVZstck6ELY0EIB863kD4qp5i6DYpQJHkwBiEo" ,
// }
var SSH_SERVER_CONFIG = &ssh.ServerConfig{
	//ServerVersion:     "SSH-2.0-OpenSSH_7.3p1 Debian-1",
	ServerVersion:     "" ,
	PublicKeyCallback: PublicKeyCallback ,
}

func PublicKeyCallback( remoteConn ssh.ConnMetadata , remoteKey ssh.PublicKey ) ( *ssh.Permissions , error ) {

	public_key := keys.PUBLIC[ USER_NUMBER_INT - 1 ]
	fmt.Println( public_key )
	parsedAuthPublicKey, _, _, _, err := ssh.ParseAuthorizedKey( public_key )
	if err != nil {
		fmt.Println("Could not parse public key")
		fmt.Println( err )
		return nil, err
	}

	// Make sure the key types match
	if remoteKey.Type() != parsedAuthPublicKey.Type() {
		fmt.Println("Key types don't match")
		return nil, errors.New("Key types do not match")
	}

	remoteKeyBytes := remoteKey.Marshal()
	authKeyBytes := parsedAuthPublicKey.Marshal()

	// Make sure the key lengths match
	if len(remoteKeyBytes) != len(authKeyBytes) {
		fmt.Println("Key lengths don't match")
		return nil, errors.New("Keys do not match")
	}

	// Make sure every byte of the key matches up
	// TODO: This should be a constant time check
	keysMatch := true
	for i, b := range remoteKeyBytes {
		if b != authKeyBytes[i] {
			keysMatch = false
		}
	}

	if keysMatch == false {
		fmt.Println("Keys don't match")
		return nil, errors.New("Keys do not match")
	}

	return nil, nil
}

func HandleRequests( reqs <-chan *ssh.Request ) {
	for req := range reqs {
		fmt.Printf( "recieved out-of-band request: %+v\n" , req )
		req.Reply(true , []byte("ok"))
	}
}

func parseDims(b []byte) (uint32, uint32) {
	w := binary.BigEndian.Uint32(b)
	h := binary.BigEndian.Uint32(b[4:])
	return w, h
}

// Winsize stores the Height and Width of a terminal.
type Winsize struct {
	Height uint16
	Width  uint16
	x      uint16 // unused
	y      uint16 // unused
}

// SetWinsize sets the size of the given pty.
func SetWinsize(fd uintptr, w, h uint32) {
	fmt.Printf("window resize %dx%d\n", w, h)
	ws := &Winsize{Width: uint16(w), Height: uint16(h)}
	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(ws)))
}

func PtyRun(c *exec.Cmd, tty *os.File) (err error) {
	defer tty.Close()
	c.Stdout = tty
	c.Stdin = tty
	c.Stderr = tty
	c.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}
	return c.Start()
}

func HandleChannels( chans <-chan ssh.NewChannel ) {
	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if t := newChannel.ChannelType(); t != "session" {
			newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			fmt.Printf("could not accept channel (%s)", err)
			continue
		}

		// allocate a terminal for this channel
		fmt.Print("creating pty...")
		// Create new pty
		f, tty, err := pty.Open()
		if err != nil {
			fmt.Printf("could not start pty (%s)", err)
			continue
		}

		var shell string
		shell = os.Getenv("SHELL")
		if shell == "" {
			shell = DEFAULT_SHELL
		}

		// Sessions have out-of-band requests such as "shell", "pty-req" and "env"
		go func(in <-chan *ssh.Request) {
			for req := range in {
				//fmt.Printf("%v %s", req.Payload, req.Payload)
				ok := false
				switch req.Type {
				case "exec":
					ok = true
					command := string(req.Payload[4 : req.Payload[3]+4])
					cmd := exec.Command(shell, []string{"-c", command}...)

					cmd.Stdout = channel
					cmd.Stderr = channel
					cmd.Stdin = channel

					err := cmd.Start()
					if err != nil {
						fmt.Printf("could not start command (%s)", err)
						continue
					}

					// teardown session
					go func() {
						_, err := cmd.Process.Wait()
						if err != nil {
							fmt.Printf("failed to exit bash (%s)", err)
						}
						channel.Close()
						fmt.Printf("session closed")
					}()
				case "shell":
					cmd := exec.Command(shell)
					cmd.Env = []string{"TERM=xterm"}
					err := PtyRun(cmd, tty)
					if err != nil {
						fmt.Printf("%s", err)
					}

					// Teardown session
					var once sync.Once
					close := func() {
						channel.Close()
						fmt.Printf("session closed")
					}

					// Pipe session to bash and visa-versa
					go func() {
						io.Copy(channel, f)
						once.Do(close)
					}()

					go func() {
						io.Copy(f, channel)
						once.Do(close)
					}()

					// We don't accept any commands (Payload),
					// only the default shell.
					if len(req.Payload) == 0 {
						ok = true
					}
				case "pty-req":
					// Responding 'ok' here will let the client
					// know we have a pty ready for input
					ok = true
					// Parse body...
					termLen := req.Payload[ 3 ]
					termEnv := string( req.Payload[ 4 : termLen + 4 ] )
					w , h := parseDims( req.Payload[ termLen + 4 : ] )
					if runtime.GOOS != "windows" {
						SetWinsize( f.Fd() , w , h )
					}
					fmt.Printf( "pty-req '%s'" , termEnv )
				case "window-change":
					w , h := parseDims( req.Payload )
					if runtime.GOOS != "windows" {
						SetWinsize( f.Fd() , w , h )
					}
					continue //no response
				}

				if !ok { fmt.Printf( "declining %s request...\n" , req.Type ) }
				req.Reply (ok , nil )
			}
		}( requests )
	}
}

func Serve( user_number_int int ) {
	SSH_SERVER_PORT = fmt.Sprint( portmap.PORTS[ user_number_int - 1][0] )
	fmt.Printf( "Starting SSH Server For User : %v , On Port : %v\n" , user_number_int , SSH_SERVER_PORT )
	USER_NUMBER_INT = user_number_int
	parsed_private_key , parsed_private_key_err := ssh.ParsePrivateKey( keys.PRIVATE[ user_number_int - 1 ] )
	if parsed_private_key_err != nil { fmt.Println( parsed_private_key_err ); fmt.Println( "here 1" ) }
	SSH_SERVER_CONFIG.AddHostKey( parsed_private_key )
	listener , err := net.Listen( "tcp4" , "0.0.0.0" + ":" + SSH_SERVER_PORT )
	if err != nil { fmt.Printf( "Failed to listen on %s (%s)\n" , SSH_SERVER_PORT , err ) }
	fmt.Printf( "listening on %s:%s\n" , "0.0.0.0" , SSH_SERVER_PORT )
	for {
		tcpConn, err := listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept incoming connection (%s)\n", err)
			continue
		}
		// Before use, a handshake must be performed on the incoming net.Conn.
		sshConn, chans, reqs, err := ssh.NewServerConn(tcpConn, SSH_SERVER_CONFIG)
		if err != nil {
			fmt.Printf("failed to handshake (%s)\n", err)
			continue
		}

		// Check remote address
		fmt.Printf("new connection from %s (%s)\n", sshConn.RemoteAddr(), sshConn.ClientVersion())

		// Print incoming out-of-band Requests
		go HandleRequests(reqs)
		// Accept all channels
		go HandleChannels(chans)
	}
}