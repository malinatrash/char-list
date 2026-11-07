package sheet

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"char-list/internal/embed"
	"github.com/signintech/gopdf"
)

type Doc struct {
	pdf  gopdf.GoPdf
	tpl  []byte
	W, H float64
	font string
}

func New(template string, w, h float64) (*Doc, error) {
	d := &Doc{tpl: embed.Sheet(template), W: w, H: h, font: "serif"}
	d.pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: w, H: h}})
	for name, file := range map[string]string{"serif": "DejaVuSerif.ttf", "bold": "DejaVuSerif-Bold.ttf", "sans": "DejaVuSans.ttf"} {
		if err := d.pdf.AddTTFFontData(name, embed.Font(file)); err != nil {
			return nil, fmt.Errorf("font %s: %w", name, err)
		}
	}
	d.pdf.SetTextColor(20, 20, 30)
	return d, nil
}

func (d *Doc) Page(n int) error {
	d.pdf.AddPage()
	var rs io.ReadSeeker = bytes.NewReader(d.tpl)
	tpl := d.pdf.ImportPageStream(&rs, n, "/MediaBox")
	d.pdf.UseImportedTemplate(tpl, 0, 0, d.W, d.H)
	return nil
}

func (d *Doc) Font(name string) { d.font = name }

func (d *Doc) set(size float64) { d.pdf.SetFont(d.font, "", size) }

func (d *Doc) Width(size float64, s string) float64 {
	d.set(size)
	w, _ := d.pdf.MeasureTextWidth(s)
	return w
}

func (d *Doc) Text(x, y, size float64, s string) {
	d.set(size)
	d.pdf.SetXY(x, y)
	d.pdf.Cell(nil, s)
}

func (d *Doc) Center(cx, y, size float64, s string) {
	d.Text(cx-d.Width(size, s)/2, y, size, s)
}

// Fit draws one line, shrinking the font until it fits the width.
func (d *Doc) Fit(x, y, w, size float64, s string) {
	for size > 4 && d.Width(size, s) > w {
		size -= 0.5
	}
	d.Text(x, y, size, s)
}

func (d *Doc) FitCenter(cx, y, w, size float64, s string) {
	for size > 4 && d.Width(size, s) > w {
		size -= 0.5
	}
	d.Center(cx, y, size, s)
}

func (d *Doc) wrap(s string, w, size float64) []string {
	d.set(size)
	var out []string
	for _, para := range strings.Split(strings.ReplaceAll(s, "\r", ""), "\n") {
		if strings.TrimSpace(para) == "" {
			out = append(out, "")
			continue
		}
		lines, _ := d.pdf.SplitTextWithWordWrap(para, w)
		out = append(out, lines...)
	}
	return out
}

// Box wraps text inside a rectangle, shrinking the font if needed; lineStep is a multiplier of size.
func (d *Doc) Box(x, y, w, h, size, lineStep float64, s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	var lines []string
	for {
		lines = d.wrap(s, w, size)
		if float64(len(lines))*size*lineStep <= h || size <= 4.5 {
			break
		}
		size -= 0.5
	}
	for i, l := range lines {
		yy := y + float64(i)*size*lineStep
		if yy+size > y+h+1 {
			break
		}
		d.Text(x, yy, size, l)
	}
}

// Lines writes text onto pre-printed ruled lines with a fixed step.
func (d *Doc) Lines(x, y, w, step, size float64, max int, s string) {
	s = strings.TrimSpace(s)
	if s == "" {
		return
	}
	lines := d.wrap(s, w, size)
	for len(lines) > max && size > 5 {
		size -= 0.5
		lines = d.wrap(s, w, size)
	}
	for i, l := range lines {
		if i >= max {
			break
		}
		d.Text(x, y+float64(i)*step, size, l)
	}
}

func (d *Doc) Write(w io.Writer) error {
	_, err := d.pdf.WriteTo(w)
	return err
}
