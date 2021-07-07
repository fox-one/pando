package color

import "math/rand"

var hexColors = []string{
	"#FF93C9",
	"#FF4D00",
	"#0BAAFF",
	"#008080",
	"#5AC18E",
	"#0066CC",
	"#FD5392",
}

func Random() string {
	return hexColors[rand.Intn(len(hexColors))]
}
