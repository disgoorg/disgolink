package lavalink

import (
	"encoding/binary"
	"io"
)

func ReadInt64(r io.Reader) (i int64, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func ReadInt32(r io.Reader) (i int32, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func ReadUInt16(r io.Reader) (i uint16, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func ReadBool(r io.Reader) (b bool, err error) {
	return b, binary.Read(r, binary.BigEndian, &b)
}

func ReadString(r io.Reader) (string, error) {
	size, err := ReadUInt16(r)
	if err != nil {
		return "", err
	}
	bytes := make([]byte, size)
	if err = binary.Read(r, binary.BigEndian, &bytes); err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ReadNullableString(r io.Reader) (*string, error) {
	b, err := ReadBool(r)
	if err != nil || !b {
		return nil, err
	}

	s, err := ReadString(r)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
