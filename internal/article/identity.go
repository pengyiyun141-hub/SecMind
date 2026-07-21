package article

import (
	"fmt"
	"time"
	"crypto/sha256"
	
)

func GenerateFileName(articleinfo ScreenedArticle)(string) {
	timeStamp := time.Now().UTC().Format("20060102")
	LinkHash := sha256.Sum256([]byte(articleinfo.Link))

	return fmt.Sprintf("%s-%s-%d-%x", timeStamp, articleinfo.Source, articleinfo.ID, LinkHash)
}