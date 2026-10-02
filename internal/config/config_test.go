package config
import("testing";"time")
func valid()Config{return Config{Version:1,Build:Build{".","Dockerfile"},Service:Service{3000},Health:Health{"/health",200},Resources:Resources{1,512},Preview:Preview{24*time.Hour,"required"}}}
func TestValidation(t *testing.T){c:=valid();if err:=c.Validate();err!=nil{t.Fatal(err)};c.Build.Context="../x";if c.Validate()==nil{t.Fatal("unsafe context accepted")}}
