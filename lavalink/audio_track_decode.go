package lavalink

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
)

type CustomTrackInfoDecoder func(info AudioTrackInfo, r io.Reader) (AudioTrack, error)

func DecodeString(str string, customTrackInfoDecoder CustomTrackInfoDecoder) (track AudioTrack, err error) {
	var data []byte
	if data, err = base64.StdEncoding.DecodeString(str); err != nil {
		return
	}

	var (
		r     = bytes.NewReader(data)
		info  AudioTrackInfo
		value int32
	)

	if value, err = ReadInt32(r); err != nil {
		return
	}

	flags := int32(int64(value) & 0xC0000000 >> 30)
	messageSize := int64(value) & 0x3FFFFFFF
	if messageSize == 0 {
		return nil, nil
	}

	var version int32
	if flags&trackInfoVersioned == 0 {
		version = 1
	} else {
		var v uint8
		if v, err = ReadUInt8(r); err != nil {
			return
		}
		version = int32(v)
	}

	if info.Title, err = ReadString(r); err != nil {
		return
	}
	if info.Author, err = ReadString(r); err != nil {
		return
	}

	var length int64
	if length, err = ReadInt64(r); err != nil {
		return
	}
	info.Length = Duration(length)

	if info.Identifier, err = ReadString(r); err != nil {
		return
	}
	if info.IsStream, err = ReadBool(r); err != nil {
		return
	}
	if version >= 2 {
		if info.URI, err = ReadNullableString(r); err != nil {
			return
		}
	}
	if info.SourceName, err = ReadString(r); err != nil {
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

	if track == nil {
		track = NewAudioTrack(info)
	}
	track.SetPosition(Duration(position))

	return track, nil
}

func ReadInt64(r io.Reader) (i int64, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func ReadInt32(r io.Reader) (i int32, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func ReadInt16(r io.Reader) (i int16, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func ReadUInt8(r io.Reader) (i uint8, err error) {
	return i, binary.Read(r, binary.BigEndian, &i)
}

func ReadBool(r io.Reader) (b bool, err error) {
	return b, binary.Read(r, binary.BigEndian, &b)
}

func ReadString(r io.Reader) (string, error) {
	size, err := ReadInt16(r)
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
