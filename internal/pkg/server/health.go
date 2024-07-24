package server

import "net/http"

// 启动一个http服务器，提供api检查泵的健康状态
func serverHealthCheck(healthAddress string, healthPath string) {
	http.HandleFunc("/"+healthPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`"status": "OK"`))
		if err != nil {
			return
		}
	})
	if err := http.ListenAndServe(healthAddress, nil); err != nil {
		panic("health check failed: " + err.Error())
	}
}
