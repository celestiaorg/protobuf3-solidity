package generator

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// SolidityVersionString is the Solidity version specifier.
const SolidityVersionString = ">=0.6.0 <8.0.0"

// SolidityABIString indicates ABIEncoderV2 use.
const SolidityABIString = "pragma experimental ABIEncoderV2;"

// Generator generates Solidity code from .proto files.
type Generator struct {
	request *pluginpb.CodeGeneratorRequest
}

// New initializes a new Generator.
func New(request *pluginpb.CodeGeneratorRequest) *Generator {
	g := new(Generator)
	g.request = request
	return g
}

// Generate generates Solidity code from the requested .proto files.
func (g *Generator) Generate() ([]*pluginpb.CodeGeneratorResponse_File, error) {
	protoFiles := g.request.GetProtoFile()
	responseFiles := make([]*pluginpb.CodeGeneratorResponse_File, len(protoFiles))

	for i := 0; i < len(protoFiles); i++ {
		responseFile, err := generateFile(protoFiles[i])
		if err != nil {
			return nil, err
		}

		responseFiles[i] = responseFile
	}

	return responseFiles, nil
}

// generateFile generates Solidity code from a single .proto file.
func generateFile(protoFile *descriptorpb.FileDescriptorProto) (*pluginpb.CodeGeneratorResponse_File, error) {
	err := checkSyntaxVersion(protoFile.GetSyntax())
	if err != nil {
		return nil, err
	}

	responseFile := &pluginpb.CodeGeneratorResponse_File{}

	b := &WriteableBuffer{}

	// TODO option for license
	b.P(fmt.Sprintf("// SPDX-License-Identifier: %s", "CC0"))
	b.P("pragma solidity " + SolidityVersionString + ";")
	b.P(SolidityABIString)
	b.P()
	b.P("import \"@lazyledger/protobuf3-solidity-lib/contracts/ProtobufLib.sol\";")
	b.P()

	for i := 0; i < len(protoFile.GetMessageType()); i++ {
		err := generateMessage(protoFile.GetMessageType()[i], b)
		if err != nil {
			return nil, err
		}
	}

	// TODO add b to response
	println(b.String())

	return responseFile, nil
}

