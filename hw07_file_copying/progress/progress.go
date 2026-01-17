package progress

import (
	"io"

	"github.com/cheggaaa/pb/v3"
)

// GetReader оборачивает исходный reader в прогресс-бар
func GetReader(source io.Reader, limit int64) (io.Reader, func()) {
	bar := pb.Full.Start64(limit)
	reader := bar.NewProxyReader(source)

	return reader, func() {
		bar.Finish()
	}
}
