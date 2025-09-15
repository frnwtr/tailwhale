package core

import "strings"

const (
    LabelEnable = "tailwhale.enable"
    LabelHost   = "tailwhale.host"
    LabelMode   = "tailwhale.mode" // values: A|B|C
)

// ParseMode maps string labels to ExposureMode.
func ParseMode(s string) ExposureMode {
    switch strings.ToUpper(strings.TrimSpace(s)) {
    case "A":
        return ModeA
    case "B":
        return ModeB
    case "C":
        return ModeC
    default:
        return ModeA
    }
}