func generateMessage(descriptor *descriptorpb.DescriptorProto, b *WriteableBuffer) error {
	structName := descriptor.GetName()
	err := checkKeyword(structName)
	if err != nil {
		return err
	}

	fields := descriptor.GetField()

	////////////////////////////////////
	// Generate struct
	////////////////////////////////////

	b.P(fmt.Sprintf("struct %s {", structName))
	b.Indent()

	fieldCount := int32(0)
	// Loop over fields
	for _, field := range fields {
		fieldNumber := field.GetNumber()
		if fieldNumber != fieldCount+1 {
			return errors.New("Field " + string(fieldNumber) + " does not increment by 1")
		}
		fieldCount++

		fieldDescriptorType := field.GetType()
		switch fieldDescriptorType {
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			return errors.New("Unsupported field type " + fieldDescriptorType.String())
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			return errors.New("Unsupported field type " + fieldDescriptorType.String())
		default:
			// Convert protobuf field type to Solidity native type
			fieldType, err := typeToSol(fieldDescriptorType)
			if err != nil {
				return err
			}
			fieldName := field.GetName()
			err = checkKeyword(fieldName)
			if err != nil {
				return err
			}

			arrayStr := ""
			if isRepeated(field.GetLabel()) {
				if isPrimitiveNumericType(fieldDescriptorType) {
					if !field.GetOptions().GetPacked() {
						return errors.New("Repeated field " + structName + "." + fieldName + " must be packed")
					}
				}
				arrayStr = "[]"
			}

			b.P(fmt.Sprintf("%s%s %s;", fieldType, arrayStr, fieldName))
		}
	}

	b.Unindent()
	b.P("}")
	b.P()

	////////////////////////////////////
	// Generate decoder
	////////////////////////////////////

	b.P(fmt.Sprintf("library %sCodec {", structName))
	b.Indent()

	// Top-level decoder function
	b.P(fmt.Sprintf("function decode(uint256 initial_pos, bytes memory buf, uint256 len) internal pure returns (bool, uint256, %s memory) {", structName))
	b.Indent()

	b.P("// Message instance")
	b.P(fmt.Sprintf("%s memory instance;", structName))
	b.P("// Previous field number")
	b.P("uint64 previous_field_number = 0;")
	b.P("// Current position in the buffer")
	b.P("uint256 pos = initial_pos;")
	b.P()

	b.P("// Sanity checks")
	b.P("if (initial_pos + len < initial_pos) {")
	b.Indent()
	b.P("return (false, pos, instance);")
	b.Unindent()
	b.P("}")
	b.P()

	b.P("while (pos - initial_pos < len) {")
	b.Indent()
	b.P("// Decode the key (field number and wire type)")
	b.P("(bool success, pos, uint64 field_number, ProtobufLib.WireType wire_type) = ProtobufLib.decode_key(pos, buf);")
	b.P("if (!success) {")
	b.Indent()
	b.P("return (false, pos, instance);")
	b.Unindent()
	b.P("}")
	b.P()

	b.P("// Check that the field number is within bounds")
	b.P(fmt.Sprintf("if (field_number > %d) {", fieldCount))
	b.Indent()
	b.P("return (false, pos, instance);")
	b.Unindent()
	b.P("}")
	b.P()

	b.P("// Check that the field number of monotonically increasing")
	b.P("if (field_number <= previous_field_number) {")
	b.Indent()
	b.P("return (false, pos, instance);")
	b.Unindent()
	b.P("}")
	b.P()

	b.P("// Check that the wire type is correct")
	b.P("success = check_key(field_number, wire_type);")
	b.P("if (!success) {")
	b.Indent()
	b.P("return (false, pos, instance);")
	b.Unindent()
	b.P("}")
	b.P()

	b.P("// Actually decode the field")
	b.P("(success, pos) = decode_field(pos, buf, instance);")
	b.P("if (!success) {")
	b.Indent()
	b.P("return (false, pos, instance);")
	b.Unindent()
	b.P("}")
	b.P()

	b.P("previous_field_number = field_number;")
	b.Unindent()
	b.P("}")
	b.P()

	b.P("// Decoding must have consumed len bytes")
	b.P("if (pos != initial_pos + len) {")
	b.Indent()
	b.P("return (false, pos, instance);")
	b.Unindent()
	b.P("}")
	b.P()

	b.P("return (true, pos, instance);")
	b.Unindent()
	b.P("}")
	b.P()

	// Check key function
	b.P("function check_key(uint64 field_number, ProtobufLib.WireType wire_type) internal pure returns (bool) {")
	b.Indent()
	for _, field := range fields {
		fieldDescriptorType := field.GetType()
		fieldNumber := field.GetNumber()

		b.P(fmt.Sprintf("if (field_number == %d) {", fieldNumber))
		b.Indent()
		wireStr, err := toSolWireType(fieldDescriptorType)
		if err != nil {
			return err
		}
		b.P(fmt.Sprintf("return wire_type == %s;", wireStr))
		b.Unindent()
		b.P("}")
		b.P()
	}

	b.P("return false;")
	b.Unindent()
	b.P("}")
	b.P()

	// Decode field dispatcher function
	b.P(fmt.Sprintf("function decode_field(uint256 initial_pos, bytes memory buf, uint256 len, uint64 field_number, %s memory instance) internal pure returns (bool, uint256) {", structName))
	b.Indent()
	b.P("uint256 pos = initial_pos;")
	b.P()

	for _, field := range fields {
		fieldNumber := field.GetNumber()

		b.P(fmt.Sprintf("if (field_number == %d) {", fieldNumber))
		b.Indent()
		b.P(fmt.Sprintf("(success, pos) = decode_%d(pos, buf, instance);", fieldNumber))
		b.P("if (!success) {")
		b.Indent()
		b.P("return (false, pos);")
		b.Unindent()
		b.P("}")
		b.P()

		b.P("return (true, pos);")
		b.Unindent()
		b.P("}")
		b.P()
	}

	b.P("return (false, pos);")
	b.Unindent()
	b.P("}")
	b.P()

	// Individual field decoders
	for _, field := range fields {
		fieldName := field.GetName()
		fieldDescriptorType := field.GetType()
		fieldType, err := typeToSol(fieldDescriptorType)
		if err != nil {
			return err
		}
		fieldNumber := field.GetNumber()

		b.P(fmt.Sprintf("function decode_%d(uint256 pos, bytes memory buf, %s memory instance) internal pure returns (bool, uint256) {", fieldNumber, structName))
		b.Indent()

		b.P("bool success;")
		b.P()

		if isRepeated(field.GetLabel()) {
			b.P(fmt.Sprintf("(success, pos, uint64 repeated_bytes = decode_uint64(pos, buf);"))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P()

			// Do one pass to count the number of elements
			b.P("while (repeated_bytes > 0) {")
			b.Indent()
			b.Unindent()
			b.P("}")
			b.P()

			// Allocated memory
			b.P()

			// Now actually parse the elements
		}

		switch fieldDescriptorType {
		case descriptorpb.FieldDescriptorProto_TYPE_INT32:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_INT64:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			b.P(fmt.Sprintf("(success, pos, %s v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_STRING:
			// TODO do this right
			b.P(fmt.Sprintf("(success, pos, %s memory v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
			// TODO do this right
			b.P(fmt.Sprintf("(success, pos, %s memory v) = decode_%s(pos, buf);", fieldType, fieldType))
			b.P("if (!success) {")
			b.Indent()
			b.P("return (false, pos);")
			b.Unindent()
			b.P("}")
			b.P(fmt.Sprintf("instance.%s = v;", fieldName))
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			return errors.New("Unsupported field type " + fieldDescriptorType.String())
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			return errors.New("Unsupported field type " + fieldDescriptorType.String())
		default:
			return errors.New("Unsupported field type " + fieldDescriptorType.String())
		}

		b.P()
		b.P("return (true, pos);")

		b.Unindent()
		b.P("}")
		b.P()
	}

	////////////////////////////////////
	// Generate encoder
	////////////////////////////////////

	b.P(fmt.Sprintf("function encode(%s memory msg) internal pure returns (bytes memory) {", structName))
	b.Indent()

	b.P("revert(\"Unimplemented feature: encoding\");")

	b.Unindent()
	b.P("}")

	b.Unindent()
	b.P("}")
	b.P()

	return nil
}

func checkSyntaxVersion(v string) error {
	if v == "proto3" {
		return nil
	}

	return errors.New("Must use syntax = \"proto3\";")
}

func checkKeyword(w string) error {
	switch w {
	case "after",
		"alias",
		"apply",
		"auto",
		"case",
		"copyof",
		"default",
		"define",
		"final",
		"immutable",
		"implements",
		"in",
		"inline",
		"let",
		"macro",
		"match",
		"mutable",
		"null",
		"of",
		"partial",
		"promise",
		"reference",
		"relocatable",
		"sealed",
		"sizeof",
		"static",
		"supports",
		"switch",
		"typedef",
		"typeof",
		"unchecked":
		return errors.New("Using Solidity keyword " + w)
	}

	return nil
}

func typeToSol(fType descriptorpb.FieldDescriptorProto_Type) (string, error) {
	s := ""

	switch fType {
	case descriptorpb.FieldDescriptorProto_TYPE_INT32:
		s = "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_INT64:
		s = "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		s = "uint32"
	case descriptorpb.FieldDescriptorProto_TYPE_UINT64:
		s = "uint64"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT32:
		s = "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		s = "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		s = "uint32"
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		s = "uint64"
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		s = "int32"
	case descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		s = "int64"
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		s = "bool"
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		s = "string"
	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		s = "bytes"
	default:
		return "", errors.New("Unsupported field type " + fType.String())
	}

	err := checkKeyword(s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func isPrimitiveNumericType(fType descriptorpb.FieldDescriptorProto_Type) bool {
	switch fType {
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return true
	}
	return false
}

func isRepeated(label descriptorpb.FieldDescriptorProto_Label) bool {
	return label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED
}

func toSolWireType(fType descriptorpb.FieldDescriptorProto_Type) (string, error) {
	switch fType {
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64,
		descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return "ProtobufLib.WireType.Varint", nil
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		return "ProtobufLib.WireType.Bits32", nil
	case descriptorpb.FieldDescriptorProto_TYPE_FIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		return "ProtobufLib.WireType.Bits64", nil
	case descriptorpb.FieldDescriptorProto_TYPE_STRING,
		descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return "ProtobufLib.WireType.LengthDelimited", nil
	}

	return "", errors.New("Unsupported field type " + fType.String())
}
