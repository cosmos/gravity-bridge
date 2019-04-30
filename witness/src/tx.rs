// This file is generated. Do not edit
// @generated

// https://github.com/Manishearth/rust-clippy/issues/702
#![allow(unknown_lints)]
#![allow(clippy)]

#![cfg_attr(rustfmt, rustfmt_skip)]

#![allow(box_pointers)]
#![allow(dead_code)]
#![allow(missing_docs)]
#![allow(non_camel_case_types)]
#![allow(non_snake_case)]
#![allow(non_upper_case_globals)]
#![allow(trivial_casts)]
#![allow(unsafe_code)]
#![allow(unused_imports)]
#![allow(unused_results)]

use protobuf::Message as Message_imported_for_functions;
use protobuf::ProtobufEnum as ProtobufEnum_imported_for_functions;

#[derive(PartialEq,Clone,Default)]
pub struct WitnessTx {
    // message fields
    pub signature: ::std::vec::Vec<u8>,
    pub sequence: i64,
    // message oneof groups
    tx: ::std::option::Option<WitnessTx_oneof_tx>,
    // special fields
    unknown_fields: ::protobuf::UnknownFields,
    cached_size: ::protobuf::CachedSize,
}

// see codegen.rs for the explanation why impl Sync explicitly
unsafe impl ::std::marker::Sync for WitnessTx {}

#[derive(Clone,PartialEq)]
pub enum WitnessTx_oneof_tx {
    lock(LockMsg),
}

impl WitnessTx {
    pub fn new() -> WitnessTx {
        ::std::default::Default::default()
    }

    pub fn default_instance() -> &'static WitnessTx {
        static mut instance: ::protobuf::lazy::Lazy<WitnessTx> = ::protobuf::lazy::Lazy {
            lock: ::protobuf::lazy::ONCE_INIT,
            ptr: 0 as *const WitnessTx,
        };
        unsafe {
            instance.get(WitnessTx::new)
        }
    }

    // .LockMsg lock = 1;

    pub fn clear_lock(&mut self) {
        self.tx = ::std::option::Option::None;
    }

    pub fn has_lock(&self) -> bool {
        match self.tx {
            ::std::option::Option::Some(WitnessTx_oneof_tx::lock(..)) => true,
            _ => false,
        }
    }

    // Param is passed by value, moved
    pub fn set_lock(&mut self, v: LockMsg) {
        self.tx = ::std::option::Option::Some(WitnessTx_oneof_tx::lock(v))
    }

    // Mutable pointer to the field.
    pub fn mut_lock(&mut self) -> &mut LockMsg {
        if let ::std::option::Option::Some(WitnessTx_oneof_tx::lock(_)) = self.tx {
        } else {
            self.tx = ::std::option::Option::Some(WitnessTx_oneof_tx::lock(LockMsg::new()));
        }
        match self.tx {
            ::std::option::Option::Some(WitnessTx_oneof_tx::lock(ref mut v)) => v,
            _ => panic!(),
        }
    }

    // Take field
    pub fn take_lock(&mut self) -> LockMsg {
        if self.has_lock() {
            match self.tx.take() {
                ::std::option::Option::Some(WitnessTx_oneof_tx::lock(v)) => v,
                _ => panic!(),
            }
        } else {
            LockMsg::new()
        }
    }

    pub fn get_lock(&self) -> &LockMsg {
        match self.tx {
            ::std::option::Option::Some(WitnessTx_oneof_tx::lock(ref v)) => v,
            _ => LockMsg::default_instance(),
        }
    }

    // bytes signature = 5;

    pub fn clear_signature(&mut self) {
        self.signature.clear();
    }

    // Param is passed by value, moved
    pub fn set_signature(&mut self, v: ::std::vec::Vec<u8>) {
        self.signature = v;
    }

    // Mutable pointer to the field.
    // If field is not initialized, it is initialized with default value first.
    pub fn mut_signature(&mut self) -> &mut ::std::vec::Vec<u8> {
        &mut self.signature
    }

    // Take field
    pub fn take_signature(&mut self) -> ::std::vec::Vec<u8> {
        ::std::mem::replace(&mut self.signature, ::std::vec::Vec::new())
    }

    pub fn get_signature(&self) -> &[u8] {
        &self.signature
    }

    fn get_signature_for_reflect(&self) -> &::std::vec::Vec<u8> {
        &self.signature
    }

    fn mut_signature_for_reflect(&mut self) -> &mut ::std::vec::Vec<u8> {
        &mut self.signature
    }

    // int64 sequence = 6;

    pub fn clear_sequence(&mut self) {
        self.sequence = 0;
    }

    // Param is passed by value, moved
    pub fn set_sequence(&mut self, v: i64) {
        self.sequence = v;
    }

    pub fn get_sequence(&self) -> i64 {
        self.sequence
    }

    fn get_sequence_for_reflect(&self) -> &i64 {
        &self.sequence
    }

    fn mut_sequence_for_reflect(&mut self) -> &mut i64 {
        &mut self.sequence
    }
}

