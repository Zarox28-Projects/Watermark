package engine

import (
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func ProcessPDF(inputPath string, text string, outputPath string) (bool, *string) {
	// Parse the watermark configuration from the text
	wmConf, err := pdfcpu.ParseTextWatermarkDetails(
		text,
		"rot:45, op:.1, scale:1",
		true,
		types.POINTS,
	)
	if err != nil {
		errStr := err.Error()
		return false, &errStr
	}

	// Add the watermark to the PDF file
	if err := api.AddWatermarksFile(inputPath, outputPath, nil, wmConf, nil); err != nil {
		errStr := err.Error()
		return false, &errStr
	}

	return true, nil
}
