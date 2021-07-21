package packet

import (
	"encoding/json"
	"strings"
)

// todo data = nil
//{"type":"publish","data":{"code":200,"message":"success","list":[{"id":1},{"id":2}]}}
type Packet struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func (p *Packet) Get(key string, defaultValue interface{}) interface{} {
	if v, ok := p.Data[key]; ok {
		return v
	}
	return defaultValue
}

func (p *Packet) GetString(name string, defaultValue string) string {
	if s, ok := p.Get(name, defaultValue).(string); ok {
		return s
	}
	return defaultValue
}

func (p *Packet) GetBoolean(name string, defaultValue bool) bool {
	if v, ok := p.Get(name, defaultValue).(bool); ok {
		return v
	}
	return defaultValue
}

func (p *Packet) Put(name string, v interface{}) {
	p.Data[name] = v
}

func (p *Packet) Reply(data map[string]interface{}) *Packet {
	p.Data = merge(p.Data, data)
	return p
}

func (p *Packet) Error(message string, code int) *Packet {
	data := make(map[string]interface{})
	data["message"] = message
	data["code"] = code
	p.Data = data
	return p
}

func (p *Packet) Build(ty string, v map[string]interface{}) *Packet {
	p.Type = ty
	p.Data = v
	return p
}

func (p *Packet) Format() []byte {
	bs, _ := json.Marshal(p)
	return bs
}

func merge(original map[string]interface{}, data map[string]interface{}) map[string]interface{} {
	for k, v := range original {
		if strings.HasPrefix(k, "[") || strings.Contains(k, "sign") {
			data[k] = v
		}
	}

	return data
}

func Parse(bs []byte) (*Packet, error) {
	var pack Packet
	err := json.Unmarshal(bs, &pack)
	if err != nil {
		return nil, err
	}

	return &pack, nil
}
