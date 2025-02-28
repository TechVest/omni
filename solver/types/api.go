package types

import (
	"encoding/json"
	"math/big"

	"github.com/omni-network/omni/lib/contracts/solvernet"
	"github.com/omni-network/omni/lib/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// JSONErrorResponse is a json response for http errors (e.g 4xx, 5xx), not used for rejections.
type JSONErrorResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// CheckRequest is the expected request body for the /api/v1/check endpoint.
//
// NOTE: Check request / response types mirror SolvertNet.OrderData, built
// specifically for EVM -> EVM orders via SolverNetInbox / Outbox contracts,
// with ERC7683 type hash matching SolverNetInbox.ORDERDATA_TYPEHASH.
//
// To support multiple order types with this api (e.g. EVM -> Solana, Solana -> EVM)
// we'd need a more generic request / response format that discriminates on
// order type hash.
type CheckRequest struct {
	SourceChainID      uint64    `json:"sourceChainId"`
	DestinationChainID uint64    `json:"destChainId"`
	FillDeadline       uint32    `json:"fillDeadline"`
	Calls              []Call    `json:"calls"`
	Expenses           []Expense `json:"expenses"`
	Deposit            AddrAmt   `json:"deposit"`
}

// CheckResponse is the response json for the /check endpoint.
type CheckResponse struct {
	Accepted          bool   `json:"accepted"`
	Rejected          bool   `json:"rejected"`
	RejectReason      string `json:"rejectReason"`
	RejectDescription string `json:"rejectDescription"`
}

// QuoteRequest is the expected request body for the /api/v1/quote endpoint.
// If deposit amount is omitted, the response will include the required deposit amount.
// If expense amount is omitted, the response will include the required expense amount.
type QuoteRequest struct {
	SourceChainID      uint64  `json:"sourceChainId"`
	DestinationChainID uint64  `json:"destChainId"`
	Deposit            AddrAmt `json:"deposit"`
	Expense            AddrAmt `json:"expense"`
}

type addrAmtJSON struct {
	Token  common.Address `json:"token"`
	Amount *hexutil.Big   `json:"amount,omitempty"`
}

// AddrAmt represents a token address and amount pair, with the amount being optional.
// If amount is nil or zero, quote response should inform the amount.
type AddrAmt struct {
	Token  common.Address
	Amount *big.Int
}

func (u AddrAmt) MarshalJSON() ([]byte, error) {
	return marshal(addrAmtJSON{
		Token:  u.Token,
		Amount: (*hexutil.Big)(u.Amount),
	})
}

func (u *AddrAmt) UnmarshalJSON(bz []byte) error {
	v := new(addrAmtJSON)
	if err := unmarshal(bz, v); err != nil {
		return err
	}

	u.Token = v.Token
	u.Amount = intOrZero(v.Amount)

	return nil
}

// QuoteResponse is the response json for the /api/v1/quote endpoint.
type QuoteResponse struct {
	Deposit AddrAmt `json:"deposit"`
	Expense AddrAmt `json:"expense"`
}

// ContractsResponse is the response json for the /api/vi/contracts endpoint.
type ContractsResponse struct {
	Portal    common.Address `json:"portal"`
	Inbox     common.Address `json:"inbox"`
	Outbox    common.Address `json:"outbox"`
	Middleman common.Address `json:"middleman"`
}

// expenseJSON is a json marshal-able solvernt.Expense.
type expenseJSON struct {
	Spender common.Address `json:"spender"`
	Token   common.Address `json:"token"`
	Amount  *hexutil.Big   `json:"amount"`
}

// Expense wraps solvernet.Expense to provide custom json marshaling.
type Expense solvernet.Expense

func (e *Expense) UnmarshalJSON(bz []byte) error {
	v := new(expenseJSON)
	if err := unmarshal(bz, v); err != nil {
		return err
	}

	e.Spender = v.Spender
	e.Token = v.Token
	e.Amount = intOrZero(v.Amount)

	return nil
}

func (e Expense) MarshalJSON() ([]byte, error) {
	return marshal(expenseJSON{
		Spender: e.Spender,
		Token:   e.Token,
		Amount:  (*hexutil.Big)(e.Amount),
	})
}

// callJSON is a json marshal-able solvernet.Call.
type callJSON struct {
	Target common.Address `json:"target"`
	Data   hexutil.Bytes  `json:"data"`
	Value  *hexutil.Big   `json:"value"`
}

// Call wraps solvernet.Call to provide custom json marshaling.
type Call solvernet.Call

func (c *Call) UnmarshalJSON(bz []byte) error {
	v := new(callJSON)
	if err := unmarshal(bz, v); err != nil {
		return err
	}

	c.Target = v.Target
	c.Value = intOrZero(v.Value)
	c.Data = v.Data

	return nil
}

func (c Call) MarshalJSON() ([]byte, error) {
	return marshal(callJSON{
		Target: c.Target,
		Value:  (*hexutil.Big)(c.Value),
		Data:   c.Data,
	})
}

func intOrZero(i *hexutil.Big) *big.Int {
	if i == nil {
		return big.NewInt(0)
	}

	return i.ToInt()
}

func marshal(v any) ([]byte, error) {
	bz, err := json.Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	return bz, nil
}

func unmarshal(bz []byte, v any) error {
	if err := json.Unmarshal(bz, v); err != nil {
		return errors.Wrap(err, "unmarshal")
	}

	return nil
}

func CallsToBindings(calls []Call) solvernet.Calls {
	var resp solvernet.Calls
	for _, c := range calls {
		resp = append(resp, solvernet.Call(c))
	}

	return resp
}

func ExpensesToBindings(expenses []Expense) solvernet.Expenses {
	var resp solvernet.Expenses
	for _, e := range expenses {
		resp = append(resp, solvernet.Expense(e))
	}

	return resp
}
