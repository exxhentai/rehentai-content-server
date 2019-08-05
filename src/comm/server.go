package comm

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	ipfs "github.com/ipfs/go-ipfs-api"
)

type Server struct {
	shell   *ipfs.Shell
	shMux   *sync.Mutex
	pinlist []*websocket.Conn
}

func NewServer() *Server {
	sh := ipfs.NewLocalShell()

	if sh == nil {
		log.Fatalln("[ERROR]", "You need to run a ipfs node first!!!")
		return nil
	} else {
		l := []*websocket.Conn{}
		return &Server{sh, &sync.Mutex{}, l}
	}
}

func (s Server) Start(port int) {
	log.Println("[INFO]", "Starting rehentai content server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
		case "POST":
			filenames := s.receive(w, r)

			// DEBUG
			fmt.Println(filenames)

			s.upload(filenames, w)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/pin", s.wshandler)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (s Server) wshandler(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	list := s.pinlist
	list = append(list, conn)
	s.pinlist = list
	go func(conn *websocket.Conn) {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			log.Println(string(msg))

		}
	}(conn)
}
