package socks

import (
	"os"
	"log"
	"fmt"
	socks5 "github.com/armon/go-socks5"
)

func Connect( socks_server_port string ) {
	creds := socks5.StaticCredentials{
		"admin": "ba771b0f20a4153ae2e2fe76b558c7be" ,
	}
	authenticator := socks5.UserPassAuthenticator{ Credentials: creds }
	conf := &socks5.Config{
			AuthMethods: []socks5.Authenticator{ authenticator } ,
			Logger: log.New( os.Stdout , "" , log.LstdFlags ) ,
	}
	server , _ := socks5.New( conf )
	fmt.Println( socks_server_port )
	server_address_string := fmt.Sprintf( "127.0.0.1:%v" , socks_server_port )
	fmt.Println( server_address_string )
	server.ListenAndServe( "tcp" , server_address_string )
}