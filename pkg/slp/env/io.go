package env

import (
	"bufio"
	"io"
	"os"
)

type ioImpl struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	writer *bufio.Writer
}

var _ IO = &ioImpl{}

func DefaultIO() IO {
	writer := bufio.NewWriter(os.Stdout)
	return &ioImpl{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		writer: writer,
	}
}

func (i *ioImpl) ReadLine() (string, error) {
	scanner := bufio.NewScanner(i.stdin)
	if scanner.Scan() {
		return scanner.Text(), nil
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", io.EOF
}

func (i *ioImpl) ReadAll() ([]byte, error) {
	return io.ReadAll(i.stdin)
}

func (i *ioImpl) Write(data []byte) (int, error) {
	return i.writer.Write(data)
}

func (i *ioImpl) WriteString(s string) (int, error) {
	return i.writer.WriteString(s)
}

func (i *ioImpl) Flush() error {
	return i.writer.Flush()
}

func (i *ioImpl) WriteError(data []byte) (int, error) {
	return i.stderr.Write(data)
}

func (i *ioImpl) WriteErrorString(s string) (int, error) {
	return i.stderr.Write([]byte(s))
}

func (i *ioImpl) SetStdin(r io.Reader) {
	i.stdin = r
}

func (i *ioImpl) SetStdout(w io.Writer) {
	i.stdout = w
	i.writer = bufio.NewWriter(w)
}

func (i *ioImpl) SetStderr(w io.Writer) {
	i.stderr = w
}
