// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: google/protobuf/empty.proto

#include "google/protobuf/empty.pb.h"

#include <algorithm>
#include "google/protobuf/io/coded_stream.h"
#include "google/protobuf/extension_set.h"
#include "google/protobuf/wire_format_lite.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/generated_message_reflection.h"
#include "google/protobuf/reflection_ops.h"
#include "google/protobuf/wire_format.h"
#include "google/protobuf/generated_message_tctable_impl.h"
// @@protoc_insertion_point(includes)

// Must be included last.
#include "google/protobuf/port_def.inc"
PROTOBUF_PRAGMA_INIT_SEG
namespace _pb = ::google::protobuf;
namespace _pbi = ::google::protobuf::internal;
namespace _fl = ::google::protobuf::internal::field_layout;
namespace google {
namespace protobuf {
      template <typename>
PROTOBUF_CONSTEXPR Empty::Empty(::_pbi::ConstantInitialized) {}
struct EmptyDefaultTypeInternal {
  PROTOBUF_CONSTEXPR EmptyDefaultTypeInternal() : _instance(::_pbi::ConstantInitialized{}) {}
  ~EmptyDefaultTypeInternal() {}
  union {
    Empty _instance;
  };
};

PROTOBUF_ATTRIBUTE_NO_DESTROY PROTOBUF_CONSTINIT PROTOBUF_EXPORT
    PROTOBUF_ATTRIBUTE_INIT_PRIORITY1 EmptyDefaultTypeInternal _Empty_default_instance_;
}  // namespace protobuf
}  // namespace google
static ::_pb::Metadata file_level_metadata_google_2fprotobuf_2fempty_2eproto[1];
static constexpr const ::_pb::EnumDescriptor**
    file_level_enum_descriptors_google_2fprotobuf_2fempty_2eproto = nullptr;
static constexpr const ::_pb::ServiceDescriptor**
    file_level_service_descriptors_google_2fprotobuf_2fempty_2eproto = nullptr;
const ::uint32_t TableStruct_google_2fprotobuf_2fempty_2eproto::offsets[] PROTOBUF_SECTION_VARIABLE(
    protodesc_cold) = {
    ~0u,  // no _has_bits_
    PROTOBUF_FIELD_OFFSET(::google::protobuf::Empty, _internal_metadata_),
    ~0u,  // no _extensions_
    ~0u,  // no _oneof_case_
    ~0u,  // no _weak_field_map_
    ~0u,  // no _inlined_string_donated_
    ~0u,  // no _split_
    ~0u,  // no sizeof(Split)
};

static const ::_pbi::MigrationSchema
    schemas[] PROTOBUF_SECTION_VARIABLE(protodesc_cold) = {
        {0, -1, -1, sizeof(::google::protobuf::Empty)},
};

static const ::_pb::Message* const file_default_instances[] = {
    &::google::protobuf::_Empty_default_instance_._instance,
};
const char descriptor_table_protodef_google_2fprotobuf_2fempty_2eproto[] PROTOBUF_SECTION_VARIABLE(protodesc_cold) = {
    "\n\033google/protobuf/empty.proto\022\017google.pr"
    "otobuf\"\007\n\005EmptyB}\n\023com.google.protobufB\n"
    "EmptyProtoP\001Z.google.golang.org/protobuf"
    "/types/known/emptypb\370\001\001\242\002\003GPB\252\002\036Google.P"
    "rotobuf.WellKnownTypesb\006proto3"
};
static ::absl::once_flag descriptor_table_google_2fprotobuf_2fempty_2eproto_once;
const ::_pbi::DescriptorTable descriptor_table_google_2fprotobuf_2fempty_2eproto = {
    false,
    false,
    190,
    descriptor_table_protodef_google_2fprotobuf_2fempty_2eproto,
    "google/protobuf/empty.proto",
    &descriptor_table_google_2fprotobuf_2fempty_2eproto_once,
    nullptr,
    0,
    1,
    schemas,
    file_default_instances,
    TableStruct_google_2fprotobuf_2fempty_2eproto::offsets,
    file_level_metadata_google_2fprotobuf_2fempty_2eproto,
    file_level_enum_descriptors_google_2fprotobuf_2fempty_2eproto,
    file_level_service_descriptors_google_2fprotobuf_2fempty_2eproto,
};

// This function exists to be marked as weak.
// It can significantly speed up compilation by breaking up LLVM's SCC
// in the .pb.cc translation units. Large translation units see a
// reduction of more than 35% of walltime for optimized builds. Without
// the weak attribute all the messages in the file, including all the
// vtables and everything they use become part of the same SCC through
// a cycle like:
// GetMetadata -> descriptor table -> default instances ->
//   vtables -> GetMetadata
// By adding a weak function here we break the connection from the
// individual vtables back into the descriptor table.
PROTOBUF_ATTRIBUTE_WEAK const ::_pbi::DescriptorTable* descriptor_table_google_2fprotobuf_2fempty_2eproto_getter() {
  return &descriptor_table_google_2fprotobuf_2fempty_2eproto;
}
// Force running AddDescriptors() at dynamic initialization time.
PROTOBUF_ATTRIBUTE_INIT_PRIORITY2
static ::_pbi::AddDescriptorsRunner dynamic_init_dummy_google_2fprotobuf_2fempty_2eproto(&descriptor_table_google_2fprotobuf_2fempty_2eproto);
namespace google {
namespace protobuf {
// ===================================================================

class Empty::_Internal {
 public:
};

Empty::Empty(::google::protobuf::Arena* arena)
    : ::google::protobuf::internal::ZeroFieldsBase(arena) {
  // @@protoc_insertion_point(arena_constructor:google.protobuf.Empty)
}
Empty::Empty(const Empty& from) : ::google::protobuf::internal::ZeroFieldsBase() {
  Empty* const _this = this;
  (void)_this;
  _internal_metadata_.MergeFrom<::google::protobuf::UnknownFieldSet>(
      from._internal_metadata_);

  // @@protoc_insertion_point(copy_constructor:google.protobuf.Empty)
}




const ::google::protobuf::Message::ClassData Empty::_class_data_ = {
    ::google::protobuf::internal::ZeroFieldsBase::CopyImpl,
    ::google::protobuf::internal::ZeroFieldsBase::MergeImpl,
};
const ::google::protobuf::Message::ClassData*Empty::GetClassData() const { return &_class_data_; }







::google::protobuf::Metadata Empty::GetMetadata() const {
  return ::_pbi::AssignDescriptors(
      &descriptor_table_google_2fprotobuf_2fempty_2eproto_getter, &descriptor_table_google_2fprotobuf_2fempty_2eproto_once,
      file_level_metadata_google_2fprotobuf_2fempty_2eproto[0]);
}
// @@protoc_insertion_point(namespace_scope)
}  // namespace protobuf
}  // namespace google
namespace google {
namespace protobuf {
template<> PROTOBUF_NOINLINE ::google::protobuf::Empty*
Arena::CreateMaybeMessage< ::google::protobuf::Empty >(Arena* arena) {
  return Arena::CreateMessageInternal< ::google::protobuf::Empty >(arena);
}
}  // namespace protobuf
}  // namespace google
// @@protoc_insertion_point(global_scope)
#include "google/protobuf/port_undef.inc"