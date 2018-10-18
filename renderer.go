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
	dc         draw.Canvas
	mx         sync.Mutex
	at         vg.Point
	sty        draw.TextStyle
	lineHeight vg.Length

	// Font, BoldFont, and ItalicFont are the names of the fonts
	// to be used for regular, bold, italic, and bold-italic text, respectively.
	// The defaults are "Helvetica", "Helvetica-Bold", "Helvetica-Oblique",
	// and "Helvetica-BoldOblique", respectively.
	Font, BoldFont, ItalicFont, BoldItalicFont string

	// Size is the base font size. The default is 12.
	Size vg.Length

	// Color is the font color.
	Color color.Color

	// PMarginTop and PMarginBottom are the margins before and
	// after paragraphs. Defaults are 0 and 0.833 text height units, respectively.
	PMarginTop, PMarginBottom float64

	// H1Scale - H6Scale are font size scaling factors for headings.
	// Defaults are 2.0, 1.5, 1.25, 1, 1, and 1, respectively.
	H1Scale, H2Scale, H3Scale, H4Scale, H5Scale, H6Scale float64

	// H*PMarginTop are the margins above headings.
	// Default values are 1, 0.833, 0.75, 0.5, 0.5, and 0.5 respectively.
	H1MarginTop, H2MarginTop, H3MarginTop, H4MarginTop, H5MarginTop, H6MarginTop float64

	// H*PMarginBottom are the margins below headings.
	// Default values are 1, 0.833, 0.75, 0.5, 0.5, and 0.5 respectively.
	H1MarginBottom, H2MarginBottom, H3MarginBottom, H4MarginBottom, H5MarginBottom, H6MarginBottom float64

	// H*Bold specifies whether the headings should be bold-face.
	// The defaults are true, true, true, true, true, and false, respectively.
	H1Bold, H2Bold, H3Bold, H4Bold, H5Bold, H6Bold bool

	// SuperscriptPosition, SubscriptPosition, and SuperSubScale
	// are the relative positions and sizes of superscripts and subscripts.
	// Defaults are +0.25, -1.25, and 0.583, respectively.
	SuperscriptPosition, SubscriptPosition, SuperSubScale float64

	// HRMarginTop and Bottom specify the spacing above and below horizontal
	// rules. Defaults are 0.833 text height units.
	HRMarginTop, HRMarginBottom float64

	// HRScale specifies the width of horizontal rules. The
	// default is 0.1 text height units.
	HRScale float64

	// HRColor specifies the color of horizontal rules. The
	// default is black.
	HRColor color.Color
}

func NewRenderer() *Renderer {
	r := &Renderer{
		Color:               color.Black,
		Size:                vg.Points(12),
		PMarginBottom:       0.833,
		SuperscriptPosition: 0.25,
		SubscriptPosition:   -1.25,
		SuperSubScale:       0.583,
	}
	r.H1Scale, r.H2Scale, r.H3Scale, r.H4Scale, r.H5Scale, r.H6Scale = 2.0, 1.5, 1.25, 1, 1, 1
	r.H1MarginTop, r.H2MarginTop, r.H3MarginTop, r.H4MarginTop, r.H5MarginTop, r.H6MarginTop =
		1, 0.833, 0.75, 0.5, 0.5, 0.5
	r.H1MarginBottom, r.H2MarginBottom, r.H3MarginBottom, r.H4MarginBottom, r.H5MarginBottom, r.H6MarginBottom =
		1, 0.833, 0.75, 0.5, 0.5, 0.5
	r.H1Bold, r.H2Bold, r.H3Bold, r.H4Bold, r.H5Bold, r.H6Bold =
		true, true, true, true, true, false
	r.Font, r.BoldFont, r.ItalicFont, r.BoldItalicFont =
		"Helvetica", "Helvetica-Bold", "Helvetica-Oblique", "Helvetica-BoldOblique"
	r.HRMarginTop, r.HRMarginBottom = 0.833, 0.833
	r.HRScale = 0.1
	r.HRColor = color.Black
	return r
}

