package model

type Node struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}
