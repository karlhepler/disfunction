package funk

type Option[Config any] func(*Config)

func ConfigWithOptions[Config any](opts []Option[Config]) Config {
	var config Config
	for _, opt := range opts {
		opt(&config)
	}
	return config
}
