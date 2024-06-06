package readbuffer

import (
	"bytes"
	"encoding/binary"
)

type Echo struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func (e *Echo) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	return binary.Read(buf, binary.BigEndian, &e)
}
