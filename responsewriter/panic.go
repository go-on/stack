package responsewriter

import (
	"fmt"
	"net/http"
)

// PanicCodes is a debugging tool to get a stack trace, if some http.StatusCode is set to one of its codes
type PanicCodes struct {
	// Codes that should trigger a panic
	Codes []int
	*Buffer
}

// WriteHeader triggers a panic, if the code is one of the Codes
// otherwise it passes to the underlying Buffer
func (p *PanicCodes) WriteHeader(code int) {

	for i := 0; i < len(p.Codes); i++ {
		if p.Codes[i] == code {
			panic(fmt.Sprintf("You told %T to panic on http code %d: Here we go!", p, code))
		}
	}

	p.Buffer.WriteHeader(code)
}

func NewPanicCodes(w http.ResponseWriter, codes ...int) *PanicCodes {
	return &PanicCodes{Codes: codes, Buffer: NewBuffer(w)}
}
