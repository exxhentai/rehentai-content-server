package comm

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	ipfs "github.com/ipfs/go-ipfs-api"
)

type Server struct {
	shell *ipfs.Shell
	shMux *sync.Mutex
}

func NewServer() *Server {
	sh := ipfs.NewLocalShell()

	if sh == nil {
		log.Fatalln("[ERROR]", "You need to run a ipfs node first!!!")
		return nil
	} else {
		return &Server{sh, &sync.Mutex{}}
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

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
