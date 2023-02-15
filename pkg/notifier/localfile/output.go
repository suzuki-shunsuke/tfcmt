package localfile

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type OutputService service

// WriteToFile Write result to file
func (f *OutputService) WriteToFile(ctx context.Context, body string, outputFile string) error {
	logE := logrus.WithFields(logrus.Fields{
		"program": "tfcmt",
	})

	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteString(body); err != nil {
		return err
	}
	logE.Debug("Output to file success")

	return err
}
