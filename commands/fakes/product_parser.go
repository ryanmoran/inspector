package fakes

import "github.com/ryanmoran/inspector/tiles"

type ProductParser struct {
	ParseCall struct {
		Returns struct {
			Product tiles.Product
			Error   error
		}
	}
}

func (p *ProductParser) Parse() (tiles.Product, error) {
	return p.ParseCall.Returns.Product, p.ParseCall.Returns.Error
}
