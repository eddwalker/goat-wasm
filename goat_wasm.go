package main

import (
	"syscall/js"
	"fmt"

  "io"
  "flag"
  "github.com/blampe/goat"
  "strings"
)

var document = js.Global().Get("document")

func getElementByID(id string) js.Value {
	return document.Call("getElementById", id)
}

func renderEditor(parent js.Value) js.Value {
	editorMarkup := `
	<div id="editor" style="display: flex; flex-flow: row wrap; margin: 0 10px 0 0">
			<textarea id="markdown" style="width: 100%; height: 200px"></textarea>
			<div id="preview" style="width: 100%;"></div>
			<button id="render" style="display:none">Render Markdown</button>
		</div>
	`
	parent.Call("insertAdjacentHTML", "beforeend", editorMarkup)
	return getElementByID("editor")
}

type MyReader struct {
    src []byte
    pos int
}

type MyWriter struct {
    dst []byte
    pos int
}

func (r *MyReader) Read(dst []byte) (n int, err error) {
    n = copy(dst, r.src[r.pos:])
    r.pos += n
    if r.pos == len(r.src) {
        return n, io.EOF
    }
    return
}

func (w *MyWriter) Write(src []byte) (n int, err error) {
    n = copy(src, w.dst[w.pos:])
    w.pos += n
    if w.pos == len(w.dst) {
        return n, io.EOF
    }
    return
}

func NewMyReader(b []byte) *MyReader { return &MyReader{b, 0} }

func NewMyWriter(b []byte) *MyWriter { return &MyWriter{b, 0} }

func main() {
	fmt.Println("wasm init.. start")

	var (
		inputFilename,
		outputFilename,
		svgColorLightScheme,
		svgColorDarkScheme string
	)

	flag.StringVar(&inputFilename, "i", "", "Input filename (default stdin)")
	flag.StringVar(&outputFilename, "o", "", "Output filename (default stdout for SVG)")
	flag.StringVar(&svgColorLightScheme, "sls", "#000000", `short for -svg-color-light-scheme`)
	flag.StringVar(&svgColorLightScheme, "svg-color-light-scheme", "#000000",
		`See help for -svg-color-dark-scheme`)
	flag.StringVar(&svgColorDarkScheme, "sds", "#FFFFFF", `short for -svg-color-dark-scheme`)
	flag.StringVar(&svgColorDarkScheme, "svg-color-dark-scheme", "#FFFFFF",
		`Goat's SVG output attempts to learn something about the background being
 drawn on top of by means of a CSS @media query, which returns a string.
 If the string is "dark", Goat draws with the color specified by
 this option; otherwise, Goat draws with the color specified by option
 -svg-color-light-scheme.

 See https://developer.mozilla.org/en-US/docs/Web/CSS/@media/prefers-color-scheme
`)
	flag.Parse()

	quit := make(chan struct{}, 0)

	editor := renderEditor(document.Get("body"))
	markdown := getElementByID("markdown")
	preview := getElementByID("preview")
	renderButton := getElementByID("markdown")
	fmt.Println("wasm point.. pre interface")
	renderButton.Set("onkeyup", js.FuncOf(func(js.Value, []js.Value) interface{} {
                outputFilename2 := new(strings.Builder)
		goat.BuildAndWriteSVG(
                   NewMyReader([]byte(markdown.Get("value").String())),
		   outputFilename2 ,
                   svgColorLightScheme,
                   svgColorDarkScheme)
                preview.Set("innerHTML", string(outputFilename2.String()))
		return nil
	}))
	renderButton.Set("onkeydown", js.FuncOf(func(js.Value, []js.Value) interface{} {
                outputFilename2 := new(strings.Builder)
		goat.BuildAndWriteSVG(
                   NewMyReader([]byte(markdown.Get("value").String())),
		   outputFilename2 ,
                   svgColorLightScheme,
                   svgColorDarkScheme)
                preview.Set("innerHTML", string(outputFilename2.String()))
		return nil
	}))

	fmt.Println("wasm point.. post interface")

	<-quit
	editor.Call("remove")
	fmt.Println("wasm init.. done")
}
