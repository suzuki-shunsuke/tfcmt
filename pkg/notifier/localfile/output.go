package localfile

import (
	"context"
	"fmt"
	"os"
)

type OutputService service

// WriteToFile Write result to file
func (f *OutputService) WriteToFile(ctx context.Context, body string, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("create a file to output the result to a file: %w", err)
	}

	defer file.Close()

	if _, err := file.WriteString(body); err != nil {
		return fmt.Errorf("write the result to a file: %w", err)
	}
	return nil
}
