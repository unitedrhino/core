package domain

type (
	DateFilterSlotReq struct {
		Code string `json:"code"`
	}
	DateFilterSlotResp struct {
		Where string `json:"where"`
	}
)
