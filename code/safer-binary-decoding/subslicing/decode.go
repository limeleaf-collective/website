package subslicing

import (
	"encoding/binary"
	"errors"
)

type Echo struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func (e *Echo) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return errors.New("invalid packet size")
	}

	e.Type = data[0]
	e.Code = data[1]

	e.Checksum = binary.BigEndian.Uint16(data[2:4])
	e.Identifier = binary.BigEndian.Uint16(data[4:6])
	e.SequenceNum = binary.BigEndian.Uint16(data[6:8])

	return nil
}
