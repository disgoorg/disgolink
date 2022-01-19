package lavalink

import (
	"bytes"
	"encoding/base64"
	"time"
)

const trackInfoVersioned int = 1
const trackInfoVersion int32 = 2

func EncodeToString(info TrackInfo) (str string, err error) {
	w := new(bytes.Buffer)

	if err = WriteInt32(w, trackInfoVersion); err != nil {
		return
	}
	if err = WriteString(w, info.Title()); err != nil {
		return
	}
	if err = WriteString(w, info.Author()); err != nil {
		return
	}
	if err = WriteInt64(w, info.Length().Milliseconds()); err != nil {
		return
	}
	if err = WriteString(w, info.Identifier()); err != nil {
		return
	}
	if err = WriteBool(w, info.IsStream()); err != nil {
		return
	}
	if err = WriteBool(w, info.URI() != nil); err != nil {
		return
	}
	if err = WriteNullableString(w, info.URI()); err != nil {
		return
	}
	if err = WriteString(w, info.SourceName()); err != nil {
		return
	}
	if err = WriteInt32(w, int32(w.Len()|trackInfoVersioned<<30)); err != nil {
		return
	}

	return base64.StdEncoding.EncodeToString(w.Bytes()), nil
}

func DecodeString(str string) (info TrackInfo, err error) {
	var data []byte
	data, err = base64.StdEncoding.DecodeString(str)
	if err != nil {
		return
	}

	r := bytes.NewReader(data)

	trackInfo := &DefaultTrackInfo{}

	value, err := ReadInt32(r)

	flags := int(value) & 0xC00000000 >> 30
	//messageSize := value & 0x3FFFFFFF

	var version uint8
	if flags&trackInfoVersioned != 0 {
		version, err = r.ReadByte()
		if err != nil {
			return
		}
		version = version & 0xFF
	} else {
		version = 1
	}

	if trackInfo.TrackTitle, err = ReadString(r); err != nil {
		return
	}

	if trackInfo.TrackAuthor, err = ReadString(r); err != nil {
		return
	}

	var length int64
	if length, err = ReadInt64(r); err != nil {
		return
	}
	trackInfo.TrackLength = time.Duration(length) * time.Millisecond

	trackInfo.TrackIdentifier, err = ReadString(r)
	if err != nil {
		return nil, err
	}

	if trackInfo.TrackIsStream, err = ReadBool(r); err != nil {
		return
	}

	if trackInfo.TrackURI, err = ReadNullableString(r); err != nil {
		return
	}

	if trackInfo.TrackSourceName, err = ReadString(r); err != nil {
		return
	}

	info = trackInfo
	return
}
