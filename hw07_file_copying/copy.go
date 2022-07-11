package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o755)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	fromFileInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}

	if offset > fromFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit > fromFileInfo.Size() {
		limit = fromFileInfo.Size()
	}

	if limit == 0 {
		limit = fromFileInfo.Size()
	}

	writer, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer writer.Close()

	fromFile.Seek(offset, 0)
	reader := io.LimitReader(fromFile, limit)

	// start new bar
	bar := pb.Full.Start64(limit)

	// create proxy reader
	barReader := bar.NewProxyReader(reader)

	// copy from proxy reader
	written, err := io.Copy(writer, barReader)
	if err != nil {
		return err
	}

	fmt.Println(written)

	// finish bar
	bar.Finish()

	return nil
}
