package mcastdiscover

type pingMsg struct {
	Port       int    `json:"port"`
	InstanceID string `json:"id"`
}
