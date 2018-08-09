package mdvg

import (
	"bytes"
	"fmt"
	"image/color"
	"strings"
	"sync"

	"github.com/miekg/mmark"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Renderer implements github.com/miekg/mmark.Renderer to render markdown
// to a gonum.org/v1/plot/vg.Canvas.
type Renderer struct {
	mmark.Renderer // We need this because some interface methods have unexported arguments.

	dc draw.Canvas
	mx sync.Mutex
	at vg.Point

	draw.TextStyle

	// PMarginTop and PMarginBottom are the margins before and
	// after paragraphs. Defaults are 0 and 10 points, respectively.
	PMarginTop, PMarginBottom vg.Length

	// SuperscriptPosition, SubscriptPosition, and SuperSubSize
	// are the relative positions and sizes of superscripts and subscripts.
	// Defaults are +0.333, -0.333, and 0.583, respectively.
	SuperscriptPosition, SubscriptPosition, SuperSubSize float64
}

func NewRenderer(font vg.Font) *Renderer {
	return &Renderer{
		TextStyle: draw.TextStyle{
			Color:  color.Black,
			Font:   font,
			XAlign: draw.XLeft,
			YAlign: draw.YTop,
		},
		PMarginBottom:       vg.Points(10),
		SuperscriptPosition: 0.333,
		SubscriptPosition:   -0.333,
		SuperSubSize:        0.583,
	}
}

// Draw renders the markdown input md to canvas dc.
func (r *Renderer) Draw(dc draw.Canvas, md []byte) {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.dc = dc
	r.at = vg.Point{X: dc.Min.X, Y: dc.Max.Y}
	mmark.Parse(md, r, 0)
}

// block-level callbacks
func (r *Renderer) BlockCode(out *bytes.Buffer, text []byte, lang string, caption []byte, subfigure bool, callouts bool) {
	panic("BlockCode is not implemented")
}
func (r *Renderer) BlockQuote(out *bytes.Buffer, text []byte, attribution []byte) {
	panic("BlockQuote is not implemented")
}
func (r *Renderer) BlockHtml(out *bytes.Buffer, text []byte) { panic("BlockHtml is not implemented") }
func (r *Renderer) CommentHtml(out *bytes.Buffer, text []byte) {
	panic("CommentHtml is not implemented")
}

// SpecialHeader is used for Abstract and Preface. The what string contains abstract or preface.
func (r *Renderer) SpecialHeader(out *bytes.Buffer, what []byte, text func() bool, id string) {
	panic("SpecialHeader is not implemented")
}

// Note is use for typesetting notes.
func (r *Renderer) Note(out *bytes.Buffer, text func() bool, id string) {
	panic("Note is not implemented")
}
func (r *Renderer) Part(out *bytes.Buffer, text func() bool, id string) {
	panic("Part is not implemented")
}
func (r *Renderer) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	panic("Header is not implemented")
}
func (r *Renderer) HRule(out *bytes.Buffer) { panic("HRule is not implemented") }
func (r *Renderer) List(out *bytes.Buffer, text func() bool, flags, start int, group []byte) {
	panic("List is not implemented")
}
func (r *Renderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	panic("ListItem is not implemented")
}
func (r *Renderer) Paragraph(out *bytes.Buffer, text func() bool, flags int) {
	r.at = vg.Point{X: 0, Y: r.at.Y - r.PMarginTop}
	ok := text()
	if !ok {
		panic("this is not ok")
	}
	r.at = vg.Point{X: 0, Y: r.at.Y - r.TextStyle.Font.Size - r.PMarginBottom}
}

func (r *Renderer) Table(out *bytes.Buffer, header []byte, body []byte, footer []byte, columnData []int, caption []byte) {
	panic("Table is not implemented")
}
func (r *Renderer) TableRow(out *bytes.Buffer, text []byte) { panic("TableRow is not implemented") }
func (r *Renderer) TableHeaderCell(out *bytes.Buffer, text []byte, flags, colspan int) {
	panic("TableHeaderCell is not implemented")
}
func (r *Renderer) TableCell(out *bytes.Buffer, text []byte, flags, colspan int) {
	panic("TableCell is not implemented")
}

func (r *Renderer) Footnotes(out *bytes.Buffer, text func() bool) {
	panic("Footnotes is not implemented")
}
func (r *Renderer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	panic("FootnoteItem is not implemented")
}
func (r *Renderer) Aside(out *bytes.Buffer, text []byte) { panic("Aside is not implemented") }
func (r *Renderer) Figure(out *bytes.Buffer, text []byte, caption []byte) {
	panic("Figure is not implemented")
}

// Span-level callbacks
func (r *Renderer) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	panic("AutoLink is not implemented")
}
func (r *Renderer) CodeSpan(out *bytes.Buffer, text []byte) { panic("CodeSpan is not implemented") }

