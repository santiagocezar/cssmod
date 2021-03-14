package cssmod

import (
	"fmt"
	"io"
	"math/rand"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

func generateClassName(filename, classname string) string {
	id := fmt.Sprintf("%06x", rand.Int()%16777216)
	return filename + "_" + classname + "_" + id
}

/*
Transform takes input from a Reader (r) with a filename.
Returns the transfromed CSS with namespaced classes, and
a map that indexes the new classes with the original
class names.
*/
func Transform(r io.Reader, filename string) (output []byte, classes map[string]string) {
	l := css.NewLexer(parse.NewInput(r))
	parsingClass := false

	classes = make(map[string]string)

	for {
		tt, text := l.Next()
		switch tt {
		case css.ErrorToken:
			if l.Err() != io.EOF {
				fmt.Printf("An error ocurred: %v", l.Err())
			}
			return
		case css.DelimToken:
			if len(text) == 1 && text[0] == 46 { // a period
				parsingClass = true
			}
			output = append(output, text...)
		case css.IdentToken:
			if parsingClass {
				parsingClass = false
				orig := string(text)
				gen := generateClassName(filename, orig)
				classes[orig] = gen
				output = append(output, gen...)
			} else {
				output = append(output, text...)
			}
		default:
			output = append(output, text...)
		}
	}

}
