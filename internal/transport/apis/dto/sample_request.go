package dto

import entitysample "go-boilerplate-clean/internal/entity/sample"

type SaveSampleRequest struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func (r *SaveSampleRequest) ToEntity() entitysample.Sample {
	return entitysample.Sample{
		Code:   r.Code,
		Name:   r.Name,
		Email:  r.Email,
		Status: r.Status,
	}
}
