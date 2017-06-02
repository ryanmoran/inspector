package fakes

import "github.com/ryanmoran/inspector/tiles"

type ProductParser struct {
	ParseCall struct {
		CallCount int
		Returns   struct {
			Product tiles.Product
			Error   error
		}
	}
}

func (p *ProductParser) Parse() (tiles.Product, error) {
	p.ParseCall.CallCount++
	return p.ParseCall.Returns.Product, p.ParseCall.Returns.Error
}
