package invoice

import (
	"embed"
	"os"
	"path/filepath"
)

//go:embed fonts/Anton/Anton-Regular.ttf
//go:embed fonts/Inter/Inter-Regular.ttf
//go:embed fonts/Inter/Inter-Bold.ttf
var fontFiles embed.FS

//Unfortunately, GC function was having trouble finding the fonts
//if we try to load them directly from the current invoice package folder
//so we need to load these font bytes in a temp directory then pass that 
//to the invoice generator function so gopdf can find them
func loadFonts() (string, error) {
    tmpDir, err := os.MkdirTemp("", "c2n-fonts")
    if err != nil {
        return "", err
    }
    files := map[string]string{
        "fonts/Anton/Anton-Regular.ttf": "Anton-Regular.ttf",
        "fonts/Inter/Inter-Regular.ttf": "Inter-Regular.ttf",
        "fonts/Inter/Inter-Bold.ttf":    "Inter-Bold.ttf",
    }
    for src, dst := range files {
        data, err := fontFiles.ReadFile(src)
        if err != nil {
            return "", err
        }
        if err := os.WriteFile(filepath.Join(tmpDir, dst), data, 0644); err != nil {
            return "", err
        }
    }
    return tmpDir, nil
}