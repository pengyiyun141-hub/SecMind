package article

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileDedup struct {
	guidseen map[string]struct{}
	linkseen	 map[string]struct{}
}

func NewFileDedup(sourceName string) (*FileDedup, error) {
	fileDedup := &FileDedup{
		guidseen: make(map[string]struct{}),
		linkseen: make(map[string]struct{}),
	}

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
				Link string `json:"link"`
			}

			err = decoder.Decode(&line)
			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, err
			}

			if line.Guid != "" {
				fileDedup.guidseen[line.Guid] = struct{}{}
			}
			if line.Link != "" {
    			fileDedup.linkseen[line.Link] = struct{}{}
			}
		}
		sourceJsonlFileHanlder.Close()
	}

	return fileDedup, nil
}

func (FileDedup *FileDedup) Filter(feedarticles []FeedArticle) []FeedArticle {
	var fresh []FeedArticle

	for _, feedArticle := range feedarticles {
		ok := FileDedup.Seen(feedArticle.Guid, feedArticle.Link)
		if ok {
			continue
		}

		FileDedup.Mark(feedArticle.Guid, feedArticle.Link)
		fresh = append(fresh, feedArticle)
	}

	return fresh
}

func (FileDedup *FileDedup) Seen(guidkey string, linkkey string) bool {
	if guidkey != "" {
        if _, ok := FileDedup.guidseen[guidkey] 
		ok {
            return true
        }
    }

    if linkkey != "" {
        if _, ok := FileDedup.linkseen[linkkey]
		ok {
            return true
        }
    }

	return false
}

func (FileDedup *FileDedup) Mark(guidkey string, linkkey string) {
	 if guidkey != "" {
        FileDedup.guidseen[guidkey] = struct{}{}
    }
    if linkkey != "" {
        FileDedup.linkseen[linkkey] = struct{}{}
    }
}
