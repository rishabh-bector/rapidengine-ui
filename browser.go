package main

import (
	"io/ioutil"
	"nuklear-golang/nk"
	"os"
)

type FileBrowser struct {
	Dir string

	fileDir string
	Files   []os.FileInfo

	Chosen func(string)

	Active bool
}

func NewFileBrowser() FileBrowser {
	cdir, _ := os.Getwd()
	return FileBrowser{
		Dir: cdir,
	}
}

func (fb *FileBrowser) Update() {
	if !fb.Active {
		return
	}

	fb.ReloadFiles()

	bounds := nk.NkRect(500, 200, 600, 800)
	if nk.NkBegin(ctx, "Browser", bounds, nk.WindowMovable|nk.WindowBorder|nk.WindowClosable) != 1 {
		nk.NkEnd(ctx)
		fb.Active = false
		return
	}

	for _, f := range fb.Files {
		nk.NkLayoutRowDynamic(ctx, 60, 1)
		nk.NkGroupBegin(ctx, f.Name(), nk.WindowNoScrollbar|nk.WindowBorder)
		nk.NkLayoutRowDynamic(ctx, 60, 1)

		if nk.NkButtonLabel(ctx, f.Name()) == 1 {
			if f.IsDir() {
				fb.Dir = f.Name()
			} else {
				fb.Chosen(f.Name())
				fb.Reset()
			}
		}

		nk.NkGroupEnd(ctx)
	}

	nk.NkEnd(ctx)
}

func (fb *FileBrowser) ReloadFiles() {
	if fb.fileDir != fb.Dir {
		fb.Files, _ = ioutil.ReadDir(fb.Dir)
		fb.fileDir = fb.Dir
	}
}

func (fb *FileBrowser) Reset() {
	cdir, _ := os.Getwd()
	fb.Active = false
	fb.Dir = cdir
}
