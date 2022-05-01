package lavalink

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
)

type CustomTrackEncoder func(track AudioTrack, w io.Writer) error

func EncodeToString(track AudioTrack, customTrackEncoder CustomTrackEncoder) (str string, err error) {
	w := new(bytes.Buffer)

	if err = WriteInt8(w, trackInfoVersion); err != nil {
		return
	}
	if err = WriteString(w, track.Info().Title); err != nil {
		return
	}
	if err = WriteString(w, track.Info().Author); err != nil {
		return
	}
	if err = WriteInt64(w, track.Info().Length.Milliseconds()); err != nil {
		return
	}
	if err = WriteString(w, track.Info().Identifier); err != nil {
		return
	}
	if err = WriteBool(w, track.Info().IsStream); err != nil {
		return
	}
	if err = WriteNullableString(w, track.Info().URI); err != nil {
		return
	}
	if err = WriteString(w, track.Info().SourceName); err != nil {
		return
	}

	if customTrackEncoder != nil {
		if err = customTrackEncoder(track, w); err != nil {
			return
		}
	}

	if err = WriteInt64(w, track.Info().Position.Milliseconds()); err != nil {
		return
	}

	output := bytes.NewBuffer(make([]byte, 0, 4+w.Len()))
	if err = WriteInt32(output, int32(w.Len())|trackInfoVersioned<<30); err != nil {
		return
	}

	if _, err = w.WriteTo(output); err != nil {
		return
	}

	return base64.StdEncoding.EncodeToString(output.Bytes()), nil
}

func WriteInt64(w io.Writer, i int64) error {
	return binary.Write(w, binary.BigEndian, i)
}

func WriteInt32(w io.Writer, i int32) error {
	return binary.Write(w, binary.BigEndian, i)
}

func WriteInt16(w io.Writer, i int16) error {
	return binary.Write(w, binary.BigEndian, i)
}

func WriteInt8(w io.Writer, i int8) error {
	return binary.Write(w, binary.BigEndian, i)
}

func WriteBool(w io.Writer, bool bool) error {
	return binary.Write(w, binary.BigEndian, bool)
}

func WriteString(w io.Writer, str string) error {
	data := []byte(str)

	if len(data) > 65535 {
		return errors.New("string too big")
	}

	if err := WriteInt16(w, int16(len(data))); err != nil {
		return err
	}
	_, err := w.Write(data) //binary.Write(w, binary.BigEndian, data)
	return err
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
