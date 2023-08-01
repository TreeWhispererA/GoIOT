package stream

import (
	"fmt"
	"io"
)

type MemoryFile struct {
	data   []byte
	cursor int
}

type MemoryStream struct {
	file *MemoryFile
}

func NewMemoryStream(data []byte) *MemoryStream {
	return &MemoryStream{file: &MemoryFile{data: data, cursor: 0}}
}

func (stream MemoryStream) Length() int {
	return len(stream.file.data)
}

func (stream MemoryStream) Position() int {
	return stream.file.cursor
}

func (stream MemoryStream) Seek(offset int64, whence int) (int64, error) {
	var newPos int64
	if whence == io.SeekStart {
		newPos = offset
	} else if whence == io.SeekCurrent {
		newPos = offset + int64(stream.file.cursor)
	} else if whence == io.SeekEnd {
		newPos = offset + int64(len(stream.file.data))
	} else {
		return int64(stream.file.cursor), fmt.Errorf("Invalid whence: %d", whence)
	}

	if newPos < 0 || newPos > int64(len(stream.file.data)) {
		return int64(stream.file.cursor), fmt.Errorf("Position out of range")
	}

	stream.file.cursor = int(newPos)
	return newPos, nil
}

func (stream MemoryStream) Read(p []byte) (n int, err error) {
	end := len(stream.file.data)
	err = error(nil)

	var i int
	for i = 0; i < len(p) && stream.file.cursor < end; i++ {
		p[i] = stream.file.data[stream.file.cursor]
		stream.file.cursor += 1
	}
	if i < len(p) {
		err = io.EOF
	}

	return i, nil
}

func (stream MemoryStream) Write(p []byte) (n int, err error) {
	end := len(stream.file.data)
	err = error(nil)

	var i int
	for i = 0; i < len(p) && stream.file.cursor < end; i++ {
		stream.file.data[stream.file.cursor] = p[i]
		stream.file.cursor += 1
	}
	if i < len(p) {
		err = io.ErrShortBuffer
	}

	return i, err
}

func CloneStream(stream io.ReadSeeker) (*MemoryStream, error) {
	var err error
	var position, length int64

	if position, err = stream.Seek(0, io.SeekCurrent); err != nil {
		return nil, err
	} else if length, err = stream.Seek(0, io.SeekEnd); err != nil {
		return nil, err
	} else if _, err = stream.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	data := make([]byte, length)
	if _, err = stream.Read(data); err != nil {
		return nil, err
	} else if _, err = stream.Seek(position, io.SeekStart); err != nil {
		return nil, err
	}

	return &MemoryStream{file: &MemoryFile{data: data, cursor: 0}}, nil
}
