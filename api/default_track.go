package api

type DefaultTrack struct {
	Base64Track *string   `json:"track"`
	TrackInfo   TrackInfo `json:"info"`
}

func (t *DefaultTrack) Track() *string {
	if t.Base64Track == nil {
		if err := t.EncodeInfo(); err != nil {
			return nil
		}
	}
	return t.Base64Track
}

func (t *DefaultTrack) Info() TrackInfo {
	if t.TrackInfo == nil {
		if err := t.DecodeInfo(); err != nil {
			return nil
		}
	}
	return t.TrackInfo
}

func (t *DefaultTrack) EncodeInfo() (err error) {
	if t.TrackInfo == nil {
		err = ErrEmptyTrackInfo
		return
	}
	t.Base64Track, err = EncodeToString(t.TrackInfo)
	return
}

func (t *DefaultTrack) DecodeInfo() (err error) {
	if t.Base64Track == nil {
		err = ErrEmptyTrack
		return
	}
	t.TrackInfo, err = DecodeString(*t.Base64Track)
	return
}

type DefaultTrackInfo struct {
	TrackIdentifier string  `json:"identifier"`
	TrackIsSeekable bool    `json:"isSeekable"`
	TrackAuthor     string  `json:"author"`
	TrackLength     int     `json:"length"`
	TrackIsStream   bool    `json:"isStream"`
	TrackPosition   int     `json:"position"`
	TrackTitle      string  `json:"title"`
	TrackURI        *string `json:"uri"`
	TrackSourceName string  `json:"sourceName"`
}

func (i *DefaultTrackInfo) Identifier() string {
	return i.TrackIdentifier
}

func (i *DefaultTrackInfo) IsSeekable() bool {
	return i.TrackIsSeekable
}

func (i *DefaultTrackInfo) Author() string {
	return i.TrackAuthor
}

func (i *DefaultTrackInfo) Length() int {
	return i.TrackLength
}

func (i *DefaultTrackInfo) IsStream() bool {
	return i.TrackIsStream
}

func (i *DefaultTrackInfo) Position() int {
	return i.TrackPosition
}

func (i *DefaultTrackInfo) Title() string {
	return i.TrackTitle
}

func (i *DefaultTrackInfo) URI() *string {
	return i.TrackURI
}

func (i *DefaultTrackInfo) SourceName() string {
	return i.TrackSourceName
}
