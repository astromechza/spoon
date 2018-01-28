package conf

func GenerateExampleConfig() *SpoonConfig {
	return &SpoonConfig{
		BasePath: "example.%(hostname)",
		Agents: []SpoonConfigAgent{
			SpoonConfigAgent{
				internalSpoonConfigAgent{
					Type:     "cmd",
					Interval: float32(30),
					Path:     ".cmd",
					Settings: map[string]interface{}{
						"cmd": []string{
							"python",
							"-c",
							"import random; print '.test.path_value', random.randint(-100, 100)",
						},
					},
					Enabled: true,
				},
			},
			SpoonConfigAgent{
				internalSpoonConfigAgent{
					Type:     "cpu",
					Interval: float32(60),
					Path:     ".cpu",
					Enabled:  true,
				},
			},
			SpoonConfigAgent{
				internalSpoonConfigAgent{
					Type:     "disk",
					Interval: float32(60),
					Path:     ".disk",
					Enabled:  true,
					Settings: map[string]interface{}{
						"device_regex": "da\\d$|disk\\d$",
					},
				},
			},
			SpoonConfigAgent{
				internalSpoonConfigAgent{
					Type:     "docker",
					Enabled:  false,
					Interval: float32(30),
					Path:     ".containers",
					Settings: map[string]interface{}{
						"container_filters": map[string]interface{}{},
					},
				},
			},
			SpoonConfigAgent{
				internalSpoonConfigAgent{
					Type:     "mem",
					Interval: float32(60),
					Path:     ".mem",
					Enabled:  true,
				},
			},
			SpoonConfigAgent{
				internalSpoonConfigAgent{
					Type:     "meta",
					Interval: float32(30),
					Path:     ".meta",
					Enabled:  true,
				},
			},
			SpoonConfigAgent{
				internalSpoonConfigAgent{
					Type:     "net",
					Interval: float32(60),
					Path:     ".net",
					Settings: map[string]interface{}{
						"nic_regex": "^e(th|n|m)\\d$",
					},
					Enabled: true,
				},
			},
			SpoonConfigAgent{
				internalSpoonConfigAgent{
					Type:     "random",
					Interval: float32(10),
					Path:     ".random",
					Enabled:  true,
				},
			},
			SpoonConfigAgent{
				internalSpoonConfigAgent{
					Type:     "time",
					Interval: float32(10),
					Path:     ".time_unix_seconds",
					Enabled:  true,
				},
			},
			SpoonConfigAgent{
				internalSpoonConfigAgent{
					Type:     "uptime",
					Interval: float32(60),
					Path:     ".uptime_seconds",
					Enabled:  true,
				},
			},
		},
		Sink: SpoonConfigSink{
			internalSpoonConfigSink{
				Type: "log",
			},
		},
	}
}
