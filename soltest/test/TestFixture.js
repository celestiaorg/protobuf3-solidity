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

        const OtherMessage = root.lookupType("OtherMessage");
        otherMessageObj = {
          otherField: 3,
        };
        const otherMessage = OtherMessage.create(otherMessageObj);

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
          optionalString: "foorbar",
          optionalBytes: Buffer.from("deadbeef", "hex"),
          optionalEnum: 1,
          optionalMessage: otherMessage,
          repeatedInt32: [-42, -41],
          repeatedInt64: [-420, -421],
          repeatedUint32: [42, 41],
          repeatedUint64: [420, 419],
          repeatedSint32: [-69, -68],
          repeatedSint64: [-690, -689],
          repeatedFixed32: [900, 899],
          repeatedFixed64: [9000, 8999],
          repeatedSfixed32: [-900, -899],
          repeatedSfixed64: [-9000, -8999],
          repeatedBool: [true, false],
          repeatedEnum: [1, 2],
          repeatedMessage: [otherMessage, otherMessage],
        };

        const message = Message.create(messageObj);
        console.log(message);
        const encoded = Message.encode(message).finish().toString("hex");
        console.log(encoded);

        const result = await instance.decode.call("0x" + encoded);
        const { 0: success, 1: decoded } = result;
        console.log(decoded);
        assert.equal(success, true);

        await instance.decode("0x" + encoded);
      });
    });
  });
});
