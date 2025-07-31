package internal

type SharedResource struct {
	value string
	write chan string
	read  chan chan string
	done  chan bool
}

func (r *SharedResource) Update(value string) {
	// fmt.Println("requested update", value)
	r.write <- value
}

func (r *SharedResource) Read() string {
	resp := make(chan string)
	r.read <- resp
	value := <-resp

	// fmt.Println("requested read", value)

	return value
}

func (r *SharedResource) Close() {
	r.done <- true
}

func NewSharedResource(value string) *SharedResource {
	resource := &SharedResource{
		value: value,
		write: make(chan string),
		read:  make(chan chan string),
		done:  make(chan bool),
	}

	go func() {
		for {
			select {
			case newValue := <-resource.write:
				resource.value = newValue
			case read := <-resource.read:
				read <- resource.value
			case <-resource.done:
				return
			}
		}
	}()

	return resource
}
