package algorithm

import (
	"net/http"
)

func extend(w http.ResponseWriter, r *http.Request){
	return
}

func StartHttpServer() {
	http.Handle("/scheduler", extend)
	http.ListenAndServe(":9000", nil)
}