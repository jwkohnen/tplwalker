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
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	for _, tld := range []string{"testdata/simple", "testdata/simple-want"} {
		for _, dir := range []string{"", "dir", "dir/sub-dir"} {
			err := os.Chmod(filepath.Join(filepath.FromSlash(tld), dir, "file-mode-640"), 0640)
			if err != nil {
				panic(err)
			}
		}
	}
	os.Exit(m.Run())
}

func TestSimpleSelfTest(t *testing.T) {
	testSimple(t, "testdata/simple-want")
}

func TestSimple(t *testing.T) {
	dst, cleanup := scratch("simple")
	defer cleanup()
	w, err := New("testdata/simple", ".tpl")
	if err != nil {
		t.Fatal(err)
	}
	w.IgnoreDir("ignore")
	err = w.WalkTemplates(dst, map[string]string{"theText": "poof"})
	if err != nil {
		t.Fatal(err.Error())
	}
	testSimple(t, dst)
}

func testSimple(t *testing.T, dst string) {
	prefices := []string{"", "dir", "dir/sub-dir"}
	for _, prefix := range prefices {
		testTextFiles(t, dst, prefix)
		testSymlink(t, dst, prefix)
		testBinaryFile(t, dst, prefix)
		testIgnoreDir(t, dst, prefix)
	}
}

func testTextFiles(t *testing.T, dst, prefix string) {
	for file, want := range map[string]string{
		"poof.txt":  "poof\n",
		"hello.txt": "hello\n",
	} {
		path := filepath.Join(dst, prefix, file)
		data, err := ioutil.ReadFile(path)
		if err != nil {
			t.Errorf("Could not read file %s: %v", path, err)
			continue
		}
		if got := string(data); want != got {
			t.Errorf("Want %s, got %s", want, got)
		}
	}
}

func testSymlink(t *testing.T, dst, prefix string) {
	var (
		link      = filepath.Join(dst, prefix, "link")
		dest      = filepath.Join(dst, prefix, "poof.txt")
		err       error
		linkLstat os.FileInfo
		linkStat  os.FileInfo
		destStat  os.FileInfo
	)
	linkLstat, err = os.Lstat(link)
	if err != nil {
		t.Errorf("Could not SYS_LSTAT the link %s: %v", link, err)
		return
	}
	if linkLstat.Mode()&os.ModeSymlink == 0 {
		t.Errorf("Link \"%s\" is not a link: %v", link, linkLstat.Mode().String())
	}
	linkStat, err = os.Stat(link)
	if err != nil {
		t.Errorf("Could not SYS_STAT the link %s: %v:", link, err)
		return
	}
	destStat, err = os.Stat(dest)
	if err != nil {
		t.Errorf("Cound not stat %s: %v", dest, err)
		return
	}
	if !os.SameFile(linkStat, destStat) {
		t.Errorf("Link %s [%#v] does not point to dest %s [%#v]", link, linkStat, dest, destStat)
	}
}

func testBinaryFile(t *testing.T, dst, prefix string) {
	data, err := ioutil.ReadFile(filepath.Join(dst, prefix, "four-binary-zeros"))
	if err != nil {
		t.Errorf("Could not read four binary zeros: %v:", err)
	}
	if !bytes.Equal(data, []byte{0x00, 0x00, 0x00, 0x00}) {
		t.Error("four-binary-zeros is not four binary zeros")
	}

	// test file modes
	mode640file := filepath.Join(dst, prefix, "file-mode-640")
	modeStat, err := os.Stat(mode640file)
	if err != nil {
		t.Errorf("Could not stat file %s: %v", mode640file, err)
		return
	}
	const want = os.FileMode(0640)
	if got := modeStat.Mode().Perm(); got != want {
		t.Errorf("File %s want mode %s, got %s", mode640file, want, got)
	}
}

func testIgnoreDir(t *testing.T, dst, prefix string) {
	ignore := filepath.Join(dst, prefix, "ignore")
	stat, err := os.Stat(ignore)
	if pathError, ok := err.(*os.PathError); ok && os.IsNotExist(pathError.Err) {
		return
	}
	if err != nil {
		t.Fatalf("%#v (%s)", err, err)
	}
	t.Errorf("there is a file %s [%s], but it shouldn't be there", ignore, stat.Mode())
}

func scratch(name string) (string, func()) {
	dir, err := ioutil.TempDir("", "tpl-walker-test-"+name+"-"+time.Now().Format(time.RFC3339)+"-")
	if err != nil {
		panic(err)
	}
	cleanup := func() {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}
	return dir, cleanup
}
