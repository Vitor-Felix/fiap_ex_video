from prometheus_client import Counter, start_http_server


processed_total = Counter(
    "video_processed_total",
    "Total de vídeos processados com sucesso",
)


def start_metrics_server():
    start_http_server(8000)
    print("📊 [WORKER] Métricas Prometheus disponíveis em :8000/metrics")
