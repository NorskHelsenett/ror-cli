package models

type ApiFilter struct {
	Skip     int `json:"skip"`
	Limit    int `json:"limit"`
	SortList []Sort
}

type Sort struct {
	By  string `json:"by"`
	Asc bool   `json:"asc"`
}

type CliConfig struct {
	Log_level        string            `json:"log_level" validate:"required"`
	Log_output       string            `json:"log_output" validate:"required"`
	Log_output_error string            `json:"log_output_error"`
	Apiconfig        ApiConfig         `json:"apiconfig" validate:"required"`
	RorAuth          RorAuthConfig     `json:"rorauth"`
	ProviderConfigs  ProviderConfigs   `json:"providerconfigs"`
	LastSession      LastSessionConfig `json:"lastsession"`
	Vim              bool              `json:"vim"`
}

type RorAuthConfig struct {
	Expiry int    `json:"expiry"`
	ApiKey string `json:"apikey"`
}

type ProviderConfigs struct {
	TanzuConfig TanzuConfig `json:"tanzuconfig"`
}

type TanzuConfig struct {
	PamToken string `json:"pamtoken"`
	Exp      int    `json:"exp"`
	PamId    string `json:"pamid"`
}
type ApiConfig struct {
	Ror     string `json:"ror" validate:"required"`
	Dex     string `json:"dex" validate:"required"`
	Vsphere string `json:"vsphere" validate:"required"`
}

type LastSessionConfig struct {
	Cluster     string `json:"cluster"`
	Environment string `json:"environment"`
	Workspace   string `json:"workspace"`
}
