// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package backup

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"

	pb "github.com/sromku/datastore-to-sql/backup/pb"
)

var (
	// ErrInvalidEntityType is returned when functions like Get or Next are
	// passed a dst or src argument of invalid type.
	ErrInvalidEntityType = errors.New("datastore: invalid entity type")
	// ErrInvalidKey is returned when an invalid key is presented.
	ErrInvalidKey = errors.New("datastore: invalid key")
	// ErrNoSuchEntity is returned when no entity was found for a given key.
	ErrNoSuchEntity = errors.New("datastore: no such entity")
)

// ErrFieldMismatch is returned when a field is to be loaded into a different
// type than the one it was stored from, or when a field is missing or
// unexported in the destination struct.
// StructType is the type of the struct pointed to by the destination argument
// passed to Get or to Iterator.Next.
type ErrFieldMismatch struct {
	StructType reflect.Type
	FieldName  string
	Reason     string
}

func (e *ErrFieldMismatch) Error() string {
	return fmt.Sprintf("datastore: cannot load field %q into a %q: %s",
		e.FieldName, e.StructType, e.Reason)
}

// protoToKey converts a Reference proto to a *Key.
func protoToKey(r *pb.Reference) (k *Key, err error) {
	appID := r.GetApp()
	namespace := r.GetNameSpace()
	for _, e := range r.Path.Element {
		k = &Key{
			kind:      e.GetType(),
			stringID:  e.GetName(),
			intID:     e.GetId(),
			parent:    k,
			appID:     appID,
			namespace: namespace,
		}
		if !k.valid() {
			return nil, ErrInvalidKey
		}
	}
	return
}

// keyToProto converts a *Key to a Reference proto.
func keyToProto(defaultAppID string, k *Key) *pb.Reference {
	appID := k.appID
	if appID == "" {
		appID = defaultAppID
	}
	n := 0
	for i := k; i != nil; i = i.parent {
		n++
	}
	e := make([]*pb.Path_Element, n)
	for i := k; i != nil; i = i.parent {
		n--
		e[n] = &pb.Path_Element{
			Type: &i.kind,
		}
		// At most one of {Name,Id} should be set.
		// Neither will be set for incomplete keys.
		if i.stringID != "" {
			e[n].Name = &i.stringID
		} else if i.intID != 0 {
			e[n].Id = &i.intID
		}
	}
	var namespace *string
	if k.namespace != "" {
		namespace = proto.String(k.namespace)
	}
	return &pb.Reference{
		App:       proto.String(appID),
		NameSpace: namespace,
		Path: &pb.Path{
			Element: e,
		},
	}
}

// It's unfortunate that the two semantically equivalent concepts pb.Reference
// and pb.PropertyValue_ReferenceValue aren't the same type. For example, the
// two have different protobuf field numbers.

// referenceValueToKey is the same as protoToKey except the input is a
// PropertyValue_ReferenceValue instead of a Reference.
func referenceValueToKey(r *pb.PropertyValue_ReferenceValue) (k *Key, err error) {
	appID := r.GetApp()
	namespace := r.GetNameSpace()
	for _, e := range r.Pathelement {
		k = &Key{
			kind:      e.GetType(),
			stringID:  e.GetName(),
			intID:     e.GetId(),
			parent:    k,
			appID:     appID,
			namespace: namespace,
		}
		if !k.valid() {
			return nil, ErrInvalidKey
		}
	}
	return
}

// keyToReferenceValue is the same as keyToProto except the output is a
// PropertyValue_ReferenceValue instead of a Reference.
func keyToReferenceValue(defaultAppID string, k *Key) *pb.PropertyValue_ReferenceValue {
	ref := keyToProto(defaultAppID, k)
	pe := make([]*pb.PropertyValue_ReferenceValue_PathElement, len(ref.Path.Element))
	for i, e := range ref.Path.Element {
		pe[i] = &pb.PropertyValue_ReferenceValue_PathElement{
			Type: e.Type,
			Id:   e.Id,
			Name: e.Name,
		}
	}
	return &pb.PropertyValue_ReferenceValue{
		App:         ref.App,
		NameSpace:   ref.NameSpace,
		Pathelement: pe,
	}
}

type multiArgType int

const (
	multiArgTypeInvalid multiArgType = iota
	multiArgTypePropertyLoadSaver
	multiArgTypeStruct
	multiArgTypeStructPtr
	multiArgTypeInterface
)
