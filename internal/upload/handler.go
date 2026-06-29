package upload

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

const maxUploadSize = 2 << 30

type Handler struct {
	manager *Manager
}

func NewHandler(manager *Manager) *Handler {
	return &Handler{manager: manager}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/upload", h.Upload)
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB in memory, rest on disk
		http.Error(w, "file too large or bad request", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := filepath.Base(header.Filename)
	savedPath, err := h.manager.Save(filename, file)

	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	hlsDir := filepath.Join(filepath.Dir(savedPath), "hls")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"path":     savedPath,
		"filename": filename,
	})

	go func() {
		if err := transcodeToHLS(savedPath, hlsDir); err != nil {
			log.Println("transcode error", err)
		}
	}()

}

func transcodeToHLS(inputPath, outputDir string) error {
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-hls_time", "6",
		"-hls_list_size", "0",
		"-f", "hls",
		filepath.Join(outputDir, "index.m3u8"),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
