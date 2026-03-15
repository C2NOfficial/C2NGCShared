package invoice

import (
	"encoding/base64"
	"os"

	"github.com/signintech/gopdf"
)

func BytesToBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func BytesToFile(b []byte) error {
	return os.WriteFile("invoice/generated/invoice.pdf", b, 0644)
}

func getTextWidth(pdf *gopdf.GoPdf, text string) float64 {
	width, _ := pdf.MeasureTextWidth(text)
	return width
}

func getTextHeight(pdf *gopdf.GoPdf, text string) float64 {
	height, _ := pdf.MeasureCellHeightByText(text)
	return height
}