package models

import (
	"fmt"
)

var Coefficients = map[string]string{
	"1.000": "2元",
	"0.500": "1元",
	"0.100": "2角",
	"0.050": "1角",
	"0.010": "2分",
	"0.001": "2厘",
}

/*
 * 返回keys
 */
func GetValidCoefficientValues() []string {
	aArr := []string{}
	for sKey, _ := range Coefficients {
		aArr = append(aArr, sKey)
	}
	return aArr
}

/*
 *根据条件获取所有系列
 */
func GetCoefficientText(fKey float64) string {
	coeKey := fmt.Sprintf("%.3f", fKey)
	return Coefficients[coeKey]
}
