package htmlvg

import (
	"os"
	"testing"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
	"gonum.org/v1/plot/vg/vgpdf"
)

func TestWrapper_Draw(t *testing.T) {
	tests := []struct {
		html          string
		filename      string
		width, height vg.Length
	}{
		{
			html:     "Hello world!",
			filename: "testdata/hello.png",
			width:    vg.Points(60),
			height:   vg.Points(20),
		},
		{
			html: `<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer interdum
odio libero, sit amet ornare tortor sollicitudin sed. Proin et porta eros. Ut id hendrerit odio.
Etiam quis augue mi. Suspendisse aliquam ligula vitae eros condimentum molestie. Integer consequat
sodales massa, vitae congue arcu lobortis ut. Nullam eros libero, tincidunt non sapien id, convallis
fringilla ante. Aliquam erat volutpat. Suspendisse dui ipsum, tristique non cursus auctor, cursus sed
eros. Morbi in eros magna.</p>

<p>Sed eu nisi ac enim congue egestas. Proin consectetur ante vitae tempus imperdiet.
Etiam ut bibendum urna. Quisque vel nulla eu dui euismod malesuada. Maecenas accumsan
ac leo quis rhoncus. Nullam egestas lectus leo, sed hendrerit erat laoreet vitae.
Mauris et tellus sagittis, laoreet dui id, mollis erat. Duis ultrices facilisis lectus et pretium.</p>`,
			filename: "testdata/lorem.png",
			width:    vg.Points(200),
			height:   vg.Points(210),
		},
		{
			html: `<p>Superscript and subscript don't currently work well with line breaks.</p>
<p>Here we try s<sup>u</sup>per<sub>s</sub>cript and s<sub>ubsc</sub>ript. H<sub>2</sub>O<sub>(2)</sub> PM<sub>2.5</sub>`,
			filename: "testdata/supersub.png",
			width:    vg.Points(150),
			height:   vg.Points(70),
		},
		{
			html: `<h1>Long Long Long Long Heading 1</h1>
<p>Paragraph</p>
<h2>Heading 2</h2>
<p>Paragraph</p>
<h3>Heading 3</h3>
<p>Paragraph</p>
<h4>Heading 4</h4>
<p>Paragraph</p>
<h5>Heading 5</h5>
<p>Paragraph</p>
<h6>Heading 6</h6>
<p>Paragraph</p>`,
			filename: "testdata/headings.png",
			width:    vg.Points(100),
			height:   vg.Points(290),
		},
		{
			html:     `<p>Here we try some <strong>bold</strong> and <em>italic</em> text.</p>`,
			filename: "testdata/bolditalic.png",
			width:    vg.Points(100),
			height:   vg.Points(25),
		},
		{
			html: `<p>Paragraph 1</p>
<hr>
<p>Paragraph 2</p>`,
			filename: "testdata/hr.png",
			width:    vg.Points(100),
			height:   vg.Points(50),
		},
		{
			html:     `<p>This is a short paragraph with a long woooooooooooord.</p>`,
			filename: "testdata/long_word.png",
			width:    vg.Points(100),
			height:   vg.Points(50),
		},

		{
			html:     `<p>Here-we're-testing-words-with-dashes-words-with-dashes-words-with-dashes.</p>`,
			filename: "testdata/dashes.png",
			width:    vg.Points(100),
			height:   vg.Points(50),
		},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			c := vgimg.New(test.width, test.height)
			dc := draw.New(c)

			r := WrapCanvas(dc)
			r.FontSize = 10
			r.Font = "Times-Roman"
			r.BoldFont = "Times-Bold"
			r.ItalicFont = "Times-Italic"
			r.BoldItalicFont = "Times-BoldItalic"

			r.FillString(vg.Font{}, vg.Point{X: dc.Min.X, Y: dc.Max.Y - r.FontSize}, test.html)

			w, err := os.Create(test.filename)
			if err != nil {
				t.Fatal(err)
			}
			pngc := vgimg.PngCanvas{Canvas: c}
			if _, err := pngc.WriteTo(w); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestWrapper_Draw_crop(t *testing.T) {
	c := vgimg.New(80, 30)
	dc := draw.New(c)
	dc = draw.Crop(dc, 20, 0, 0, -10)

	r := WrapCanvas(dc)
	r.FontSize = 10
	r.Font = "Times-Roman"
	r.BoldFont = "Times-Bold"
	r.ItalicFont = "Times-Italic"
	r.BoldItalicFont = "Times-BoldItalic"

	r.FillString(vg.Font{}, vg.Point{X: dc.Min.X, Y: dc.Max.Y - r.FontSize}, "hello world!")

	w, err := os.Create("testdata/hello_crop.png")
	if err != nil {
		t.Fatal(err)
	}
	pngc := vgimg.PngCanvas{Canvas: c}
	if _, err := pngc.WriteTo(w); err != nil {
		t.Fatal(err)
	}
}

func TestWrapper_pdf(t *testing.T) {
	c := vgpdf.New(100, 80)
	dc := draw.New(c)

	r := WrapCanvas(dc)

	r.FillString(vg.Font{}, vg.Point{X: dc.Min.X, Y: dc.Max.Y - r.FontSize}, `
<h1>Heading 1</h1>
<p>Here is some <strong>bold</strong> and <em>italic</em> text.</p>
`)

	w, err := os.Create("testdata/pdf_test.pdf")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.WriteTo(w); err != nil {
		t.Fatal(err)
	}
}

func TestWrapper_Plot(t *testing.T) {
	p, err := plot.New()
	if err != nil {
		t.Fatal(err)
	}
	l, err := plotter.NewLine(plotter.XYs{{X: 0, Y: 0}, {X: 1, Y: 1}})
	if err != nil {
		t.Fatal(err)
	}
	p.Add(l)
	p.Title.Text = "<em>This</em> is the <strong>Title</strong>"
	p.X.Label.Text = "H<sub>2</sub>O"
	p.Y.Label.Text = "x<sub>y</sub><sup>z</sup>x<sup>z</sup>"

	// With the current implementation, we have to manually adjust the
	// alignment using guess-and-check. This is because the HTML tags
	// get counted as part of the text width.
	p.Title.TextStyle.XAlign = -0.17
	p.Y.Label.TextStyle.XAlign = -0.05
	p.X.Label.TextStyle.XAlign = -0.1

	c := vgimg.New(100, 100)
	dc := draw.New(c)
	h := WrapCanvas(dc)
	h.WrapLines = false
	hc := draw.New(h)
	p.Draw(hc)

	w, err := os.Create("testdata/plot.png")
	if err != nil {
		t.Fatal(err)
	}
	pngc := vgimg.PngCanvas{Canvas: c}
	if _, err := pngc.WriteTo(w); err != nil {
		t.Fatal(err)
	}
}
