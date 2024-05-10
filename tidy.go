package docconv

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"golang.org/x/net/html"
)

// Tidy attempts to tidy up XML.
// Errors & warnings are deliberately suppressed as underlying tools
// throw warnings very easily.
/*func Tidy(r io.Reader, xmlIn bool) ([]byte, error) {
	f, err := os.CreateTemp(os.TempDir(), "docconv")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())
	io.Copy(f, r)

	var output []byte
	if xmlIn {
		output, err = exec.Command("tidy", "-xml", "-numeric", "-asxml", "-quiet", "-utf8", f.Name()).Output()
	} else {
		output, err = exec.Command("tidy", "-numeric", "-asxml", "-quiet", "-utf8", f.Name()).Output()
	}

	if err != nil && err.Error() != "exit status 1" {
		return nil, err
	}
	return output, nil
}*/

// TidyHTML cleans and formats the provided HTML content
func TidyHTML(r io.Reader) ([]byte, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var buf bytes.Buffer
	err = html.Render(&buf, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to render HTML: %w", err)
	}

	return buf.Bytes(), nil
}

// TidyXML cleans and formats the provided XML content
func TidyXML(r io.Reader) ([]byte, error) {
	decoder := xml.NewDecoder(r)
	var buffer bytes.Buffer
	encoder := xml.NewEncoder(&buffer)
	encoder.Indent("", "  ")

	// Decode and re-encode to pretty print
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to decode XML: %w", err)
		}

		err = encoder.EncodeToken(token)
		if err != nil {
			return nil, fmt.Errorf("failed to encode XML: %w", err)
		}
	}
	err := encoder.Flush()
	if err != nil {
		return nil, fmt.Errorf("failed to flush encoder: %w", err)
	}

	return buffer.Bytes(), nil
}

func Tidy(r io.Reader, xmlIn bool) ([]byte, error) {
	if xmlIn {
		return TidyXML(r)
	}
	return TidyHTML(r)
}
