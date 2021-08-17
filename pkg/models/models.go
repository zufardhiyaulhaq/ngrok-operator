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
		Count  int `json:"count"`
		Gauge  int `json:"gauge"`
		Rate1  int `json:"rate1"`
		Rate5  int `json:"rate5"`
		Rate15 int `json:"rate15"`
		P50    int `json:"p50"`
		P90    int `json:"p90"`
		P95    int `json:"p95"`
		P99    int `json:"p99"`
	} `json:"conns"`
	HTTP struct {
		Count  int `json:"count"`
		Rate1  int `json:"rate1"`
		Rate5  int `json:"rate5"`
		Rate15 int `json:"rate15"`
		P50    int `json:"p50"`
		P90    int `json:"p90"`
		P95    int `json:"p95"`
		P99    int `json:"p99"`
	} `json:"http"`
}
