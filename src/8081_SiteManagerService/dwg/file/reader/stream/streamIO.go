package stream

import (
	"io"

	"golang.org/x/text/encoding/charmap"
)

type StreamIO struct {
	encoding *charmap.Charmap
	stream   io.ReadSeeker
}

func (this *StreamIO) GetPosition() (int64, error) {
	return this.stream.Seek(0, io.SeekCurrent)
}

func (this *StreamIO) SetPosition(pos int64) (int64, error) {
	return this.stream.Seek(pos, io.SeekStart)
}

func (this *StreamIO) GetLength() (int64, error) {
	var position, length int64
	var err error
	if position, err = this.GetPosition(); err != nil {
		return 0, err
	} else if length, err = this.stream.Seek(0, io.SeekEnd); err != nil {
		return length, err
	}
	position, err = this.SetPosition(position)

	return length, err
}
