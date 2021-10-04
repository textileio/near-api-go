package types

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"

	uint128 "lukechampine.com/uint128"
)

var (
	tenPower24 = uint128.From64(uint64(math.Pow10(12))).Mul64(uint64(math.Pow10(12)))
	zeroNEAR   = Balance(uint128.From64(0))
)

// Balance holds amount of yoctoNEAR
type Balance uint128.Uint128

func (bal *Balance) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	val := big.Int{}
	if _, ok := val.SetString(s, 10); !ok {
		return fmt.Errorf("unable to parse '%s'", s)
	}

	*bal = Balance(uint128.FromBig(&val))

	return nil
}

func (bal Balance) MarshalJSON() ([]byte, error) {
	return json.Marshal(bal.String())
}

func (bal Balance) String() string {
	return uint128.Uint128(bal).String()
}

// Convenience funcs
func (bal Balance) Div64(div uint64) Balance {
	return Balance(uint128.Uint128(bal).Div64(div))
}

// TODO
func NEARToYocto(near uint64) Balance {
	if near == 0 {
		return zeroNEAR
	}

	return Balance(uint128.From64(near).Mul(tenPower24))
}

// TODO
func YoctoToNEAR(yocto Balance) uint64 {
	div := uint128.Uint128(yocto).Div(tenPower24)
	if h := div.Hi; h != 0 {
		panic(fmt.Errorf("yocto div failed, remaining: %d", h))
	}

	return div.Lo
}

func scaleToYocto(f *big.Float) (r *big.Int) {
	// Convert reference 1 NEAR to big.Float
	base := new(big.Float).SetPrec(128).SetInt(uint128.Uint128(NEARToYocto(1)).Big())

	// Multiply base using the supplied float
	// XXX: small precision issues here will haunt me forever
	bigf2 := new(big.Float).SetPrec(128).SetMode(big.ToZero).Mul(base, f)

	// Convert it to big.Int
	r, _ = bigf2.Int(nil)
	return
}

// TODO
func BalanceFromFloat(f float64) (bal Balance) {
	bigf := big.NewFloat(f)
	bal = Balance(uint128.FromBig(scaleToYocto(bigf)))
	return
}

// TODO
func BalanceFromString(s string) (bal Balance, err error) {
	var bigf *big.Float
	bigf, _, err = big.ParseFloat(s, 10, 128, big.ToZero)
	if err != nil {
		return
	}

	bal = Balance(uint128.FromBig(scaleToYocto(bigf)))
	return
}
