package article

import ()//"path/filepath"

type FileDedup struct {
	seen map[string]struct{}
}

/*func NewFileDedup(sourceName string) (*FileDedup, error) {
	fileDedup := &FileDedup{seen: make(map[string]struct{})}

	sourceJsonlPath := filepath.Join("internal", "data", "pool", sourceName, "*.jsonl")
	sourceJsonlFiles, err := filepath.Glob(sourceJsonlPath)
	if err != nil {
		return nil, fmt.Errorf("NewFileDedup(),sourceJsonlFiles文件获取失败：%w\n", err)
	}

	for _, sourceJsonlFile := range sourceJsonlFiles {

	}	


}

func (FileDedup *FileDedup) Filter(feedarticles []FeedArticle) ([]FeedArticle) {
	
}

func (FileDedup *FileDedup) Seen(key string) (bool) {

}

func (FileDedup *FileDedup) Mark(key string) {

}*/