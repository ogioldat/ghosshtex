package ghosshtex

type EditorText string

type Editor struct {
	text EditorText

	onInput []func(input string)
	onLoad  []func()
	onClose []func()
}

func NewEditor() *Editor {
	return &Editor{
		text:    "",
		onInput: []func(input string){},
		onLoad:  []func(){},
		onClose: []func(){},
	}
}

func (ed *Editor) Close() {
	for _, handler := range ed.onClose {
		go handler()
	}
}

func (ed *Editor) Handle(session *Session) {
	tui := NewEditorTUI(session.channel, session.channel)
	tui.Run()
}

func (ed *Editor) OnConnect() {}