impl ::protobuf::Message for WitnessTx {
    fn is_initialized(&self) -> bool {
        if let Some(WitnessTx_oneof_tx::lock(ref v)) = self.tx {
            if !v.is_initialized() {
                return false;
            }
        }
        true
    }

    fn merge_from(&mut self, is: &mut ::protobuf::CodedInputStream) -> ::protobuf::ProtobufResult<()> {
        while !is.eof()? {
            let (field_number, wire_type) = is.read_tag_unpack()?;
            match field_number {
                1 => {
                    if wire_type != ::protobuf::wire_format::WireTypeLengthDelimited {
                        return ::std::result::Result::Err(::protobuf::rt::unexpected_wire_type(wire_type));
                    }
                    self.tx = ::std::option::Option::Some(WitnessTx_oneof_tx::lock(is.read_message()?));
                },
                5 => {
                    ::protobuf::rt::read_singular_proto3_bytes_into(wire_type, is, &mut self.signature)?;
                },
                6 => {
                    if wire_type != ::protobuf::wire_format::WireTypeVarint {
                        return ::std::result::Result::Err(::protobuf::rt::unexpected_wire_type(wire_type));
                    }
                    let tmp = is.read_int64()?;
                    self.sequence = tmp;
                },
                _ => {
                    ::protobuf::rt::read_unknown_or_skip_group(field_number, wire_type, is, self.mut_unknown_fields())?;
                },
            };
        }
        ::std::result::Result::Ok(())
    }

    // Compute sizes of nested messages
    #[allow(unused_variables)]
    fn compute_size(&self) -> u32 {
        let mut my_size = 0;
        if !self.signature.is_empty() {
            my_size += ::protobuf::rt::bytes_size(5, &self.signature);
        }
        if self.sequence != 0 {
            my_size += ::protobuf::rt::value_size(6, self.sequence, ::protobuf::wire_format::WireTypeVarint);
        }
        if let ::std::option::Option::Some(ref v) = self.tx {
            match v {
                &WitnessTx_oneof_tx::lock(ref v) => {
                    let len = v.compute_size();
                    my_size += 1 + ::protobuf::rt::compute_raw_varint32_size(len) + len;
                },
            };
        }
        my_size += ::protobuf::rt::unknown_fields_size(self.get_unknown_fields());
        self.cached_size.set(my_size);
        my_size
    }

