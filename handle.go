package sftpd

import "strconv"

type FileOpenArgs struct {
	name  string
	flags uint32
	attr  *Attr
}

type handles struct {
	f map[string]*FileOpenArgs
	d map[string]string
	c int64
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
	} else if k[0] == 'd' {
		delete(h.d, k)
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
