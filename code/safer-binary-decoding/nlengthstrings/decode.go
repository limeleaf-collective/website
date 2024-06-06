package nlengthstrings

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Message struct {
	Identifier uint16
	Hostname   string
}

func (m *Message) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	// Decode Identifer the same we did for others above in Echo
	if err := binary.Read(buf, binary.BigEndian, &m.Identifier); err != nil {
		return err
	}

	// Read the HostnameLen into a temporary variable
	var hostnameLen uint16
	if err := binary.Read(buf, binary.BigEndian, &hostnameLen); err != nil {
		return err
	}

	// Make a slice with sizing it with the value of the temporary variable
	// from above and read the next n bytes into it. Finally, we convert
	// that []byte into string and set the field on the Message.
	hostname := make([]byte, hostnameLen)
	n, err := buf.Read(hostname)
	if err != nil {
		return err
	}
	if n != int(hostnameLen) {
		return io.ErrUnexpectedEOF
	}
	m.Hostname = string(hostname)

	return nil
}
