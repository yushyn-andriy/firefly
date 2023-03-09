package config

type ReplMode int

const (
	_ ReplMode = iota
	INTERACTIVE
	FROM_FILE
)

type Config struct {
	Mode         ReplMode
	Debug        bool
	CompilerMode bool
}
