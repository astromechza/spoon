package conf

import (
    "encoding/json"
    "io/ioutil"
)

// SpoonConfigAgent is a sub structure of SpoonConfig
type SpoonConfigAgent struct {
    Type string
    Interval float32
    Path string
    Settings map[string]interface{}
}

// SpoonConfigLog is a sub structure of SpoonConfig
type SpoonConfigLog struct {
    Path string
    RotateSize int64 `json:"rotate_size"`
}

// SpoonConfig is the definition of the json config structure
type SpoonConfig struct {
    Logging SpoonConfigLog
    Agents []SpoonConfigAgent
}

// Load the config information from the file on disk
func Load(path *string) (*SpoonConfig, error) {

    // first read all bytes from file
    data, err := ioutil.ReadFile(*path)
    if err != nil {
        return nil, err
    }

    // now parse config object out
    var cfg SpoonConfig
    err = json.Unmarshal(data, &cfg)
    if err != nil {
        return nil, err
    }

    // and return
    return &cfg, nil
}
