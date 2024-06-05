package repeatedfields

import (
	"bytes"
	"encoding/binary"
)

type Message struct {
	Identifier uint16
	Ports      []uint16
}

func (m *Message) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	if err := binary.Read(buf, binary.BigEndian, &m.Identifier); err != nil {
		return err
	}

	// Read the HostnameLen into a temporary variable
	var numPorts uint16
	if err := binary.Read(buf, binary.BigEndian, &numPorts); err != nil {
		return err
	}

	// Make a slice with sizing it with the value of the temporary variable
	// from above and read the next n fields into it.
	m.Ports = make([]uint16, numPorts)
	for n := range numPorts {
		if err := binary.Read(buf, binary.BigEndian, &m.Ports[n]); err != nil {
			return err
		}
	}

	return nil
}
