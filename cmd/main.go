package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/iamhectorosa/iamhectorsosa/internal/headless"
	"github.com/iamhectorosa/iamhectorsosa/internal/templates"
)

const (
	RootPath   = ".dist"
	AssetsPath = "assets"
	Addr       = "3000"
)

func main() {
	mode := flag.String("mode", "", "Mode of operation: 'build' to build static files.")
	flag.Parse()

	switch *mode {
	case "build":
		build()
	default:
		build()
		serve()
	}
}

func build() {
	if err := os.MkdirAll(RootPath, 0755); err != nil {
		log.Fatalf("failed to create output directory\n%v", err)
	}

	name := path.Join(RootPath, "index.html")
	file, err := os.Create(name)
	if err != nil {
		log.Fatalf("failed to create output file\n%v", err)
	}
	defer file.Close()

	templs := templates.New()

	config, err := templates.ParseYaml("data.yml")
	if err != nil {
		log.Fatalf("error reading YAML, err=%v", err)
	}

	if err := templs.Render(file, "home", config); err != nil {
		log.Fatalf("couldn't render file\n%v", err)
	}

	faviconName := path.Join(RootPath, "favicon.svg")
	faviconFile, err := os.Create(faviconName)
	if err != nil {
		log.Fatalf("failed to create output file\n%v", err)
	}
	defer faviconFile.Close()

	if err := templs.Render(faviconFile, "favicon", struct{}{}); err != nil {
		log.Fatalf("couldn't render file\n%v", err)
	}

	tmpFile, err := os.CreateTemp("", "og-*.html")
	if err != nil {
		log.Fatalf("failed to create temp file\n%v", err)
	}
	defer os.Remove(tmpFile.Name())

	cssPath, err := filepath.Abs(path.Join(RootPath, "css", "tailwind.css"))
	if err != nil {
		log.Fatalf("failed to get absolute path\n%v", err)
	}

	ogConfig := struct {
		CSSPath string
		templates.Config
	}{
		"file://" + cssPath,
		config,
	}

	if err := templs.Render(tmpFile, "opengraph-image", ogConfig); err != nil {
		log.Fatalf("couldn't render file\n%v", err)
	}
	defer tmpFile.Close()

	buf, err := headless.Capture(tmpFile.Name(), 1200, 640, 2.0)
	if err != nil {
		log.Fatalf("failed to capture screenshot\n%v", err)
	}

	if err := os.WriteFile(path.Join(AssetsPath, "opengraph-image.png"), buf, 0644); err != nil {
		log.Fatalf("failed to save screenshot\n%v", err)
	}
}

func serve() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(RootPath)))

	assetsPath := fmt.Sprintf("/%s/", AssetsPath)
	mux.Handle(assetsPath, http.StripPrefix(assetsPath, http.FileServer(http.Dir(AssetsPath))))

	fmt.Printf(
		"Listening on http://localhost:%s\n",
		Addr,
	)

	http.ListenAndServe(":"+Addr, mux)
}
