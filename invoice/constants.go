package invoice

var (
	colorRed        = [3]uint8{255, 61, 34}
	colorBlack      = [3]uint8{0, 0, 0}
	colorWhite      = [3]uint8{255, 255, 255}
	colorGray       = [3]uint8{240, 240, 240}
	colorGrayBorder = [3]uint8{220, 220, 220}
	colorGrayText   = [3]uint8{150, 150, 150}
)

const (
	pageWidth    = 595.28
	pageHeight   = 841.89
	marginLeft   = 50
	marginRight  = pageWidth - marginLeft
	marginTop    = 80
	marginBottom = pageHeight - marginTop
)
