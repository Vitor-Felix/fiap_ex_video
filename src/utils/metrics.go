package utils

import "github.com/prometheus/client_golang/prometheus"

var (
	UploadsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "video_uploads_total",
			Help: "Total de uploads recebidos",
		},
	)

	ProcessingErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "video_processing_errors_total",
			Help: "Total de erros durante o processamento",
		},
	)

	ProcessedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "video_processed_total",
			Help: "Total de vídeos processados com sucesso",
		},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(UploadsTotal)
	prometheus.MustRegister(ProcessingErrorsTotal)
	prometheus.MustRegister(ProcessedTotal)
}
