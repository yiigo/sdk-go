package wechat

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

// ValueToXML value to xml
func ValueToXML(vals V) (string, error) {
	var builder strings.Builder

	builder.WriteString("<xml>")
	for k, v := range vals {
		builder.WriteString("<" + k + ">")
		if err := xml.EscapeText(&builder, []byte(v)); err != nil {
			return "", err
		}
		builder.WriteString("</" + k + ">")
	}
	builder.WriteString("</xml>")

	return builder.String(), nil
}

// XMLToValue xml to value
func XMLToValue(b []byte) (V, error) {
	m := make(V)

	xmlReader := bytes.NewReader(b)

	var (
		d     = xml.NewDecoder(xmlReader)
		tk    xml.Token
		depth = 0 // current xml.Token depth
		key   string
		buf   bytes.Buffer
		err   error
	)

	d.Strict = false

	for {
		tk, err = d.Token()
		if err != nil {
			if err == io.EOF {
				return m, nil
			}
			return nil, err
		}
		switch v := tk.(type) {
		case xml.StartElement:
			depth++
			switch depth {
			case 2:
				key = v.Name.Local
				buf.Reset()
			case 3:
				if err = d.Skip(); err != nil {
					return nil, err
				}
				depth--
				key = "" // key == "" indicates that the node with depth==2 has children
			}
		case xml.CharData:
			if depth == 2 && key != "" {
				buf.Write(v)
			}
		case xml.EndElement:
			if depth == 2 && key != "" {
				m[key] = buf.String()
			}
			depth--
		}
	}
}
