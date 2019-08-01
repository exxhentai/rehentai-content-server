package comm

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

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

func (s Server) upload(filenames []string, w http.ResponseWriter) {
	dir := os.TempDir()

	hash_ch := make(chan string, len(filenames))

	for _, fn := range filenames {
		go func(filename string, hash_ch chan<- string) {
			folder, err := decompress(dir, filename)
			if err != nil {
				log.Println("[ERROR]", err)
				hash_ch <- ""
				return
			}

			s.shMux.Lock()
			hash, err := s.shell.AddDir(folder)
			s.shMux.Unlock()

			if err != nil {
				log.Println("[ERROR]", err)
				hash_ch <- ""
				return
			}

			hashAndFolder := hash + " " + path.Base(folder)

			// DEBUG
			fmt.Println(hashAndFolder)

			hash_ch <- hashAndFolder

			// remove uploaded folder
			os.RemoveAll(folder)
		}(fn, hash_ch)
	}

	for i := 0; i < len(filenames); i++ {
		if hash := <-hash_ch; hash != "" {
			w.Write([]byte(hash + "\n"))
		}
	}
}

func decompress(dir, filename string) (string, error) {
	re := regexp.MustCompile("\\.(png|jpg|jpeg|bmp)\\z")

	reader, err := zip.OpenReader(path.Join(dir, filename))
	if err != nil {
		return "", err
	}

	defer reader.Close()

	var baseFolder string

	for _, f := range reader.File {
		filepath := path.Join(dir, f.Name)

		if !strings.HasPrefix(filepath, path.Clean(dir)+string(os.PathSeparator)) {
			err = errors.New("illegal file path: " + filepath)
			return "", err
		}

		// create folder
		if f.FileInfo().IsDir() {
			os.MkdirAll(filepath, os.ModePerm)

			if l := len(baseFolder); l == 0 || len(filepath) < l {
				baseFolder = filepath
			}

			continue
		}

		// exclude non-image extension
		if !re.MatchString(path.Ext(f.Name)) {
			continue
		}

		// create file
		err = os.MkdirAll(path.Dir(filepath), os.ModePerm)
		if err != nil {
			return "", err
		}

		dp, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", err
		}

		sp, err := f.Open()
		if err != nil {
			return "", err
		}

		_, err = io.Copy(dp, sp)

		dp.Close()
		sp.Close()

		if err != nil {
			return "", err
		}
	}

	// remove decompressed zip
	os.Remove(path.Join(dir, filename))

	return baseFolder, nil
}
