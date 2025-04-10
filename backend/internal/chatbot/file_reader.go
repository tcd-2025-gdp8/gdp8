package chatbot

import (
	"bytes"
	"os/exec"
	"strings"
	"unicode"
)

func ReadPDF(path string) (string, error) {
	out, err := exec.Command("pdftotext", "-layout", path, "-").Output()
	if err != nil {
		return "", err
	}
	raw := string(out)
	clean := normaliseText(raw)
	return clean, nil
}

func normaliseText(text string) string {
	var buf bytes.Buffer
	for _, r := range text {
		if isAllowedRune(r) {
			buf.WriteRune(r)
		}
	}
	return strings.Join(strings.Fields(buf.String()), " ")
}

func isAllowedRune(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsPunct(r) || unicode.IsSpace(r) {
		return true
	}
	return false
}
