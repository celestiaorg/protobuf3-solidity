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
				arrayStr = "[]"
				if !field.GetOptions().GetPacked() {
					// TODO repeated is false for primitive numeric, why?
					// return errors.New("Repeated field " + structName + "." + fieldName + " must be packed")
				}
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

	b.P(fmt.Sprintf("function decode(bytes memory buf) internal pure returns (bool, %s memory) {", structName))
	b.Indent()

	b.P("// Message instance")
	b.P(fmt.Sprintf("%s memory instance;", structName))
	b.P("// Current field number")
	b.P("uint64 current_field_number;")
	b.P("// Current position in the buffer")
	b.P("uint256 pos;")
	b.P()

	b.P("while (pos < buf.length) {")
	b.Indent()
	b.P("(bool success, pos, uint64 field_number, ProtobufLib.WireType wire_type) = ProtobufLib.decode_key(pos, buf);")
	b.P("if (!success) {")
	b.Indent()
	b.P("return (false, instance);")
	b.Unindent()
	b.P("}")
	b.P()

	b.P("// Check that the field number is bounded")
	b.P(fmt.Sprintf("if (field_number > %d) {", fieldCount))
	b.Indent()
	b.P("return (false, instance);")
	b.Unindent()
	b.P("}")
	b.P()

	b.Unindent()
	b.P("}")

	for _, field := range fields {
		fieldName := field.GetName()
		fieldDescriptorType := field.GetType()
		fieldType, err := typeToSol(fieldDescriptorType)
		if err != nil {
			return err
		}
		fieldNumber := field.GetNumber()

		b.P()
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
	}

	b.Unindent()
	b.P("}")

	////////////////////////////////////
	// Generate encoder
	////////////////////////////////////

	b.P()
	b.P(fmt.Sprintf("function encode(%s memory msg) internal pure returns (bytes memory) {", structName))
	b.Indent()

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

func isRepeated(label descriptorpb.FieldDescriptorProto_Label) bool {
	return label == descriptorpb.FieldDescriptorProto_LABEL_REPEATED
}