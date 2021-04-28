package api

type DefaultTrack struct {
	Track_ *string           `json:"track"`
	Info_  *DefaultTrackInfo `json:"info"`
}

func (t *DefaultTrack) Track() *string {
	if t.Track_ == nil {
		if err := t.EncodeInfo(); err != nil {
			return nil
		}
	}
	return t.Track_
}

func (t *DefaultTrack) Info() TrackInfo {
	if t.Info_ == nil {
		if err := t.DecodeInfo(); err != nil {
			return nil
		}
	}
	return t.Info_
}

func (t *DefaultTrack) EncodeInfo() (err error) {
	if t.Info_ == nil {
		err = ErrEmptyTrackInfo
		return
	}
	t.Track_, err = EncodeToString(t.Info_)
	return
}

func (t *DefaultTrack) DecodeInfo() (err error) {
	if t.Track_ == nil {
		err = ErrEmptyTrack
		return
	}
	t.Info_, err = DecodeString(*t.Track_)
	return
}

type DefaultTrackInfo struct {
	Identifier_ string  `json:"identifier"`
	IsSeekable_ bool    `json:"isSeekable"`
	Author_     string  `json:"author"`
	Length_     int     `json:"length"`
	IsStream_   bool    `json:"isStream"`
	Position_   int     `json:"position"`
	Title_      string  `json:"title"`
	URI_        *string `json:"uri"`
	SourceName_ string  `json:"sourceName"`
}

func (i *DefaultTrackInfo) Identifier() string {
	return i.Identifier_
}

func (i *DefaultTrackInfo) IsSeekable() bool {
	return i.IsSeekable_
}

func (i *DefaultTrackInfo) Author() string {
	return i.Author_
}

func (i *DefaultTrackInfo) Length() int {
	return i.Length_
}

func (i *DefaultTrackInfo) IsStream() bool {
	return i.IsStream_
}

func (i *DefaultTrackInfo) Position() int {
	return i.Position_
}

func (i *DefaultTrackInfo) Title() string {
	return i.Title_
}

func (i *DefaultTrackInfo) URI() *string {
	return i.URI_
}

func (i *DefaultTrackInfo) SourceName() string {
	return i.SourceName_
}
