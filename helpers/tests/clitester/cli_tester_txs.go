package clitester

import (
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (ct *CLITester) TxCurrenciesIssue(recipientAddr, fromAddr, symbol string, amount sdk.Int, decimals int8, issueID string) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"currencies",
		fromAddr,
		"ms-issue-currency",
		symbol,
		amount.String(),
		strconv.Itoa(int(decimals)),
		recipientAddr,
		issueID)

	return r
}

func (ct *CLITester) TxCurrenciesDestroy(recipientAddr, fromAddr, symbol string, amount sdk.Int) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"currencies",
		fromAddr,
		"destroy-currency",
		ct.ChainID,
		symbol,
		amount.String(),
		recipientAddr)

	return r
}

func (ct *CLITester) TxOracleAddAsset(nomineeAddress, assetCode string, oracleAddresses ...string) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"oracle",
		"",
		"add-asset",
		nomineeAddress,
		assetCode,
		strings.Join(oracleAddresses, ","))

	return r
}

func (ct *CLITester) TxPoaAddValidator(fromAddr, address, ethAddress, issueId string) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"poa",
		fromAddr,
		"ms-add-validator",
		address,
		ethAddress,
		issueId)

	return r
}

func (ct *CLITester) TxPoaRemoveValidator(fromAddr, address, issueId string) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"poa",
		fromAddr,
		"ms-remove-validator",
		address,
		issueId)

	return r
}

func (ct *CLITester) TxPoaReplaceValidator(fromAddr, targetAddress, address, ethAddress, issueId string) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"poa",
		fromAddr,
		"ms-replace-validator",
		targetAddress,
		address,
		ethAddress,
		issueId)

	return r
}

func (ct *CLITester) TxOracleSetAsset(nomineeAddress, assetCode string, oracleAddresses ...string) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"oracle",
		"",
		"set-asset",
		nomineeAddress,
		assetCode,
		strings.Join(oracleAddresses, ","))

	return r
}

func (ct *CLITester) TxOracleAddOracle(nomineeAddress, assetCode string, oracleAddress string) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"oracle",
		"",
		"add-oracle",
		nomineeAddress,
		assetCode,
		oracleAddress)

	return r
}

func (ct *CLITester) TxOracleSetOracles(nomineeAddress, assetCode string, oracleAddresses ...string) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"oracle",
		"",
		"set-oracles",
		nomineeAddress,
		assetCode,
		strings.Join(oracleAddresses, ","))

	return r
}

func (ct *CLITester) TxOraclePostPrice(nomineeAddress, assetCode string, price sdk.Int, receivedAt time.Time) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"oracle",
		"",
		"postprice",
		nomineeAddress,
		assetCode,
		price.String(),
		strconv.FormatInt(receivedAt.Unix(), 10))

	return r
}

func (ct *CLITester) TxMultiSigConfirmCall(fromAddress string, callID uint64) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"multisig",
		fromAddress,
		"confirm-call",
		strconv.FormatUint(callID, 10))

	return r
}

func (ct *CLITester) TxMultiSigRevokeConfirm(fromAddress string, callID uint64) *TxRequest {
	r := ct.newTxRequest()
	r.SetCmd(
		"multisig",
		fromAddress,
		"revoke-confirm",
		strconv.FormatUint(callID, 10))

	return r
}

func (ct *CLITester) TxVmExecuteScript(fromAddress, filePath string, args ...string) *TxRequest {
	cmdArgs := make([]string, 0, 2+len(args))
	cmdArgs = append(cmdArgs, "execute-script")
	cmdArgs = append(cmdArgs, filePath)
	cmdArgs = append(cmdArgs, args...)

	r := ct.newTxRequest()
	r.SetCmd(
		"vm",
		fromAddress,
		cmdArgs...)
	r.cmd.AddArg("compiler", ct.vmCompilerAddress)

	return r
}

func (ct *CLITester) TxVmDeployModule(fromAddress, filePath string) *TxRequest {
	cmdArgs := []string{
		"deploy-module",
		filePath,
	}

	r := ct.newTxRequest()
	r.SetCmd(
		"vm",
		fromAddress,
		cmdArgs...)

	return r
}
