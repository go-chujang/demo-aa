package logx

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-chujang/demo-aa/common/utils/ternary"
)

const (
	mb                   = 1024 * 1024
	defaultMaxFileCount  = 50
	defaultMaxFileSizeMb = 10
)

var _ io.Writer = (*FileLogger)(nil)

type FileLogger struct {
	mu sync.Mutex

	fp    *os.File
	mode  fs.FileMode
	date  time.Time
	size  int64
	count int

	filePath string
	dir      string
	base     string
	ext      string
	prefix   string

	maxFileCount int
	maxFileSize  int64
}

func New(logPath string, maxFileCount, maxFileSizeMb int) (*FileLogger, error) {
	if logPath == "" {
		return nil, fmt.Errorf("empty logPath")
	}
	var (
		dir    = filepath.Dir(logPath)
		base   = filepath.Base(logPath)
		ext    = filepath.Ext(base)
		prefix = base[:len(base)-len(ext)]
		logger = &FileLogger{
			filePath:     logPath,
			mode:         fs.FileMode(0644),
			maxFileCount: ternary.Cond(maxFileCount == 0, defaultMaxFileCount, maxFileCount),
			maxFileSize:  int64(ternary.Cond(maxFileSizeMb == 0, defaultMaxFileSizeMb, maxFileSizeMb)),
			dir:          dir,
			base:         base,
			ext:          ext,
			prefix:       prefix,
		}
	)
	if err := os.MkdirAll(logger.dir, fs.FileMode(0744)); err != nil {
		return nil, err
	}
	if err := logger.newFile(); err != nil {
		return nil, err
	}
	return logger, nil
}

func (w *FileLogger) Write(p []byte) (n int, err error) {
	inputLen := int64(len(p))

	// ~~~ first round ~~~
	// pre-process
	switch {
	case w == nil, inputLen == 0: // skip write
		return 0, nil
	case inputLen > w.maxFileSize: // oversize log
		return 0, fmt.Errorf("log length : %d, max file size : %d", inputLen, w.maxFileSize)
	case w.fp == nil:
		if err = w.newFile(); err != nil {
			return 0, err
		}
	}

	// ~~~ second round ~~~
	// check-rotate
	w.mu.Lock()
	defer w.mu.Unlock()

	now := logDateYmd(time.Now())
	bRotate := false
	if now.After(w.date) || inputLen+w.size > w.maxFileSize {
		bRotate = true
	}

	// ~~~ final round ~~~
	// rotate, wirte
	if bRotate {
		w.autoRemove()
		if err = w.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = w.fp.Write(p)
	if err == nil {
		w.size += int64(n)
		w.date = now
	}
	return
}

func (w *FileLogger) fileList() []string {
	var files []string
	matches, _ := filepath.Glob(w.dir + "/" + w.prefix + "*" + w.ext)
	for _, s := range matches {
		fi, _ := os.Stat(s)

		if !fi.IsDir() {
			files = append(files, s)
		}
	}
	return files
}

func (w *FileLogger) newFile() error {
	var err error
	w.fp, err = os.OpenFile(w.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, w.mode)
	if err != nil {
		return err
	}
	fi, err := w.fp.Stat()
	if err != nil {
		return err
	}
	w.size = fi.Size()
	w.count = len(w.fileList())
	w.date = logDateYmd(fi.ModTime())
	return nil
}

func (w *FileLogger) autoRemove() {
	fl := w.fileList()
	count := len(fl)
	if w.maxFileCount == 0 || w.maxFileCount >= count {
		return
	}
	os.Remove(fl[0])
}

func (w *FileLogger) rotateFileName() string {
	t := logDateYmd(time.Now())
	date := subDate(t)

	count := 0
	fl := w.fileList()
	for _, f := range fl {
		if strings.Contains(f, date) {
			count = count + 1
		}
	}
	return filepath.Join(
		w.dir,
		fmt.Sprintf("%s.%s.%d%s",
			w.prefix, date, count, w.ext,
		))
}

func (w *FileLogger) close() {
	if w.fp != nil {
		w.fp.Close()
	}
}

func (w *FileLogger) rotate() error {
	if _, err := os.Stat(w.filePath); err != nil {
		return err
	}
	w.close()

	if err := os.Rename(w.filePath, w.rotateFileName()); err != nil {
		return err
	}
	return w.newFile()
}

func subDate(date time.Time) string {
	y, m, d := date.AddDate(0, 0, -1).Date()
	return fmt.Sprintf("%d-%02d-%02d", y, m, d)
}

func logDateYmd(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}
