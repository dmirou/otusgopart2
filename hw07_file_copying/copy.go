package main

import (
	"errors"
	"fmt"
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

const (
	progressBarWidth   = 64
	oneHundredPercents = 100
)

// Copy copies limit bytes from fromPath to toPath with offset.
func Copy(fromPath string, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("can't open file %s: %w", fromPath, err)
	}

	defer srcFile.Close()

	sizeToCopy, err := getSizeToCopy(srcFile, offset, limit)
	if err != nil {
		return err
	}

	name := path.Base(fromPath)

	container := mpb.New(mpb.WithWidth(progressBarWidth))

	bar := container.AddBar(sizeToCopy,
		mpb.BarStyle("[=>-|"),
		mpb.PrependDecorators(
			decor.CountersKibiByte("% .2f / % .2f "),
			decor.OnComplete(decor.Name(name, decor.WC{W: len(name), C: decor.DextraSpace}), "done!"),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)

	defer container.Wait()

	bufSize := sizeToCopy / oneHundredPercents
	buf := make([]byte, bufSize)

	dstFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("can't create file %s: %w", toPath, err)
	}

	defer dstFile.Close()

	for i := 1; i < oneHundredPercents; i++ {
		if err := copyPortion(srcFile, dstFile, &buf, offset); err != nil {
			return err
		}

		bar.IncrBy(int(bufSize))
		offset += bufSize
	}

	lastBufSize := bufSize + sizeToCopy%oneHundredPercents
	lastBuf := make([]byte, lastBufSize)

	if err := copyPortion(srcFile, dstFile, &lastBuf, offset); err != nil {
		return err
	}

	bar.IncrBy(int(lastBufSize))

	return nil
}

// getSizeToCopy returns count of bytes to copy according with offset and limit.
// It returns ErrUnsupportedFile if file is dir or it has zero size,
// ErrOffsetExceedsFileSize if offset exceeds file size.
func getSizeToCopy(srcFile *os.File, offset, limit int64) (int64, error) {
	info, err := srcFile.Stat()
	if err != nil {
		return 0, fmt.Errorf("can't get file info for %s: %w", srcFile.Name(), err)
	}

	if info.IsDir() {
		return 0, ErrUnsupportedFile
	}

	if info.Size() == 0 {
		return 0, ErrUnsupportedFile
	}

	if offset > info.Size() {
		return 0, ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		return info.Size() - offset, nil
	}

	size := limit

	if offset+limit > info.Size() {
		size = info.Size() - offset
	}

	return size, nil
}

// copyPortion copies len(buf) bytes from srcFile into dstFile.
func copyPortion(srcFile io.ReaderAt, dstFile io.Writer, buf *[]byte, offset int64) error {
	_, err := srcFile.ReadAt(*buf, offset)
	if err != nil && errors.Is(err, io.EOF) {
		return fmt.Errorf("can't read from file with offset: %w", err)
	}

	_, err = dstFile.Write(*buf)
	if err != nil {
		return fmt.Errorf("can't write buffer to file: %w", err)
	}

	return nil
}
