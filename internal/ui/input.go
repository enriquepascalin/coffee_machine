package ui

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type Input struct {
	reader *bufio.Reader
}

func NewInput(reader io.Reader) *Input {
	return &Input{
		reader: bufio.NewReader(reader),
	}
}

func (i *Input) ReadLine() (string, error) {
	line, err := i.reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && len(line) > 0 {
			return strings.TrimSpace(line), nil
		}

		return "", err
	}

	return strings.TrimSpace(line), nil
}

func (i *Input) ReadInt() (int, error) {
	line, err := i.ReadLine()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(line))
}
