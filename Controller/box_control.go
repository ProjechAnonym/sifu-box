package controller

type box_url struct {
	url   string
	proxy bool
	label string
}
type box_ruleset struct {
	label string
	value map[string]interface{}
}
type Box_config struct {
	Url     []box_url     `json:"url"`
	Ruleset []box_ruleset `json:"rule_set"`
}

func Add_items(box_config Box_config) {

	
}