package api

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
)

const trackInfoVersioned int32 = 1

type Track struct {
	Track string     `json:"track"`
	Info  *TrackInfo `json:"info"`
}

type TrackInfo struct {
	Identifier string `json:"identifier"`
	IsSeekable bool   `json:"isSeekable"`
	Author     string `json:"author"`
	Length     int    `json:"length"`
	IsStream   bool   `json:"isStream"`
	Position   int    `json:"position"`
	Title      string `json:"title"`
	URI        string `json:"uri"`
	SourceName string `json:"sourceName"`
}

func (t *Track) DecodeInfo() (err error) {
	t.Info, err = DecodeString(t.Track)
	if err != nil {
		return
	}
	return
}

// DecodeString thx to https://github.com/foxbot/gavalink/blob/master/decoder.go
func DecodeString(str string) (info *TrackInfo, err error) {

	var data []byte
	data, err = base64.StdEncoding.DecodeString(str)
	if err != nil {
		return
	}

	r := bytes.NewReader(data)

	info = &TrackInfo{}

	var value uint8
	if err = binary.Read(r, binary.LittleEndian, &value); err != nil {
		return
	}

	flags := int32(int64(value) & 0xC00000000)

	var ignore [2]byte
	if err = binary.Read(r, binary.LittleEndian, &ignore); err != nil {
		return
	}

	var version uint8
	if flags&trackInfoVersioned == 0 {
		version = 1
	} else {
		if err = binary.Read(r, binary.LittleEndian, &version); err != nil {
			return
		}
	}

	if err = binary.Read(r, binary.LittleEndian, &ignore); err != nil {
		return nil, err
	}

	info.Title, err = readStr(r)
	if err != nil {
		return
	}

	info.Author, err = readStr(r)
	if err != nil {
		return
	}

	var length uint64
	if err = binary.Read(r, binary.BigEndian, &length); err != nil {
		return
	}
	info.Length = int(length)

	info.Identifier, err = readStr(r)
	if err != nil {
		return nil, err
	}

	var isStream uint8
	if err = binary.Read(r, binary.LittleEndian, &isStream); err != nil {
		return
	}
	info.IsStream = isStream == 1
	info.IsSeekable = !info.IsStream

	var hasURI uint8
	if err := binary.Read(r, binary.LittleEndian, &hasURI); err != nil {
		return nil, err
	}

	if hasURI == 1 {
		info.URI, err = readStr(r)
		if err != nil {
			return
		}
	} else {
		_, err = readStr(r)
		if err != nil {
			return
		}
	}

	info.SourceName, err = readStr(r)
	if err != nil {
		return
	}

	return info, nil
}

func readStr(r io.Reader) (string, error) {
	var size uint16
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return "", err
	}
	buf := make([]byte, size)
	if err := binary.Read(r, binary.BigEndian, &buf); err != nil {
		return "", err
	}
	return string(buf), nil
}
