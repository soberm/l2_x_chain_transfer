package operator

import (
	"math/big"
)

func Div(_ *big.Int, inputs []*big.Int, outputs []*big.Int) error {
	result := big.NewInt(0)
	result.Div(inputs[0], inputs[1])
	outputs[0] = result
	return nil
}