// CalloutText is called when a callout is seen in the text. Id is the text
// seen between < and > and ids references the callout counter(s) in the code.
func (r *Renderer) CalloutText(out *bytes.Buffer, id string, ids []string) {
	panic("CalloutText is not implemented")
}

// Called when a callout is seen in a code block. Index is the callout counter, id
// is the number seen between < and >.
func (r *Renderer) CalloutCode(out *bytes.Buffer, index, id string) {
	panic("CalloutCode is not implemented")
}
func (r *Renderer) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	panic("DoubleEmphasis is not implemented")
}
func (r *Renderer) Emphasis(out *bytes.Buffer, text []byte) { panic("Emphasis is not implemented") }
func (r *Renderer) Subscript(out *bytes.Buffer, text []byte) {
	r.at.Y -= r.TextStyle.Font.Size * vg.Length(r.SubscriptPosition)
	ts := r.TextStyle
	ts.Font.Size *= vg.Length(r.SuperSubSize)
	r.writeLines(string(text), ts)
}
func (r *Renderer) Superscript(out *bytes.Buffer, text []byte) {
	r.at.Y += r.TextStyle.Font.Size * vg.Length(r.SuperscriptPosition)
	ts := r.TextStyle
	ts.Font.Size *= vg.Length(r.SuperSubSize)
	//r.at.Y -= r.TextStyle.Font.Size * vg.Length(r.SuperscriptPosition)
}
func (r *Renderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte, subfigure bool) {
	panic("Image is not implemented")
}
func (r *Renderer) LineBreak(out *bytes.Buffer) { panic("LineBreak is not implemented") }
func (r *Renderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	panic("Link is not implemented")
}
func (r *Renderer) RawHtmlTag(out *bytes.Buffer, tag []byte) { panic("RawHtmlTag is not implemented") }
func (r *Renderer) TripleEmphasis(out *bytes.Buffer, text []byte) {
	panic("TripleEmphasis is not implemented")
}
func (r *Renderer) StrikeThrough(out *bytes.Buffer, text []byte) {
	panic("StrikeThrough is not implemented")
}
func (r *Renderer) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	panic("FootnoteRef is not implemented")
}
func (r *Renderer) Index(out *bytes.Buffer, primary, secondary []byte, prim bool) {
	panic("Index is not implemented")
}
func (r *Renderer) Citation(out *bytes.Buffer, link, title []byte) {
	panic("Citation is not implemented")
}
func (r *Renderer) Abbreviation(out *bytes.Buffer, abbr, title []byte) {
	panic("Abbreviation is not implemented")
}
func (r *Renderer) Example(out *bytes.Buffer, index int) { panic("Example is not implemented") }
func (r *Renderer) Math(out *bytes.Buffer, text []byte, display bool) {
	panic("Math is not implemented")
}

// Low-level callbacks
func (r *Renderer) Entity(out *bytes.Buffer, entity []byte) { panic("Entity is not implemented") }
func (r *Renderer) NormalText(out *bytes.Buffer, text []byte) {
	r.writeLines(string(text), r.TextStyle)
}

// writeLines writes the given text to the canvas, inserting line breaks
// as necessary.
func (r *Renderer) writeLines(text string, sty draw.TextStyle) {
	fmt.Println("bbb", text)
	splitFunc := func(r rune) bool {
		return r == ' ' || r == '-' // Function for choosing possible line breaks.
	}

	str := strings.Replace(text, "\n", " ", -1)

	var lineStart int
	var line string
	for {
		nextBreak := -1
		if len(str) > 1 {
			nextBreak = strings.IndexFunc(str[lineStart+len(line)+1:], splitFunc)
		}
		lineEnd := lineStart + len(line) + 1 + nextBreak
		if nextBreak == -1 {
			out := str[lineStart:]
			if r.at.X == 0 { // Remove any trailing space at the beginning of a line.
				out = strings.TrimLeft(out, " ")
			}
			r.dc.FillText(sty, r.at, out)
			r.at.X += sty.Width(out)
			break
		} else if sty.Font.Width(str[lineStart:lineEnd]) > r.dc.Max.X-r.at.X {
			// If we go to the next break, will the line be too long? If so,
			// insert a line break.
			lineStart += len(line)
			if r.at.X == 0 { // Remove any trailing space at the beginning of a line.
				line = strings.TrimLeft(line, " ")
			}
			r.dc.FillText(sty, r.at, line)
			r.at.X = 0
			r.at.Y -= sty.Font.Size
			line = ""
		} else {
			line = str[lineStart:lineEnd]
		}
	}
}

// Header and footer
func (r *Renderer) DocumentHeader(out *bytes.Buffer, start bool) {
	// Don't currently do anything for the header.
}
func (r *Renderer) DocumentFooter(out *bytes.Buffer, start bool) {
	// Don't currently do anything for the footer.
}

// Frontmatter, mainmatter or backmatter
func (r *Renderer) DocumentMatter(out *bytes.Buffer, matter int) {
	panic("DocumentMatter is not implemented")
}
