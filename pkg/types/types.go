package types

type License struct {
	MainVersion     string   `json:"mainVersion"`
	Feature         []string `json:"feature"`
	ExpireTimestamp int64    `json:"expireTimestamp"`
}
