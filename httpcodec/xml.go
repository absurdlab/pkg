package httpcodec

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
)

// EncodeXML encodes payload as XML into writer.
func EncodeXML(writer io.Writer, payload interface{}) error {
	return xml.NewEncoder(writer).Encode(payload)
}

// DecodeXML decodes response body as XML into destination. The decoding charset
// is read from the XML attribute. The response body is closed afterwards.
func DecodeXML(destination interface{}) Decoder {
	return func(resp *http.Response) error {
		defer func() {
			_ = resp.Body.Close()
		}()
		decoder := xml.NewDecoder(resp.Body)
		decoder.CharsetReader = charset.NewReaderLabel
		return decoder.Decode(destination)
	}
}
