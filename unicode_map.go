package gopdf

import (
	"bytes"
	"fmt"
)

type UnicodeMap struct {
	buffer             bytes.Buffer
	PtrToSubsetFontObj *SubsetFontObj
}

func (u *UnicodeMap) Init(funcGetRoot func() *GoPdf) {}

func (u *UnicodeMap) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	u.PtrToSubsetFontObj = ptr
}

func (u *UnicodeMap) Build() error {
	u.buffer.Write(u.pdfToUnicodeMap().Bytes())
	return nil
}

func (u *UnicodeMap) GetType() string {
	return "Unicode"
}

func (u *UnicodeMap) GetObjBuff() *bytes.Buffer {
	return &u.buffer
}

func (u *UnicodeMap) pdfToUnicodeMap() *bytes.Buffer {
	//stream
	characterToGlyphIndex := u.PtrToSubsetFontObj.CharacterToGlyphIndex
	prefix :=
		"/CIDInit /ProcSet findresource begin\n" +
			"12 dict begin\n" +
			"begincmap\n" +
			"/CIDSystemInfo << /Registry (Adobe)/Ordering (UCS)/Supplement 0>> def\n" +
			"/CMapName /Adobe-Identity-UCS def /CMapType 2 def\n"
	suffix := "endcmap CMapName currentdict /CMap defineresource pop end end"

	glyphIndexToCharacter := make(map[int]rune)
	lowIndex := 65536
	hiIndex := -1
	for k, v := range characterToGlyphIndex {
		index := int(v)
		if index < lowIndex {
			lowIndex = index
		}
		if index > hiIndex {
			hiIndex = index
		}
		glyphIndexToCharacter[index] = k
	}

	var buff bytes.Buffer
	buff.WriteString(prefix)
	buff.WriteString("1 begincodespacerange\n")
	buff.WriteString(fmt.Sprintf("<%04X><%04X>\n", lowIndex, hiIndex))
	buff.WriteString("endcodespacerange\n")
	buff.WriteString(fmt.Sprintf("%d beginbfrange\n", len(glyphIndexToCharacter)))
	for k, v := range glyphIndexToCharacter {
		buff.WriteString(fmt.Sprintf("<%04X><%04X><%04X>\n", k, k, v))
	}
	buff.WriteString("endbfrange\n")
	buff.WriteString(suffix)
	buff.WriteString("\n")

	length := buff.Len()
	var streambuff bytes.Buffer
	streambuff.WriteString("<<\n")
	streambuff.WriteString(fmt.Sprintf("/Length %d\n", length))
	streambuff.WriteString(">>\n")
	streambuff.WriteString("stream\n")
	streambuff.Write(buff.Bytes())
	streambuff.WriteString("endstream\n")

	return &streambuff
}
