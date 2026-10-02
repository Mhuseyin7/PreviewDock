package github
import("crypto/hmac";"crypto/sha256";"fmt";"testing")
func TestVerifySignature(t *testing.T){b:=[]byte("payload");m:=hmac.New(sha256.New,[]byte("secret"));m.Write(b);if !VerifySignature("secret",fmt.Sprintf("sha256=%x",m.Sum(nil)),b){t.Fatal("valid signature rejected")};if VerifySignature("bad","sha256=00",b){t.Fatal("invalid accepted")}}
