package utils

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

func BoolToOuiNon(b bool) string {
	if b {
		return "Oui"
	}
	return "Non"
}

func GenerateFileName(debut, fin time.Time) string {
	debutFormat := debut.Format("2006-01-02")
	finFormat := fin.Format("2006-01-02")

	return fmt.Sprintf("pointages_%s_%s.xlsx", debutFormat, finFormat)
}

func CreateBorder() []excelize.Border {
	return []excelize.Border{
		{Type: "left", Color: "#1E1F24", Style: 1},
		{Type: "bottom", Color: "#1E1F24", Style: 1},
		{Type: "right", Color: "#1E1F24", Style: 1},
		{Type: "top", Color: "#1E1F24", Style: 1},
	}
}

func FormatHeure(heure float64) string {
	hh := int(heure)
	mm := int((heure - float64(hh)) * 60)
	return fmt.Sprintf("%02d:%02d", hh, mm)
}
