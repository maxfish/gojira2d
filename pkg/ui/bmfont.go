package ui

// Bitmap font loader. The format is the one described here
// http://www.angelcode.com/products/bmfont/doc/file_format.html

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// BmChar holds information about a single character, see
// http://www.angelcode.com/products/bmfont/doc/file_format.html
type BmChar struct {
	id             int32
	x              int
	y              int
	width          int
	height         int
	offsetX        int
	offsetY        int
	advanceX       int
	pageIndex      int
	textureChannel int
	kernings       map[int32]int
}

// BmFont holds all information about the font, see
// http://www.angelcode.com/products/bmfont/doc/file_format.html
type BmFont struct {
	pageFiles      map[int]string
	charactersList []*BmChar
	Characters     map[int32]*BmChar

	// Info
	face          string
	size          int
	bold          bool
	italic        bool
	charset       string
	unicode       bool
	stretchH      int
	smooth        bool
	superSampling int
	padding       [4]int // top, right, bottom, left
	spacing       [2]int // x, y
	// Common
	lineHeight int
	base       int
	pageWidth  int
	pageHeight int
	packed     bool
	numPages   int
	// Chars
	charactersCount int
}

// NewBmFontFromFile parse the font data out of a file
func NewBmFontFromFile(fileName string) *BmFont {
	f := &BmFont{}

	f.pageFiles = make(map[int]string)
	f.Characters = make(map[int32]*BmChar)

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(fileContent), "\n")
	for _, line := range lines {
		section, keyValues := f.tokenizeLine(line)
		switch section {
		case "info":
			f.parseInfoSection(keyValues)
		case "common":
			f.parseCommonSection(keyValues)
		case "page":
			f.parsePageSection(keyValues)
		case "char":
			f.parseCharSection(keyValues)
		case "kerning":
			f.parseKerningSection(keyValues)
		}
	}

	return f
}

func (f *BmFont) parseInfoSection(keyValues map[string]string) {
	f.face = keyValues["face"]
	f.size, _ = strconv.Atoi(keyValues["size"])
	f.bold, _ = strconv.ParseBool(keyValues["bold"])
	f.italic, _ = strconv.ParseBool(keyValues["italic"])
	f.unicode, _ = strconv.ParseBool(keyValues["unicode"])
	f.stretchH, _ = strconv.Atoi(keyValues["stretchH"])
	f.smooth, _ = strconv.ParseBool(keyValues["smooth"])
	f.superSampling, _ = strconv.Atoi(keyValues["aa"])

	// Padding
	paddingStrings := strings.Split(keyValues["padding"], ",")
	for i := 0; i < 4; i++ {
		f.padding[i], _ = strconv.Atoi(paddingStrings[i])
	}
	// Spacing
	spacingStrings := strings.Split(keyValues["spacing"], ",")
	for i := 0; i < 2; i++ {
		f.spacing[i], _ = strconv.Atoi(spacingStrings[i])
	}
}

func (f *BmFont) parseCommonSection(keyValues map[string]string) {
	f.lineHeight, _ = strconv.Atoi(keyValues["lineHeight"])
	f.base, _ = strconv.Atoi(keyValues["base"])
	f.pageWidth, _ = strconv.Atoi(keyValues["scaleW"])
	f.pageHeight, _ = strconv.Atoi(keyValues["scaleH"])
	f.packed, _ = strconv.ParseBool(keyValues["packed"])
	f.numPages, _ = strconv.Atoi(keyValues["pages"])
}

func (f *BmFont) parsePageSection(keyValues map[string]string) {
	id, _ := strconv.Atoi(keyValues["id"])
	f.pageFiles[id] = keyValues["file"]
}

func (f *BmFont) parseCharSection(keyValues map[string]string) {
	c := &BmChar{}
	id, _ := strconv.Atoi(keyValues["id"])
	c.id = int32(id)
	c.x, _ = strconv.Atoi(keyValues["x"])
	c.y, _ = strconv.Atoi(keyValues["y"])
	c.width, _ = strconv.Atoi(keyValues["width"])
	c.height, _ = strconv.Atoi(keyValues["height"])
	c.offsetX, _ = strconv.Atoi(keyValues["xoffset"])
	c.offsetY, _ = strconv.Atoi(keyValues["yoffset"])
	c.advanceX, _ = strconv.Atoi(keyValues["xadvance"])
	c.pageIndex, _ = strconv.Atoi(keyValues["page"])
	c.textureChannel, _ = strconv.Atoi(keyValues["chnl"])

	c.kernings = make(map[int32]int)
	f.Characters[c.id] = c
}

func (f *BmFont) parseKerningSection(keyValues map[string]string) {
	first, _ := strconv.Atoi(keyValues["first"])
	second, _ := strconv.Atoi(keyValues["second"])
	amount, _ := strconv.Atoi(keyValues["amount"])

	char, ok := f.Characters[int32(second)]
	if !ok {
		fmt.Printf("Kerning parse error: char %v not found", first)
	}
	char.kernings[int32(first)] = amount
}

var bmSectionRex = regexp.MustCompile("^(\\w+) ")
var bmKeyValueRex = regexp.MustCompile("(\\w+)=\"?([\\w\\s ,._\\-]*)\"?( |$|\")")

func (f *BmFont) tokenizeLine(line string) (string, map[string]string) {
	sectionMatches := bmSectionRex.FindStringSubmatch(line)
	if sectionMatches == nil {
		return "", nil
	}
	sectionName := sectionMatches[1]
	data := bmKeyValueRex.FindAllStringSubmatch(line, -1)

	keyValues := make(map[string]string)
	for _, kv := range data {
		k := kv[1]
		v := strings.Trim(kv[2], " ")
		keyValues[k] = v
	}

	return sectionName, keyValues
}
