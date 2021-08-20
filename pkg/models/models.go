package models

type TunnelsConfiguration struct {
	Tunnels []Tunnel `json:"tunnels"`
	URI     string   `json:"uri"`
}

type Tunnel struct {
	Name      string        `json:"name"`
	URI       string        `json:"uri"`
	PublicURL string        `json:"public_url"`
	Proto     string        `json:"proto"`
	Config    TunnelConfig  `json:"config"`
	Metrics   TunnelMetrics `json:"metrics"`
}

type TunnelConfig struct {
	Addr    string `json:"addr"`
	Inspect bool   `json:"inspect"`
}

type TunnelMetrics struct {
	Conns struct {
		Count  float64 `json:"count"`
		Gauge  float64 `json:"gauge"`
		Rate1  float64 `json:"rate1"`
		Rate5  float64 `json:"rate5"`
		Rate15 float64 `json:"rate15"`
		P50    float64 `json:"p50"`
		P90    float64 `json:"p90"`
		P95    float64 `json:"p95"`
		P99    float64 `json:"p99"`
	} `json:"conns"`
	HTTP struct {
		Count  float64 `json:"count"`
		Rate1  float64 `json:"rate1"`
		Rate5  float64 `json:"rate5"`
		Rate15 float64 `json:"rate15"`
		P50    float64 `json:"p50"`
		P90    float64 `json:"p90"`
		P95    float64 `json:"p95"`
		P99    float64 `json:"p99"`
	} `json:"http"`
}
