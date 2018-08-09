package mdvg

import (
	"os"
	"testing"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func TestRenderer_Draw(t *testing.T) {
	tests := []struct {
		md            string
		filename      string
		width, height vg.Length
	}{
		{
			md:       "Hello world!",
			filename: "testdata/hello.png",
			width:    vg.Points(60),
			height:   vg.Points(20),
		},
		{
			md: `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer interdum
odio libero, sit amet ornare tortor sollicitudin sed. Proin et porta eros. Ut id hendrerit odio.
Etiam quis augue mi. Suspendisse aliquam ligula vitae eros condimentum molestie. Integer consequat
sodales massa, vitae congue arcu lobortis ut. Nullam eros libero, tincidunt non sapien id, convallis
fringilla ante. Aliquam erat volutpat. Suspendisse dui ipsum, tristique non cursus auctor, cursus sed
eros. Morbi in eros magna.

Sed eu nisi ac enim congue egestas. Proin consectetur ante vitae tempus imperdiet.
Etiam ut bibendum urna. Quisque vel nulla eu dui euismod malesuada. Maecenas accumsan
ac leo quis rhoncus. Nullam egestas lectus leo, sed hendrerit erat laoreet vitae.
Mauris et tellus sagittis, laoreet dui id, mollis erat. Duis ultrices facilisis lectus et pretium.`,
			filename: "testdata/lorem.png",
			width:    vg.Points(200),
			height:   vg.Points(210),
		},
		{
			md:       "Here we try s^u^per", //~s~cript and s~ubsc~ript. H~2~O^(2)^ PM~2.5~",
			filename: "testdata/supersub.png",
			width:    vg.Points(150),
			height:   vg.Points(100),
		},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			c := vgimg.New(test.width, test.height)
			dc := draw.New(c)

			font, err := vg.MakeFont("Times-Roman", 10)
			if err != nil {
				t.Fatal(err)
			}

			r := NewRenderer(font)

			r.Draw(dc, []byte(test.md))

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
