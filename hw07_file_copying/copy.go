package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fi, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fi.Close()
	fiStat, err := fi.Stat()
	if err != nil {
		return err
	}
	fiSize := fiStat.Size()
	if fiSize == 0 {
		return ErrUnsupportedFile
	}
	if offset > fiSize {
		return ErrOffsetExceedsFileSize
	}
	fi.Seek(offset, io.SeekStart)
	fo, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fo.Close()
	if limit == 0 {
		limit = fiSize
	}
	countBytes := min(limit, fiSize-offset)
	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	barReader := bar.NewProxyReader(fi)
	_, err = io.CopyN(fo, barReader, countBytes)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	return nil
}
