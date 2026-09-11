package frontend

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const ssrOutlet = "<!--ssr-outlet-->"

type Handler struct {
	renderer   *renderer
	template   string
	fileServer http.Handler
	clientDir  string
	artifacts  *artifactStore
}

func FrontendHandler(
	rendererProcess *renderer,
	clientDirectory string,
	rendererBundlePath string,
	artifactDirectory string,
) (*Handler, error) {
	indexPath := path.Join(clientDirectory, "index.html")

	content, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("read frontend template: %w", err)
	}

	template := string(content)

	if !strings.Contains(template, ssrOutlet) {
		return nil, fmt.Errorf("frontend template is missing SSR outlet")
	}

	rendererBundle, err := os.ReadFile(
		rendererBundlePath,
	)
	if err != nil {
		return nil, err
	}

	version := buildVersion(
		template,
		rendererBundle,
	)

	return &Handler{
		renderer:   rendererProcess,
		template:   template,
		clientDir:  clientDirectory,
		fileServer: http.FileServer(http.Dir(clientDirectory)),
		artifacts: newArtifactStore(artifactDirectory,
			version),
	}, nil
}

func (h *Handler) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	filePath := path.Join(
		h.clientDir,
		path.Clean(r.URL.Path),
	)

	if info, err := os.Stat(filePath); err == nil && info.Mode().IsRegular() {
		h.fileServer.ServeHTTP(w, r)
		return
	}

	h.serveRenderedPage(w, r)
}

func (h *Handler) serveRenderedPage(
	w http.ResponseWriter,
	r *http.Request,
) {
	artifactRoute := path.Clean(r.URL.Path)

	cacheable := IsArtifactRoute(
		artifactRoute,
		r.URL.RawQuery,
	)

	if cacheable {
		artifact, found, err := h.artifacts.read(
			artifactRoute,
		)

		if err != nil {
			http.Error(
				w,
				"Could not read rendered artifact",
				http.StatusInternalServerError,
			)
			return
		}

		if found {
			writeHTML(
				w,
				http.StatusOK,
				artifact,
			)
			return
		}
	}

	result, err := h.renderer.Render(
		r.URL.RequestURI(),
	)

	if err != nil {
		http.Error(
			w,
			"Could not render page",
			http.StatusInternalServerError,
		)
		return
	}

	document := strings.Replace(
		h.template,
		ssrOutlet,
		result.HTML,
		1,
	)

	if cacheable && result.StatusCode == http.StatusOK {
		if err := h.artifacts.write(
			artifactRoute,
			document,
		); err != nil {
			log.Printf(
				"could not store render artifact: %v",
				err,
			)
		}
	}

	writeHTML(
		w,
		result.StatusCode,
		[]byte(document),
	)
}
