package subslicing_test

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sbd/nlengthstrings"
	"sbd/readbuffer"
	"sbd/repeatedfields"
	"sbd/subslicing"
	"testing"
)

type subtest struct {
	expected encoding.BinaryUnmarshaler
	actual   encoding.BinaryUnmarshaler
	err      error
}

func TestEcho(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		subtests []subtest
	}{
		{
			name: "success",
			data: []byte{
				0x08,       // Type
				0x00,       // Code
				0x30, 0x39, // Checksum
				0x00, 0x0d, // Identifier
				0x00, 0x04, // SequencyNum
			},
			subtests: []subtest{
				{
					expected: &subslicing.Echo{
						Type:        8,
						Code:        0,
						Checksum:    12345,
						Identifier:  13,
						SequenceNum: 4,
					},
					actual: &subslicing.Echo{},
				},
				{
					expected: &readbuffer.Echo{
						Type:        8,
						Code:        0,
						Checksum:    12345,
						Identifier:  13,
						SequenceNum: 4,
					},
					actual: &readbuffer.Echo{},
				},
			},
		},
		{
			name: "error",
			data: []byte{
				0x08,       // Type
				0x00,       // Code
				0x30, 0x39, // Checksum
				0x02, // Invalid Identifier
			},
			subtests: []subtest{
				{
					actual:   &subslicing.Echo{},
					expected: &subslicing.Echo{},
					err:      errors.New("invalid packet size"),
				},
				{
					actual: &readbuffer.Echo{},
					expected: &readbuffer.Echo{
						Type:     8,
						Code:     0,
						Checksum: 12345,
					},
					err: io.ErrUnexpectedEOF,
				},
			},
		},
	}

	for _, test := range tests {
		for _, subtest := range test.subtests {
			t.Run(fmt.Sprintf("%s-%T", test.name, subtest.expected), func(t *testing.T) {
				err := subtest.actual.UnmarshalBinary(test.data)
				if err != nil && err.Error() != subtest.err.Error() {
					t.Error("unexpected error", err, subtest.err)
				}
				if !reflect.DeepEqual(subtest.actual, subtest.expected) {
					t.Errorf("expected %v, got %v", subtest.expected, subtest.actual)
				}
			})
		}
	}
}

func TestNLengthString(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected nlengthstrings.Message
		err      error
	}{
		{
			name: "success",
			data: []byte{
				0x00, 0x05, // Identifier
				0x00, 0x06, // Length of next string
				'f', 'o', 'o', 'b', 'a', 'r', // String "foobar"
			},
			expected: nlengthstrings.Message{
				Identifier: 5,
				Hostname:   "foobar",
			},
		},
		{
			name: "length too long",
			data: []byte{
				0x00, 0x05, // Identifier
				0x00, 0x12, // Length of next string
				'f', 'o', 'o', 'b', 'a', 'r', // String "foobar"
			},
			expected: nlengthstrings.Message{
				Identifier: 5,
			},
			err: io.ErrUnexpectedEOF,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var message nlengthstrings.Message

			err := message.UnmarshalBinary(test.data)
			if err != nil && err.Error() != test.err.Error() {
				t.Error("unexpected error", err, test.err)
			}
			if !reflect.DeepEqual(message, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, message)
			}
		})
	}
}

func TestRepeated(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected repeatedfields.Message
		err      error
	}{
		{
			name: "success",
			data: []byte{
				0x00, 0x05, // Identifier
				0x00, 0x03, // Number of ports
				0x00, 0x16,
				0x00, 0x50,
				0x0b, 0xb8,
			},
			expected: repeatedfields.Message{
				Identifier: 5,
				Ports:      []uint16{22, 80, 3000},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var message repeatedfields.Message

			err := message.UnmarshalBinary(test.data)
			if err != nil && err.Error() != test.err.Error() {
				t.Error("unexpected error", err, test.err)
			}
			if !reflect.DeepEqual(message, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, message)
			}
		})
	}
}
