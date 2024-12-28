package dig

import (
	"context"
	"crypto/sha256"
	"io"

	"github.com/hsfzxjy/sdxtra/internal/dig/helper"
)

type Writer interface {
	helper.Helper
	io.Writer
	Context() context.Context
	CopyFromFile(filepath string) error
	new() Writer
}

type writer struct {
	filecache *fileHashCache
	ctx       context.Context
	helper.H
}

func newWriter(ctx context.Context) writer {
	return writer{
		ctx:       ctx,
		H:         helper.H{Hash: sha256.New()},
		filecache: fileHashCacheInstance,
	}
}

func NewWriter(ctx context.Context) Writer {
	w := newWriter(ctx)
	return &w
}

func (w *writer) new() Writer {
	n := newWriter(w.ctx)
	return &n
}

func (w *writer) Context() context.Context {
	return w.ctx
}

func (w *writer) CopyFromFile(filepath string) error {
	hasher := sha256.New()
	hsh, err := w.filecache.Get(filepath, hasher)
	if err != nil {
		return err
	}
	w.Hash.Write(hsh)
	return nil
}
