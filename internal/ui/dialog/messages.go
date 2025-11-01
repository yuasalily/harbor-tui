package dialog

type OpenConfirmDialogMsg struct {
	Title   string
	Body    string
	Hint    string
	Payload any
}

type DialogResultMsg struct {
	Confirmed bool
	Payload   any
}
