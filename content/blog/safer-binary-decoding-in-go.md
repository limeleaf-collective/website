+++
title = "Safer Binary Decoding in Go"
date = 2024-05-30
authors = ["Blain Smith"]

[taxonomies]
tags = ["Engineering"]

[extra]
feature_photo = ""
feature_photo_alt = ""
+++

Go is a popular language choice for building web services. Typically, when
building those web services, we end up encoding/decoding JSON as the
data format. The [`encoding/json`][0] package provides a safe way to turn
JSON payloads into Go structs, and vice versa. 

However, if we need to
handle raw `[]byte` that follow a binary encoding format that
is not self-describing, we need to do a bit more work and
implement the [`encoding.BinaryMarshaler`][1] and [`encoding.BinaryUnmarshaler`][2]
directly. Since we're dealing with `[]byte`, we need to respect slice
bounds to avoid triggering a `panic()` and crashing our service. 

Let's look at the two ways we can decode data into Go structs and compare how
one way will be safer than the other while yielding the same result. As an
added bonus, we'll end up with easier to understand code.

<!-- more -->

## ICMP Packet Format

To make this a bit more fun and challenging, lets decode a packet format
we all tend to take for granted, [Internet Control Message Protocol or
ICMP][3]. If you have every used [`ping`][4], then you've sent ICMP packets in
order to test and measure a remote host.

Reading the ICMP specification, we see that the packet format is
as follows for Echo and Echo Reply messages only, since they are more
interesting than some of the other ICMP payloads with fewer fields.

- Type (1 byte)
- Code (1 byte)
- Checksum (2 bytes)
- ID (2 bytes)
- Sequence Number (2 bytes)

Another way to read these payloads is to look at the protocol layout
in ASCII, like you see in the ICMP RFC:

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|     Type      |     Code      |          Checksum             |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|           Identifier          |        Sequence Number        |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|     Data ...
+-+-+-+-+-
```

In this view, we can visualize the bounds of the fields by number of
bits. For example, we see that `Code` occupied bits 8 to 15 which is 8
total bits or 1 byte.

## Defining the Struct and Decoding Function

Now that we know the layout of the binary data, we can define our Go
struct. Notice we are defining explicit sizes to our fields and not
using regular `int` types. We explicitly want `uint8` or `uint16` types
for our fields so we can properly decode it in binary.

```go
type Echo struct {
  Type        uint8
  Code        uint8
  Checksum    uint16
  Identifier  uint16
  SequenceNum uint16
}
```

Next, we can set up the stub decoding function. We define the function
with this name and signature to implement the `encoding.UnmarshalBinary`
interface to make this flexible throughout the Go standard library.

```go
func (e *Echo) UnmarshalBinary(buf []byte) error {
  // Read through buf and assign to Echo fields
  return nil
}
```

## Decoding by Subslicing the `[]byte`

First, we are going to perform the decoding by subslicing the `[]byte`
into our required `uint8` and `uint16` fields.

```go
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
```

We can assign `Type` and `Code` fields directly since a `byte` is just
a `uint8`. You can see this in the [type definition for `byte`][5]. However, 
since the rest of our fields are `uint16`, we need to
take 2 `byte`s and convert both of them at the same time to a single
`uint16` value using [big endian][6].

Looking back at our protocol format, we see that `Checksum` occupies bits
16 to 32 which is a total of 16 bits or 2 bytes. That is why we are
taking bytes 2 up to, but not including 4 (2 and 3 only). We follow the
same pattern for the rest of the `uint16` fields.

You might be thinking this is a perfectly good way to decode a byte
slice into a struct since we're also checking the length of the incoming
slice to ensure we won't slice out of bounds and cause a `panic`. While
you would be correct, it is a bit awkward to reason about the slicing
semantics as you're scanning the code. The `buf[2:4]` isn't very clear
if your goal is to scan the code and understand it quickly. What if this
payload had 5-10 times as many fields? Slicing could lead to incorrect
field bounds. Fortunately, there is a clearer and safer way to perform
the same decoding.

## Deocoding with `bytes.Buffer` and `binary.Read`

For the clearer and safer method, we use [`bytes.Buffer`][7] and
[`binary.Read`][8]. Converting our `BinaryUnmarshal` function now looks
like this.

```go
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
```

Now we have a much clearer and safer decoding function that reads from
the `bytes.Buffer` that wraps the `data []byte`, and since our struct
fields contain the correct sizes (`uint8` and `uint16`), we no longer
have to keep track of the subslicing indicies. Even though we have much
more of the dreaded `if err != nil` Go error checking, we end up with
`panic` free code since `binary.Read` will safely return an `io.EOF`
should something bad happen. We'll end up catching that error and
returning it to the caller.

## Bonus: Decoding the Entire `struct` at Once

Since `Echo`'s fields are all fixed sizes, we can pass the entire thing
into `binary.Read` and get the same result as above, but with a single
call.

```go
func (e *Echo) UnmarshalBinary(data []byte) error {
  buf := bytes.NewBuffer(data)

  return binary.Read(buf, binary.BigEndian, &e)
}
```

This method works great for structs that have all fixed sized fields and
the binary data on the wire matches our struct, but that will not
always be the case. Sometimes we'll have to account for n-length
strings and n-repeated fields.

## Variable Length Strings

For these examples, let's assume we have a new message that is encoded as
binary.

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|   Identifier  |  HostnameLen  |          Hostname ...
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-
```

Since we have a variable length string `Hostname`, the length is
first encoded before the string data itself to tell us how many bytes
we should read that represent the hostname.

Knowing this, we can now set up our type to store the relevant fields
like we did with `Echo`.

```go
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
```

This length-prefixed string format is very common for binary encodings.
The only thing we need to know for sure is what type the length is
encoded as so we know what kind of temporary variable to use.

## Repeated Fields

Repeated fields can work in a similar fashion by size-prefixing the
number of times the field occurs in the data.

For these examples, lets assume we have a new message that is encoded as
binary.

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|   Identifier  |    NumPorts   |      Port0      |    Port1    |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|      Port2    |      Port3    |      PortN ...
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-
```

Now we can decode this data with the following strategy.

```go
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
```

Since we're still using `binary.Read` and managing any returned `error`, 
we can safely handle situations where the data we're decoding indicates 
there should be 5 ports, but only provides 2. We'll catch the `io.EOF` 
error and return accordingly instead of needing to handle subslicing 
and potential `panic`s being thrown.

## Conclusion

We've seen that using [`bytes.Buffer`][7] and [`binary.Read`][8] leads
to safer reading of `[]byte` and clearer code for other
engineers on your team who might read it later.

So, the next time you need to decode binary data and have a specification
to follow, try your hand at these techniques to make your code safe,
correct, readable, and composable.

The [complete source code][9] is avaialable along with supporting tests.

[0]: https://pkg.go.dev/encoding/json
[1]: https://pkg.go.dev/encoding#BinaryMarshaler
[2]: https://pkg.go.dev/encoding#BinaryUnmarshaler
[3]: https://www.rfc-editor.org/rfc/rfc792
[4]: https://www.man7.org/linux/man-pages/man8/ping.8.html
[5]: https://pkg.go.dev/builtin#byte
[6]: https://en.wikipedia.org/wiki/Endianness
[7]: https://pkg.go.dev/bytes#Buffer
[8]: https://pkg.go.dev/encoding/binary#Read
[9]: https://github.com/limeleaf-collective/limeleaf.net/tree/main/code/safer-binary-decoding
