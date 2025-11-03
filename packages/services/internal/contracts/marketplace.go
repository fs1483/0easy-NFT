// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IMarketplaceOrder is an auto generated low-level Go binding around an user-defined struct.
type IMarketplaceOrder struct {
	Maker        common.Address
	Nft          common.Address
	TokenId      *big.Int
	PaymentToken common.Address
	Price        *big.Int
	Expiry       *big.Int
	Nonce        *big.Int
	Side         uint8
}

// OeasyMarketplaceMetaData contains all meta data concerning the OeasyMarketplace contract.
var OeasyMarketplaceMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelOrder\",\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"consumedNonces\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"eip712Domain\",\"inputs\":[],\"outputs\":[{\"name\":\"fields\",\"type\":\"bytes1\",\"internalType\":\"bytes1\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"verifyingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"salt\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"extensions\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"executeTrade\",\"inputs\":[{\"name\":\"makerOrder\",\"type\":\"tuple\",\"internalType\":\"structIMarketplace.Order\",\"components\":[{\"name\":\"maker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nft\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"paymentToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"side\",\"type\":\"uint8\",\"internalType\":\"enumIMarketplace.OrderSide\"}]},{\"name\":\"taker\",\"type\":\"tuple\",\"internalType\":\"structIMarketplace.Order\",\"components\":[{\"name\":\"maker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nft\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"paymentToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"side\",\"type\":\"uint8\",\"internalType\":\"enumIMarketplace.OrderSide\"}]},{\"name\":\"makerSignature\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeBps\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeRecipient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hashOrder\",\"inputs\":[{\"name\":\"order\",\"type\":\"tuple\",\"internalType\":\"structIMarketplace.Order\",\"components\":[{\"name\":\"maker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nft\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"paymentToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"price\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"side\",\"type\":\"uint8\",\"internalType\":\"enumIMarketplace.OrderSide\"}]}],\"outputs\":[{\"name\":\"digest\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFeeConfiguration\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeBps_\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"EIP712DomainChanged\",\"inputs\":[],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeUpdated\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"feeBps\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OrderCancelled\",\"inputs\":[{\"name\":\"maker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TradeExecuted\",\"inputs\":[{\"name\":\"maker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"taker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"nft\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"paymentToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"side\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIMarketplace.OrderSide\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureLength\",\"inputs\":[{\"name\":\"length\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ECDSAInvalidSignatureS\",\"inputs\":[{\"name\":\"s\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFeeConfiguration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidOrder\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MakerCannotBeTaker\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NonceConsumed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"UnsupportedSide\",\"inputs\":[]}]",
}

// OeasyMarketplaceABI is the input ABI used to generate the binding from.
// Deprecated: Use OeasyMarketplaceMetaData.ABI instead.
var OeasyMarketplaceABI = OeasyMarketplaceMetaData.ABI

// OeasyMarketplace is an auto generated Go binding around an Ethereum contract.
type OeasyMarketplace struct {
	OeasyMarketplaceCaller     // Read-only binding to the contract
	OeasyMarketplaceTransactor // Write-only binding to the contract
	OeasyMarketplaceFilterer   // Log filterer for contract events
}

// OeasyMarketplaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type OeasyMarketplaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OeasyMarketplaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OeasyMarketplaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OeasyMarketplaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OeasyMarketplaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OeasyMarketplaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OeasyMarketplaceSession struct {
	Contract     *OeasyMarketplace // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OeasyMarketplaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OeasyMarketplaceCallerSession struct {
	Contract *OeasyMarketplaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// OeasyMarketplaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OeasyMarketplaceTransactorSession struct {
	Contract     *OeasyMarketplaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// OeasyMarketplaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type OeasyMarketplaceRaw struct {
	Contract *OeasyMarketplace // Generic contract binding to access the raw methods on
}

// OeasyMarketplaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OeasyMarketplaceCallerRaw struct {
	Contract *OeasyMarketplaceCaller // Generic read-only contract binding to access the raw methods on
}

// OeasyMarketplaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OeasyMarketplaceTransactorRaw struct {
	Contract *OeasyMarketplaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOeasyMarketplace creates a new instance of OeasyMarketplace, bound to a specific deployed contract.
