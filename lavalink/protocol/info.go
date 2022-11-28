package protocol

type Info struct {
	Version        Version  `json:"version"`
	BuildTime      Time     `json:"buildTime"`
	Git            Git      `json:"git"`
	JVM            string   `json:"jvm"`
	Lavaplayer     string   `json:"lavaplayer"`
	SourceManagers []string `json:"sourceManagers"`
	Filters        []string `json:"filters"`
	Plugins        []Plugin `json:"plugins"`
}

type Version struct {
	Semver     string `json:"semver"`
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	PreRelease string `json:"preRelease"`
}

type Git struct {
	Branch     string `json:"branch"`
	Commit     string `json:"commit"`
	CommitTime Time   `json:"commitTime"`
}
