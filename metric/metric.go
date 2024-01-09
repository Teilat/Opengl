package metric

type MetricsProvider interface {
}

type metricsProvider struct {
	enabled bool
}

func Init() MetricsProvider {
	return &metricsProvider{}
}
