package entity


type Wallet struct{
	Id string `json:"id"`
	OwnedBy string `json:"owned_by"`
	Status string `json:"status"`
	EnabledAt *string `json:"enabled_at,omitempty"`	
	DisabledAt *string `json:"disabled_at,omitempty"`	
	Balance float64 `json:"balance"`
}