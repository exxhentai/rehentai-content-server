package comm

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	ipfs "github.com/ipfs/go-ipfs-api"
)

type Server struct {
	shell *ipfs.Shell
}

func NewServer() *Server {
	sh := ipfs.NewLocalShell()

	if sh == nil {
		log.Fatalln("[ERROR]", "You need to run a ipfs node first!!!")
		return nil
	} else {
		return &Server{sh}
	}
}

func (s Server) Start(port int) {
	log.Println("[INFO]", "Starting rehentai content server")

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			reader, err := r.MultipartReader()
			if err != nil {
				log.Println("[ERROR]", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}

				if part.FileName() == "" {
					continue
				}

				fp, err := os.Create(path.Join("/tmp", part.FileName()))
				if err != nil {
					log.Println("[ERROR]", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if _, err := io.Copy(fp, part); err != nil {
					log.Println("[ERROR]", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		log.Println("[INFO]", "File from", r.RemoteAddr, "upload successful")
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
