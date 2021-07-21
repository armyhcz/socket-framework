package socket_framework

import (
	"net/http"

	"socket-framework/alive"
	"socket-framework/processor"

	"github.com/gorilla/mux"
	"nhooyr.io/websocket"
)

var acceptOptions *websocket.AcceptOptions

func init() {
	acceptOptions = &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	}
}

var consumers *processor.Consumers

func Start(path string, listen string, c *processor.Consumers) error {
	consumers = c

	router := mux.NewRouter()
	router.HandleFunc(path, Server)

	return http.ListenAndServe(listen, router)
}

func Server(writer http.ResponseWriter, request *http.Request) {
	if conn, err := websocket.Accept(writer, request, acceptOptions); err == nil {

		go func() {
			ctx, sig := alive.KeepAlive()

			for {
				_, bs, err := conn.Read(ctx)
				if err != nil {
					_ = conn.Close(websocket.StatusNormalClosure, "")
					return
				}

				sig <- struct{}{}
				pack, err := consumers.Start(ctx, bs)
				if err != nil || pack.Type == processor.TypeCloseEvent {
					_ = conn.Close(websocket.StatusBadGateway, "")
					return
				}
				err = conn.Write(ctx, websocket.MessageText, pack.Format())
				if err != nil {
					_ = conn.Close(websocket.StatusNormalClosure, "")
					return
				}
			}
		}()

	}
}
