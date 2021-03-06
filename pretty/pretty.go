package pretty

import "fmt"

type IconType int

const (
	None IconType = iota
	Info
	Failure
	Launch
	Success
)

var icons = map[IconType]string{
	None:    "  ",
	Info:    "âšī¸  ",
	Failure: "â ",
	Launch:  "đ ",
	Success: "â ",
}

func Printf(icon IconType, format string, args ...interface{}) {
	// TODO: add fall back for terminal that don't support unicode emojis
	fmt.Printf("%s ", icons[icon])
	fmt.Printf(format, args...)
	fmt.Println()
}
