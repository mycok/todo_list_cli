package colors

type Color string

const (
	Green  Color = "\u001b[32;1m"
	White        = "\u001b[37;1m"
	Yellow       = "\u001b[33;1m"
	Reset  Color = "\u001b[0m"
)
