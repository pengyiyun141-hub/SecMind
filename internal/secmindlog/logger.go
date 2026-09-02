package secmindlog

import (
	"fmt"
	"path/filepath"
	"log/slog"
	"os"
	"secmind/configs"
)

type Options struct {
	Level  string
	Format string
	Output string
	File   string 
}

func Init(opts configs.LogConfigs) error {

	err := os.MkdirAll(filepath.Dir(opts.File), 0755)
	if err != nil {
		return fmt.Errorf("创建目录失败：%v", err)
	}

	logFile, err := os.OpenFile(opts.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("打开文件%s失败：%v", opts.File, err)
	}

	handler := slog.NewJSONHandler(logFile, nil)
	Logger := slog.New(handler)
	slog.SetDefault(Logger)

	return nil
}
