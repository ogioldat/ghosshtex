# Ghosshtex 🚀

A collaborative text editor over SSH, written in Go.

## Project Goals 🎯

- [ ] Learn Go
- [ ] Use SSH to support a multi-user text editor
- [ ] Use BubbleTea for TUI

## Getting Started 🏁

### Prerequisites

- Go 1.23 or newer
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Bubbles](https://github.com/charmbracelet/bubbles) (installed via Go modules)
- An SSH client (e.g., `ssh`)

### Setup

1. **Generate SSH host keys:**
   ```sh
   ssh-keygen -t rsa -b 4096 -f id_rsa
   ```
   (The private key `id_rsa` must be in the project root.)

2. **Build and run the server:**
   ```sh
   go run main.go
   ```

3. **Connect with SSH:**
   ```sh
   ssh -p 2022 <username>@localhost
   ```
   (🔑 Public key authentication is accepted for any key.)

## Development 🛠️

- To run tests:
  ```sh
  go test ./ghosshtex
  ```

- To debug, use the provided VS Code launch configuration.

## Roadmap 🗺️

- [ ] Implement shared document editing across sessions
- [ ] Add authentication and user management
- [ ] Improve TUI features (syntax highlighting, file saving, etc.)
- [ ] Add logging and monitoring