    fn write_to_with_cached_sizes(&self, os: &mut ::protobuf::CodedOutputStream) -> ::protobuf::ProtobufResult<()> {
        if !self.signature.is_empty() {
            os.write_bytes(5, &self.signature)?;
        }
        if self.sequence != 0 {
            os.write_int64(6, self.sequence)?;
        }
        if let ::std::option::Option::Some(ref v) = self.tx {
            match v {
                &WitnessTx_oneof_tx::lock(ref v) => {
                    os.write_tag(1, ::protobuf::wire_format::WireTypeLengthDelimited)?;
                    os.write_raw_varint32(v.get_cached_size())?;
                    v.write_to_with_cached_sizes(os)?;
                },
            };
        }
        os.write_unknown_fields(self.get_unknown_fields())?;
        ::std::result::Result::Ok(())
    }

    fn get_cached_size(&self) -> u32 {
        self.cached_size.get()
    }

    fn get_unknown_fields(&self) -> &::protobuf::UnknownFields {
        &self.unknown_fields
    }

    fn mut_unknown_fields(&mut self) -> &mut ::protobuf::UnknownFields {
        &mut self.unknown_fields
    }

    fn as_any(&self) -> &::std::any::Any {
        self as &::std::any::Any
    }
    fn as_any_mut(&mut self) -> &mut ::std::any::Any {
        self as &mut ::std::any::Any
    }
    fn into_any(self: Box<Self>) -> ::std::boxed::Box<::std::any::Any> {
        self
    }

    fn descriptor(&self) -> &'static ::protobuf::reflect::MessageDescriptor {
        ::protobuf::MessageStatic::descriptor_static(None::<Self>)
    }
}

impl ::protobuf::MessageStatic for WitnessTx {
    fn new() -> WitnessTx {
        WitnessTx::new()
    }

    fn descriptor_static(_: ::std::option::Option<WitnessTx>) -> &'static ::protobuf::reflect::MessageDescriptor {
        static mut descriptor: ::protobuf::lazy::Lazy<::protobuf::reflect::MessageDescriptor> = ::protobuf::lazy::Lazy {
            lock: ::protobuf::lazy::ONCE_INIT,
            ptr: 0 as *const ::protobuf::reflect::MessageDescriptor,
        };
        unsafe {
            descriptor.get(|| {
                let mut fields = ::std::vec::Vec::new();
                fields.push(::protobuf::reflect::accessor::make_singular_message_accessor::<_, LockMsg>(
                    "lock",
                    WitnessTx::has_lock,
                    WitnessTx::get_lock,
                ));
                fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeBytes>(
                    "signature",
                    WitnessTx::get_signature_for_reflect,
                    WitnessTx::mut_signature_for_reflect,
                ));
                fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeInt64>(
                    "sequence",
                    WitnessTx::get_sequence_for_reflect,
                    WitnessTx::mut_sequence_for_reflect,
                ));
                ::protobuf::reflect::MessageDescriptor::new::<WitnessTx>(
                    "WitnessTx",
                    fields,
                    file_descriptor_proto()
                )
            })
        }
    }
}

impl ::protobuf::Clear for WitnessTx {
    fn clear(&mut self) {
        self.clear_lock();
        self.clear_signature();
        self.clear_sequence();
        self.unknown_fields.clear();
    }
}

impl ::std::fmt::Debug for WitnessTx {
    fn fmt(&self, f: &mut ::std::fmt::Formatter) -> ::std::fmt::Result {
        ::protobuf::text_format::fmt(self, f)
    }
}

impl ::protobuf::reflect::ProtobufValue for WitnessTx {
    fn as_ref(&self) -> ::protobuf::reflect::ProtobufValueRef {
        ::protobuf::reflect::ProtobufValueRef::Message(self)
    }
}

#[derive(PartialEq,Clone,Default)]
pub struct LockMsg {
    // message fields
    pub dest: ::std::vec::Vec<u8>,
    pub value: u64,
    pub token: ::std::vec::Vec<u8>,
    pub nonce: u64,
    // special fields
    unknown_fields: ::protobuf::UnknownFields,
    cached_size: ::protobuf::CachedSize,
}

// see codegen.rs for the explanation why impl Sync explicitly
unsafe impl ::std::marker::Sync for LockMsg {}

