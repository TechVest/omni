// SPDX-License-Identifier: GPL-3.0-only
pragma solidity =0.8.24;

import { OwnableRoles } from "solady/src/auth/OwnableRoles.sol";
import { ReentrancyGuard } from "solady/src/utils/ReentrancyGuard.sol";
import { Initializable } from "solady/src/utils/Initializable.sol";
import { XAppBase } from "core/src/pkg/XAppBase.sol";

import { SafeTransferLib } from "solady/src/utils/SafeTransferLib.sol";
import { ConfLevel } from "core/src/libraries/ConfLevel.sol";
import { TypeMax } from "core/src/libraries/TypeMax.sol";
import { Solve } from "./Solve.sol";

import { ISolveInbox } from "./interfaces/ISolveInbox.sol";

/**
 * @title SolveOutbox
 * @notice Entrypoint for fulfillments of user solve requests.
 */
contract SolveOutbox is OwnableRoles, ReentrancyGuard, Initializable, XAppBase {
    using SafeTransferLib for address;

    error CallFailed();
    error CallNotAllowed();
    error AreadyFulfilled();
    error WrongDestChain();
    error IncorrectPreReqs();
    error InsufficientFee();

    event AllowedCallSet(address indexed target, bytes4 indexed selector, bool allowed);

    /**
     * @notice Emitted when a request is fulfilled.
     * @param reqId       ID of the request.
     * @param callHash    Hash of the call executed.
     * @param solvedBy    Address of the solver.
     */
    event Fulfilled(bytes32 indexed reqId, bytes32 indexed callHash, address indexed solvedBy);

    /**
     * @notice Role for solvers.
     * @dev _ROLE_0 evaluates to '1'.
     */
    uint256 internal constant SOLVER = _ROLE_0;

    /**
     * @notice Gas limit for SolveInbox.markFulfilled callback.
     */
    uint64 internal constant MARK_FULFILLED_GAS_LIMIT = 100_000;

    /**
     * @notice Stubbed calldata for SolveInbox.markFulfilled. Used to estimate the gas cost.
     * @dev Type maxes used to ensure no non-zero bytes in fee estimation.
     */
    bytes internal constant MARK_FULFILLED_STUB_CDATA =
        abi.encodeCall(ISolveInbox.markFulfilled, (TypeMax.Bytes32, TypeMax.Bytes32));

    /**
     * @notice Address of the inbox contract.
     */
    address internal _inbox;

    /**
     * @notice Mapping of allowed calls per contract.
     */
    mapping(address target => mapping(bytes4 selector => bool)) public allowedCalls;

    /**
     * @notice Mapping of fulfilled calls.
     * @dev callHash used to prevent duplicate fulfillment.
     */
    mapping(bytes32 callHash => bool fulfilled) public fulfilledCalls;

    constructor() {
        _disableInitializers();
    }

    /**
     * @notice Initialize the contract's owner and solver.
     * @dev Used instead of constructor as we want to use the transparent upgradeable proxy pattern.
     * @param owner_  Address of the owner.
     * @param solver_ Address of the solver.
     */
    function initialize(address owner_, address solver_, address omni_, address inbox_) external initializer {
        _initializeOwner(owner_);
        _grantRoles(solver_, SOLVER);
        _setOmniPortal(omni_);
        _inbox = inbox_;
    }

    /**
     * @notice Calculate the message passing fee for a fulfill call.
     * @param sourceChainId ID of the source chain.
     */
    function fulfillFee(uint64 sourceChainId) public view returns (uint256) {
        return feeFor(sourceChainId, MARK_FULFILLED_STUB_CDATA, MARK_FULFILLED_GAS_LIMIT);
    }

    /**
     * @notice Check if a call has been fulfilled.
     * @param reqId          ID of the request.
     * @param sourceChainId  ID of the source chain.
     * @param call           Details of the call executed.
     */
    function didFulfill(bytes32 reqId, uint64 sourceChainId, Solve.Call calldata call) external view returns (bool) {
        return fulfilledCalls[_callHash(reqId, sourceChainId, call)];
    }

    /**
     * @notice Set an allowed call for a target contract.
     * @param target    Address of the target contract.
     * @param selector  4-byte selector of the function to allow.
     * @param allowed   Whether the call is allowed.
     */
    function setAllowedCall(address target, bytes4 selector, bool allowed) external onlyOwner {
        allowedCalls[target][selector] = allowed;
        emit AllowedCallSet(target, selector, allowed);
    }

    /**
     * @notice Fulfill a request.
     * @param reqId          ID of the request.
     * @param sourceChainId  ID of the source chain.
     * @param call           Details of the call to execute.
     * @param prereqs        Pre-requisite token deposits required by the call.
     */
    function fulfill(
        bytes32 reqId,
        uint64 sourceChainId,
        Solve.Call calldata call,
        Solve.TokenPrereq[] calldata prereqs
    ) external payable onlyRoles(SOLVER) nonReentrant {
        if (call.destChainId != block.chainid) revert WrongDestChain();
        if (!allowedCalls[call.target][bytes4(call.data)]) revert CallNotAllowed();

        // Check if the call has already been fulfilled
        bytes32 callHash = _callHash(reqId, sourceChainId, call);
        if (fulfilledCalls[callHash]) revert AreadyFulfilled();

        // Process token pre-requisites
        // Record pre-call balances (checked against post-call balances)
        uint256[] memory prereqBalances = new uint256[](prereqs.length);
        for (uint256 i; i < prereqs.length; ++i) {
            prereqBalances[i] = prereqs[i].token.balanceOf(address(this));
            prereqs[i].token.safeTransferFrom(msg.sender, address(this), prereqs[i].amount);
            prereqs[i].token.safeApprove(prereqs[i].spender, prereqs[i].amount);
        }

        // Execute the call
        (bool success,) = payable(call.target).call{ value: call.value }(call.data);
        if (!success) revert CallFailed();

        // Require post-call balances matches pre-call balances
        // Ensures token pre-requisites match the call's required transfers
        for (uint256 i; i < prereqs.length; ++i) {
            if (prereqs[i].token.balanceOf(address(this)) != prereqBalances[i]) revert IncorrectPreReqs();
        }

        // Mark the call as fulfilled on inbox
        bytes memory xcalldata = abi.encodeCall(ISolveInbox.markFulfilled, (reqId, callHash));
        uint256 fee = xcall(sourceChainId, ConfLevel.Finalized, _inbox, xcalldata, MARK_FULFILLED_GAS_LIMIT);
        if (msg.value - call.value < fee) revert InsufficientFee();

        emit Fulfilled(reqId, callHash, msg.sender);
    }

    /**
     * @dev Returns call hash. Used to discern fullfilment.
     */
    function _callHash(bytes32 id, uint64 sourceChainId, Solve.Call calldata call) internal pure returns (bytes32) {
        return keccak256(abi.encode(id, sourceChainId, call));
    }
}
