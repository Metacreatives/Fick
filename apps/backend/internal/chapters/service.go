package chapters

import (
	"encoding/json"
	"fick/backend/internal/store"
)

func findChapterByWorkAndNumber(workID, chapterNumber string) (Chapter, bool) {
	var chapters []Chapter

	if err := json.Unmarshal([]byte(store.Chapters), &chapters); err != nil {
		return Chapter{}, false
	}

	for _, chapter := range chapters {
		if chapter.WorkID == workID && chapter.Number == chapterNumber {
			return chapter, true
		}
	}

	return Chapter{}, false
}
