package parliament

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

func randomHexColor() string {
	return hexColors[rand.Intn(len(hexColors))]
}
