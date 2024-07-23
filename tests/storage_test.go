/*
 * @PackageName: tests
 * @Description:
 * @Author: gabbymrh
 * @Date: 2024-07-18 10:18:04
 * @LastModifiedBy: gabbymrh
 * @LastModifiedAt: 2024-07-18 10:18:04
 */

package tests

import (
	"bytes"
	"easy_dfs/pkg/filesystem"
	"testing"
)

func TestSaveAndLoadFile(t *testing.T) {
	storage := filesystem.FileSystemStorage{BaseDir: "./testdata"}
	filename := "testdir/testfile.txt"
	data := []byte("Hello, world!")

	// Save file
	if err := storage.Save(filename, bytes.NewReader(data)); err != nil {
		t.Fatalf("Failed to save file: %v", err)
	}

	// Load file
	reader, err := storage.Load(filename)
	if err != nil {
		t.Fatalf("Failed to load file: %v", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	if buf.String() != string(data) {
		t.Fatalf("Expected %s but got %s", string(data), buf.String())
	}

	// Delete file
	if err := storage.Delete(filename); err != nil {
		t.Fatalf("Failed to delete file: %v", err)
	}
}
