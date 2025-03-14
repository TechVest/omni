package bridging

import (
	"fmt"
	"math/big"
	"time"

	"github.com/omni-network/omni/contracts/bindings"
	"github.com/omni-network/omni/e2e/app/eoa"
	"github.com/omni-network/omni/lib/contracts/solvernet"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/netconf"
	"github.com/omni-network/omni/monitor/flowgen/types"
	"github.com/omni-network/omni/solver/app"

	"github.com/ethereum/go-ethereum/common"
)

// NewJob instantiates the job that bridges native ETH.
func NewJob(
	network netconf.ID,
	srcChain,
	dstChain uint64,
	role eoa.Role,
	token common.Address,
	amount *big.Int,
) (types.Job, error) {
	owner := eoa.MustAddress(network, role)
	data, err := orderData(owner, srcChain, dstChain, token, amount)
	if err != nil {
		return types.Job{}, errors.Wrap(err, "new job")
	}

	namer := netconf.ChainNamer(network)

	return types.Job{
		Name:    fmt.Sprintf("Bridging (%v->%v)", namer(srcChain), namer(dstChain)),
		Cadence: 30 * time.Minute,
		Network: network,

		SrcChain: srcChain,
		DstChain: dstChain,

		Owner: owner,

		OrderData: data,
	}, nil
}

// OrderData returns the order data required to do the job.
func orderData(
	owner common.Address,
	srcChain, dstChain uint64,
	tokenAddr common.Address,
	amount *big.Int,
) (bindings.SolverNetOrderData, error) {
	token, ok := app.AllTokens().Find(srcChain, tokenAddr)
	if !ok {
		return bindings.SolverNetOrderData{}, errors.New("token not found")
	}
	// Tokens that will be deposited to the user on the destination chain.
	expense := app.TokenAmt{Token: token, Amount: amount}

	depositWithFee, err := app.QuoteDeposit(expense.Token, expense)
	if err != nil {
		return bindings.SolverNetOrderData{}, errors.Wrap(err, "quote expense")
	}

	orderData := bindings.SolverNetOrderData{
		Owner:       owner,
		DestChainId: dstChain,
		Deposit: solvernet.Deposit{
			Token:  depositWithFee.Token.Address,
			Amount: depositWithFee.Amount,
		},
		Expenses: []solvernet.Expense{},
		Calls: []bindings.SolverNetCall{
			{
				Target: owner,
				Value:  expense.Amount,
			},
		},
	}

	return orderData, nil
}
