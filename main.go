package main

import (
	"ghosshtex.com/ghosshtex"
)

func main() {
	server := ghosshtex.NewSSHServer()
	editor := ghosshtex.NewEditor()

	server.SetHandler(editor)

	server.Start()
}
