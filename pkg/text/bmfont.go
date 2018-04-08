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
	id             int
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
}

type BmFont struct {
	pageFiles      map[int]string
	charactersList []BmChar
	Characters     map[string]BmChar

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
}

func NewBmFontFromFile(fileName string) *BmFont {
	f := &BmFont{}

	f.pageFiles = make(map[int]string)
	f.Characters = make(map[string]BmChar)

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
}

func (f *BmFont) parsePageSection(keyValues map[string]string) {
	id, _ := strconv.Atoi(keyValues["id"])
	f.pageFiles[id] = keyValues["file"]
}

func (f *BmFont) parseCharSection(keyValues map[string]string) {
	c := BmChar{}
	c.id, _ = strconv.Atoi(keyValues["id"])
	c.x, _ = strconv.Atoi(keyValues["x"])
	c.y, _ = strconv.Atoi(keyValues["y"])
	c.width, _ = strconv.Atoi(keyValues["width"])
	c.height, _ = strconv.Atoi(keyValues["height"])
	c.offsetX, _ = strconv.Atoi(keyValues["xoffset"])
	c.offsetY, _ = strconv.Atoi(keyValues["yoffset"])
	c.advanceX, _ = strconv.Atoi(keyValues["xadvance"])
	c.pageIndex, _ = strconv.Atoi(keyValues["page"])
	c.textureChannel, _ = strconv.Atoi(keyValues["chnl"])

	if letter, ok := keyValues["letter"]; ok {
		c.letter = letter
	} else {
		c.letter = string(c.id)
	}

	if c.letter == "space" {
		c.letter = " "
	}
	f.Characters[c.letter] = c
}

func (f *BmFont) tokenizeLine(line string) (string, map[string]string) {
	sectionRex := regexp.MustCompile("^(\\w+) ")
	keyValueRex := regexp.MustCompile("(\\w+)=\"?([\\w\\s ,._-]*)\"?[ |$|\"]")
	sectionMatches := sectionRex.FindStringSubmatch(line)
	if sectionMatches == nil {
		return "", nil
	}
	sectionName := sectionMatches[1]
	data := keyValueRex.FindAllStringSubmatch(line, -1)

	keyValues := make(map[string]string)
	for _, kv := range data {
		k := kv[1]
		v := strings.Trim(kv[2], " ")
		keyValues[k] = v
	}

	return sectionName, keyValues
}
