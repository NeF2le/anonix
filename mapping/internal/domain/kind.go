package domain

type Kind struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	RussianName string `json:"russian_name"`
	AccessLevel int32  `json:"access_level"`
	Mask        string `json:"mask"`
	ShortName   string `json:"short_name"`
}
