package htmlvg

import (
	"bytes"
	"fmt"
	"image/color"
	"strings"
	"sync"

	"golang.org/x/net/html"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// Renderer to renders HTML to a gonum.org/v1/plot/vg.Canvas.
type Renderer struct {
	dc draw.Canvas
	mx sync.Mutex
	at vg.Point

	// TextStyle specifies the text style. The default is black, left-aligned
	// text. Left-alignment is currently the only supported alignment.
	draw.TextStyle

	// PMarginTop and PMarginBottom are the margins before and
	// after paragraphs. Defaults are 0 and 10 points, respectively.
	PMarginTop, PMarginBottom vg.Length

	// SuperscriptPosition, SubscriptPosition, and SuperSubSize
	// are the relative positions and sizes of superscripts and subscripts.
	// Defaults are +0.25, -1.25, and 0.583, respectively.
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
		SuperscriptPosition: 0.25,
		SubscriptPosition:   -1.25,
		SuperSubSize:        0.583,
	}
}

// Draw renders the HTML input to canvas dc.
func (r *Renderer) Draw(dc draw.Canvas, HTML []byte) error {
	r.mx.Lock()
	defer r.mx.Unlock()
	r.dc = dc
	r.at = vg.Point{X: dc.Min.X, Y: dc.Max.Y}
	doc, err := html.Parse(bytes.NewBuffer(HTML))
	if err != nil {
		return fmt.Errorf("htmlvg: %v", err)
	}
	return r.draw(doc)
}

func (r *Renderer) draw(n *html.Node) error {
	switch n.Type {
	case html.ErrorNode:
		return fmt.Errorf("htmlvg: node error: %+v", n)
	case html.TextNode:
		return r.text(n)
	case html.DocumentNode, html.DoctypeNode, html.CommentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := r.draw(c); err != nil {
				return err
			}
		}
	case html.ElementNode:
		if err := r.element(n); err != nil {
			return err
		}
	default:
		panic(fmt.Errorf("invalid node type %v", n.Type))
	}
	return nil
}

// element renders an HTML element.
func (r *Renderer) element(e *html.Node) error {
	switch e.Data {
	case "p":
		return r.paragraph(e)
	case "sup":
		return r.subsuperscript(e, vg.Length(r.SuperscriptPosition))
	case "sub":
		return r.subsuperscript(e, vg.Length(r.SubscriptPosition))
	case "html", "head", "body":
		for c := e.FirstChild; c != nil; c = c.NextSibling {
			if err := r.draw(c); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("htmlvg: '%s' not implemented", e.Data)
	}
}

// paragraph renders an HTML p element.
func (r *Renderer) paragraph(p *html.Node) error {
	r.at = vg.Point{X: 0, Y: r.at.Y - r.PMarginTop}
	for c := p.FirstChild; c != nil; c = c.NextSibling {
		if err := r.draw(c); err != nil {
			return err
		}
	}
	r.at = vg.Point{X: 0, Y: r.at.Y - r.TextStyle.Font.Size - r.PMarginBottom}
	return nil
}

// text renders HTML normal text.
func (r *Renderer) text(t *html.Node) error {
	r.writeLines(t.Data, r.TextStyle)
	return nil
}

// subscript renders subscript text.
func (r *Renderer) subscript(s *html.Node) error {
	r.TextStyle.Font.Size *= vg.Length(r.SuperSubSize)
	r.at.Y += r.TextStyle.Font.Size * vg.Length(r.SubscriptPosition)
	for c := s.FirstChild; c != nil; c = c.NextSibling {
		if err := r.draw(c); err != nil {
			return err
		}
	}
	r.at.Y = r.TextStyle.Font.Size * vg.Length(r.SubscriptPosition)
	r.TextStyle.Font.Size /= vg.Length(r.SuperSubSize)
	return nil
}

// subsuperscript renders superscript or subscript text.
func (r *Renderer) subsuperscript(s *html.Node, position vg.Length) error {
	r.TextStyle.Font.Size *= vg.Length(r.SuperSubSize)
	r.at.Y += r.TextStyle.Font.Size * position
	for c := s.FirstChild; c != nil; c = c.NextSibling {
		if err := r.draw(c); err != nil {
			return err
		}
	}
	r.at.Y -= r.TextStyle.Font.Size * position
	r.TextStyle.Font.Size /= vg.Length(r.SuperSubSize)
	return nil
}

// writeLines writes the given text to the canvas, inserting line breaks
// as necessary.
func (r *Renderer) writeLines(text string, sty draw.TextStyle) {
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
			if r.at.X == r.dc.Min.X { // Remove any trailing space at the beginning of a line.
				out = strings.TrimLeft(out, " ")
			}
			r.dc.FillText(sty, r.at, out)
			r.at.X += sty.Width(out)
			break
		} else if sty.Font.Width(str[lineStart:lineEnd]) > r.dc.Max.X-r.at.X {
			// If we go to the next break, will the line be too long? If so,
			// insert a line break.
			lineStart += len(line)
			if r.at.X == r.dc.Min.X { // Remove any trailing space at the beginning of a line.
				line = strings.TrimLeft(line, " ")
			}
			r.dc.FillText(sty, r.at, line)
			r.at.X = r.dc.Min.X
			r.at.Y -= sty.Font.Size
			line = ""
		} else {
			line = str[lineStart:lineEnd]
		}
	}
}
