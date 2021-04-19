package api

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
	"log"
)

const trackInfoVersioned int = 1
const trackInfoVersion int32 = 2

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

func (t *Track) EncodeInfo() (err error) {
	t.Track, err = EncodeToString(t.Info)
	return
}

func EncodeToString(info *TrackInfo) (str string, err error) {
	w := new(bytes.Buffer)

	if err = binary.Write(w, binary.BigEndian, trackInfoVersion); err != nil {
		return
	}
	if err = writeStr(w, info.Title); err != nil {
		return
	}
	if err = writeStr(w, info.Author); err != nil {
		return
	}
	if err = binary.Write(w, binary.LittleEndian, uint64(info.Length)); err != nil {
		return
	}
	if err = writeStr(w, info.Identifier); err != nil {
		return
	}
	if err = binary.Write(w, binary.LittleEndian, info.IsStream); err != nil {
		return
	}
	if err = writeStr(w, info.URI); err != nil {
		return
	}
	if err = writeStr(w, info.SourceName); err != nil {
		return
	}
	if err = binary.Write(w, binary.LittleEndian, uint64(info.Position)); err != nil {
		return
	}


	buf := new(bytes.Buffer)
	log.Println("BRUG: ", trackInfoVersioned << 30)
	_ = binary.Write(buf, binary.LittleEndian, int32(w.Len() | trackInfoVersioned << 30))
	buf.Write(w.Bytes())

	log.Println("actual: ", buf.Bytes())
	log.Println("actual: ", string(buf.Bytes()))

	str = base64.StdEncoding.EncodeToString(buf.Bytes())
	return
}

func writeStr(w io.Writer, str string) (err error) {
	data := []byte(str)

	if err = binary.Write(w, binary.BigEndian, uint16(len(data))); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, data); err != nil {
		return
	}
	return
}

func (t *Track) DecodeInfo() (err error) {
	t.Info, err = DecodeString(t.Track)
	return
}

// DecodeString thx to https://github.com/foxbot/gavalink/blob/master/decoder.go
func DecodeString(str string) (info *TrackInfo, err error) {

	var data []byte
	data, err = base64.StdEncoding.DecodeString(str)
	if err != nil {
		return
	}

	log.Println("data: ", data)
	log.Println("data: ", string(data))
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
	log.Println("ignore: ", ignore)
	log.Println("value: ", value)

	var version uint8
	if flags&int32(trackInfoVersioned) == 0 {
		println("flags&trackInfoVersioned == 0")
		version = 1
	} else {
		if err = binary.Read(r, binary.LittleEndian, &version); err != nil {
			return
		}
	}

	if err = binary.Read(r, binary.LittleEndian, &ignore); err != nil {
		return nil, err
	}
	log.Println("ignore: ", ignore)

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
	if err = binary.Read(r, binary.LittleEndian, &hasURI); err != nil {
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
