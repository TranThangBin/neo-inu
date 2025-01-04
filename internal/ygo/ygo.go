package ygo

type Response struct {
	Data []*YgoCard     `json:"data"`
	Meta map[string]any `json:"meta"`
}

type YgoCard struct {
	CardImages []*YgoCardImage `json:"card_images"`
}

type YgoCardImage struct {
	Id              int    `json:"id"`
	ImageUrl        string `json:"image_url"`
	ImageUrlSmall   string `json:"image_url_small"`
	ImageUrlCropped string `json:"image_url_cropped"`
}
