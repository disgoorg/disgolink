package lavalink

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
)

type CustomTrackEncoder func(info AudioTrack, w io.Writer) error

func EncodeToString(track AudioTrack, customTrackEncoder CustomTrackEncoder) (str string, err error) {
	w := new(bytes.Buffer)

	if err = WriteInt32(w, trackInfoVersion); err != nil {
		return
	}
	if err = WriteString(w, track.Info().Title()); err != nil {
		return
	}
	if err = WriteString(w, track.Info().Author()); err != nil {
		return
	}
	if err = WriteInt64(w, track.Info().Length().Milliseconds()); err != nil {
		return
	}
	if err = WriteString(w, track.Info().Identifier()); err != nil {
		return
	}
	if err = WriteBool(w, track.Info().IsStream()); err != nil {
		return
	}
	if err = WriteBool(w, track.Info().URI() != nil); err != nil {
		return
	}
	if err = WriteNullableString(w, track.Info().URI()); err != nil {
		return
	}
	if err = WriteString(w, track.Info().SourceName()); err != nil {
		return
	}

	if customTrackEncoder != nil {
		if err = customTrackEncoder(track, w); err != nil {
			return
		}
	}

	if err = WriteInt64(w, track.Info().Position().Milliseconds()); err != nil {
		return
	}
	if err = WriteInt32(w, int32(w.Len()|trackInfoVersioned<<30)); err != nil {
		return
	}

	return base64.StdEncoding.EncodeToString(w.Bytes()), nil
}

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
