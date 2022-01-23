package lavalink

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
	"time"
)

type CustomTrackInfoDecoder func(info AudioTrackInfo, r io.Reader) (AudioTrack, error)

func DecodeString(str string, customTrackInfoDecoder CustomTrackInfoDecoder) (track AudioTrack, err error) {
	var data []byte
	if data, err = base64.StdEncoding.DecodeString(str); err != nil {
		return
	}

	r := bytes.NewReader(data)

	info := &DefaultAudioTrackInfo{}
	var value int32
	if value, err = ReadInt32(r); err != nil {
		return
	}
	flags := int(value) & 0xC00000000 >> 30
	//messageSize := value & 0x3FFFFFFF

	var version uint8
	if flags&trackInfoVersioned != 0 {
		if version, err = r.ReadByte(); err != nil {
			return
		}
		version = version & 0xFF
	} else {
		version = 1
	}

	if info.TrackTitle, err = ReadString(r); err != nil {
		return
	}

	if info.TrackAuthor, err = ReadString(r); err != nil {
		return
	}

	var length int64
	if length, err = ReadInt64(r); err != nil {
		return
	}
	info.TrackLength = time.Duration(length) * time.Millisecond

	if info.TrackIdentifier, err = ReadString(r); err != nil {
		return
	}

	if info.TrackIsStream, err = ReadBool(r); err != nil {
		return
	}

	if info.TrackURI, err = ReadNullableString(r); err != nil {
		return
	}

	if info.TrackSourceName, err = ReadString(r); err != nil {
		return
	}

	if customTrackInfoDecoder != nil {
		if track, err = customTrackInfoDecoder(info, r); err != nil {
			return
		}
	}

	var position int64
	if position, err = ReadInt64(r); err != nil {
		return
	}
	info.TrackPosition = time.Duration(position) * time.Millisecond

	return NewAudioTrack(str, info), nil
}

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
	b := make([]byte, size)
	if err = binary.Read(r, binary.BigEndian, &b); err != nil {
		return "", err
	}
	return string(b), nil
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
