package deployments
import("crypto/rand";"encoding/hex";"fmt";"time")
type State string
const(Queued State="QUEUED";Validating State="VALIDATING";Building State="BUILDING";Provisioning State="PROVISIONING";Starting State="STARTING";Verifying State="VERIFYING";Ready State="READY";Updating State="UPDATING";Failed State="FAILED";Stopping State="STOPPING";Deleting State="DELETING";Deleted State="DELETED";Expired State="EXPIRED")
type Deployment struct{ID string `json:"id"`;Repository string `json:"repository"`;PullRequest int `json:"pullRequest"`;CommitSHA string `json:"commitSha"`;State State `json:"state"`;CreatedAt time.Time `json:"createdAt"`}
func New(repo string,pr int,sha string)Deployment{b:=make([]byte,16);_,_=rand.Read(b);return Deployment{ID:hex.EncodeToString(b),Repository:repo,PullRequest:pr,CommitSHA:sha,State:Queued,CreatedAt:time.Now().UTC()}}
func CanTransition(from,to State)bool{allowed:=map[State]map[State]bool{Queued:{Validating:true,Failed:true},Validating:{Building:true,Failed:true},Building:{Provisioning:true,Failed:true},Provisioning:{Starting:true,Failed:true},Starting:{Verifying:true,Failed:true},Verifying:{Ready:true,Failed:true},Ready:{Updating:true,Stopping:true,Deleting:true,Expired:true},Updating:{Building:true,Failed:true},Failed:{Deleting:true},Stopping:{Deleting:true},Deleting:{Deleted:true}};return allowed[from][to]}
func (d *Deployment)Transition(to State)error{if !CanTransition(d.State,to){return fmt.Errorf("invalid transition %s -> %s",d.State,to)};d.State=to;return nil}
