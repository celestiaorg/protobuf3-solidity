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
          optionalInt32: -42,
          optionalInt64: -420,
          optionalUint32: 42,
          optionalUint64: 420,
          optionalSint32: -69,
          optionalSint64: -690,
          optionalFixed32: 900,
          optionalFixed64: 9000,
          optionalSfixed32: -900,
          optionalSfixed64: -9000,
          optionalBool: true,
          // optionsString: "foorbar",
          // optionsBytes: "0xdeadbeef",
        };

        const message = Message.create(messageObj);
        const encoded = Message.encode(message).finish().toString("hex");
        console.log(encoded);

        const result = await instance.decode.call("0x" + encoded);
        const { 0: success, 1: decoded } = result;
        console.log(decoded);
        assert.equal(success, true);

        // await instance.decode("0x" + encoded);
      });
    });
  });
});
