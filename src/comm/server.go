package comm

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
		case "POST":
			filenames := s.receive(w, r)
			fmt.Println(filenames)
			s.decompress(filenames)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (s Server) receive(w http.ResponseWriter, r *http.Request) (filenames []string) {
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

		if filename := part.FileName(); filename == "" {
			continue
		} else if ext := path.Ext(filename); ext != ".zip" {
			http.Error(w, "invalid file extension ("+ext+")", http.StatusUnprocessableEntity)
			return
		} else if !func() bool {
			for i := range filenames {
				if filenames[i] == filename {
					return true
				}
			}

			return false
		}() {
			filenames = append(filenames, filename)
		}

		fp, err := os.Create(path.Join(os.TempDir(), part.FileName()))
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

	return
}

func (s Server) decompress(filenames []string) {
	dest := os.TempDir()

	for _, fn := range filenames {
		go func(filename string) {
			reader, err := zip.OpenReader(path.Join(dest, filename))
			if err != nil {
				log.Println("[ERROR]", err)
				// TODO: return fail message to user
				return
			}

			defer reader.Close()

			for _, f := range reader.File {
				filepath := path.Join(dest, f.Name)

				if !strings.HasPrefix(filepath, path.Clean(dest)+string(os.PathSeparator)) {
					log.Println("[ERROR]", "illegal file path:", filepath)
					// TODO: return fail message to user
					return
				}

				// create folder
				if f.FileInfo().IsDir() {
					fmt.Println(f.Name)
					os.MkdirAll(filepath, os.ModePerm)
					continue
				}

				// create file
				if err := os.MkdirAll(path.Dir(filepath), os.ModePerm); err != nil {
					log.Println("[ERROR]", err)
					// TODO: return fail message to user
					return
				}

				dp, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					log.Println("[ERROR]", err)
					// TODO: return fail message to user
					return
				}

				sp, err := f.Open()
				if err != nil {
					log.Println("[ERROR]", err)
					// TODO: return fail message to user
					return
				}

				_, err = io.Copy(dp, sp)

				dp.Close()
				sp.Close()

				if err != nil {
					log.Println("[ERROR]", err)
					// TODO: return fail message to user
					return
				}
			}
		}(fn)
	}
}
