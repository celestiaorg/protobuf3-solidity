// SPDX-License-Identifier: Apache-2.0
pragma solidity >=0.6.0 <8.0.0;
pragma experimental ABIEncoderV2;

import "@lazyledger/protobuf3-solidity-lib/contracts/ProtobufLib.sol";
import "./all_features.proto.sol";
import "./top.proto.sol";

contract TestFixture {
    // Functions are not pure so that we can measure gas

    function decode(bytes memory buf) public returns (bool, Message memory) {
        (bool success, uint64 pos, Message memory instance) = MessageCodec.decode(0, buf, uint64(buf.length));

        return (success, instance);
    }

    // function encode(Message memory instance) public returns (bytes memory) {
    //     return MessageCodec.encode(instance);
    // }
}
