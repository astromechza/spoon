package conf

import "encoding/json"

// SpoonConfig is the definition of the json config structure
type SpoonConfig struct {
	BasePath string             `json:"base_path"`
	Agents   []SpoonConfigAgent `json:"agents"`
	Sink     SpoonConfigSink    `json:"sink"`
}

type internalSpoonConfigAgent struct {
	Enabled     bool            `json:"enabled"`
	Type        string          `json:"type"`
	Interval    float32         `json:"interval"`
	Path        string          `json:"path"`
	SettingsRaw json.RawMessage `json:"settings,omitempty"`
	Settings    interface{}     `json:"-"`
}

// SpoonConfigAgent is a sub structure of SpoonConfig
type SpoonConfigAgent struct {
	internalSpoonConfigAgent
}

func (s *SpoonConfigAgent) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(s.Settings)
	if err != nil {
		return nil, err
	}
	s.SettingsRaw = json.RawMessage(b)
	return json.Marshal(s.internalSpoonConfigAgent)
}

type internalSpoonConfigSink struct {
	Type        string          `json:"type"`
	SettingsRaw json.RawMessage `json:"settings,omitempty"`
	Settings    interface{}     `json:"-"`
}

// SpoonConfigSink is a sub structure of SpoonConfig
type SpoonConfigSink struct {
	internalSpoonConfigSink
}

func (s *SpoonConfigSink) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(s.Settings)
	if err != nil {
		return nil, err
	}
	s.SettingsRaw = json.RawMessage(b)
	return json.Marshal(s.internalSpoonConfigSink)
}
