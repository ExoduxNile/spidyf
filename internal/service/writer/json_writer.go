// internal/service/writer/json_writer.go - New JSON writer implementation
package writer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/twiny/spidy/v2/internal/pkg/spider/v1"
)

type JSONWriter struct {
	l    *sync.Mutex
	f    *os.File
	enc  *json.Encoder
}

func NewJSONWriter(dir string) (*JSONWriter, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	fp := filepath.Join(dir, time.Now().Format("2006-01-02")+"_domains.json")
	f, err := os.Create(fp)
	if err != nil {
		return nil, err
	}

	return &JSONWriter{
		l:   &sync.Mutex{},
		f:   f,
		enc: json.NewEncoder(f),
	}, nil
}

func (j *JSONWriter) Write(d *spider.Domain) error {
	j.l.Lock()
	defer j.l.Unlock()
	return j.enc.Encode(d)
}

func (j *JSONWriter) Close() error {
	return j.f.Close()
}
