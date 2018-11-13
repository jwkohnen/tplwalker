//    This file is part of tplwalker.
//
//    tplwalker is free software: you can redistribute it and/or modify it
//    under the terms of the GNU General Public License as published by the
//    Free Software Foundation, either version 3 of the License, or (at your
//    option) any later version.
//
//    tplwalker is distributed in the hope that it will be useful, but WITHOUT
//    ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
//    FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
//    more details.
//
//    You should have received a copy of the GNU General Public License along
//    with tplwalker.  If not, see <http://www.gnu.org/licenses/>.

package tplwalker

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateWalker interface {
	WalkTemplates(dst string, data interface{}) error
	IgnoreDir(dir ...string)
}

var _ TemplateWalker = (*TplWalker)(nil)

type TplWalker struct {
	source     string
	skipDirSet map[string]struct{}
	suffix     string
}

func New(source, suffix string) (*TplWalker, error) {
	t := &TplWalker{source: source, suffix: suffix}
	err := t.walkerInit()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t TplWalker) WalkTemplates(dst string, data interface{}) error {
	return filepath.Walk(t.source, func(path string, stat os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if stat.IsDir() {
			if _, ok := t.skipDirSet[stat.Name()]; ok {
				return filepath.SkipDir
			}
		}
		relPath, err := filepath.Rel(t.source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, relPath)
		switch mode := stat.Mode(); {
		case mode.IsDir():
			return os.MkdirAll(target, mode)
		case mode.IsRegular():
			if strings.HasSuffix(path, t.suffix) {
				target = target[:len(target)-len(t.suffix)]
				return execute(target, path, mode, data)
			}
			return copyFile(target, path, mode)
		case mode&os.ModeSymlink != 0:
			return symlink(target, path)
		default:
			return fmt.Errorf("unsupported file type %s: %s", path, stat.Mode())
		}
	})
}

func (t *TplWalker) IgnoreDir(dirs ...string) {
	for _, dir := range dirs {
		t.skipDirSet[dir] = struct{}{}
	}
}

func execute(target, path string, mode os.FileMode, data interface{}) error {
	tpl, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	dst, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer func() { err = syncClose(dst, err) }()

	err = tpl.Execute(dst, data)
	defer func() { err = syncDir(filepath.Dir(target), err) }()

	return err
}

func syncDir(path string, inError error) (err error) {
	err = inError

	dir, errOpen := os.OpenFile(path, os.O_RDONLY, 0)
	if errOpen != nil {
		if err == nil {
			err = errOpen
		}
		return err
	}

	return syncClose(dir, err)
}

func syncClose(path *os.File, inError error) (err error) {
	err = inError

	errSync := path.Sync()
	if errSync != nil {
		if err == nil {
			err = errSync
		}
		return err
	}

	errClose := path.Close()
	if errClose != nil {
		if err == nil {
			err = errClose
		}
		return err
	}

	return inError
}

func copyFile(target string, path string, mode os.FileMode) error {
	dst, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer func() {
		if errDst := syncClose(dst, err); err == nil {
			err = errDst
		}
	}()

	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := src.Close(); err == nil {
			err = errClose
		}
	}()

	_, err = io.Copy(dst, src)
	return err
}

func (t *TplWalker) walkerInit() error {
	t.skipDirSet = make(map[string]struct{})
	return nil
}

func symlink(target, path string) error {
	linksToPath, err := os.Readlink(path)
	if err != nil {
		return err
	}
	return os.Symlink(linksToPath, target)
}