// Draw renders the HTML input to canvas dc.
// It returns the canvas coordinates of the cursor after drawing.
func (r *Renderer) Draw(dc draw.Canvas, HTML []byte) (vg.Point, error) {
	f, err := vg.MakeFont(r.Font, r.Size)
	if err != nil {
		return r.at, err
	}
	r.sty = draw.TextStyle{
		Font:   f,
		Color:  r.Color,
		XAlign: draw.XLeft,
		YAlign: draw.YTop,
	}

	r.mx.Lock()
	defer r.mx.Unlock()
	r.dc = dc
	r.at = vg.Point{X: dc.Min.X, Y: dc.Max.Y}
	doc, err := html.Parse(bytes.NewBuffer(HTML))
	if err != nil {
		return r.at, fmt.Errorf("htmlvg: %v", err)
	}
	return r.draw(doc)
}

func (r *Renderer) draw(n *html.Node) (vg.Point, error) {
	switch n.Type {
	case html.ErrorNode:
		return r.at, fmt.Errorf("htmlvg: node error: %+v", n)
	case html.TextNode:
		return r.text(n)
	case html.DocumentNode, html.DoctypeNode, html.CommentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if at, err := r.draw(c); err != nil {
				return at, err
			}
		}
	case html.ElementNode:
		if at, err := r.element(n); err != nil {
			return at, err
		}
	default:
		panic(fmt.Errorf("invalid node type %v", n.Type))
	}
	return r.at, nil
}

// element renders an HTML element.
func (r *Renderer) element(e *html.Node) (vg.Point, error) {
	switch e.Data {
	case "p":
		return r.paragraph(e)
	case "h1":
		return r.heading(e, r.H1Scale, r.H1MarginTop, r.H1MarginBottom, r.H1Bold)
	case "h2":
		return r.heading(e, r.H2Scale, r.H2MarginTop, r.H2MarginBottom, r.H2Bold)
	case "h3":
		return r.heading(e, r.H3Scale, r.H3MarginTop, r.H3MarginBottom, r.H3Bold)
	case "h4":
		return r.heading(e, r.H4Scale, r.H4MarginTop, r.H4MarginBottom, r.H4Bold)
	case "h5":
		return r.heading(e, r.H5Scale, r.H5MarginTop, r.H5MarginBottom, r.H5Bold)
	case "h6":
		return r.heading(e, r.H6Scale, r.H6MarginTop, r.H6MarginBottom, r.H6Bold)
	case "strong", "b":
		return r.newFont(e, r.BoldFont)
	case "em", "i":
		return r.newFont(e, r.ItalicFont)
	case "hr":
		return r.hr()
	case "sup":
		return r.subsuperscript(e, vg.Length(r.SuperscriptPosition))
	case "sub":
		return r.subsuperscript(e, vg.Length(r.SubscriptPosition))
	case "html", "head", "body":
		for c := e.FirstChild; c != nil; c = c.NextSibling {
			if at, err := r.draw(c); err != nil {
				return at, err
			}
		}
		return r.at, nil
	default:
		return r.at, fmt.Errorf("htmlvg: '%s' not implemented", e.Data)
	}
}

// paragraph renders an HTML p element.
func (r *Renderer) paragraph(p *html.Node) (vg.Point, error) {
	r.at = vg.Point{X: r.dc.Min.X, Y: r.at.Y - r.Size*vg.Length(r.PMarginTop)}
	r.lineHeight = r.sty.Font.Size
	for c := p.FirstChild; c != nil; c = c.NextSibling {
		if at, err := r.draw(c); err != nil {
			return at, err
		}
	}
	r.at = vg.Point{X: r.dc.Min.X, Y: r.at.Y - r.Size*(1+vg.Length(r.PMarginBottom))}
	return r.at, nil
}

// text renders HTML normal text.
func (r *Renderer) text(t *html.Node) (vg.Point, error) {
	r.writeLines(t.Data, r.sty)
	return r.at, nil
}

