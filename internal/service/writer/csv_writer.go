// internal/service/writer/csv_writer.go - Updated CSV writer with full domain export
package writer

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/twiny/spidy/v2/internal/pkg/spider/v1"
)

type CSVWriter struct {
	l *sync.Mutex
	f *os.File
	w *csv.Writer
}

func NewCSVWriter(dir string) (*CSVWriter, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	fp := filepath.Join(dir, time.Now().Format("2006-01-02")+"_domains.csv")
	f, err := os.Create(fp)
	if err != nil {
		return nil, err
	}

	w := csv.NewWriter(f)
	if err := w.Write([]string{"URL", "Name", "TLD", "Status"}); err != nil {
		return nil, err
	}

	return &CSVWriter{
		l: &sync.Mutex{},
		f: f,
		w: w,
	}, nil
}

func (c *CSVWriter) Write(d *spider.Domain) error {
	c.l.Lock()
	defer func() {
		c.w.Flush()
		c.l.Unlock()
	}()
	return c.w.Write(d.CSVRow())
}

func (c *CSVWriter) Close() error {
	return c.f.Close()
}
