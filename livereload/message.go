package livereload

type (
	hello struct {
		Command    string   `json:"command"`
		Protocols  []string `json:"protocols"`
		ServerName string   `json:"serverName"`
	}

	alert struct {
		Command string `json:"command"`
		Message string `json:"message"`
	}

	reload struct {
		Command string `json:"command"`
		Path    string `json:"path"`    // as full as possible/known, absolute path preferred, file name only is OK
		LiveCSS bool   `json:"liveCSS"` // false to disable live CSS refresh
	}
)
