package htmlvg

import (
	"os"
	"testing"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func TestRenderer_Draw(t *testing.T) {
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
<p>Here we try s<sup>u</sup>per<sub>s</sub>cript and s<sub>ubsc</sub>ript. H<sub>2</sub>O<sup>(2)</sub> PM<sub>2.5</sub>`,
			filename: "testdata/supersub.png",
			width:    vg.Points(150),
			height:   vg.Points(60),
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
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			c := vgimg.New(test.width, test.height)
			dc := draw.New(c)

			r := NewRenderer()
			r.Size = 10
			r.Font = "Times-Roman"
			r.BoldFont = "Times-Bold"
			r.ItalicFont = "Times-Italic"
			r.BoldItalicFont = "Times-BoldItalic"

			if _, err := r.Draw(dc, []byte(test.html)); err != nil {
				t.Fatal(err)
			}

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

func TestRenderer_Draw_crop(t *testing.T) {
	c := vgimg.New(80, 30)
	dc := draw.New(c)
	dc = draw.Crop(dc, 20, 0, 0, -10)

	r := NewRenderer()
	r.Size = 10
	r.Font = "Times-Roman"
	r.BoldFont = "Times-Bold"
	r.ItalicFont = "Times-Italic"
	r.BoldItalicFont = "Times-BoldItalic"

	if _, err := r.Draw(dc, []byte("hello world!")); err != nil {
		t.Fatal(err)
	}

	w, err := os.Create("testdata/hello_crop.png")
	if err != nil {
		t.Fatal(err)
	}
	pngc := vgimg.PngCanvas{Canvas: c}
	if _, err := pngc.WriteTo(w); err != nil {
		t.Fatal(err)
	}
}
