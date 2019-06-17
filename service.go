package yan

import (
	"net/http"
)

func handler(res http.ResponseWriter, resq http.Request) error {
	return nil
}

func start() {
	server := new(http.Server)
	server.Addr = "127.0.0.1:8000"
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
