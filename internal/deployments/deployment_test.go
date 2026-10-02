package deployments
import "testing"
func TestTransitions(t *testing.T){d:=New("o/r",1,"abc");for _,s:=range []State{Validating,Building,Provisioning,Starting,Verifying,Ready}{if err:=d.Transition(s);err!=nil{t.Fatal(err)}};if d.Transition(Building)==nil{t.Fatal("expected invalid transition")}}
