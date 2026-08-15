package generators

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"math/big"
	"strconv"
	"strings"
)

func RunGUID() {
	run("guid", "GUID Generator", func(cmdry.Request) (cmdry.View, error) {
		b := make([]byte, 16)
		if _, e := rand.Read(b); e != nil {
			return cmdry.View{}, e
		}
		b[6] = b[6]&0x0f | 0x40
		b[8] = b[8]&0x3f | 0x80
		s := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
		return result("New GUID", s), nil
	}, nil)
}
func RunHash() {
	form := func() cmdry.View {
		return cmdry.View{Title: "Hash Generator", Components: []cmdry.Component{{Type: "form", Title: "Hash text", Action: "generate", Submit: "Generate hash", Fields: []cmdry.Field{{Name: "input", Label: "Text", Type: "textarea", Required: true}, {Name: "algorithm", Label: "Algorithm", Type: "select", Value: "sha256", Options: []cmdry.Option{{Value: "md5", Label: "MD5"}, {Value: "sha1", Label: "SHA-1"}, {Value: "sha256", Label: "SHA-256"}, {Value: "sha512", Label: "SHA-512"}}}}}}}
	}
	run("hash", "Hash Generator", func(r cmdry.Request) (cmdry.View, error) { return form(), nil }, func(r cmdry.Request) (cmdry.View, error) {
		v := fmt.Sprint(r.Params["input"])
		a := fmt.Sprint(r.Params["algorithm"])
		var sum string
		switch a {
		case "md5":
			x := md5.Sum([]byte(v))
			sum = fmt.Sprintf("%x", x)
		case "sha1":
			x := sha1.Sum([]byte(v))
			sum = fmt.Sprintf("%x", x)
		case "sha256":
			x := sha256.Sum256([]byte(v))
			sum = fmt.Sprintf("%x", x)
		case "sha512":
			x := sha512.Sum512([]byte(v))
			sum = fmt.Sprintf("%x", x)
		default:
			return cmdry.View{}, fmt.Errorf("unsupported algorithm")
		}
		return result(strings.ToUpper(a)+" hash", sum), nil
	})
}
func RunPassword() {
	form := func() cmdry.View {
		return cmdry.View{Title: "Password Generator", Components: []cmdry.Component{{Type: "form", Title: "Password options", Action: "generate", Submit: "Generate password", Fields: []cmdry.Field{{Name: "min", Label: "Minimum length", Type: "number", Value: "16", Min: "4", Max: "256", Required: true}, {Name: "max", Label: "Maximum length", Type: "number", Value: "24", Min: "4", Max: "256", Required: true}, {Name: "lower", Label: "Lowercase", Type: "checkbox", Value: "true"}, {Name: "upper", Label: "Uppercase", Type: "checkbox", Value: "true"}, {Name: "numbers", Label: "Numbers", Type: "checkbox", Value: "true"}, {Name: "symbols", Label: "Symbols", Type: "checkbox"}}}}}
	}
	run("password", "Password Generator", func(r cmdry.Request) (cmdry.View, error) { return form(), nil }, func(r cmdry.Request) (cmdry.View, error) {
		min, e := strconv.Atoi(fmt.Sprint(r.Params["min"]))
		if e != nil {
			return cmdry.View{}, fmt.Errorf("minimum length must be a number")
		}
		max, e := strconv.Atoi(fmt.Sprint(r.Params["max"]))
		if e != nil || min < 4 || max > 256 || min > max {
			return cmdry.View{}, fmt.Errorf("choose lengths from 4 to 256 with minimum no greater than maximum")
		}
		sets := []string{}
		if r.Params["lower"] == "true" {
			sets = append(sets, "abcdefghijkmnopqrstuvwxyz")
		}
		if r.Params["upper"] == "true" {
			sets = append(sets, "ABCDEFGHJKLMNPQRSTUVWXYZ")
		}
		if r.Params["numbers"] == "true" {
			sets = append(sets, "23456789")
		}
		if r.Params["symbols"] == "true" {
			sets = append(sets, "!@#$%^&*()-_=+")
		}
		if len(sets) == 0 {
			return cmdry.View{}, fmt.Errorf("select at least one character set")
		}
		chars := strings.Join(sets, "")
		length := min + secureInt(max-min+1)
		out := make([]byte, length)
		for i := range out {
			out[i] = chars[secureInt(len(chars))]
		}
		return result("Generated password", string(out)), nil
	})
}
func secureInt(n int) int {
	v, e := rand.Int(rand.Reader, bigInt(n))
	if e != nil {
		panic(e)
	}
	return int(v.Int64())
}
func bigInt(n int) *big.Int { return big.NewInt(int64(n)) }
func result(title, value string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "code", Title: title, Text: value}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Generate another"}}}}}
}
func run(id, name string, overview cmdry.Handler, generate cmdry.Handler) {
	actions := map[string]cmdry.Handler{"overview": overview}
	manifestActions := []cmdry.Action{{ID: "overview", Name: "New value", Method: "read"}}
	if generate == nil {
		actions["overview"] = overview
	} else {
		actions["generate"] = generate
		manifestActions = append(manifestActions, cmdry.Action{ID: "generate", Name: "Generate", Method: "write"})
	}
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: "Generate a local value.", Category: "developer", Pages: []cmdry.Page{{ID: "overview", Name: "Generate", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: manifestActions}, Actions: actions})
}