func NewOeasyMarketplace(address common.Address, backend bind.ContractBackend) (*OeasyMarketplace, error) {
	contract, err := bindOeasyMarketplace(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplace{OeasyMarketplaceCaller: OeasyMarketplaceCaller{contract: contract}, OeasyMarketplaceTransactor: OeasyMarketplaceTransactor{contract: contract}, OeasyMarketplaceFilterer: OeasyMarketplaceFilterer{contract: contract}}, nil
}

// NewOeasyMarketplaceCaller creates a new read-only instance of OeasyMarketplace, bound to a specific deployed contract.
func NewOeasyMarketplaceCaller(address common.Address, caller bind.ContractCaller) (*OeasyMarketplaceCaller, error) {
	contract, err := bindOeasyMarketplace(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceCaller{contract: contract}, nil
}

// NewOeasyMarketplaceTransactor creates a new write-only instance of OeasyMarketplace, bound to a specific deployed contract.
func NewOeasyMarketplaceTransactor(address common.Address, transactor bind.ContractTransactor) (*OeasyMarketplaceTransactor, error) {
	contract, err := bindOeasyMarketplace(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceTransactor{contract: contract}, nil
}

// NewOeasyMarketplaceFilterer creates a new log filterer instance of OeasyMarketplace, bound to a specific deployed contract.
func NewOeasyMarketplaceFilterer(address common.Address, filterer bind.ContractFilterer) (*OeasyMarketplaceFilterer, error) {
	contract, err := bindOeasyMarketplace(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceFilterer{contract: contract}, nil
}

// bindOeasyMarketplace binds a generic wrapper to an already deployed contract.
func bindOeasyMarketplace(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OeasyMarketplaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OeasyMarketplace *OeasyMarketplaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OeasyMarketplace.Contract.OeasyMarketplaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OeasyMarketplace *OeasyMarketplaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.OeasyMarketplaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OeasyMarketplace *OeasyMarketplaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.OeasyMarketplaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OeasyMarketplace *OeasyMarketplaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OeasyMarketplace.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OeasyMarketplace *OeasyMarketplaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OeasyMarketplace *OeasyMarketplaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_OeasyMarketplace *OeasyMarketplaceCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _OeasyMarketplace.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_OeasyMarketplace *OeasyMarketplaceSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _OeasyMarketplace.Contract.UPGRADEINTERFACEVERSION(&_OeasyMarketplace.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_OeasyMarketplace *OeasyMarketplaceCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _OeasyMarketplace.Contract.UPGRADEINTERFACEVERSION(&_OeasyMarketplace.CallOpts)
}

// ConsumedNonces is a free data retrieval call binding the contract method 0x4593eb63.
//
// Solidity: function consumedNonces(address , uint256 ) view returns(bool)
func (_OeasyMarketplace *OeasyMarketplaceCaller) ConsumedNonces(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var out []interface{}
	err := _OeasyMarketplace.contract.Call(opts, &out, "consumedNonces", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ConsumedNonces is a free data retrieval call binding the contract method 0x4593eb63.
//
// Solidity: function consumedNonces(address , uint256 ) view returns(bool)
func (_OeasyMarketplace *OeasyMarketplaceSession) ConsumedNonces(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _OeasyMarketplace.Contract.ConsumedNonces(&_OeasyMarketplace.CallOpts, arg0, arg1)
}

// ConsumedNonces is a free data retrieval call binding the contract method 0x4593eb63.
//
// Solidity: function consumedNonces(address , uint256 ) view returns(bool)
func (_OeasyMarketplace *OeasyMarketplaceCallerSession) ConsumedNonces(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _OeasyMarketplace.Contract.ConsumedNonces(&_OeasyMarketplace.CallOpts, arg0, arg1)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_OeasyMarketplace *OeasyMarketplaceCaller) Eip712Domain(opts *bind.CallOpts) (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	var out []interface{}
	err := _OeasyMarketplace.contract.Call(opts, &out, "eip712Domain")

	outstruct := new(struct {
		Fields            [1]byte
		Name              string
		Version           string
		ChainId           *big.Int
		VerifyingContract common.Address
		Salt              [32]byte
		Extensions        []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Fields = *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)
	outstruct.Name = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Version = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.VerifyingContract = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Salt = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Extensions = *abi.ConvertType(out[6], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_OeasyMarketplace *OeasyMarketplaceSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _OeasyMarketplace.Contract.Eip712Domain(&_OeasyMarketplace.CallOpts)
}

// Eip712Domain is a free data retrieval call binding the contract method 0x84b0196e.
//
// Solidity: function eip712Domain() view returns(bytes1 fields, string name, string version, uint256 chainId, address verifyingContract, bytes32 salt, uint256[] extensions)
func (_OeasyMarketplace *OeasyMarketplaceCallerSession) Eip712Domain() (struct {
	Fields            [1]byte
	Name              string
	Version           string
	ChainId           *big.Int
	VerifyingContract common.Address
	Salt              [32]byte
	Extensions        []*big.Int
}, error) {
	return _OeasyMarketplace.Contract.Eip712Domain(&_OeasyMarketplace.CallOpts)
}

// FeeBps is a free data retrieval call binding the contract method 0x24a9d853.
//
// Solidity: function feeBps() view returns(uint96)
func (_OeasyMarketplace *OeasyMarketplaceCaller) FeeBps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _OeasyMarketplace.contract.Call(opts, &out, "feeBps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeBps is a free data retrieval call binding the contract method 0x24a9d853.
//
// Solidity: function feeBps() view returns(uint96)
func (_OeasyMarketplace *OeasyMarketplaceSession) FeeBps() (*big.Int, error) {
	return _OeasyMarketplace.Contract.FeeBps(&_OeasyMarketplace.CallOpts)
}

// FeeBps is a free data retrieval call binding the contract method 0x24a9d853.
//
// Solidity: function feeBps() view returns(uint96)
func (_OeasyMarketplace *OeasyMarketplaceCallerSession) FeeBps() (*big.Int, error) {
	return _OeasyMarketplace.Contract.FeeBps(&_OeasyMarketplace.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_OeasyMarketplace *OeasyMarketplaceCaller) FeeRecipient(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OeasyMarketplace.contract.Call(opts, &out, "feeRecipient")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_OeasyMarketplace *OeasyMarketplaceSession) FeeRecipient() (common.Address, error) {
	return _OeasyMarketplace.Contract.FeeRecipient(&_OeasyMarketplace.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_OeasyMarketplace *OeasyMarketplaceCallerSession) FeeRecipient() (common.Address, error) {
	return _OeasyMarketplace.Contract.FeeRecipient(&_OeasyMarketplace.CallOpts)
}

// HashOrder is a free data retrieval call binding the contract method 0x5a893e74.
//
// Solidity: function hashOrder((address,address,uint256,address,uint256,uint256,uint256,uint8) order) view returns(bytes32 digest)
func (_OeasyMarketplace *OeasyMarketplaceCaller) HashOrder(opts *bind.CallOpts, order IMarketplaceOrder) ([32]byte, error) {
	var out []interface{}
	err := _OeasyMarketplace.contract.Call(opts, &out, "hashOrder", order)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashOrder is a free data retrieval call binding the contract method 0x5a893e74.
//
// Solidity: function hashOrder((address,address,uint256,address,uint256,uint256,uint256,uint8) order) view returns(bytes32 digest)
func (_OeasyMarketplace *OeasyMarketplaceSession) HashOrder(order IMarketplaceOrder) ([32]byte, error) {
	return _OeasyMarketplace.Contract.HashOrder(&_OeasyMarketplace.CallOpts, order)
}

// HashOrder is a free data retrieval call binding the contract method 0x5a893e74.
//
// Solidity: function hashOrder((address,address,uint256,address,uint256,uint256,uint256,uint8) order) view returns(bytes32 digest)
func (_OeasyMarketplace *OeasyMarketplaceCallerSession) HashOrder(order IMarketplaceOrder) ([32]byte, error) {
	return _OeasyMarketplace.Contract.HashOrder(&_OeasyMarketplace.CallOpts, order)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OeasyMarketplace *OeasyMarketplaceCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OeasyMarketplace.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OeasyMarketplace *OeasyMarketplaceSession) Owner() (common.Address, error) {
	return _OeasyMarketplace.Contract.Owner(&_OeasyMarketplace.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OeasyMarketplace *OeasyMarketplaceCallerSession) Owner() (common.Address, error) {
	return _OeasyMarketplace.Contract.Owner(&_OeasyMarketplace.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_OeasyMarketplace *OeasyMarketplaceCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _OeasyMarketplace.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_OeasyMarketplace *OeasyMarketplaceSession) Paused() (bool, error) {
	return _OeasyMarketplace.Contract.Paused(&_OeasyMarketplace.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_OeasyMarketplace *OeasyMarketplaceCallerSession) Paused() (bool, error) {
	return _OeasyMarketplace.Contract.Paused(&_OeasyMarketplace.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_OeasyMarketplace *OeasyMarketplaceCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _OeasyMarketplace.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_OeasyMarketplace *OeasyMarketplaceSession) ProxiableUUID() ([32]byte, error) {
	return _OeasyMarketplace.Contract.ProxiableUUID(&_OeasyMarketplace.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_OeasyMarketplace *OeasyMarketplaceCallerSession) ProxiableUUID() ([32]byte, error) {
	return _OeasyMarketplace.Contract.ProxiableUUID(&_OeasyMarketplace.CallOpts)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x514fcac7.
//
// Solidity: function cancelOrder(uint256 nonce) returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactor) CancelOrder(opts *bind.TransactOpts, nonce *big.Int) (*types.Transaction, error) {
	return _OeasyMarketplace.contract.Transact(opts, "cancelOrder", nonce)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x514fcac7.
//
// Solidity: function cancelOrder(uint256 nonce) returns()
func (_OeasyMarketplace *OeasyMarketplaceSession) CancelOrder(nonce *big.Int) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.CancelOrder(&_OeasyMarketplace.TransactOpts, nonce)
}

// CancelOrder is a paid mutator transaction binding the contract method 0x514fcac7.
//
// Solidity: function cancelOrder(uint256 nonce) returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactorSession) CancelOrder(nonce *big.Int) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.CancelOrder(&_OeasyMarketplace.TransactOpts, nonce)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0xad7c1676.
//
// Solidity: function executeTrade((address,address,uint256,address,uint256,uint256,uint256,uint8) makerOrder, (address,address,uint256,address,uint256,uint256,uint256,uint8) taker, bytes makerSignature) returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactor) ExecuteTrade(opts *bind.TransactOpts, makerOrder IMarketplaceOrder, taker IMarketplaceOrder, makerSignature []byte) (*types.Transaction, error) {
	return _OeasyMarketplace.contract.Transact(opts, "executeTrade", makerOrder, taker, makerSignature)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0xad7c1676.
//
// Solidity: function executeTrade((address,address,uint256,address,uint256,uint256,uint256,uint8) makerOrder, (address,address,uint256,address,uint256,uint256,uint256,uint8) taker, bytes makerSignature) returns()
func (_OeasyMarketplace *OeasyMarketplaceSession) ExecuteTrade(makerOrder IMarketplaceOrder, taker IMarketplaceOrder, makerSignature []byte) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.ExecuteTrade(&_OeasyMarketplace.TransactOpts, makerOrder, taker, makerSignature)
}

// ExecuteTrade is a paid mutator transaction binding the contract method 0xad7c1676.
//
// Solidity: function executeTrade((address,address,uint256,address,uint256,uint256,uint256,uint8) makerOrder, (address,address,uint256,address,uint256,uint256,uint256,uint8) taker, bytes makerSignature) returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactorSession) ExecuteTrade(makerOrder IMarketplaceOrder, taker IMarketplaceOrder, makerSignature []byte) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.ExecuteTrade(&_OeasyMarketplace.TransactOpts, makerOrder, taker, makerSignature)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactor) Initialize(opts *bind.TransactOpts, owner_ common.Address) (*types.Transaction, error) {
	return _OeasyMarketplace.contract.Transact(opts, "initialize", owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_OeasyMarketplace *OeasyMarketplaceSession) Initialize(owner_ common.Address) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.Initialize(&_OeasyMarketplace.TransactOpts, owner_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address owner_) returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactorSession) Initialize(owner_ common.Address) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.Initialize(&_OeasyMarketplace.TransactOpts, owner_)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OeasyMarketplace.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_OeasyMarketplace *OeasyMarketplaceSession) Pause() (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.Pause(&_OeasyMarketplace.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactorSession) Pause() (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.Pause(&_OeasyMarketplace.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OeasyMarketplace.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OeasyMarketplace *OeasyMarketplaceSession) RenounceOwnership() (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.RenounceOwnership(&_OeasyMarketplace.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.RenounceOwnership(&_OeasyMarketplace.TransactOpts)
}

// SetFeeConfiguration is a paid mutator transaction binding the contract method 0xff282391.
//
// Solidity: function setFeeConfiguration(address recipient, uint96 feeBps_) returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactor) SetFeeConfiguration(opts *bind.TransactOpts, recipient common.Address, feeBps_ *big.Int) (*types.Transaction, error) {
	return _OeasyMarketplace.contract.Transact(opts, "setFeeConfiguration", recipient, feeBps_)
}

// SetFeeConfiguration is a paid mutator transaction binding the contract method 0xff282391.
//
// Solidity: function setFeeConfiguration(address recipient, uint96 feeBps_) returns()
func (_OeasyMarketplace *OeasyMarketplaceSession) SetFeeConfiguration(recipient common.Address, feeBps_ *big.Int) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.SetFeeConfiguration(&_OeasyMarketplace.TransactOpts, recipient, feeBps_)
}

// SetFeeConfiguration is a paid mutator transaction binding the contract method 0xff282391.
//
// Solidity: function setFeeConfiguration(address recipient, uint96 feeBps_) returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactorSession) SetFeeConfiguration(recipient common.Address, feeBps_ *big.Int) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.SetFeeConfiguration(&_OeasyMarketplace.TransactOpts, recipient, feeBps_)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _OeasyMarketplace.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OeasyMarketplace *OeasyMarketplaceSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.TransferOwnership(&_OeasyMarketplace.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.TransferOwnership(&_OeasyMarketplace.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OeasyMarketplace.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_OeasyMarketplace *OeasyMarketplaceSession) Unpause() (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.Unpause(&_OeasyMarketplace.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactorSession) Unpause() (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.Unpause(&_OeasyMarketplace.TransactOpts)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _OeasyMarketplace.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_OeasyMarketplace *OeasyMarketplaceSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.UpgradeToAndCall(&_OeasyMarketplace.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_OeasyMarketplace *OeasyMarketplaceTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _OeasyMarketplace.Contract.UpgradeToAndCall(&_OeasyMarketplace.TransactOpts, newImplementation, data)
}

// OeasyMarketplaceEIP712DomainChangedIterator is returned from FilterEIP712DomainChanged and is used to iterate over the raw logs and unpacked data for EIP712DomainChanged events raised by the OeasyMarketplace contract.
type OeasyMarketplaceEIP712DomainChangedIterator struct {
	Event *OeasyMarketplaceEIP712DomainChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OeasyMarketplaceEIP712DomainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OeasyMarketplaceEIP712DomainChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OeasyMarketplaceEIP712DomainChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OeasyMarketplaceEIP712DomainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OeasyMarketplaceEIP712DomainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OeasyMarketplaceEIP712DomainChanged represents a EIP712DomainChanged event raised by the OeasyMarketplace contract.
type OeasyMarketplaceEIP712DomainChanged struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEIP712DomainChanged is a free log retrieval operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_OeasyMarketplace *OeasyMarketplaceFilterer) FilterEIP712DomainChanged(opts *bind.FilterOpts) (*OeasyMarketplaceEIP712DomainChangedIterator, error) {

	logs, sub, err := _OeasyMarketplace.contract.FilterLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceEIP712DomainChangedIterator{contract: _OeasyMarketplace.contract, event: "EIP712DomainChanged", logs: logs, sub: sub}, nil
}

// WatchEIP712DomainChanged is a free log subscription operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_OeasyMarketplace *OeasyMarketplaceFilterer) WatchEIP712DomainChanged(opts *bind.WatchOpts, sink chan<- *OeasyMarketplaceEIP712DomainChanged) (event.Subscription, error) {

	logs, sub, err := _OeasyMarketplace.contract.WatchLogs(opts, "EIP712DomainChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OeasyMarketplaceEIP712DomainChanged)
				if err := _OeasyMarketplace.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEIP712DomainChanged is a log parse operation binding the contract event 0x0a6387c9ea3628b88a633bb4f3b151770f70085117a15f9bf3787cda53f13d31.
//
// Solidity: event EIP712DomainChanged()
func (_OeasyMarketplace *OeasyMarketplaceFilterer) ParseEIP712DomainChanged(log types.Log) (*OeasyMarketplaceEIP712DomainChanged, error) {
	event := new(OeasyMarketplaceEIP712DomainChanged)
	if err := _OeasyMarketplace.contract.UnpackLog(event, "EIP712DomainChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OeasyMarketplaceFeeUpdatedIterator is returned from FilterFeeUpdated and is used to iterate over the raw logs and unpacked data for FeeUpdated events raised by the OeasyMarketplace contract.
type OeasyMarketplaceFeeUpdatedIterator struct {
	Event *OeasyMarketplaceFeeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OeasyMarketplaceFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OeasyMarketplaceFeeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OeasyMarketplaceFeeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OeasyMarketplaceFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OeasyMarketplaceFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OeasyMarketplaceFeeUpdated represents a FeeUpdated event raised by the OeasyMarketplace contract.
type OeasyMarketplaceFeeUpdated struct {
	Recipient common.Address
	FeeBps    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFeeUpdated is a free log retrieval operation binding the contract event 0xcfadc6b7fb7eadbc2f0418cc2e7b2a4f37d3c39d702c9a36bea28de07fb31f66.
//
// Solidity: event FeeUpdated(address indexed recipient, uint96 feeBps)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) FilterFeeUpdated(opts *bind.FilterOpts, recipient []common.Address) (*OeasyMarketplaceFeeUpdatedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _OeasyMarketplace.contract.FilterLogs(opts, "FeeUpdated", recipientRule)
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceFeeUpdatedIterator{contract: _OeasyMarketplace.contract, event: "FeeUpdated", logs: logs, sub: sub}, nil
}

// WatchFeeUpdated is a free log subscription operation binding the contract event 0xcfadc6b7fb7eadbc2f0418cc2e7b2a4f37d3c39d702c9a36bea28de07fb31f66.
//
// Solidity: event FeeUpdated(address indexed recipient, uint96 feeBps)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) WatchFeeUpdated(opts *bind.WatchOpts, sink chan<- *OeasyMarketplaceFeeUpdated, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _OeasyMarketplace.contract.WatchLogs(opts, "FeeUpdated", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OeasyMarketplaceFeeUpdated)
				if err := _OeasyMarketplace.contract.UnpackLog(event, "FeeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFeeUpdated is a log parse operation binding the contract event 0xcfadc6b7fb7eadbc2f0418cc2e7b2a4f37d3c39d702c9a36bea28de07fb31f66.
//
// Solidity: event FeeUpdated(address indexed recipient, uint96 feeBps)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) ParseFeeUpdated(log types.Log) (*OeasyMarketplaceFeeUpdated, error) {
	event := new(OeasyMarketplaceFeeUpdated)
	if err := _OeasyMarketplace.contract.UnpackLog(event, "FeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OeasyMarketplaceInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the OeasyMarketplace contract.
type OeasyMarketplaceInitializedIterator struct {
	Event *OeasyMarketplaceInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OeasyMarketplaceInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OeasyMarketplaceInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OeasyMarketplaceInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OeasyMarketplaceInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OeasyMarketplaceInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OeasyMarketplaceInitialized represents a Initialized event raised by the OeasyMarketplace contract.
type OeasyMarketplaceInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) FilterInitialized(opts *bind.FilterOpts) (*OeasyMarketplaceInitializedIterator, error) {

	logs, sub, err := _OeasyMarketplace.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceInitializedIterator{contract: _OeasyMarketplace.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *OeasyMarketplaceInitialized) (event.Subscription, error) {

	logs, sub, err := _OeasyMarketplace.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OeasyMarketplaceInitialized)
				if err := _OeasyMarketplace.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) ParseInitialized(log types.Log) (*OeasyMarketplaceInitialized, error) {
	event := new(OeasyMarketplaceInitialized)
	if err := _OeasyMarketplace.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OeasyMarketplaceOrderCancelledIterator is returned from FilterOrderCancelled and is used to iterate over the raw logs and unpacked data for OrderCancelled events raised by the OeasyMarketplace contract.
type OeasyMarketplaceOrderCancelledIterator struct {
	Event *OeasyMarketplaceOrderCancelled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OeasyMarketplaceOrderCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OeasyMarketplaceOrderCancelled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OeasyMarketplaceOrderCancelled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OeasyMarketplaceOrderCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OeasyMarketplaceOrderCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OeasyMarketplaceOrderCancelled represents a OrderCancelled event raised by the OeasyMarketplace contract.
type OeasyMarketplaceOrderCancelled struct {
	Maker common.Address
	Nonce *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterOrderCancelled is a free log retrieval operation binding the contract event 0xdd003742fb214507783ce004fe55f5ac14f89c6de4a7cd7487e47eb091c62265.
//
// Solidity: event OrderCancelled(address indexed maker, uint256 indexed nonce)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) FilterOrderCancelled(opts *bind.FilterOpts, maker []common.Address, nonce []*big.Int) (*OeasyMarketplaceOrderCancelledIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _OeasyMarketplace.contract.FilterLogs(opts, "OrderCancelled", makerRule, nonceRule)
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceOrderCancelledIterator{contract: _OeasyMarketplace.contract, event: "OrderCancelled", logs: logs, sub: sub}, nil
}

// WatchOrderCancelled is a free log subscription operation binding the contract event 0xdd003742fb214507783ce004fe55f5ac14f89c6de4a7cd7487e47eb091c62265.
//
// Solidity: event OrderCancelled(address indexed maker, uint256 indexed nonce)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) WatchOrderCancelled(opts *bind.WatchOpts, sink chan<- *OeasyMarketplaceOrderCancelled, maker []common.Address, nonce []*big.Int) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var nonceRule []interface{}
	for _, nonceItem := range nonce {
		nonceRule = append(nonceRule, nonceItem)
	}

	logs, sub, err := _OeasyMarketplace.contract.WatchLogs(opts, "OrderCancelled", makerRule, nonceRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OeasyMarketplaceOrderCancelled)
				if err := _OeasyMarketplace.contract.UnpackLog(event, "OrderCancelled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOrderCancelled is a log parse operation binding the contract event 0xdd003742fb214507783ce004fe55f5ac14f89c6de4a7cd7487e47eb091c62265.
//
// Solidity: event OrderCancelled(address indexed maker, uint256 indexed nonce)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) ParseOrderCancelled(log types.Log) (*OeasyMarketplaceOrderCancelled, error) {
	event := new(OeasyMarketplaceOrderCancelled)
	if err := _OeasyMarketplace.contract.UnpackLog(event, "OrderCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OeasyMarketplaceOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OeasyMarketplace contract.
type OeasyMarketplaceOwnershipTransferredIterator struct {
	Event *OeasyMarketplaceOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OeasyMarketplaceOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OeasyMarketplaceOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OeasyMarketplaceOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OeasyMarketplaceOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OeasyMarketplaceOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OeasyMarketplaceOwnershipTransferred represents a OwnershipTransferred event raised by the OeasyMarketplace contract.
type OeasyMarketplaceOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OeasyMarketplaceOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OeasyMarketplace.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceOwnershipTransferredIterator{contract: _OeasyMarketplace.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OeasyMarketplaceOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OeasyMarketplace.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OeasyMarketplaceOwnershipTransferred)
				if err := _OeasyMarketplace.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) ParseOwnershipTransferred(log types.Log) (*OeasyMarketplaceOwnershipTransferred, error) {
	event := new(OeasyMarketplaceOwnershipTransferred)
	if err := _OeasyMarketplace.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OeasyMarketplacePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the OeasyMarketplace contract.
type OeasyMarketplacePausedIterator struct {
	Event *OeasyMarketplacePaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OeasyMarketplacePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OeasyMarketplacePaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OeasyMarketplacePaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OeasyMarketplacePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OeasyMarketplacePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OeasyMarketplacePaused represents a Paused event raised by the OeasyMarketplace contract.
type OeasyMarketplacePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) FilterPaused(opts *bind.FilterOpts) (*OeasyMarketplacePausedIterator, error) {

	logs, sub, err := _OeasyMarketplace.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplacePausedIterator{contract: _OeasyMarketplace.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *OeasyMarketplacePaused) (event.Subscription, error) {

	logs, sub, err := _OeasyMarketplace.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OeasyMarketplacePaused)
				if err := _OeasyMarketplace.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) ParsePaused(log types.Log) (*OeasyMarketplacePaused, error) {
	event := new(OeasyMarketplacePaused)
	if err := _OeasyMarketplace.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OeasyMarketplaceTradeExecutedIterator is returned from FilterTradeExecuted and is used to iterate over the raw logs and unpacked data for TradeExecuted events raised by the OeasyMarketplace contract.
type OeasyMarketplaceTradeExecutedIterator struct {
	Event *OeasyMarketplaceTradeExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OeasyMarketplaceTradeExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OeasyMarketplaceTradeExecuted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OeasyMarketplaceTradeExecuted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OeasyMarketplaceTradeExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OeasyMarketplaceTradeExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OeasyMarketplaceTradeExecuted represents a TradeExecuted event raised by the OeasyMarketplace contract.
type OeasyMarketplaceTradeExecuted struct {
	Maker        common.Address
	Taker        common.Address
	Nft          common.Address
	TokenId      *big.Int
	PaymentToken common.Address
	Price        *big.Int
	Side         uint8
	Fee          *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTradeExecuted is a free log retrieval operation binding the contract event 0xbe2552ec99ac4dff73b91d536b01bfb876d9613954afcbe8524eebc6a8275d0f.
//
// Solidity: event TradeExecuted(address indexed maker, address indexed taker, address indexed nft, uint256 tokenId, address paymentToken, uint256 price, uint8 side, uint256 fee)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) FilterTradeExecuted(opts *bind.FilterOpts, maker []common.Address, taker []common.Address, nft []common.Address) (*OeasyMarketplaceTradeExecutedIterator, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}
	var nftRule []interface{}
	for _, nftItem := range nft {
		nftRule = append(nftRule, nftItem)
	}

	logs, sub, err := _OeasyMarketplace.contract.FilterLogs(opts, "TradeExecuted", makerRule, takerRule, nftRule)
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceTradeExecutedIterator{contract: _OeasyMarketplace.contract, event: "TradeExecuted", logs: logs, sub: sub}, nil
}

// WatchTradeExecuted is a free log subscription operation binding the contract event 0xbe2552ec99ac4dff73b91d536b01bfb876d9613954afcbe8524eebc6a8275d0f.
//
// Solidity: event TradeExecuted(address indexed maker, address indexed taker, address indexed nft, uint256 tokenId, address paymentToken, uint256 price, uint8 side, uint256 fee)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) WatchTradeExecuted(opts *bind.WatchOpts, sink chan<- *OeasyMarketplaceTradeExecuted, maker []common.Address, taker []common.Address, nft []common.Address) (event.Subscription, error) {

	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}
	var nftRule []interface{}
	for _, nftItem := range nft {
		nftRule = append(nftRule, nftItem)
	}

	logs, sub, err := _OeasyMarketplace.contract.WatchLogs(opts, "TradeExecuted", makerRule, takerRule, nftRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OeasyMarketplaceTradeExecuted)
				if err := _OeasyMarketplace.contract.UnpackLog(event, "TradeExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTradeExecuted is a log parse operation binding the contract event 0xbe2552ec99ac4dff73b91d536b01bfb876d9613954afcbe8524eebc6a8275d0f.
//
// Solidity: event TradeExecuted(address indexed maker, address indexed taker, address indexed nft, uint256 tokenId, address paymentToken, uint256 price, uint8 side, uint256 fee)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) ParseTradeExecuted(log types.Log) (*OeasyMarketplaceTradeExecuted, error) {
	event := new(OeasyMarketplaceTradeExecuted)
	if err := _OeasyMarketplace.contract.UnpackLog(event, "TradeExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OeasyMarketplaceUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the OeasyMarketplace contract.
type OeasyMarketplaceUnpausedIterator struct {
	Event *OeasyMarketplaceUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OeasyMarketplaceUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OeasyMarketplaceUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OeasyMarketplaceUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OeasyMarketplaceUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OeasyMarketplaceUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OeasyMarketplaceUnpaused represents a Unpaused event raised by the OeasyMarketplace contract.
type OeasyMarketplaceUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) FilterUnpaused(opts *bind.FilterOpts) (*OeasyMarketplaceUnpausedIterator, error) {

	logs, sub, err := _OeasyMarketplace.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceUnpausedIterator{contract: _OeasyMarketplace.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *OeasyMarketplaceUnpaused) (event.Subscription, error) {

	logs, sub, err := _OeasyMarketplace.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OeasyMarketplaceUnpaused)
				if err := _OeasyMarketplace.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) ParseUnpaused(log types.Log) (*OeasyMarketplaceUnpaused, error) {
	event := new(OeasyMarketplaceUnpaused)
	if err := _OeasyMarketplace.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OeasyMarketplaceUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the OeasyMarketplace contract.
type OeasyMarketplaceUpgradedIterator struct {
	Event *OeasyMarketplaceUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OeasyMarketplaceUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OeasyMarketplaceUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OeasyMarketplaceUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OeasyMarketplaceUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OeasyMarketplaceUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OeasyMarketplaceUpgraded represents a Upgraded event raised by the OeasyMarketplace contract.
type OeasyMarketplaceUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*OeasyMarketplaceUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _OeasyMarketplace.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &OeasyMarketplaceUpgradedIterator{contract: _OeasyMarketplace.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *OeasyMarketplaceUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _OeasyMarketplace.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OeasyMarketplaceUpgraded)
				if err := _OeasyMarketplace.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_OeasyMarketplace *OeasyMarketplaceFilterer) ParseUpgraded(log types.Log) (*OeasyMarketplaceUpgraded, error) {
	event := new(OeasyMarketplaceUpgraded)
	if err := _OeasyMarketplace.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
