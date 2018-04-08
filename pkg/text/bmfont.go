package text
// Bitmap font loader. The format is the one described here
// http://www.angelcode.com/products/bmfont/doc/file_format.html

import (
	"io/ioutil"
	"log"
	"strings"
	"regexp"
	"strconv"
)

type BmChar struct {
	id             int32
	letter         string
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

	// page-size scaled
	f32x           float32
	f32y           float32
	f32width       float32
	f32height      float32

	// line-height scaled
	f32lineWidth   float32
	f32lineHeight  float32
	f32offsetX     float32
	f32offsetY     float32
	f32advanceX    float32
	f32kernings    map[int32]float32
}

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
	padding       [4]int
	spacing       [2]int
	// Common
	lineHeight int
	base       int
	pageWidth  int
	pageHeight int
	packed     bool
	numPages   int
	// Chars
	charactersCount int
	// Pre-calculated f32 scale values
	f32scaleLine  float32
	f32scaleW     float32
	f32scaleH     float32
}

func NewBmFontFromFile(fileName string) *BmFont {
	f := &BmFont{}

	f.pageFiles = make(map[int]string)
	f.Characters = make(map[int32]*BmChar)

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln(err)
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
	//f.padding = keyValues["padding"]
	//f.spacing = keyValues["spacing"]
}

func (f *BmFont) parseCommonSection(keyValues map[string]string) {
	f.lineHeight, _ = strconv.Atoi(keyValues["lineHeight"])
	f.base, _ = strconv.Atoi(keyValues["base"])
	f.pageWidth, _ = strconv.Atoi(keyValues["scaleW"])
	f.pageHeight, _ = strconv.Atoi(keyValues["scaleH"])
	f.packed, _ = strconv.ParseBool(keyValues["packed"])
	f.numPages, _ = strconv.Atoi(keyValues["pages"])

	f.f32scaleLine = 1.0/float32(f.lineHeight)
	f.f32scaleH = 1.0/float32(f.pageHeight)
	f.f32scaleW = 1.0/float32(f.pageWidth)
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

	c.f32x = float32(c.x) * f.f32scaleW
	c.f32y = float32(c.y) * f.f32scaleH
	c.f32width = float32(c.width) * f.f32scaleW
	c.f32height = float32(c.height) * f.f32scaleH

	c.f32lineWidth = float32(c.width) * f.f32scaleLine
	c.f32lineHeight = float32(c.height) * f.f32scaleLine
	c.f32offsetX = float32(c.offsetX) * f.f32scaleLine
	c.f32offsetY = float32(c.offsetY) * f.f32scaleLine
	c.f32advanceX = float32(c.advanceX) * f.f32scaleLine

	c.kernings = make(map[int32]int)
	c.f32kernings = make(map[int32]float32)

	if letter, ok := keyValues["letter"]; ok {
		c.letter = letter
	} else {
		c.letter = string(c.id)
	}

	if c.letter == "space" {
		c.letter = " "
	}
	f.Characters[c.id] = c
}

func (f *BmFont) parseKerningSection(keyValues map[string]string) {
	first, err := strconv.Atoi(keyValues["first"])
	second, err := strconv.Atoi(keyValues["second"])
	amount, err := strconv.Atoi(keyValues["amount"])
	if err != nil {
		log.Printf("Error parsing kerning: %v", err)
		return
	}

	char, ok := f.Characters[int32(first)]
	if !ok {
		log.Printf("Kerning parse error: char %v not found", first)
	}
	char.kernings[int32(second)] = amount
	char.f32kernings[int32(second)] = float32(amount)*f.f32scaleLine
}

var BmSectionRex = regexp.MustCompile("^(\\w+) ")
var BmKeyValueRex = regexp.MustCompile("(\\w+)=\"?([\\w\\s ,._\\-]*)\"?( |$|\")")

func (f *BmFont) tokenizeLine(line string) (string, map[string]string) {
	sectionMatches := BmSectionRex.FindStringSubmatch(line)
	if sectionMatches == nil {
		return "", nil
	}
	sectionName := sectionMatches[1]
	data := BmKeyValueRex.FindAllStringSubmatch(line, -1)

	keyValues := make(map[string]string)
	for _, kv := range data {
		k := kv[1]
		v := strings.Trim(kv[2], " ")
		keyValues[k] = v
	}

	return sectionName, keyValues
}
