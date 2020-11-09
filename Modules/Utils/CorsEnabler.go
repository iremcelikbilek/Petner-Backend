package Utils

import "net/http"

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "content-type")
	(*w).Header().Set("Access-Control-Allow-Headers", "token")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}
