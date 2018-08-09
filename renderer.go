package mdvg

import "bytes"

// Renderer implements github.com/miekg/mmark.Renderer to render markdown
// to a gonum.org/v1/plot/vg.Canvas.
type Renderer struct {
}

// block-level callbacks
func (r *Renderer) BlockCode(out *bytes.Buffer, text []byte, lang string, caption []byte, subfigure bool, callouts bool) {
}
func (r *Renderer) BlockQuote(out *bytes.Buffer, text []byte, attribution []byte) {}
func (r *Renderer) BlockHtml(out *bytes.Buffer, text []byte)                      {}
func (r *Renderer) CommentHtml(out *bytes.Buffer, text []byte)                    {}

// SpecialHeader is used for Abstract and Preface. The what string contains abstract or preface.
func (r *Renderer) SpecialHeader(out *bytes.Buffer, what []byte, text func() bool, id string) {}

// Note is use for typesetting notes.
func (r *Renderer) Note(out *bytes.Buffer, text func() bool, id string)                      {}
func (r *Renderer) Part(out *bytes.Buffer, text func() bool, id string)                      {}
func (r *Renderer) Header(out *bytes.Buffer, text func() bool, level int, id string)         {}
func (r *Renderer) HRule(out *bytes.Buffer)                                                  {}
func (r *Renderer) List(out *bytes.Buffer, text func() bool, flags, start int, group []byte) {}
func (r *Renderer) ListItem(out *bytes.Buffer, text []byte, flags int)                       {}
func (r *Renderer) Paragraph(out *bytes.Buffer, text func() bool, flags int)                 {}

func (r *Renderer) Table(out *bytes.Buffer, header []byte, body []byte, footer []byte, columnData []int, caption []byte) {
}
func (r *Renderer) TableRow(out *bytes.Buffer, text []byte)                            {}
func (r *Renderer) TableHeaderCell(out *bytes.Buffer, text []byte, flags, colspan int) {}
func (r *Renderer) TableCell(out *bytes.Buffer, text []byte, flags, colspan int)       {}

func (r *Renderer) Footnotes(out *bytes.Buffer, text func() bool)                {}
func (r *Renderer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {}
func (r *Renderer) TitleBlockTOML(out *bytes.Buffer, data *title)                {}
func (r *Renderer) Aside(out *bytes.Buffer, text []byte)                         {}
func (r *Renderer) Figure(out *bytes.Buffer, text []byte, caption []byte)        {}

// Span-level callbacks
func (r *Renderer) AutoLink(out *bytes.Buffer, link []byte, kind int) {}
func (r *Renderer) CodeSpan(out *bytes.Buffer, text []byte)           {}

// CalloutText is called when a callout is seen in the text. Id is the text
// seen between < and > and ids references the callout counter(s) in the code.
func (r *Renderer) CalloutText(out *bytes.Buffer, id string, ids []string) {}

// Called when a callout is seen in a code block. Index is the callout counter, id
// is the number seen between < and >.
func (r *Renderer) CalloutCode(out *bytes.Buffer, index, id string)                                {}
func (r *Renderer) DoubleEmphasis(out *bytes.Buffer, text []byte)                                  {}
func (r *Renderer) Emphasis(out *bytes.Buffer, text []byte)                                        {}
func (r *Renderer) Subscript(out *bytes.Buffer, text []byte)                                       {}
func (r *Renderer) Superscript(out *bytes.Buffer, text []byte)                                     {}
func (r *Renderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte, subfigure bool) {}
func (r *Renderer) LineBreak(out *bytes.Buffer)                                                    {}
func (r *Renderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte)              {}
func (r *Renderer) RawHtmlTag(out *bytes.Buffer, tag []byte)                                       {}
func (r *Renderer) TripleEmphasis(out *bytes.Buffer, text []byte)                                  {}
func (r *Renderer) StrikeThrough(out *bytes.Buffer, text []byte)                                   {}
func (r *Renderer) FootnoteRef(out *bytes.Buffer, ref []byte, id int)                              {}
func (r *Renderer) Index(out *bytes.Buffer, primary, secondary []byte, prim bool)                  {}
func (r *Renderer) Citation(out *bytes.Buffer, link, title []byte)                                 {}
func (r *Renderer) Abbreviation(out *bytes.Buffer, abbr, title []byte)                             {}
func (r *Renderer) Example(out *bytes.Buffer, index int)                                           {}
func (r *Renderer) Math(out *bytes.Buffer, text []byte, display bool)                              {}

// Low-level callbacks
func (r *Renderer) Entity(out *bytes.Buffer, entity []byte)   {}
func (r *Renderer) NormalText(out *bytes.Buffer, text []byte) {}

// Header and footer
func (r *Renderer) DocumentHeader(out *bytes.Buffer, start bool) {}
func (r *Renderer) DocumentFooter(out *bytes.Buffer, start bool) {}

// Frontmatter, mainmatter or backmatter
func (r *Renderer) DocumentMatter(out *bytes.Buffer, matter int)                 {}
func (r *Renderer) References(out *bytes.Buffer, citations map[string]*citation) {}

// Helper functions
func (r *Renderer) Flags() int {}

// Attr returns the inline attribute.
func (r *Renderer) Attr() *inlineAttr {}

// SetAttr set the inline attribute.
func (r *Renderer) SetAttr(*inlineAttr) {}

// AttrString return the string representation of this inline attribute.
func (r *Renderer) AttrString(*inlineAttr) string {}
