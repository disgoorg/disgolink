package lavalink

import (
	"encoding/binary"
	"io"
)

func WriteInt64(w io.Writer, i int64) error {
	return binary.Write(w, binary.BigEndian, i)
}

func WriteInt32(w io.Writer, i int32) error {
	return binary.Write(w, binary.BigEndian, i)
}

func WriteUInt16(w io.Writer, i uint16) error {
	return binary.Write(w, binary.BigEndian, i)
}

func WriteBool(w io.Writer, bool bool) (err error) {
	var bInt uint8
	if bool {
		bInt = 1
	} else {
		bInt = 0
	}

	if err = binary.Write(w, binary.BigEndian, bInt); err != nil {
		return
	}
	return
}

func WriteString(w io.Writer, str string) (err error) {
	data := []byte(str)

	if err = WriteUInt16(w, uint16(len(data))); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, data); err != nil {
		return
	}
	return
}

func WriteNullableString(w io.Writer, str *string) error {
	if err := WriteBool(w, str != nil); err != nil {
		return err
	}
	if str != nil {
		return WriteString(w, *str)
	}
	return nil
}
