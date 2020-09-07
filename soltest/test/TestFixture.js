const protobuf = require("protobufjs");
const truffleAssert = require("truffle-assertions");

const TestFixture = artifacts.require("TestFixture");

contract("TestFixture", async (accounts) => {
  describe("constructor", async () => {
    it("should deploy", async () => {
      await TestFixture.deployed();
    });
  });

  //////////////////////////////////////
  // NOTICE
  // Tests call functions twice, once to run and another to measure gas.
  //////////////////////////////////////

  describe("decode", async () => {
    describe("passing", async () => {
      it("all features", async () => {
        const instance = await TestFixture.deployed();

        const root = await protobuf.load("../test/pass/all_features/all_features.proto");
        const Message = root.lookupType("Message");
        const messageObj = {
          optional_int32: -42,
          optional_int64: -420,
          optional_uint32: 42,
          optional_uint64: 420,
          optional_sint32: -69,
          optional_sint64: -690,
          optional_fixed32: 900,
          optional_fixed64: 9000,
          optional_sfixed32: -900,
          optional_sfixed64: -9000,
          optional_bool: true,
          optional_string: "foobar",
          optional_bytes: "0xdeadbeef",
        };
        const message = Message.create(messageObj);
        const encoded = Message.encode(message).finish().toString("hex");

        const result = await instance.decode.call("0x" + encoded);
        const { 0: success, 1: decoded } = result;
        assert.equal(success, true);

        // await instance.decode("0x" + encoded);
      });
    });
  });
});
