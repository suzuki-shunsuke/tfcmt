package local_file

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type LocalFileService service

// Post posts comment
func (f *LocalFileService) Post(ctx context.Context, body string, OutputFile string) error {
	logE := logrus.WithFields(logrus.Fields{
		"program": "tfcmt",
	})

	
	file, err := os.Create(OutputFile)
	
	if err != nil {
		return err
	}


	defer file.Close()
	
	_, err2 := file.WriteString(body)
	
	if err2 != nil {
		return err2
	}
	logE.Debug("Output to file success")
	
	return err
}
