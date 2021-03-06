package colors

// Color type represents a single color string value.
type Color string

const (
	Green   Color = "\u001b[32;1m"
	White   Color = "\u001b[37;1m"
	Yellow  Color = "\u001b[33;1m"
	Cyan    Color = "\u001b[36;1m"
	Magenta Color = "\u001b[35;1m"
	Red     Color = "\u001b[31;1m"
	Reset   Color = "\u001b[0m"
)
