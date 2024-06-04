package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	exec "golang.org/x/sys/execabs"
)

//go:embed go.patch
var goPatch []byte

func CreateCustomizedGOROOT() error {
	customizedGOROOT, err := xdg.DataFile("goos-js2-go/customized-goroot")
	if err != nil {
		return err
	}
	log.Printf("customizedGOROOT=%#+v", customizedGOROOT)

	err = os.MkdirAll(customizedGOROOT, 0755)
	if err != nil {
		return err
	}

	// Pseudo-copy the GOROOT directory (the native original Go toolchain installation) using as many symlinks
	// as possible to our new custom GOROOT directory. This new custom GOROOT directory will have go.patch applied
	// which will change some of the src/... files but none of the src/cmd/... files.
	err = filepath.WalkDir(filepath.Join(runtime.GOROOT()), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		pathRel, err := filepath.Rel(runtime.GOROOT(), path)
		if err != nil {
			return err
		}
		// log.Printf("pathRel=%#+v", pathRel)
		if pathRel == "." {
			return nil
		}

		// If it's src/, copy all the contents (except src/cmd/ which is symlinked) and skip further crawling.
		if pathRel == "src" {
			err = os.Mkdir(filepath.Join(customizedGOROOT, pathRel), 0755)
			if err != nil && !errors.Is(err, fs.ErrExist) {
				return err
			}

			err = filepath.WalkDir(filepath.Join(runtime.GOROOT(), pathRel), func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				pathRel, err := filepath.Rel(runtime.GOROOT(), path)
				if err != nil {
					return err
				}
				// log.Printf("pathRel=%#+v", pathRel)
				if pathRel == "." {
					return nil
				}

				if pathRel == filepath.Join("src", "cmd") || pathRel == filepath.Join("src", "vendor") || pathRel == filepath.Join("src", "testdata") {
					err = os.Symlink(path, filepath.Join(customizedGOROOT, pathRel))
					if err != nil && !errors.Is(err, fs.ErrExist) {
						return err
					}
					if d.IsDir() {
						return filepath.SkipDir
					} else {
						return nil
					}
				} else if d.IsDir() {
					err = os.Mkdir(filepath.Join(customizedGOROOT, pathRel), 0755)
					if err != nil && !errors.Is(err, fs.ErrExist) {
						return err
					} else {
						return nil
					}
				} else if len(filepath.SplitList(pathRel)) == 1 {
					err = os.Symlink(path, filepath.Join(customizedGOROOT, pathRel))
					if err != nil && !errors.Is(err, fs.ErrExist) {
						return err
					} else {
						return nil
					}
				} else {
					return os.WriteFile(filepath.Join(customizedGOROOT, pathRel), goPatch, 0644)
				}
			})
			if err != nil {
				return err
			}
			return filepath.SkipDir
		} else {
			err = os.Symlink(path, filepath.Join(customizedGOROOT, pathRel))
			if err != nil && !errors.Is(err, fs.ErrExist) {
				return err
			}
			if d.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
	})
	if err != nil {
		return err
	}

	goPatch2, err := xdg.DataFile("goos-js2-go/go.patch")
	if err != nil {
		return err
	}
	err = os.WriteFile(goPatch2, goPatch, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "apply", goPatch2)
	cmd.Dir = customizedGOROOT
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cmd.Output() failed: %w\noutput:\n%s", err, out)
	}
	log.Printf("out=%#+v", out)

	return nil
}
