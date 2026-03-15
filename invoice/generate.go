package invoice

import (
	"log"
	"os"
	"path/filepath"

	"github.com/signintech/gopdf"
)

func Generate(data *PDFData) []byte {

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	
	fontDir, err := loadFonts()
	if err != nil {
	    log.Println("Font load error:", err)
	}
	defer os.RemoveAll(fontDir)
	
	pdf.AddTTFFont("Anton", filepath.Join(fontDir, "Anton-Regular.ttf"))
	pdf.AddTTFFont("Inter", filepath.Join(fontDir, "Inter-Regular.ttf"))
	pdf.AddTTFFont("InterBold", filepath.Join(fontDir, "Inter-Bold.ttf"))

	pdf.AddPage()

	// ─────────────────────────────────────────
	// HEADER — Brand name (left) + INVOICE (right)
	// ─────────────────────────────────────────
	pdf.SetFont("Anton", "", 28)
	pdf.SetX(marginLeft)
	pdf.SetY(marginTop)
	pdf.SetTextColor(colorRed[0], colorRed[1], colorRed[2])
	pdf.Text(data.CompanyTradeName)

	pdf.SetFont("InterBold", "", 28)
	pdf.SetTextColor(colorBlack[0], colorBlack[1], colorBlack[2])
	pdf.SetX(marginRight - getTextWidth(pdf, "INVOICE"))
	pdf.Text("INVOICE")

	// Red horizontal divider line
	pdf.SetX(marginLeft)
	pdf.SetY(pdf.GetY() + 15)
	pdf.SetLineWidth(3)
	pdf.SetStrokeColor(colorRed[0], colorRed[1], colorRed[2])
	pdf.Line(pdf.GetX(), pdf.GetY(), marginRight, pdf.GetY())

	// ─────────────────────────────────────────
	// INVOICE META TABLE — Invoice No, Date, Order Ref, Place of Supply
	// ─────────────────────────────────────────
	pdf.SetY(pdf.GetY() + 10)
	pdf.SetX(marginLeft + 15)

	tableWidths := getEqualTableWidths(pdf.GetX(), marginRight-15, 4)
	totalTableWidth := getTotalTableWidth(tableWidths)
	initialX, initialY := pdf.GetX(), pdf.GetY()

	// Header row — labels
	height := drawTableHeader(pdf, &tableHeaderStyle{
		text:               [1][]string{{"INVOICE NO.", "DATE", "ORDER REF.", "PLACE OF SUPPLY"}},
		cellWidths:         tableWidths,
		fontFamily:         "Inter",
		fontStyle:          "",
		fontSize:           7,
		textColor:          colorGrayText,
		BorderColor:        colorGrayBorder,
		FillColor:          colorGray,
		lineWidth:          0.5,
		height:             20,
		textMarginFromEdge: 5,
		width:              totalTableWidth,
	})

	// Body row — values
	pdf.SetXY(initialX, initialY+height)
	totalTableHeight := drawTableBody(pdf, &tableBodyStyle{
		text:               [][]string{{data.InvNo, data.InvDate, data.OrderNo, data.CompanyState}},
		cellWidths:         tableWidths,
		fontFamily:         "InterBold",
		fontStyle:          "",
		fontSize:           8,
		textColor:          colorBlack,
		lineWidth:          0.5,
		BorderColor:        colorGrayBorder,
		textMarginFromEdge: 5,
		width:              totalTableWidth,
		tableHeaderHeight:  height,
		rowTextTopPadding:  4,
	}, true)

	pdf.SetXY(marginLeft, initialY+totalTableHeight+10)

	// ─────────────────────────────────────────
	// FROM / BILL TO — Company + Customer details
	// ─────────────────────────────────────────
	pdf.SetStrokeColor(colorGrayBorder[0], colorGrayBorder[1], colorGrayBorder[2])
	pdf.SetFillColor(colorGray[0], colorGray[1], colorGray[2])

	var boxHeight float64 = 120
	var boxWidth float64 = marginLeft + marginRight

	// Outer box
	pdf.Rectangle(marginLeft, pdf.GetY(), marginRight, pdf.GetY()+boxHeight, "FD", 0, 0)

	// Center divider line
	pdf.Line((boxWidth)/2, pdf.GetY(), (boxWidth)/2, pdf.GetY()+boxHeight)

	pdf.SetXY(marginLeft+10, pdf.GetY()+15)
	initialX = pdf.GetX()
	initialY = pdf.GetY()

	// ── Company (FROM) details
	pdf.SetFont("Inter", "", 7)
	pdf.SetTextColor(colorGrayText[0], colorGrayText[1], colorGrayText[2])
	pdf.Text("FROM")

	pdf.SetXY(initialX, pdf.GetY()+15)
	pdf.SetFont("InterBold", "", 8)
	pdf.SetTextColor(colorBlack[0], colorBlack[1], colorBlack[2])
	pdf.Text(data.CompanyName)

	pdf.SetFont("Inter", "", 7)
	pdf.SetTextColor(colorGrayText[0], colorGrayText[1], colorGrayText[2])

	pdf.SetXY(initialX, pdf.GetY()+10)
	pdf.Text("Trading as " + data.CompanyTradeName)

	pdf.SetXY(initialX, pdf.GetY()+10)
	pdf.Text(data.CompanyAddressLine)

	pdf.SetXY(initialX, pdf.GetY()+10)
	pdf.Text(data.CompanyCity + ", " + data.CompanyState + " - " + data.CompanyPincode)

	pdf.SetXY(initialX, pdf.GetY()+15)
	pdf.Text("Email: " + data.CompanyContactEmail)

	pdf.SetXY(initialX, pdf.GetY()+15)
	pdf.Text("GSTIN: " + data.CompanyGSTIN)

	// ── Customer (BILL TO) details
	rightSideX := initialX + boxWidth/2 - marginLeft
	pdf.SetXY(rightSideX, initialY)
	pdf.SetFont("Inter", "", 7)
	pdf.SetTextColor(colorGrayText[0], colorGrayText[1], colorGrayText[2])
	pdf.Text("BILL TO")

	pdf.SetXY(rightSideX, pdf.GetY()+15)
	pdf.SetFont("InterBold", "", 8)
	pdf.SetTextColor(colorBlack[0], colorBlack[1], colorBlack[2])
	pdf.Text(data.CustomerName)

	pdf.SetFont("Inter", "", 7)
	pdf.SetTextColor(colorGrayText[0], colorGrayText[1], colorGrayText[2])

	pdf.SetXY(rightSideX, pdf.GetY()+10)
	pdf.Text(data.CustomerAddressLine)

	pdf.SetXY(rightSideX, pdf.GetY()+10)
	pdf.Text(data.CustomerCity + ", " + data.CustomerState + " - " + data.CustomerPincode)

	pdf.SetXY(rightSideX, pdf.GetY()+15)
	pdf.Text("Ph: " + data.CustomerPhone)

	pdf.SetXY(rightSideX, pdf.GetY()+15)
	pdf.Text("Email: " + data.CustomerEmail)

	// ─────────────────────────────────────────
	// ITEMS TABLE — #, Description, HSN, Qty, Unit Price, Total
	// ─────────────────────────────────────────
	pdf.SetXY(marginLeft, initialY+boxHeight-5)
	initialX = pdf.GetX()
	initialY = pdf.GetY()

	tableWidths = []float64{20, 210, 50, 50, 80, 85}
	totalTableWidth = getTotalTableWidth(tableWidths)

	// Header row — column labels (white text on black background)
	height = drawTableHeader(pdf, &tableHeaderStyle{
		text:               [1][]string{{"#", "DESCRIPTION", "HSN", "QTY", "UNIT PRICE", "TOTAL"}},
		cellWidths:         tableWidths,
		fontFamily:         "Inter",
		fontStyle:          "",
		fontSize:           7,
		textColor:          colorWhite,
		BorderColor:        colorGrayBorder,
		FillColor:          colorBlack,
		lineWidth:          0.5,
		height:             25,
		textMarginFromEdge: 5,
		width:              totalTableWidth,
	})

	// Body rows — order items
	pdf.SetXY(initialX, initialY+height)
	totalTableHeight = drawTableBody(pdf, &tableBodyStyle{
		text:               data.Items,
		cellWidths:         tableWidths,
		fontFamily:         "Inter",
		fontStyle:          "",
		fontSize:           7,
		textColor:          colorBlack,
		lineWidth:          0.25,
		BorderColor:        colorGrayBorder,
		textMarginFromEdge: 5,
		width:              totalTableWidth,
		tableHeaderHeight:  height,
		rowTextTopPadding:  4,
	}, false)

	// ─────────────────────────────────────────
	// NOTES + ORDER SUMMARY — side by side
	// ─────────────────────────────────────────
	pdf.SetXY(marginLeft+10, initialY+totalTableHeight+10)
	initialX = pdf.GetX()
	initialY = pdf.GetY()

	// Notes box (left side)
	pdf.SetStrokeColor(colorGrayBorder[0], colorGrayBorder[1], colorGrayBorder[2])
	pdf.SetFillColor(colorGray[0], colorGray[1], colorGray[2])
	boxHeight = 70
	boxWidth = marginRight/2 + 15
	pdf.Rectangle(marginLeft, pdf.GetY(), boxWidth, pdf.GetY()+boxHeight, "FD", 0, 0)

	pdf.SetXY(initialX, initialY+15)
	pdf.SetFont("Inter", "", 7)
	pdf.SetTextColor(colorGrayText[0], colorGrayText[1], colorGrayText[2])
	pdf.Text("NOTES")

	pdf.SetXY(initialX, pdf.GetY()+15)
	pdf.Text("• Thank you for shopping with " + data.CompanyTradeName + "!")

	pdf.SetXY(initialX, pdf.GetY()+10)
	pdf.Text("• All sales are final. Refer to the return and refund policy")
	pdf.SetXY(initialX + 3.5, pdf.GetY()+8)
	pdf.Text(" on the website for further details..")

	pdf.SetXY(initialX, pdf.GetY()+10)
	pdf.Text("• For queries: " + data.CompanyContactEmail)

	// Order summary (right side)
	pdf.SetXY(initialX+boxWidth, initialY+8)
	initialX = pdf.GetX()
	initialY = pdf.GetY()

	const summaryTextToValueGap = 142.0

	pdf.SetFont("Inter", "", 8)
	pdf.SetTextColor(colorBlack[0], colorBlack[1], colorBlack[2])

	// Labels column
	pdf.SetXY(initialX, initialY)
	pdf.Text("Subtotal (Excluding GST)")

	pdf.SetXY(initialX, pdf.GetY()+20)
	// Show CGST + SGST for same state, IGST for interstate
	if data.IGST == "" {
		pdf.Text("CGST (2.5%)")
		pdf.SetXY(initialX, pdf.GetY()+20)
		pdf.Text("SGST (2.5%)")
	} else {
		pdf.Text("IGST (5%)")
	}

	pdf.SetXY(initialX, pdf.GetY()+20)
	pdf.Text("Shipping")

	// Values column
	pdf.SetXY(initialX+summaryTextToValueGap, initialY)
	pdf.Text(data.Subtotal)

	pdf.SetXY(initialX+summaryTextToValueGap, pdf.GetY()+20)
	if data.IGST == "" {
		pdf.Text(data.CGST)
		pdf.SetXY(initialX+summaryTextToValueGap, pdf.GetY()+20)
		pdf.Text(data.SGST)
	} else {
		pdf.Text(data.IGST)
	}

	pdf.SetXY(initialX+summaryTextToValueGap, pdf.GetY()+20)
	pdf.Text(data.Shipping)

	// Total row — double lines above and below, red bold text
	pdf.SetXY(initialX, pdf.GetY()+10)
	pdf.SetStrokeColor(colorBlack[0], colorBlack[1], colorBlack[2])
	pdf.SetLineWidth(1.5)
	pdf.Line(pdf.GetX(), pdf.GetY(), pdf.GetX()+summaryTextToValueGap+55, pdf.GetY())

	pdf.SetXY(initialX, pdf.GetY()+12)
	pdf.SetFont("InterBold", "", 9)
	pdf.SetTextColor(colorRed[0], colorRed[1], colorRed[2])
	pdf.Text("TOTAL PAID")

	pdf.SetX(initialX + summaryTextToValueGap)
	pdf.Text(data.Total)

	pdf.SetXY(initialX, pdf.GetY()+8)
	pdf.SetStrokeColor(colorBlack[0], colorBlack[1], colorBlack[2])
	pdf.SetLineWidth(1.5)
	pdf.Line(pdf.GetX(), pdf.GetY(), pdf.GetX()+summaryTextToValueGap+55, pdf.GetY())

	// ─────────────────────────────────────────
	// FOOTER — divider line + company info + legal note
	// ─────────────────────────────────────────
	pdf.SetXY(marginLeft, marginBottom)
	pdf.SetStrokeColor(colorGrayBorder[0], colorGrayBorder[1], colorGrayBorder[2])
	pdf.SetLineWidth(1)
	pdf.Line(pdf.GetX(), pdf.GetY(), marginRight, pdf.GetY())

	pdf.SetXY((marginRight/2)-72, pdf.GetY()+13)
	pdf.SetFont("Inter", "", 7)
	pdf.SetTextColor(colorGrayText[0], colorGrayText[1], colorGrayText[2])
	pdf.Text(data.CompanyTradeName + " | " + data.CompanyWebsite + " | " + data.CompanyContactEmail + " | " + data.CompanyInstagram)

	pdf.SetXY((marginRight/2)-72, pdf.GetY()+18)
	pdf.Text("This is a computer generated invoice. No signature is required.")

	return pdf.GetBytesPdf()
}