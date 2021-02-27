package model

type BaseSysSetting struct {
	BlogName     string `json:"blog_name"`
	Username     string `json:"username"`
	Password     string `json:"password,omitempty"`
	GithubName   string `json:"github_name"`
	IsShowGithub bool   `json:"is_show_github"`
	IsShowBook   bool   `json:"is_show_book"`
	SavePath     string `json:"save_path"`
}

type SysSetting struct {
	BaseSysSetting
	LogoName string `json:"logo_name"`
	InitTime int64  `json:"init_time"`
	Salt     string `json:"salt,omitempty"`
}
