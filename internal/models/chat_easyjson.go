// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson9b8f5552DecodeYulaInternalModels(in *jlexer.Lexer, out *Message) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "info":
			(out.MI).UnmarshalEasyJSON(in)
		case "message":
			out.Msg = string(in.String())
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9b8f5552EncodeYulaInternalModels(out *jwriter.Writer, in Message) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"info\":"
		out.RawString(prefix[1:])
		(in.MI).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Msg))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9b8f5552EncodeYulaInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9b8f5552EncodeYulaInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9b8f5552DecodeYulaInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9b8f5552DecodeYulaInternalModels(l, v)
}
func easyjson9b8f5552DecodeYulaInternalModels1(in *jlexer.Lexer, out *IMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "from":
			out.IdFrom = int64(in.Int64())
		case "to":
			out.IdTo = int64(in.Int64())
		case "adv":
			out.IdAdv = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9b8f5552EncodeYulaInternalModels1(out *jwriter.Writer, in IMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"from\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.IdFrom))
	}
	{
		const prefix string = ",\"to\":"
		out.RawString(prefix)
		out.Int64(int64(in.IdTo))
	}
	{
		const prefix string = ",\"adv\":"
		out.RawString(prefix)
		out.Int64(int64(in.IdAdv))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v IMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9b8f5552EncodeYulaInternalModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v IMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9b8f5552EncodeYulaInternalModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *IMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9b8f5552DecodeYulaInternalModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *IMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9b8f5552DecodeYulaInternalModels1(l, v)
}
func easyjson9b8f5552DecodeYulaInternalModels2(in *jlexer.Lexer, out *IDialog) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user1":
			out.Id1 = int64(in.Int64())
		case "user2":
			out.Id2 = int64(in.Int64())
		case "adv":
			out.IdAdv = int64(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9b8f5552EncodeYulaInternalModels2(out *jwriter.Writer, in IDialog) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user1\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Id1))
	}
	{
		const prefix string = ",\"user2\":"
		out.RawString(prefix)
		out.Int64(int64(in.Id2))
	}
	{
		const prefix string = ",\"adv\":"
		out.RawString(prefix)
		out.Int64(int64(in.IdAdv))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v IDialog) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9b8f5552EncodeYulaInternalModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v IDialog) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9b8f5552EncodeYulaInternalModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *IDialog) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9b8f5552DecodeYulaInternalModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *IDialog) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9b8f5552DecodeYulaInternalModels2(l, v)
}
func easyjson9b8f5552DecodeYulaInternalModels3(in *jlexer.Lexer, out *Dialog) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "info":
			(out.DI).UnmarshalEasyJSON(in)
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9b8f5552EncodeYulaInternalModels3(out *jwriter.Writer, in Dialog) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"info\":"
		out.RawString(prefix[1:])
		(in.DI).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Dialog) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9b8f5552EncodeYulaInternalModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Dialog) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9b8f5552EncodeYulaInternalModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Dialog) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9b8f5552DecodeYulaInternalModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Dialog) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9b8f5552DecodeYulaInternalModels3(l, v)
}
