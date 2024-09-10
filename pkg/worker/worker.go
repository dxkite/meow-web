package worker

import "io"

type Worker struct {
	Name        string
	WorkPath    string
	Restart     int
	Environment []string
	Command     []string

	output io.Reader
	input  io.Writer
}
