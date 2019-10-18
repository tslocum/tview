package tview

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gdamore/tcell"
)

type textViewTestCase struct {
	colors     bool
	regions    bool
	scrollable bool
	wrap       bool
	wordwrap   bool
}

func (c *textViewTestCase) String() string {
	return fmt.Sprintf("%sColor%sRegion%sScroll%sWrap%sWordWrap", tvl(c.colors), tvl(c.regions), tvl(c.scrollable), tvl(c.wrap), tvl(c.wordwrap))
}

const randomDataSize = 512

var (
	randomData        = generateRandomData()
	textViewTestCases = generateTestCases()
)

func TestTextViewWrite(t *testing.T) {
	t.Parallel()

	for _, c := range textViewTestCases {
		c := c // Capture

		t.Run(c.String(), func(t *testing.T) {
			t.Parallel()

			var (
				tv, _, err = testTextView(c)
				n          int
			)
			if err != nil {
				t.Error(err)
			}

			n, err = tv.Write(randomData)
			if err != nil {
				t.Errorf("failed to write (successfully wrote %d) bytes: %s", n, err)
			} else if n != randomDataSize {
				t.Errorf("failed to write: expected to write %d bytes, wrote %d", randomDataSize, n)
			}

			tv.Clear()
		})
	}
}

func BenchmarkTextViewWrite(b *testing.B) {
	for _, c := range textViewTestCases {
		c := c // Capture

		b.Run(c.String(), func(b *testing.B) {
			var (
				tv, _, err = testTextView(c)
				n          int
			)
			if err != nil {
				b.Error(err)
			}

			tv.Clear()

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				n, err = tv.Write(randomData)
				if err != nil {
					b.Errorf("failed to write (successfully wrote %d) bytes: %s", n, err)
				} else if n != randomDataSize {
					b.Errorf("failed to write: expected to write %d bytes, wrote %d", randomDataSize, n)
				}

				tv.Clear()
			}
		})
	}
}

func TestTextViewDraw(t *testing.T) {
	t.Parallel()

	for _, c := range textViewTestCases {
		c := c // Capture

		t.Run(c.String(), func(t *testing.T) {
			t.Parallel()

			var (
				tv, sc, err = testTextView(c)
				n           int
			)
			if err != nil {
				t.Error(err)
			}

			n, err = tv.Write(randomData)
			if err != nil {
				t.Errorf("failed to write (successfully wrote %d) bytes: %s", n, err)
			} else if n != randomDataSize {
				t.Errorf("failed to write: expected to write %d bytes, wrote %d", randomDataSize, n)
			}

			tv.Draw(sc)
		})
	}
}

func BenchmarkTextViewDraw(b *testing.B) {
	for _, c := range textViewTestCases {
		c := c // Capture

		b.Run(c.String(), func(b *testing.B) {
			var (
				tv, sc, err = testTextView(c)
				n           int
			)
			if err != nil {
				b.Error(err)
			}

			n, err = tv.Write(randomData)
			if err != nil {
				b.Errorf("failed to write (successfully wrote %d) bytes: %s", n, err)
			} else if n != randomDataSize {
				b.Errorf("failed to write: expected to write %d bytes, wrote %d", randomDataSize, n)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				tv.Draw(sc)
			}
		})
	}
}

func generateTestCases() []*textViewTestCase {
	var cases []*textViewTestCase

	colors := false
	for i := 0; i < 2; i++ {
		if i == 1 {
			colors = true
		}

		regions := false
		for i := 0; i < 2; i++ {
			if i == 1 {
				regions = true
			}

			scrollable := false
			for i := 0; i < 2; i++ {
				if i == 1 {
					scrollable = true
				}

				wrap := false
				for i := 0; i < 2; i++ {
					if i == 1 {
						wrap = true
					}

					wordwrap := false
					for i := 0; i < 2; i++ {
						if i == 1 {
							wordwrap = true
						}

						cases = append(cases, &textViewTestCase{colors: colors, regions: regions, scrollable: scrollable, wrap: wrap, wordwrap: wordwrap})
					}
				}
			}
		}
	}

	return cases
}

func generateRandomData() []byte {
	var (
		b bytes.Buffer
		r = 33
	)

	for i := 0; i < randomDataSize; i++ {
		if i%80 == 0 && i <= 160 {
			b.WriteRune('\n')
		} else if i%7 == 0 {
			b.WriteRune(' ')
		} else {
			b.WriteRune(rune(r))
		}

		r++
		if r == 127 {
			r = 33
		}
	}

	return b.Bytes()
}

func tvc(tv *TextView, c *textViewTestCase) *TextView {
	return tv.SetDynamicColors(c.colors).SetRegions(c.regions).SetScrollable(c.scrollable).SetWrap(c.wrap).SetWordWrap(c.wordwrap)
}

func tvl(v bool) string {
	if v {
		return "Y"
	}

	return "N"
}

func testTextView(c *textViewTestCase) (*TextView, tcell.Screen, error) {
	tv := NewTextView()

	sc := tcell.NewSimulationScreen("UTF-8")
	sc.SetSize(80, 24)

	err := sc.Init()
	if err != nil {
		return nil, nil, err
	}

	app := NewApplication()
	app.screen = sc
	app.SetRoot(tv, true)

	return tvc(tv, c), sc, nil
}
