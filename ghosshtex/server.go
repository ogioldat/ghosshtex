package ghosshtex

import (
	"errors"
	"log"
	"log/slog"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSHServer struct {
	handler SessionHandler
	config  ssh.ServerConfig
	wg      *sync.WaitGroup
}

type User struct {
	name string
	host net.Addr
}

type Session struct {
	channel ssh.Channel
	user    User
}

type SessionHandler interface {
	Handle(*Session)
	OnConnect()
}

func NewSSHServer(handler SessionHandler) SSHServer {
	var wg sync.WaitGroup

	algos := ssh.SupportedAlgorithms()
	config := ssh.ServerConfig{
		Config:                  ssh.Config{KeyExchanges: algos.KeyExchanges, Ciphers: algos.Ciphers, MACs: algos.MACs},
		PublicKeyAuthAlgorithms: algos.PublicKeyAuths,
		PublicKeyCallback: func(conn ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			return &ssh.Permissions{
				// Record the public key used for authentication.
				Extensions: map[string]string{
					"pubkey-fp": ssh.FingerprintSHA256(pubKey),
				},
			}, nil

			// TODO: Check authorized keys
			// return nil, fmt.Errorf("unknown public key for %q", c.User())
		},
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			return nil, errors.New("password auth not supported")
		},
	}

	privateBytes, err := os.ReadFile("id_rsa")
	if err != nil {
		slog.Error("failed to load private key: ", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		slog.Error("failed to parse private key: ", err)
	}
	config.AddHostKey(private)

	return SSHServer{
		config:  config,
		handler: handler,
		wg:      &wg,
	}
}

func (server *SSHServer) Start() error {
	server.OnStart()

	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		return err
	}

	server.OnListen()

	for {
		newConn, err := listener.Accept()
		if err != nil {
			log.Println("unable to accept new conn:", err)
			continue
		}
		go server.handleConnection(newConn)
	}
}

func (server *SSHServer) handleConnection(newConn net.Conn) {
	sshConn, sshChan, _, err := ssh.NewServerConn(newConn, &server.config)
	if err != nil {
		log.Println("failed to handshake:", err)
		return
	}
	for newChannel := range sshChan {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, _, err := newChannel.Accept()
		if err != nil {
			log.Println("could not accept channel:", err)
			continue
		}
		user := User{
			name: sshConn.Conn.User(),
			host: sshConn.Conn.RemoteAddr(),
		}
		session := Session{
			channel: channel,
			user:    user,
		}
		server.handler.Handle(&session)
	}
}

func (server *SSHServer) OnStart() {
	slog.Info("Server started")
}

func (server *SSHServer) OnConnect() {
	slog.Info("Server started")
}

func (server *SSHServer) OnListen() {
	slog.Info("Server listening")
}

func (server *SSHServer) OnError() {
	slog.Error("Error occurred")
}

func (server *SSHServer) OnExit() {
	slog.Info("Server exited")
}
