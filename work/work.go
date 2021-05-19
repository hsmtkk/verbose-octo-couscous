package work

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hsmtkk/verbose-octo-couscous/http"
)

type Worker interface {
	Run()
}

type workerImpl struct {
	id          int
	accessor    http.Accessor
	dataSrcChan <-chan string
}

func New(id int, accessor http.Accessor, dataSrcChan <-chan string) Worker {
	return &workerImpl{id: id, accessor: accessor, dataSrcChan: dataSrcChan}
}

func (w *workerImpl) Run() {
	for dataSrc := range w.dataSrcChan {
		fmt.Printf("worker %d is handling %s\n", w.id, dataSrc)

		thumbnailURL, err := w.accessor.GetDataSrc(dataSrc)
		if err != nil {
			log.Printf("failed to get data-src; %s; %s", dataSrc, err.Error())
			continue
		}
		thumbnailBytes, err := w.accessor.GetThumbnail(thumbnailURL)
		if err != nil {
			log.Printf("failed to get thumbnail; %s; %s", thumbnailURL, err.Error())
			continue
		}
		name, err := w.getFileName(thumbnailURL)
		if err != nil {
			log.Printf("failed to get thumbnail name; %s; %s", thumbnailURL, err.Error())
			continue
		}
		path := filepath.Join("photo", name)
		if err := os.WriteFile(path, thumbnailBytes, 0644); err != nil {
			log.Printf("failed to save thumbnail; %s; %s", path, err.Error())
			continue
		}
	}
}

func (w *workerImpl) getFileName(thumnailURL string) (string, error) {
	parsed, err := url.Parse(thumnailURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL; %s; %w", thumnailURL, err)
	}
	elems := strings.Split(parsed.Path, "/")
	name := elems[len(elems)-1]
	return name, nil

}
