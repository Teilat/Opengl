package metric

type MetricsProvider interface {
}

type metricsProvider struct {
}

func Init() MetricsProvider {
	return &metricsProvider{}
}
