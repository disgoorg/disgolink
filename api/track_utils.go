package api

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
)

const trackInfoVersioned int = 1
const trackInfoVersion int8 = 2

func EncodeToString(info TrackInfo) (str *string, err error) {
	w := new(bytes.Buffer)

	if err = binary.Write(w, binary.BigEndian, trackInfoVersion); err != nil {
		return
	}
	if err = writeStr(w, info.Title()); err != nil {
		return
	}
	if err = writeStr(w, info.Author()); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint64(info.Length())); err != nil {
		return
	}
	if err = writeStr(w, info.Identifier()); err != nil {
		return
	}
	if err = writeBool(w, info.IsStream()); err != nil {
		return
	}
	if err = writeBool(w, info.URI() != nil); err != nil {
		return
	}
	if info.URI() != nil {
		if err = writeStr(w, *info.URI()); err != nil {
			return
		}
	}

	if err = writeStr(w, info.SourceName()); err != nil {
		return
	}
	if err = binary.Write(w, binary.BigEndian, uint64(info.Position())); err != nil {
		return
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, int32(w.Len()|trackInfoVersioned<<30))
	if err != nil {
		return
	}
	buf.Write(w.Bytes())

	str2 := base64.StdEncoding.EncodeToString(buf.Bytes())
	str = &str2
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

func writeBool(w io.Writer, bool bool) (err error) {
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

// DecodeString thx to https://github.com/foxbot/gavalink/blob/master/decoder.go
func DecodeString(str string) (info *DefaultTrackInfo, err error) {

	var data []byte
	data, err = base64.StdEncoding.DecodeString(str)
	if err != nil {
		return
	}

	r := bytes.NewReader(data)

	trackInfo := &DefaultTrackInfo{}

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
	if flags&int32(trackInfoVersioned) == 0 {
		version = 1
	} else {
		if err = binary.Read(r, binary.LittleEndian, &version); err != nil {
			return
		}
	}

	if err = binary.Read(r, binary.LittleEndian, &ignore); err != nil {
		return nil, err
	}

	trackInfo.Title_, err = readStr(r)
	if err != nil {
		return
	}

	trackInfo.Author_, err = readStr(r)
	if err != nil {
		return
	}

	var length uint64
	if err = binary.Read(r, binary.BigEndian, &length); err != nil {
		return
	}
	trackInfo.Length_ = int(length)

	trackInfo.Identifier_, err = readStr(r)
	if err != nil {
		return nil, err
	}

	var isStream uint8
	if err = binary.Read(r, binary.LittleEndian, &isStream); err != nil {
		return
	}
	trackInfo.IsStream_ = isStream == 1
	trackInfo.IsSeekable_ = !trackInfo.IsStream()

	var hasURI uint8
	if err = binary.Read(r, binary.LittleEndian, &hasURI); err != nil {
		return nil, err
	}

	if hasURI == 1 {
		var uri string
		uri, err = readStr(r)
		trackInfo.URI_ = &uri
		if err != nil {
			return
		}
	} else {
		trackInfo.URI_ = nil
		_, err = readStr(r)
		if err != nil {
			return
		}
	}

	trackInfo.SourceName_, err = readStr(r)
	if err != nil {
		return
	}

	var position uint64
	if err = binary.Read(r, binary.BigEndian, &position); err != nil {
		return
	}
	trackInfo.Position_ = int(position)

	info = trackInfo
	return
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

func DefaultTracksToTracks(defaultTracks []*DefaultTrack) []Track {
	tracks := make([]Track, len(defaultTracks))
	for i := 0; i < len(defaultTracks); i++ {
		tracks[i] = defaultTracks[i]
	}
	return tracks
}