impl LockMsg {
    pub fn new() -> LockMsg {
        ::std::default::Default::default()
    }

    pub fn default_instance() -> &'static LockMsg {
        static mut instance: ::protobuf::lazy::Lazy<LockMsg> = ::protobuf::lazy::Lazy {
            lock: ::protobuf::lazy::ONCE_INIT,
            ptr: 0 as *const LockMsg,
        };
        unsafe {
            instance.get(LockMsg::new)
        }
    }

    // bytes dest = 1;

    pub fn clear_dest(&mut self) {
        self.dest.clear();
    }

    // Param is passed by value, moved
    pub fn set_dest(&mut self, v: ::std::vec::Vec<u8>) {
        self.dest = v;
    }

    // Mutable pointer to the field.
    // If field is not initialized, it is initialized with default value first.
    pub fn mut_dest(&mut self) -> &mut ::std::vec::Vec<u8> {
        &mut self.dest
    }

    // Take field
    pub fn take_dest(&mut self) -> ::std::vec::Vec<u8> {
        ::std::mem::replace(&mut self.dest, ::std::vec::Vec::new())
    }

    pub fn get_dest(&self) -> &[u8] {
        &self.dest
    }

    fn get_dest_for_reflect(&self) -> &::std::vec::Vec<u8> {
        &self.dest
    }

    fn mut_dest_for_reflect(&mut self) -> &mut ::std::vec::Vec<u8> {
        &mut self.dest
    }

    // uint64 value = 2;

    pub fn clear_value(&mut self) {
        self.value = 0;
    }

    // Param is passed by value, moved
    pub fn set_value(&mut self, v: u64) {
        self.value = v;
    }

    pub fn get_value(&self) -> u64 {
        self.value
    }

    fn get_value_for_reflect(&self) -> &u64 {
        &self.value
    }

    fn mut_value_for_reflect(&mut self) -> &mut u64 {
        &mut self.value
    }

    // bytes token = 3;

    pub fn clear_token(&mut self) {
        self.token.clear();
    }

    // Param is passed by value, moved
    pub fn set_token(&mut self, v: ::std::vec::Vec<u8>) {
        self.token = v;
    }

    // Mutable pointer to the field.
    // If field is not initialized, it is initialized with default value first.
    pub fn mut_token(&mut self) -> &mut ::std::vec::Vec<u8> {
        &mut self.token
    }

    // Take field
    pub fn take_token(&mut self) -> ::std::vec::Vec<u8> {
        ::std::mem::replace(&mut self.token, ::std::vec::Vec::new())
    }

    pub fn get_token(&self) -> &[u8] {
        &self.token
    }

    fn get_token_for_reflect(&self) -> &::std::vec::Vec<u8> {
        &self.token
    }

    fn mut_token_for_reflect(&mut self) -> &mut ::std::vec::Vec<u8> {
        &mut self.token
    }

    // uint64 nonce = 5;

    pub fn clear_nonce(&mut self) {
        self.nonce = 0;
    }

    // Param is passed by value, moved
    pub fn set_nonce(&mut self, v: u64) {
        self.nonce = v;
    }

    pub fn get_nonce(&self) -> u64 {
        self.nonce
    }

    fn get_nonce_for_reflect(&self) -> &u64 {
        &self.nonce
    }

    fn mut_nonce_for_reflect(&mut self) -> &mut u64 {
        &mut self.nonce
    }
}

impl ::protobuf::Message for LockMsg {
    fn is_initialized(&self) -> bool {
        true
    }

