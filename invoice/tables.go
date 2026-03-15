package invoice

import (
	"math"

	"github.com/signintech/gopdf"
)

type tableHeaderStyle struct {
	text                   [1][]string //table header can have no more than 1 row
	cellWidths             []float64
	fontFamily             string
	fontStyle              string
	fontSize               float64
	textColor              [3]uint8
	BorderColor, FillColor [3]uint8
	lineWidth              float64
	textMarginFromEdge     float64
	height                 float64 //entire table height
	width                  float64 //entire table width
}

type tableBodyStyle struct {
	text                   [][]string
	cellWidths             []float64
	fontFamily             string
	fontStyle              string
	fontSize               float64
	textColor              [3]uint8
	BorderColor, FillColor [3]uint8
	textMarginFromEdge     float64
	lineWidth              float64
	rowTextTopPadding      float64
	tableHeaderHeight      float64
	width                  float64 //entire table width
}

func drawTableHeader(pdf *gopdf.GoPdf, tableHeaderStyle *tableHeaderStyle) float64 {
	pdf.SetFillColor(tableHeaderStyle.FillColor[0], tableHeaderStyle.FillColor[1], tableHeaderStyle.FillColor[2])
	pdf.SetStrokeColor(tableHeaderStyle.BorderColor[0], tableHeaderStyle.BorderColor[1], tableHeaderStyle.BorderColor[2])
	pdf.SetLineWidth(tableHeaderStyle.lineWidth)

	// Draw the outer rectangle
	pdf.Rectangle(pdf.GetX(), pdf.GetY(), pdf.GetX()+tableHeaderStyle.width, pdf.GetY()+tableHeaderStyle.height, "FD", 0, 0)

	// Setup the font to get appropriate text width
	pdf.SetFont(tableHeaderStyle.fontFamily, tableHeaderStyle.fontStyle, tableHeaderStyle.fontSize)
	pdf.SetTextColor(tableHeaderStyle.textColor[0], tableHeaderStyle.textColor[1], tableHeaderStyle.textColor[2])

	textHeight := getTextHeight(pdf, "A")
	startX := pdf.GetX() + tableHeaderStyle.textMarginFromEdge //some padding so it doesn't stick to the left edge
	//starting at center of the table
	startY := pdf.GetY() + ((tableHeaderStyle.height - textHeight) / 2)

	//Draw the inner text now
	for i, text := range tableHeaderStyle.text[0] {
		pdf.SetXY(startX, startY)
		pdf.MultiCell(&gopdf.Rect{
			W: tableHeaderStyle.cellWidths[i],
			H: textHeight,
		}, text)
		startX += tableHeaderStyle.cellWidths[i] + tableHeaderStyle.textMarginFromEdge //some padding
	}
	return tableHeaderStyle.height
}

func drawTableBody(pdf *gopdf.GoPdf, tableBodyStyle *tableBodyStyle, drawCellVerticalBorders bool) float64 {
	pdf.SetFont(tableBodyStyle.fontFamily, tableBodyStyle.fontStyle, tableBodyStyle.fontSize)
	pdf.SetTextColor(tableBodyStyle.textColor[0], tableBodyStyle.textColor[1], tableBodyStyle.textColor[2])
	textHeight := getTextHeight(pdf, "A")
	var tableHeight float64 = 0.0

	var initialX, initialY = pdf.GetX(), pdf.GetY()
	startX := pdf.GetX() + tableBodyStyle.textMarginFromEdge
	startY := pdf.GetY() + tableBodyStyle.rowTextTopPadding //some padding

	//Draw the text now
	for _, row := range tableBodyStyle.text {
		for j, text := range row {
			pdf.SetXY(startX, startY)
			pdf.MultiCell(&gopdf.Rect{
				W: tableBodyStyle.cellWidths[j],
				H: textHeight,
			}, text)
			startX += tableBodyStyle.cellWidths[j] + tableBodyStyle.textMarginFromEdge //some padding
		}
		startX = initialX + tableBodyStyle.textMarginFromEdge //Reset the x position to beginning
		startY += textHeight + 2*(tableBodyStyle.rowTextTopPadding)
		tableHeight += textHeight + 2*(tableBodyStyle.rowTextTopPadding) //need some padding on bottom

		//Draw the row rectangle
		pdf.SetStrokeColor(tableBodyStyle.BorderColor[0], tableBodyStyle.BorderColor[1], tableBodyStyle.BorderColor[2])
		pdf.SetLineWidth(tableBodyStyle.lineWidth)
		pdf.Rectangle(initialX, initialY, initialX+tableBodyStyle.width, initialY+tableHeight, "D", 0, 0)
	}

	totalTableHeight := tableBodyStyle.tableHeaderHeight + tableHeight
	
	//Vertical borders for each column
	if drawCellVerticalBorders {
		curX := initialX
		y1 := initialY - tableBodyStyle.tableHeaderHeight
		y2 := (initialY - tableBodyStyle.tableHeaderHeight) + totalTableHeight

		for i, width := range tableBodyStyle.cellWidths {
			curX += width
			// no need to draw after last cell
			if i != len(tableBodyStyle.cellWidths)-1 {
				pdf.Line(curX, y1, curX, y2)
			}
		}
	}
	return totalTableHeight
}

func getEqualTableWidths(x0, x1 float64, arraySize int) []float64 {
	totalAvailable := math.Abs(x1 - x0)
	tableWidths := make([]float64, arraySize)
	for i := range tableWidths {
		tableWidths[i] = totalAvailable / float64(arraySize)
	}
	return tableWidths
}

func getTotalTableWidth(cellWidths []float64) float64 {
	total := 0.0
	for _, w := range cellWidths {
		total += w
	}
	return total
}
