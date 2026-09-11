package frontend

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
)

func buildVersion(
	template string,
	rendererBundle []byte,
) string {
	hash := sha256.New()

	_, _ = hash.Write([]byte(template))
	_, _ = hash.Write([]byte{0})
	_, _ = hash.Write(rendererBundle)

	sum := hash.Sum(nil)

	return hex.EncodeToString(sum[:8])
}

func writeHTML(
	w http.ResponseWriter,
	statusCode int,
	document []byte,
) {
	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.Header().Set(
		"Cache-Control",
		"no-cache",
	)

	w.WriteHeader(statusCode)

	_, _ = w.Write(document)
}

func IsArtifactRoute(
	route string,
	rawQuery string,
) bool {
	if rawQuery != "" {
		return false
	}

	parts := strings.Split(
		strings.Trim(route, "/"),
		"/",
	)

	if len(parts) == 2 &&
		parts[0] == "works" &&
		parts[1] != "" {
		return true
	}

	if len(parts) == 4 &&
		parts[0] == "works" &&
		parts[1] != "" &&
		parts[2] == "chapters" &&
		parts[3] != "" {
		return true
	}

	return false
}
