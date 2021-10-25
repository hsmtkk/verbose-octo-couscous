package main

import (
	"log"
	"sync"

	"github.com/hsmtkk/verbose-octo-couscous/env"
	"github.com/hsmtkk/verbose-octo-couscous/html"
	"github.com/hsmtkk/verbose-octo-couscous/http"
	"github.com/hsmtkk/verbose-octo-couscous/work"
	"github.com/spf13/cobra"
)

const parallelNumber = 4

var command = &cobra.Command{
	Use:  "en-photo-thumbnail url",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(args[0])
	},
}

func main() {
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(url string) {
	username, password, err := env.GetUsernamePassword()
	if err != nil {
		log.Fatal(err)
	}
	accessor, err := http.New()
	if err != nil {
		log.Fatal(err)
	}
	if err := accessor.Login(username, password); err != nil {
		log.Fatal(err)
	}
	albumHTML, err := accessor.GetAlbum(url)
	if err != nil {
		log.Fatal(err)
	}
	dataSrcs, err := html.SelectDataSrc(albumHTML)
	if err != nil {
		log.Fatal(err)
	}

	dataSrcChan := make(chan string)
	var dataSrcWait sync.WaitGroup

	for i := 0; i < parallelNumber; i++ {
		worker := work.New(i, accessor, dataSrcChan)
		dataSrcWait.Add(1)
		go func() {
			defer dataSrcWait.Done()
			worker.Run()
		}()
	}

	for _, dataSrc := range dataSrcs {
		dataSrcChan <- dataSrc
	}
	close(dataSrcChan)

	dataSrcWait.Wait()
}
