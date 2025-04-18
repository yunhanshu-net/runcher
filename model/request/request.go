package request

import (
	"encoding/json"
	"fmt"
	"github.com/yunhanshu-net/runcher/model"
)

type Request struct {
	UUID       string              `json:"uuid"`
	TraceID    string              `json:"trace_id"`
	Route      string              `json:"route"`
	Method     string              `json:"method"`
	Headers    map[string]string   `json:"headers"`
	Body       interface{}         `json:"body"`        //请求json
	BodyString string              `json:"body_string"` //请求json
	FileMap    map[string][]string `json:"file_map"`
}

type Ping struct {
}
type Runner struct {
	Command         string `json:"command"`
	WorkPath        string `json:"work_path"`
	Name            string `json:"name"`
	User            string `json:"user"`
	Version         string `json:"version"`
	RequestJsonPath string `json:"request_json_path"`
}
type RunnerRequest struct {
	UUID            string                 `json:"uuid"`
	Timeout         int                    `json:"timeout"`
	Runner          *model.Runner          `json:"runner"`
	TransportConfig *TransportConfig       `json:"transport_config"`
	Metadata        map[string]interface{} `json:"metadata"`
	//Request         *Request               `json:"request"`
	Request *Request    `json:"request"`
	Body    interface{} `json:"body"`
}

type TransportConfig struct {
	IdleTime int                    `json:"idle_time"`
	Type     string                 `json:"type"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (r *RunnerRequest) GetSubject() string {
	return fmt.Sprintf("runner.%s.%s.%s.run", r.Runner.User, r.Runner.Name, r.Runner.Version)
}
func (r *RunnerRequest) Bytes() []byte {
	jsonBytes, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return jsonBytes
}
