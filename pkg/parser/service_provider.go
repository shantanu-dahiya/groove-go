package parser

type ServiceProvider struct {
	Id    int    `json:"id"`
	Host  string `json:"host"`
	Port  int    `json:"port"`
}
