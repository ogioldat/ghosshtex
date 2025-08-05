package main

import (
	"ghosshtex.com/ghosshtex"
)

func main() {
	editor := ghosshtex.NewEditor()
	server := ghosshtex.NewSSHServer(editor)
	defer server.OnExit()

	if err := server.Start(); err != nil {
		server.OnError()
		panic(err) // or log.Fatal(err)
	}
}
