package main

import (
	"io/ioutil"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("UnsupportedFile/RandomSize", func(t *testing.T) {
		err := Copy("/dev/urandom", "./out.txt", 3, 0)
		if err != ErrUnsupportedFile {
			t.Fatalf(
				"unexpected error in Copy: %v, expected: %v",
				err, ErrUnsupportedFile,
			)
		}
	})

	t.Run("UnsupportedFile/Dir", func(t *testing.T) {
		dir, err := ioutil.TempDir("/tmp", "test-copy-*")
		if err != nil {
			t.Fatalf("Unexpected error in TempDir: %v", err)
		}

		err = Copy(dir, "./out.txt", 3, 0)
		if err != ErrUnsupportedFile {
			t.Fatalf(
				"unexpected error in Copy: %v, expected: %v",
				err, ErrUnsupportedFile,
			)
		}
	})

	t.Run("OffsetExceedsFileSize", func(t *testing.T) {
		err := Copy("./testdata/input_size2.txt", "./out.txt", 3, 0)
		if err != ErrOffsetExceedsFileSize {
			t.Fatalf(
				"unexpected error in Copy: %v, expected: %v",
				err, ErrOffsetExceedsFileSize,
			)
		}
	})
}
