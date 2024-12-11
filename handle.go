package sftpd

import (
	"io"
	"strconv"
)

type FileOpenArgs struct {
	name  string
	flags uint32
	attr  *Attr
}

type DirReader struct {
	attrs []NamedAttr
	pos   int
}

func (d *DirReader) Readdir(count int) ([]NamedAttr, error) {
	if d.pos >= len(d.attrs) {
		return nil, io.EOF
	}
	var nextPos int
	if d.pos+count > len(d.attrs) {
		nextPos = len(d.attrs)
	} else {
		nextPos = d.pos + count
	}
	ret := d.attrs[d.pos:nextPos]
	d.pos = nextPos
	return ret, nil
}

func (d *DirReader) Close() error {
	return nil
}

type handles struct {
	f  map[string]*FileOpenArgs
	d  map[string]string
	fw map[string]io.WriteCloser
	fr map[string]io.ReadCloser
	dr map[string]Dir
	c  int64
}

func (h *handles) init() {
	h.f = map[string]*FileOpenArgs{}
	h.d = map[string]string{}
}

func (h *handles) closeHandle(k string) {
	if k == "" {
		return
	}
	if k[0] == 'f' {
		delete(h.f, k)
		if c, ok := h.fw[k]; ok {
			_ = c.Close()
			delete(h.fw, k)
		}
		if c, ok := h.fr[k]; ok {
			_ = c.Close()
			delete(h.fr, k)
		}
	} else if k[0] == 'd' {
		delete(h.d, k)
		if c, ok := h.dr[k]; ok {
			_ = c.Close()
			delete(h.dr, k)
		}
	}
}

func (h *handles) nfiles() int { return len(h.f) }

func (h *handles) ndir() int { return len(h.d) }

func (h *handles) newFile(f *FileOpenArgs) string {
	h.c++
	k := "f" + strconv.FormatInt(h.c, 16)
	h.f[k] = f
	return k
}
func (h *handles) newDir(f string) string {
	h.c++
	k := "d" + strconv.FormatInt(h.c, 16)
	h.d[k] = f
	return k
}
func (h *handles) getFile(n string) *FileOpenArgs {
	return h.f[n]
}
func (h *handles) getDir(n string) string {
	return h.d[n]
}
