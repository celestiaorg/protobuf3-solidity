const protobuf = require("protobufjs");
const truffleAssert = require("truffle-assertions");

const TestFixture = artifacts.require("TestFixture");

const AllFeaturesProtoFile = "../test/pass/all_features/all_features.proto";

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

        const root = await protobuf.load(AllFeaturesProtoFile);

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
          repeatedInt32: ["-42", "-41"],
          repeatedInt64: ["-420", "-421"],
          repeatedUint32: ["42", "41"],
          repeatedUint64: ["420", "419"],
          repeatedSint32: ["-69", "-68"],
          repeatedSint64: ["-690", "-689"],
          repeatedFixed32: ["900", "899"],
          repeatedFixed64: ["9000", "8999"],
          repeatedSfixed32: ["-900", "-899"],
          repeatedSfixed64: ["-9000", "-8999"],
          repeatedBool: [true, false],
          repeatedEnum: ["1", "2"],
          repeatedMessage: [otherMessage, otherMessage],
        };

        const message = Message.create(messageObj);
        const encoded = Message.encode(message).finish().toString("hex");

        const result = await instance.decode.call("0x" + encoded);
        const { 0: success, 1: decoded } = result;
        assert.equal(success, true);
        assert.equal(decoded.optional_int32, messageObj.optionalInt32);
        assert.equal(decoded.optional_int64, messageObj.optionalInt64);
        assert.equal(decoded.optional_uint32, messageObj.optionalUint32);
        assert.equal(decoded.optional_uint64, messageObj.optionalUint64);
        assert.equal(decoded.optional_sint32, messageObj.optionalSint32);
        assert.equal(decoded.optional_sint64, messageObj.optionalSint64);
        assert.equal(decoded.optional_fixed32, messageObj.optionalFixed32);
        assert.equal(decoded.optional_fixed64, messageObj.optionalFixed64);
        assert.equal(decoded.optional_sfixed32, messageObj.optionalSfixed32);
        assert.equal(decoded.optional_sfixed64, messageObj.optionalSfixed64);
        assert.equal(decoded.optional_bool, messageObj.optionalBool);
        assert.equal(decoded.optional_string, messageObj.optionalString);
        assert.equal(decoded.optional_bytes.slice(2), messageObj.optionalBytes.toString("hex"));
        assert.equal(decoded.optional_enum, messageObj.optionalEnum);
        assert.equal(decoded.optional_message.other_field, messageObj.optionalMessage.otherField);
        assert.deepStrictEqual(decoded.repeated_int32, messageObj.repeatedInt32);
        assert.deepStrictEqual(decoded.repeated_int64, messageObj.repeatedInt64);
        assert.deepStrictEqual(decoded.repeated_uint32, messageObj.repeatedUint32);
        assert.deepStrictEqual(decoded.repeated_uint64, messageObj.repeatedUint64);
        assert.deepStrictEqual(decoded.repeated_sint32, messageObj.repeatedSint32);
        assert.deepStrictEqual(decoded.repeated_sint64, messageObj.repeatedSint64);
        assert.deepStrictEqual(decoded.repeated_fixed32, messageObj.repeatedFixed32);
        assert.deepStrictEqual(decoded.repeated_fixed64, messageObj.repeatedFixed64);
        assert.deepStrictEqual(decoded.repeated_sfixed32, messageObj.repeatedSfixed32);
        assert.deepStrictEqual(decoded.repeated_sfixed64, messageObj.repeatedSfixed64);
        assert.deepStrictEqual(decoded.repeated_bool, messageObj.repeatedBool);
        assert.deepStrictEqual(decoded.repeated_enum, messageObj.repeatedEnum);
        assert.equal(decoded.repeated_message.length, messageObj.repeatedMessage.length);
        assert.equal(decoded.repeated_message[0].other_field, messageObj.repeatedMessage[0].otherField);
        assert.equal(decoded.repeated_message[1].other_field, messageObj.repeatedMessage[1].otherField);

        await instance.decode("0x" + encoded);
      });

      it("field start >1", async () => {
        const instance = await TestFixture.deployed();

        const root = await protobuf.load(AllFeaturesProtoFile);

        const Message = root.lookupType("Message");
        const messageObj = {
          optionalUint64: 420,
        };

        const message = Message.create(messageObj);
        const encoded = Message.encode(message).finish().toString("hex");

        const result = await instance.decode.call("0x" + encoded);
        const { 0: success, 1: decoded } = result;
        assert.equal(success, true);
        assert.equal(decoded.optional_uint64, messageObj.optionalUint64);

        await instance.decode("0x" + encoded);
      });
    });

    describe("failing", async () => {
      it("fields out of order", async () => {
        const instance = await TestFixture.deployed();

        const messageObj = {
          optionalUint32: 1, // 1801
          optionalUint64: 1, // 2001
        };

        const encoded = "20011801";

        const result = await instance.decode.call("0x" + encoded);
        const { 0: success, 1: decoded } = result;
        assert.equal(success, false);
      });

      it("repeated not-repeated field", async () => {
        const instance = await TestFixture.deployed();

        const messageObj = {
          optionalUint64: 1, // 2001
        };

        const encoded = "20012001";

        const result = await instance.decode.call("0x" + encoded);
        const { 0: success, 1: decoded } = result;
        assert.equal(success, false);
      });

      it("included default value", async () => {
        const instance = await TestFixture.deployed();

        const messageObj = {
          optionalUint32: 0, // 1800
          optionalUint64: 1, // 2001
        };

        const encoded = "18002001";

        const result = await instance.decode.call("0x" + encoded);
        const { 0: success, 1: decoded } = result;
        assert.equal(success, false);
      });

      it("extra data", async () => {
        const instance = await TestFixture.deployed();

        const messageObj = {
          optionalUint32: 1, // 1801
          optionalUint64: 1, // 2001
        };

        const encoded = "18012001deadbeef";

        const result = await instance.decode.call("0x" + encoded);
        const { 0: success, 1: decoded } = result;
        assert.equal(success, false);
      });
    });
  });

  describe("encode", async () => {});
});
