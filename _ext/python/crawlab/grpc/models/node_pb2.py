# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: models/node.proto
"""Generated protocol buffer code."""
from google.protobuf.internal import builder as _builder
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x11models/node.proto\x12\x04grpc\"\xcd\x01\n\x04Node\x12\x0b\n\x03_id\x18\x01 \x01(\t\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\n\n\x02ip\x18\x03 \x01(\t\x12\x0c\n\x04port\x18\x05 \x01(\t\x12\x0b\n\x03mac\x18\x06 \x01(\t\x12\x10\n\x08hostname\x18\x07 \x01(\t\x12\x13\n\x0b\x64\x65scription\x18\x08 \x01(\t\x12\x0b\n\x03key\x18\t \x01(\t\x12\x11\n\tis_master\x18\x0b \x01(\x08\x12\x11\n\tupdate_ts\x18\x0c \x01(\t\x12\x11\n\tcreate_ts\x18\r \x01(\t\x12\x16\n\x0eupdate_ts_unix\x18\x0e \x01(\x03\x42\x08Z\x06.;grpcb\x06proto3')

_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, globals())
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'models.node_pb2', globals())
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'Z\006.;grpc'
  _NODE._serialized_start=28
  _NODE._serialized_end=233
# @@protoc_insertion_point(module_scope)
