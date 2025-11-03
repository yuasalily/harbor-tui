package dialog

type OpenDialogMsg struct {
	Title   string
	Body    string
	Hint    string
	Payload any
	Kind    DialogKind
}

type DialogResultMsg struct {
	Confirmed bool
	Payload   any
}

type DialogKind int

const (
	DialogKindConfirm DialogKind = iota
	DialogKindError
	DialogKindInfo
)