    fn merge_from(&mut self, is: &mut ::protobuf::CodedInputStream) -> ::protobuf::ProtobufResult<()> {
        while !is.eof()? {
            let (field_number, wire_type) = is.read_tag_unpack()?;
            match field_number {
                1 => {
                    ::protobuf::rt::read_singular_proto3_bytes_into(wire_type, is, &mut self.dest)?;
                },
                2 => {
                    if wire_type != ::protobuf::wire_format::WireTypeVarint {
                        return ::std::result::Result::Err(::protobuf::rt::unexpected_wire_type(wire_type));
                    }
                    let tmp = is.read_uint64()?;
                    self.value = tmp;
                },
                3 => {
                    ::protobuf::rt::read_singular_proto3_bytes_into(wire_type, is, &mut self.token)?;
                },
                5 => {
                    if wire_type != ::protobuf::wire_format::WireTypeVarint {
                        return ::std::result::Result::Err(::protobuf::rt::unexpected_wire_type(wire_type));
                    }
                    let tmp = is.read_uint64()?;
                    self.nonce = tmp;
                },
                _ => {
                    ::protobuf::rt::read_unknown_or_skip_group(field_number, wire_type, is, self.mut_unknown_fields())?;
                },
            };
        }
        ::std::result::Result::Ok(())
    }

    // Compute sizes of nested messages
    #[allow(unused_variables)]
    fn compute_size(&self) -> u32 {
        let mut my_size = 0;
        if !self.dest.is_empty() {
            my_size += ::protobuf::rt::bytes_size(1, &self.dest);
        }
        if self.value != 0 {
            my_size += ::protobuf::rt::value_size(2, self.value, ::protobuf::wire_format::WireTypeVarint);
        }
        if !self.token.is_empty() {
            my_size += ::protobuf::rt::bytes_size(3, &self.token);
        }
        if self.nonce != 0 {
            my_size += ::protobuf::rt::value_size(5, self.nonce, ::protobuf::wire_format::WireTypeVarint);
        }
        my_size += ::protobuf::rt::unknown_fields_size(self.get_unknown_fields());
        self.cached_size.set(my_size);
        my_size
    }

    fn write_to_with_cached_sizes(&self, os: &mut ::protobuf::CodedOutputStream) -> ::protobuf::ProtobufResult<()> {
        if !self.dest.is_empty() {
            os.write_bytes(1, &self.dest)?;
        }
        if self.value != 0 {
            os.write_uint64(2, self.value)?;
        }
        if !self.token.is_empty() {
            os.write_bytes(3, &self.token)?;
        }
        if self.nonce != 0 {
            os.write_uint64(5, self.nonce)?;
        }
        os.write_unknown_fields(self.get_unknown_fields())?;
        ::std::result::Result::Ok(())
    }

    fn get_cached_size(&self) -> u32 {
        self.cached_size.get()
    }

    fn get_unknown_fields(&self) -> &::protobuf::UnknownFields {
        &self.unknown_fields
    }

    fn mut_unknown_fields(&mut self) -> &mut ::protobuf::UnknownFields {
        &mut self.unknown_fields
    }

    fn as_any(&self) -> &::std::any::Any {
        self as &::std::any::Any
    }
    fn as_any_mut(&mut self) -> &mut ::std::any::Any {
        self as &mut ::std::any::Any
    }
    fn into_any(self: Box<Self>) -> ::std::boxed::Box<::std::any::Any> {
        self
    }

    fn descriptor(&self) -> &'static ::protobuf::reflect::MessageDescriptor {
        ::protobuf::MessageStatic::descriptor_static(None::<Self>)
    }
}

impl ::protobuf::MessageStatic for LockMsg {
    fn new() -> LockMsg {
        LockMsg::new()
    }

