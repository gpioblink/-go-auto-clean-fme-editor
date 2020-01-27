package main

import (
	"github.com/go-chi/chi"
	"github.com/gpioblink/go-auto-clean-fme-editor/pkg/common/cmd"
	editorApp "github.com/gpioblink/go-auto-clean-fme-editor/pkg/editor/application"
	editorInfraLyric "github.com/gpioblink/go-auto-clean-fme-editor/pkg/editor/infrastructure/lyric"
	editorInterfaceLyric "github.com/gpioblink/go-auto-clean-fme-editor/pkg/editor/interfaces/private/intraprocess"
	editorInterfaceHttp "github.com/gpioblink/go-auto-clean-fme-editor/pkg/editor/interfaces/public/http"
	fmeApp "github.com/gpioblink/go-auto-clean-fme-editor/pkg/fme/application"
	fmeInfraFme "github.com/gpioblink/go-auto-clean-fme-editor/pkg/fme/infrastructure/fme"
	fmeInfraLyric "github.com/gpioblink/go-auto-clean-fme-editor/pkg/fme/infrastructure/lyric"
	fmeInterfaceHttp "github.com/gpioblink/go-auto-clean-fme-editor/pkg/fme/interfaces/public/http"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting api server")
	ctx := cmd.Context()

	router := createMonolith()

	go func() {
		if err := http.ListenAndServe(":5814", router); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Monolith is listening on 5814")

	<-ctx.Done()
	log.Println("Closing monolith")

	// TODO: httpサーバーの終了処理 (go-chiの使い方変わった？

}

func createMonolith() *chi.Mux {
	editorRepo := editorInfraLyric.NewMemoryRepository()
	editorService := editorApp.NewLyricService(editorRepo, editorRepo)
	editorIntraprocessInterface := editorInterfaceLyric.NewLyricsInterface(editorService)

	fmeRepo := fmeInfraFme.NewMemoryRepository()
	fmeService := fmeApp.NewFmeService(fmeInfraLyric.NewIntraprocessService(editorIntraprocessInterface), fmeRepo)

	r := cmd.CreateRouter()
	editorInterfaceHttp.AddRoutes(r, editorService, editorRepo)
	fmeInterfaceHttp.AddRoutes(r, fmeService, fmeRepo)

	return r
}
