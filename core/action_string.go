// Code generated by "stringer -type Action -trimprefix Action"; DO NOT EDIT.

package core

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ActionSys-0]
	_ = x[ActionSysWithdraw-1]
	_ = x[ActionSysProperty-2]
	_ = x[ActionProposal-10]
	_ = x[ActionProposalMake-11]
	_ = x[ActionProposalShout-12]
	_ = x[ActionProposalVote-13]
	_ = x[ActionCat-20]
	_ = x[ActionCatCreate-21]
	_ = x[ActionCatSupply-22]
	_ = x[ActionCatEdit-23]
	_ = x[ActionCatFold-24]
	_ = x[ActionCatMove-25]
	_ = x[ActionCatGain-26]
	_ = x[ActionCatFill-27]
	_ = x[ActionVat-30]
	_ = x[ActionVatOpen-31]
	_ = x[ActionVatDeposit-32]
	_ = x[ActionVatWithdraw-33]
	_ = x[ActionVatPayback-34]
	_ = x[ActionVatGenerate-35]
	_ = x[ActionFlip-40]
	_ = x[ActionFlipKick-41]
	_ = x[ActionFlipBid-42]
	_ = x[ActionFlipDeal-43]
	_ = x[ActionOracle-50]
	_ = x[ActionOracleCreate-51]
	_ = x[ActionOracleEdit-52]
	_ = x[ActionOraclePoke-53]
	_ = x[ActionOracleRely-54]
	_ = x[ActionOracleDeny-55]
}

const (
	_Action_name_0 = "SysSysWithdrawSysProperty"
	_Action_name_1 = "ProposalProposalMakeProposalShoutProposalVote"
	_Action_name_2 = "CatCatCreateCatSupplyCatEditCatFoldCatMoveCatGainCatFill"
	_Action_name_3 = "VatVatOpenVatDepositVatWithdrawVatPaybackVatGenerate"
	_Action_name_4 = "FlipFlipKickFlipBidFlipDeal"
	_Action_name_5 = "OracleOracleCreateOracleEditOraclePokeOracleRelyOracleDeny"
)

var (
	_Action_index_0 = [...]uint8{0, 3, 14, 25}
	_Action_index_1 = [...]uint8{0, 8, 20, 33, 45}
	_Action_index_2 = [...]uint8{0, 3, 12, 21, 28, 35, 42, 49, 56}
	_Action_index_3 = [...]uint8{0, 3, 10, 20, 31, 41, 52}
	_Action_index_4 = [...]uint8{0, 4, 12, 19, 27}
	_Action_index_5 = [...]uint8{0, 6, 18, 28, 38, 48, 58}
)

func (i Action) String() string {
	switch {
	case 0 <= i && i <= 2:
		return _Action_name_0[_Action_index_0[i]:_Action_index_0[i+1]]
	case 10 <= i && i <= 13:
		i -= 10
		return _Action_name_1[_Action_index_1[i]:_Action_index_1[i+1]]
	case 20 <= i && i <= 27:
		i -= 20
		return _Action_name_2[_Action_index_2[i]:_Action_index_2[i+1]]
	case 30 <= i && i <= 35:
		i -= 30
		return _Action_name_3[_Action_index_3[i]:_Action_index_3[i+1]]
	case 40 <= i && i <= 43:
		i -= 40
		return _Action_name_4[_Action_index_4[i]:_Action_index_4[i+1]]
	case 50 <= i && i <= 55:
		i -= 50
		return _Action_name_5[_Action_index_5[i]:_Action_index_5[i+1]]
	default:
		return "Action(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