// subsuperscript renders superscript or subscript text.
func (r *Renderer) subsuperscript(s *html.Node, position vg.Length) (vg.Point, error) {
	r.sty.Font.Size *= vg.Length(r.SuperSubScale)
	r.at.Y += r.sty.Font.Size * position
	for c := s.FirstChild; c != nil; c = c.NextSibling {
		if at, err := r.draw(c); err != nil {
			return at, err
		}
	}
	r.at.Y -= r.sty.Font.Size * position
	r.sty.Font.Size /= vg.Length(r.SuperSubScale)
	return r.at, nil
}

func (r *Renderer) heading(h *html.Node, scale, marginTop, marginBottom float64, bold bool) (vg.Point, error) {
	if bold {
		f := r.sty.Font
		if err := r.sty.Font.SetName(r.BoldFont); err != nil {
			return r.at, err
		}
		defer func() {
			r.sty.Font = f
		}()
	}
	r.at.X = r.dc.Min.X
	r.at.Y -= r.Size * vg.Length(marginTop)
	r.sty.Font.Size *= vg.Length(scale)
	r.lineHeight = r.sty.Font.Size
	for c := h.FirstChild; c != nil; c = c.NextSibling {
		if at, err := r.draw(c); err != nil {
			return at, err
		}
	}
	r.at.Y -= r.sty.Font.Size * vg.Length(marginBottom)
	r.sty.Font.Size /= vg.Length(scale)
	r.at.X = r.dc.Min.X
	r.at.Y -= r.Size * vg.Length(marginBottom)
	return r.at, nil
}

func (r *Renderer) hr() (vg.Point, error) {
	r.at.Y -= r.Size * vg.Length(r.HRMarginTop)
	r.dc.StrokeLine2(draw.LineStyle{
		Color: r.HRColor,
		Width: r.Size * vg.Length(r.HRScale),
	}, r.dc.Min.X, r.at.Y, r.dc.Max.X, r.at.Y)
	r.at.Y -= r.Size * vg.Length(r.HRMarginBottom)
	return r.at, nil
}

// newFont temporarily changes the font.
func (r *Renderer) newFont(n *html.Node, font string) (vg.Point, error) {
	f := r.sty.Font
	if err := r.sty.Font.SetName(font); err != nil {
		return r.at, err
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if at, err := r.draw(c); err != nil {
			return at, err
		}
	}
	r.sty.Font = f
	return r.at, nil
}

// writeLines writes the given text to the canvas, inserting line breaks
// as necessary.
func (r *Renderer) writeLines(text string, sty draw.TextStyle) {
	splitFunc := func(r rune) bool {
		return r == ' ' || r == '-' // Function for choosing possible line breaks.
	}

	str := strings.Replace(text, " \n ", " ", -1)
	str = strings.Replace(str, " \n", " ", -1)
	str = strings.Replace(str, "\n ", " ", -1)
	str = strings.Replace(str, "\n", " ", -1)

	var lineStart int
	var line string
	for {
		nextBreak := -1
		if len(str) > 1 {
			nextBreak = strings.IndexFunc(str[lineStart+len(line)+1:], splitFunc)
		}
		var lineEnd int
		if nextBreak == -1 {
			lineEnd = len(str)
		} else {
			lineEnd = lineStart + len(line) + 1 + nextBreak
		}

		if sty.Font.Width(str[lineStart:lineEnd]) > r.dc.Max.X-r.at.X {
			// If we go to the next break, will the line be too long? If so,
			// insert a line break.
			lineStart += len(line)
			if r.at.X == r.dc.Min.X { // Remove any trailing space at the beginning of a line.
				line = strings.TrimLeft(line, " ")
			}
			r.dc.FillText(sty, r.at, line)
			r.newLine()
			line = ""
		} else {
			line = str[lineStart:lineEnd]
		}
		if nextBreak == -1 {
			out := str[lineStart:]
			if r.at.X == r.dc.Min.X { // Remove any trailing space at the beginning of a line.
				out = strings.TrimLeft(out, " ")
			}
			r.dc.FillText(sty, r.at, out)
			r.at.X += sty.Width(out)
			break
		}
	}
}

func (r *Renderer) newLine() {
	r.at.X = r.dc.Min.X
	r.at.Y -= r.lineHeight
}
