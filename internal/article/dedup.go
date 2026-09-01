package article

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileDedup struct {
	seen map[string]struct{}
}

func NewFileDedup(sourceName string) (*FileDedup, error) {
	fileDedup := &FileDedup{seen: make(map[string]struct{})}

	sourceJsonlPath := filepath.Join("internal", "data", "pool", sourceName, "*.jsonl")
	sourceJsonlFiles, err := filepath.Glob(sourceJsonlPath)
	if err != nil {
		return nil, fmt.Errorf("[NewFileDedup()]:sourceJsonlFiles文件获取失败：%w\n", err)
	}

	for _, sourceJsonlFile := range sourceJsonlFiles {

		sourceJsonlFileHanlder, err := os.Open(sourceJsonlFile)
		if err != nil {
			return nil, fmt.Errorf("[NewFileDedup()]:打开sourceJsonlFile文件失败：%w\n", err)
		}

		decoder := json.NewDecoder(sourceJsonlFileHanlder)

		for {
			var line struct {
				Guid string `json:"guid"`
			}

			err = decoder.Decode(&line)
			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, err
			}

			if line.Guid != "" {
				fileDedup.seen[line.Guid] = struct{}{}
			}
		}
		sourceJsonlFileHanlder.Close()
	}

	return fileDedup, nil
}

func (FileDedup *FileDedup) Filter(feedarticles []FeedArticle) []FeedArticle {
	var fresh []FeedArticle

	for _, feedArticle := range feedarticles {
		ok := FileDedup.Seen(feedArticle.Guid)
		if ok {
			continue
		}

		FileDedup.Mark(feedArticle.Guid)
		fresh = append(fresh, feedArticle)
	}

	return fresh
}

func (FileDedup *FileDedup) Seen(key string) bool {
	_, ok := FileDedup.seen[key]

	return ok
}

func (FileDedup *FileDedup) Mark(key string) {
	FileDedup.seen[key] = struct{}{}
}
