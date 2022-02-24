package lavalink

const (
	trackInfoVersioned int32 = 1
	trackInfoVersion   int32 = 2
)

type AudioTrackInfo struct {
	Identifier string   `json:"identifier"`
	Author     string   `json:"author"`
	Length     Duration `json:"length"`
	IsStream   bool     `json:"isStream"`
	Title      string   `json:"title"`
	URI        *string  `json:"uri"`
	SourceName string   `json:"sourceName"`
	Position   Duration `json:"position"`
}

func NewAudioTrack(info AudioTrackInfo) AudioTrack {
	return &BasicAudioTrack{
		AudioTrackInfo: info,
	}
}

type AudioTrack interface {
	Info() AudioTrackInfo
	SetPosition(position Duration)
	UserData() interface{}
	SetUserData(interface{})
	Clone() AudioTrack
}

type BasicAudioTrack struct {
	AudioTrackInfo AudioTrackInfo `json:"info"`
	userData       interface{}
}

func (t *BasicAudioTrack) Info() AudioTrackInfo {
	return t.AudioTrackInfo
}

func (t *BasicAudioTrack) SetPosition(position Duration) {
	t.AudioTrackInfo.Position = position
}

func (t *BasicAudioTrack) SetUserData(userData interface{}) {
	t.userData = userData
}

func (t *BasicAudioTrack) UserData() interface{} {
	return t.userData
}

func (t *BasicAudioTrack) Clone() AudioTrack {
	info := t.AudioTrackInfo
	info.Position = 0
	return &BasicAudioTrack{
		AudioTrackInfo: info,
		userData:       t.userData,
	}
}