    fn descriptor_static(_: ::std::option::Option<LockMsg>) -> &'static ::protobuf::reflect::MessageDescriptor {
        static mut descriptor: ::protobuf::lazy::Lazy<::protobuf::reflect::MessageDescriptor> = ::protobuf::lazy::Lazy {
            lock: ::protobuf::lazy::ONCE_INIT,
            ptr: 0 as *const ::protobuf::reflect::MessageDescriptor,
        };
        unsafe {
            descriptor.get(|| {
                let mut fields = ::std::vec::Vec::new();
                fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeBytes>(
                    "dest",
                    LockMsg::get_dest_for_reflect,
                    LockMsg::mut_dest_for_reflect,
                ));
                fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeUint64>(
                    "value",
                    LockMsg::get_value_for_reflect,
                    LockMsg::mut_value_for_reflect,
                ));
                fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeBytes>(
                    "token",
                    LockMsg::get_token_for_reflect,
                    LockMsg::mut_token_for_reflect,
                ));
                fields.push(::protobuf::reflect::accessor::make_simple_field_accessor::<_, ::protobuf::types::ProtobufTypeUint64>(
                    "nonce",
                    LockMsg::get_nonce_for_reflect,
                    LockMsg::mut_nonce_for_reflect,
                ));
                ::protobuf::reflect::MessageDescriptor::new::<LockMsg>(
                    "LockMsg",
                    fields,
                    file_descriptor_proto()
                )
            })
        }
    }
}

impl ::protobuf::Clear for LockMsg {
    fn clear(&mut self) {
        self.clear_dest();
        self.clear_value();
        self.clear_token();
        self.clear_nonce();
        self.unknown_fields.clear();
    }
}

impl ::std::fmt::Debug for LockMsg {
    fn fmt(&self, f: &mut ::std::fmt::Formatter) -> ::std::fmt::Result {
        ::protobuf::text_format::fmt(self, f)
    }
}

impl ::protobuf::reflect::ProtobufValue for LockMsg {
    fn as_ref(&self) -> ::protobuf::reflect::ProtobufValueRef {
        ::protobuf::reflect::ProtobufValueRef::Message(self)
    }
}

