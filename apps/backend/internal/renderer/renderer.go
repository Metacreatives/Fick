package renderer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type RenderResult struct {
	HTML       string
	StatusCode int
}

type renderRequest struct {
	URL string `json:"url"`
}

type renderResponse struct {
	HTML       string `json:"html"`
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
}

type Renderer struct {
	mu sync.Mutex

	cmd     *exec.Cmd
	stdin   io.WriteCloser
	encoder *json.Encoder
	decoder *json.Decoder

	closed bool
}

func Start(scriptPath, apiOrigin string) (*Renderer, error) {
	cmd := exec.Command("node", scriptPath, apiOrigin)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("open renderer stdin: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdin.Close()
		return nil, fmt.Errorf("open renderer stdout: %w", err)
	}

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return nil, fmt.Errorf("start renderer: %w", err)
	}

	return &Renderer{
		cmd:     cmd,
		stdin:   stdin,
		encoder: json.NewEncoder(stdin),
		decoder: json.NewDecoder(stdout),
	}, nil
}

func (r *Renderer) Render(url string) (RenderResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return RenderResult{}, errors.New("renderer is closed")
	}

	if url == "" {
		return RenderResult{}, errors.New("render URL is empty")
	}

	if err := r.encoder.Encode(renderRequest{
		URL: url,
	}); err != nil {
		return RenderResult{}, fmt.Errorf(
			"send render request: %w",
			err,
		)
	}

	var response renderResponse

	if err := r.decoder.Decode(&response); err != nil {
		return RenderResult{}, fmt.Errorf(
			"read render response: %w",
			err,
		)
	}

	if response.Error != "" {
		return RenderResult{}, fmt.Errorf(
			"renderer: %s",
			response.Error,
		)
	}

	if response.StatusCode < 100 || response.StatusCode > 599 {
		return RenderResult{}, fmt.Errorf(
			"renderer returned invalid status code %d",
			response.StatusCode,
		)
	}

	return RenderResult{
		HTML:       response.HTML,
		StatusCode: response.StatusCode,
	}, nil
}

func (r *Renderer) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}

	r.closed = true

	closeError := r.stdin.Close()
	waitError := r.cmd.Wait()

	return errors.Join(closeError, waitError)
}
