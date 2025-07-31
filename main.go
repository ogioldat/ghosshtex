package main

import (
	"errors"
	"log"
	"net"
	"os"
	"sync"

	"ghosshtex.com/internal"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type SSHServer struct{}
type SSHServerRequestHandler interface {
	Handle()
}

func (server *SSHServer) Init() ssh.ServerConfig {
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
		log.Fatal("failed to load private key: ", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("failed to parse private key: ", err)
	}
	config.AddHostKey(private)

	return config
}

func (server *SSHServer) Listen() net.Listener {
	listener, err := net.Listen("tcp", "0.0.0.0:2022")

	log.Println("listening")

	if err != nil {
		log.Fatal(err)
	}

	return listener
}

func main() {
	server := SSHServer{}
	config := server.Init()
	listener := server.Listen()

	var wg sync.WaitGroup
	defer wg.Wait()

	sharedText := internal.NewSharedResource("hello!")

	for {
		newConn, err := listener.Accept()

		if err != nil {
			log.Fatalln("unable to accept new conn", err)
		}

		go func() {
			_, sshChan, _, err := ssh.NewServerConn(newConn, &config)

			if err != nil {
				log.Fatalln("failed to handshake: ", err)
			}

			for newChannel := range sshChan {
				if newChannel.ChannelType() != "session" {
					newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
					continue
				}
				channel, _, err := newChannel.Accept()
				if err != nil {
					log.Fatalf("could not accept channel: %v", err)
				}

				term := term.NewTerminal(channel, ">")

				wg.Add(1)
				go func() {
					defer func() {
						channel.Close()
						wg.Done()
					}()

					text := sharedText.Read()

					for {
						term.Write([]byte(text))

						line, err := term.ReadLine()
						if err != nil {
							break
						}

						sharedText.Update(line)
						text = sharedText.Read()
						log.Println("shared text:", text)
					}
				}()
			}
		}()
	}
}
