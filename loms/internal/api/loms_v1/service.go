package loms_v1

import (
	desc "route256/loms/pkg/loms_v1"
)

type Implementation struct {
	desc.UnimplementedLOMSV1Server
}

func NewLomsV1() *Implementation {
	return &Implementation{
		desc.UnimplementedLOMSV1Server{},
	}
}
