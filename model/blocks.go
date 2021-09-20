package model

type Block struct {
	Id     int     `json:"id"`
	Label  string  `json:"label"`
	Parent int     `json:"parent"`
	Type   string  `json:"type"`
	PosX   float64 `json:"x"`
	PosY   float64 `json:"y"`
}

type Connection struct {
	Id          int `json:"id"`
	SourceId    int `json:"source_id"`
	SourceRoute int `json:"source_route"`
	TargetId    int `json:"target_id"`
	TargetRoute int `json:"target_route"`
}

// type Link struct {
//     Id       int `json:"id"`
//     SourceId int `json:"id"`
//     BlockId  int `json:"id"`
// }

// type Source struct {
//     Id         int               `json:"id"`
//     Label      string            `json:"label"`
//     Type       string            `json:"type"`
//     Parent     int               `json:"parent"`
//     Parameters map[string]string `json:"params"`
//     PosX       float64           `json:"x"`
//     PosY       float64           `json:"y"`
// }

// type Group struct {
// 	Id       int     `json:"id"`
// 	Group    int     `json:"parent"`
// 	Children []int   `json:"children"`
// 	Label    string  `json:"label"`
// 	PosX     float64 `json:"x"`
// 	PosY     float64 `json:"y"`
// }