static file_descriptor_proto_data: &'static [u8] = b"\
    \n\x08tx.proto\"k\n\tWitnessTx\x12\x1e\n\x04lock\x18\x01\x20\x01(\x0b2\
    \x08.LockMsgH\0R\x04lock\x12\x1c\n\tsignature\x18\x05\x20\x01(\x0cR\tsig\
    nature\x12\x1a\n\x08sequence\x18\x06\x20\x01(\x03R\x08sequenceB\x04\n\
    \x02tx\"_\n\x07LockMsg\x12\x12\n\x04dest\x18\x01\x20\x01(\x0cR\x04dest\
    \x12\x14\n\x05value\x18\x02\x20\x01(\x04R\x05value\x12\x14\n\x05token\
    \x18\x03\x20\x01(\x0cR\x05token\x12\x14\n\x05nonce\x18\x05\x20\x01(\x04R\
    \x05nonceJ\xca\x04\n\x06\x12\x04\0\0\x10\x01\n\x08\n\x01\x0c\x12\x03\0\0\
    \x12\n\n\n\x02\x04\0\x12\x04\x02\0\x08\x01\n\n\n\x03\x04\0\x01\x12\x03\
    \x02\x08\x11\n\x0c\n\x04\x04\0\x08\0\x12\x04\x03\x04\x05\x05\n\x0c\n\x05\
    \x04\0\x08\0\x01\x12\x03\x03\n\x0c\n\x0b\n\x04\x04\0\x02\0\x12\x03\x04\
    \x08\x19\n\x0c\n\x05\x04\0\x02\0\x06\x12\x03\x04\x08\x0f\n\x0c\n\x05\x04\
    \0\x02\0\x01\x12\x03\x04\x10\x14\n\x0c\n\x05\x04\0\x02\0\x03\x12\x03\x04\
    \x17\x18\n\x0b\n\x04\x04\0\x02\x01\x12\x03\x06\x04\x18\n\r\n\x05\x04\0\
    \x02\x01\x04\x12\x04\x06\x04\x05\x05\n\x0c\n\x05\x04\0\x02\x01\x05\x12\
    \x03\x06\x04\t\n\x0c\n\x05\x04\0\x02\x01\x01\x12\x03\x06\n\x13\n\x0c\n\
    \x05\x04\0\x02\x01\x03\x12\x03\x06\x16\x17\n\x0b\n\x04\x04\0\x02\x02\x12\
    \x03\x07\x04\x17\n\r\n\x05\x04\0\x02\x02\x04\x12\x04\x07\x04\x06\x18\n\
    \x0c\n\x05\x04\0\x02\x02\x05\x12\x03\x07\x04\t\n\x0c\n\x05\x04\0\x02\x02\
    \x01\x12\x03\x07\n\x12\n\x0c\n\x05\x04\0\x02\x02\x03\x12\x03\x07\x15\x16\
    \n\n\n\x02\x04\x01\x12\x04\n\0\x10\x01\n\n\n\x03\x04\x01\x01\x12\x03\n\
    \x08\x0f\n\x0b\n\x04\x04\x01\x02\0\x12\x03\x0b\x04\x13\n\r\n\x05\x04\x01\
    \x02\0\x04\x12\x04\x0b\x04\n\x11\n\x0c\n\x05\x04\x01\x02\0\x05\x12\x03\
    \x0b\x04\t\n\x0c\n\x05\x04\x01\x02\0\x01\x12\x03\x0b\n\x0e\n\x0c\n\x05\
    \x04\x01\x02\0\x03\x12\x03\x0b\x11\x12\n\x0b\n\x04\x04\x01\x02\x01\x12\
    \x03\x0c\x04\x15\n\r\n\x05\x04\x01\x02\x01\x04\x12\x04\x0c\x04\x0b\x13\n\
    \x0c\n\x05\x04\x01\x02\x01\x05\x12\x03\x0c\x04\n\n\x0c\n\x05\x04\x01\x02\
    \x01\x01\x12\x03\x0c\x0b\x10\n\x0c\n\x05\x04\x01\x02\x01\x03\x12\x03\x0c\
    \x13\x14\n\x0b\n\x04\x04\x01\x02\x02\x12\x03\r\x04\x14\n\r\n\x05\x04\x01\
    \x02\x02\x04\x12\x04\r\x04\x0c\x15\n\x0c\n\x05\x04\x01\x02\x02\x05\x12\
    \x03\r\x04\t\n\x0c\n\x05\x04\x01\x02\x02\x01\x12\x03\r\n\x0f\n\x0c\n\x05\
    \x04\x01\x02\x02\x03\x12\x03\r\x12\x13\n\x1c\n\x04\x04\x01\x02\x03\x12\
    \x03\x0f\x04\x15\x1a\x0f\x20\x20bytes\x20chain;\n\n\r\n\x05\x04\x01\x02\
    \x03\x04\x12\x04\x0f\x04\r\x14\n\x0c\n\x05\x04\x01\x02\x03\x05\x12\x03\
    \x0f\x04\n\n\x0c\n\x05\x04\x01\x02\x03\x01\x12\x03\x0f\x0b\x10\n\x0c\n\
    \x05\x04\x01\x02\x03\x03\x12\x03\x0f\x13\x14b\x06proto3\
";

static mut file_descriptor_proto_lazy: ::protobuf::lazy::Lazy<::protobuf::descriptor::FileDescriptorProto> = ::protobuf::lazy::Lazy {
    lock: ::protobuf::lazy::ONCE_INIT,
    ptr: 0 as *const ::protobuf::descriptor::FileDescriptorProto,
};

fn parse_descriptor_proto() -> ::protobuf::descriptor::FileDescriptorProto {
    ::protobuf::parse_from_bytes(file_descriptor_proto_data).unwrap()
}

pub fn file_descriptor_proto() -> &'static ::protobuf::descriptor::FileDescriptorProto {
    unsafe {
        file_descriptor_proto_lazy.get(|| {
            parse_descriptor_proto()
        })
    }
}
