package main

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Copy copies limit bytes from fromPath to toPath with offset.
func Copy(fromPath string, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		return ErrUnsupportedFile
	}

	if info.Size() == 0 {
		return ErrUnsupportedFile
	}

	if offset > info.Size() {
		return ErrOffsetExceedsFileSize
	}

	var actualLimit int64

	if limit == 0 {
		actualLimit = info.Size()
	} else {
		actualLimit = limit

		if offset+limit > info.Size() {
			actualLimit = info.Size() - offset
		}
	}

	name := path.Base(fromPath)

	container := mpb.New(mpb.WithWidth(64))

	bar := container.AddBar(actualLimit,
		mpb.BarStyle("[=>-|"),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f "),
			decor.OnComplete(decor.Name(name, decor.WC{W: len(name), C: decor.DextraSpace}), "done!"),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)

	bufSize := actualLimit / 100
	buf := make([]byte, bufSize)

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer dstFile.Close()

	for i := 1; i < 100; i++ {
		if err := copyPortion(srcFile, dstFile, &buf, offset); err != nil {
			return err
		}
		bar.IncrBy(int(bufSize))
		offset += bufSize
	}

	lastBufSize := bufSize + actualLimit%100
	lastBuf := make([]byte, lastBufSize)

	if err := copyPortion(srcFile, dstFile, &lastBuf, offset); err != nil {
		return err
	}
	bar.IncrBy(int(lastBufSize))

	container.Wait()

	return nil
}

// copyPortion copies len(buf) bytes from srcFile into dstFile.
func copyPortion(srcFile io.ReaderAt, dstFile io.Writer, buf *[]byte, offset int64) error {
	_, err := srcFile.ReadAt(*buf, offset)
	if err != nil && err != io.EOF {
		return err
	}

	_, err = dstFile.Write(*buf)
	if err != nil {
		return err
	}

	return nil
}
