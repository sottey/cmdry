package webpriority2

import (
	"strings"
	"testing"
	"time"
)

func TestWebUtilities(t *testing.T) {
	v, e := BasicAuth("ada", "secret")
	if e != nil || v != "Basic YWRhOnNlY3JldA==" {
		t.Fatal(v, e)
	}
	v, e = Entity("<x>", "encode")
	if e != nil || v != "&lt;x&gt;" {
		t.Fatal(v, e)
	}
	jwt, e := InspectJWT("eyJhbGciOiJub25lIn0.eyJzdWIiOiJhZGEifQ.")
	if e != nil || !strings.Contains(jwt.Claims, "ada") {
		t.Fatal(jwt, e)
	}
	otp, e := OTP("JBSWY3DPEHPK3PXP", time.Unix(0, 0))
	if e != nil || len(otp) != 6 {
		t.Fatal(otp, e)
	}
	safe, e := SafeLink("https://nam01.safelinks.protection.outlook.com/?url=https%3A%2F%2Fexample.com%2Fx")
	if e != nil || safe != "https://example.com/x" {
		t.Fatal(safe, e)
	}
}
