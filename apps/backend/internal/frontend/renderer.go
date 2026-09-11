package frontend

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type renderResult struct {
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

type renderer struct {
	mu sync.Mutex

	cmd     *exec.Cmd
	stdin   io.WriteCloser
	encoder *json.Encoder
	decoder *json.Decoder

	closed bool
}

func Start(scriptPath, apiOrigin string) (*renderer, error) {
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

	return &renderer{
		cmd:     cmd,
		stdin:   stdin,
		encoder: json.NewEncoder(stdin),
		decoder: json.NewDecoder(stdout),
	}, nil
}

func (r *renderer) Render(url string) (renderResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return renderResult{}, errors.New("renderer is closed")
	}

	if url == "" {
		return renderResult{}, errors.New("render URL is empty")
	}

	if err := r.encoder.Encode(renderRequest{
		URL: url,
	}); err != nil {
		return renderResult{}, fmt.Errorf(
			"send render request: %w",
			err,
		)
	}

	var response renderResponse

	if err := r.decoder.Decode(&response); err != nil {
		return renderResult{}, fmt.Errorf(
			"read render response: %w",
			err,
		)
	}

	if response.Error != "" {
		return renderResult{}, fmt.Errorf(
			"renderer: %s",
			response.Error,
		)
	}

	if response.StatusCode < 100 || response.StatusCode > 599 {
		return renderResult{}, fmt.Errorf(
			"renderer returned invalid status code %d",
			response.StatusCode,
		)
	}

	return renderResult{
		HTML:       response.HTML,
		StatusCode: response.StatusCode,
	}, nil
}

func (r *renderer) Close() error {
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
