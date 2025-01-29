// SPDX-License-Identifier: GPL-3.0-only
pragma solidity 0.8.26;

import "../TestBase.sol";

contract GeneralBridgeTest is TestBase {
    function test_initialize_reverts() public {
        address impl = address(new Bridge());
        bridgeWithLockbox = Bridge(address(new TransparentUpgradeableProxy(impl, admin, "")));

        // `admin` cannot be zero address.
        vm.expectRevert(IBridge.ZeroAddress.selector);
        bridgeWithLockbox.initialize({
            admin_: address(0),
            pauser_: address(0),
            omni_: address(0),
            token_: address(0),
            lockbox_: address(0)
        });

        // `pauser` cannot be zero address.
        vm.expectRevert(IBridge.ZeroAddress.selector);
        bridgeWithLockbox.initialize({
            admin_: admin,
            pauser_: address(0),
            omni_: address(0),
            token_: address(0),
            lockbox_: address(0)
        });

        // `omni` cannot be zero address.
        vm.expectRevert(IBridge.ZeroAddress.selector);
        bridgeWithLockbox.initialize({
            admin_: admin,
            pauser_: pauser,
            omni_: address(0),
            token_: address(0),
            lockbox_: address(0)
        });

        // `token` cannot be zero address.
        vm.expectRevert(IBridge.ZeroAddress.selector);
        bridgeWithLockbox.initialize({
            admin_: admin,
            pauser_: pauser,
            omni_: address(omni),
            token_: address(0),
            lockbox_: address(0)
        });

        // Initialization works with all fields populated.
        bridgeWithLockbox.initialize({
            admin_: admin,
            pauser_: pauser,
            omni_: address(omni),
            token_: address(token),
            lockbox_: address(lockbox)
        });

        // Initialization can also occur without a `lockbox` if none is present.
        bridgeNoLockbox = Bridge(address(new TransparentUpgradeableProxy(impl, admin, "")));
        bridgeNoLockbox.initialize({
            admin_: admin,
            pauser_: pauser,
            omni_: address(omni),
            token_: address(token),
            lockbox_: address(0)
        });
    }

    function test_setRoutes_reverts() public prank(admin) {
        uint64[] memory chainIds = new uint64[](1);
        chainIds[0] = DEST_CHAIN_ID + 1;
        address[] memory bridgeAddrs;

        vm.expectRevert(IBridge.ArrayLengthMismatch.selector);
        bridgeWithLockbox.setRoutes(chainIds, bridgeAddrs);
        bridgeAddrs = new address[](1);

        vm.expectRevert(IBridge.ZeroAddress.selector);
        bridgeWithLockbox.setRoutes(chainIds, bridgeAddrs);
        bridgeAddrs[0] = makeAddr("newBridge");

        // Configuration successful with valid inputs.
        bridgeWithLockbox.setRoutes(chainIds, bridgeAddrs);
    }

    function test_setRoutes_succeeds() public prank(admin) {
        uint64[] memory chainIds = new uint64[](1);
        chainIds[0] = DEST_CHAIN_ID + 1;
        address[] memory bridgeAddrs = new address[](1);
        bridgeAddrs[0] = makeAddr("newBridge");

        bridgeWithLockbox.setRoutes(chainIds, bridgeAddrs);
    }
}
