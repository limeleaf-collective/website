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

	if err := binary.Read(buf, binary.BigEndian, &e.Type); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &e.Code); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &e.Checksum); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &e.Identifier); err != nil {
		return err
	}

	if err := binary.Read(buf, binary.BigEndian, &e.SequenceNum); err != nil {
		return err
	}

	return nil
}
