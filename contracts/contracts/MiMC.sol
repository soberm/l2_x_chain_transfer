// Copyright (c) 2018 HarryR
// SPDX-License-Identifier: LGPL-3.0+

//Adapted to use exponent of 5
pragma solidity ^0.8.19;

/**
 * Implements MiMC-p/p over the altBN scalar field used by zkSNARKs
 *
 * See: https://eprint.iacr.org/2016/492.pdf
 *
 * Round constants are generated in sequence from a seed
 */
library MiMC {
    function getScalarField() internal pure returns (uint256) {
        return
            0x30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001;
    }

    function encipher(
        uint256 in_x,
        uint256 in_k
    ) public pure returns (uint256 out_x) {
        return MiMCpe5(in_x, in_k, uint256(keccak256("seed")), 110);
    }

    /**
     * MiMC-p/p with exponent of 5
     *
     * Recommended at least 46 rounds, for a polynomial degree of 2^126
     */
    function MiMCpe5(
        uint256 in_x,
        uint256 in_k,
        uint256 in_seed,
        uint256 round_count
    ) internal pure returns (uint256 out_x) {
        assembly {
            if lt(round_count, 1) {
                revert(0, 0)
            }

        // Initialise round constants, k will be hashed
            let c := mload(0x40)
            mstore(0x40, add(c, 32))
            mstore(c, in_seed)

            let
            localQ
            := 0x30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001
            let t
            let a

        // Further n-2 subsequent rounds include a round constant
            for {
                let i := round_count
            } gt(i, 0) {
                i := sub(i, 1)
            } {
            // c = H(c)
                mstore(c, keccak256(c, 32))

            // x = pow(x + c_i, 7, p) + k
                t := addmod(addmod(in_x, mload(c), localQ), in_k, localQ) // t = x + c_i + k
                a := mulmod(t, t, localQ) // t^2
                in_x := mulmod(mulmod(a, a, localQ), t, localQ) // t^5
            }

        // Result adds key again as blinding factor
            out_x := addmod(in_x, in_k, localQ)
        }
    }

    function MiMCpe5_mp(
        uint256[] memory in_x,
        uint256 in_k,
        uint256 in_seed,
        uint256 round_count
    ) internal pure returns (uint256) {
        uint256 r = in_k;
        uint256 localQ = 0x30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001;

        for (uint256 i = 0; i < in_x.length; i++) {
            r =
                (r + in_x[i] + MiMCpe5(in_x[i], r, in_seed, round_count)) %
                localQ;
        }

        return r;
    }

    function hash(
        uint256[] memory in_msgs,
        uint256 in_key
    ) public pure returns (uint256) {
        return MiMCpe5_mp(in_msgs, in_key, uint256(keccak256("seed")), 110);
    }

    function hash(uint256[] memory in_msgs) public pure returns (uint256) {
        return hash(in_msgs, 0);
    }
}