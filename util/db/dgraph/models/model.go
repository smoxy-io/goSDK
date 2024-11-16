package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgraph-io/dgo/v230/protos/api"
	goerrors "github.com/go-errors/errors"
	db "github.com/smoxy-io/goSDK/util/db/dgraph"
	str "github.com/smoxy-io/goSDK/util/strings"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"reflect"
	"strconv"
	"strings"
)

const (
	DgraphTypePredicate = "dgraph.type"

	ModelIdField = "id"
)

type Scalar interface {
	NQuadValue() string
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

const (
	// Example "increment" query:
	//
	// query {
	//   counter(func: has(counter.val)) {
	//     old: A as counter.val
	//     new: B as math(A+1)
	//   }
	// }
	// mutation @if(gt(len(A), 0)) {
	//   set {
	//     uid(A) <counter.val> val(B) .
	//   }
	// }
	incrementQuery = `
query Q($id: string!){
  counter(func: uid($id), first: 1) @filter(has(%model%.%field%)) {
    old: A as %model%.%field%
    new: B as math(A+%delta%)
  }
}
`

	incrementSetNQuads = `
uid(A) <%model%.%field%> val(B) .
`
)

type IdQueryExtractor[T proto.Message] func(resp *api.Response) (T, error)

type QueryParams[T any] struct {
	Query string
	Res   T
	Vars  map[string]string
}

func NewQueryParams[T any](query string, vars map[string]string, result T) *QueryParams[T] {
	return &QueryParams[T]{
		Query: query,
		Res:   result,
		Vars:  vars,
	}
}

// Get gets a model by ID using reflection. EXPERIMENTAL
func Get[T proto.Message](ctx context.Context, id string) (T, error) {
	var resp T

	if id == "" {
		return resp, errors.New("id is required")
	}

	dgraph, cErr, ctx := db.GetClient(ctx)

	if cErr != nil {
		return resp, cErr
	}

	txn := db.GetTxn(ctx)
	extTxn := txn != nil

	if !extTxn {
		txn = dgraph.NewReadOnlyTxn()
		defer txn.Discard(ctx)
	}

	vars := map[string]string{
		"$id": id,
	}

	qResp, qErr := txn.QueryWithVars(ctx, buildGetByIdQuery(resp), vars)

	if qErr != nil {
		return resp, qErr
	}

	q := struct {
		Q []T `json:"q"`
	}{}

	if err := json.Unmarshal(qResp.Json, &q); err != nil {
		return resp, err
	}

	if len(q.Q) < 1 {
		return resp, errors.New("not found")
	}

	return q.Q[0], nil
}

func buildGetByIdQuery[T proto.Message](m T) string {
	q := `
query Q($id: string!) {
  q(func: uid($id), first: 1) {
    %%FIELDS%%
  }
}
`

	fields := []string{}

	v := reflect.ValueOf(m)
	vk := v.Kind()

	if vk == reflect.Pointer || vk == reflect.Interface {
		// dereference the pointer or interface
		v = v.Elem()
		vk = v.Kind()
	}

	if vk == reflect.Struct {
		// convert each struct field into a dql field line
		buildQuery(v, vk, fields)
	}

	q = strings.ReplaceAll(q, "%%FIELDS%%", strings.Join(fields, "\n"))

	return q
}

func buildQuery(v reflect.Value, vk reflect.Kind, flds []string) {
	typeName := v.Type().Name()

	if vk == reflect.Pointer || vk == reflect.Interface {
		// dereference the pointer or interface
		v = v.Elem()
		vk = v.Kind()
	}

	if vk != reflect.Struct {
		return
	}

	for _, field := range reflect.VisibleFields(v.Type()) {
		if !field.IsExported() {
			// ignore unexported fields
			continue
		}

		fName := field.Name
		predName := str.FirstToLower(fName)

		vf := v.FieldByName(fName)
		vfk := vf.Kind()

		if vfk == reflect.Pointer || vfk == reflect.Interface {
			// dereference the pointer or interface
			// necessary for handling nullable fields as they are pointers in the model struct
			vf = vf.Elem()
			vfk = vf.Kind()
		}

		switch strings.ToLower(fName) {
		case "id", "uid":
			// this is the uid
			//   it's already in the query
			flds = append(flds, "id: uid")
			continue
		default:
			switch vfk {
			case reflect.Slice, reflect.Array:
				// this is a list of relationships
				// dql query is similar to single relationship
				vfR := vf.Elem()
				vkR := vfR.Kind()

				if vkR == reflect.Pointer || vkR == reflect.Interface {
					vfR = vfR.Elem()
					vkR = vfR.Kind()
				}

				if vkR != reflect.Struct {
					// no fields to iterate, treat as a standard field
					flds = append(flds, predName+": "+typeName+"."+predName)
					continue
				}

				f := []string{"id: uid"}

				buildQuery(vfR, vkR, f)

				flds = append(flds, predName+": "+typeName+"."+predName+" {\n"+strings.Join(f, "\n")+"\n}")

				continue
			case reflect.Struct:
				// this is a single relationship
				f := []string{"id: uid"}

				buildQuery(vf, vk, f)

				flds = append(flds, predName+": "+typeName+"."+predName+" {\n"+strings.Join(f, "\n")+"\n}")

				continue
			case reflect.Invalid:
				// invalid type. ignore it
				continue
			default:
				// add the field to the list
				flds = append(flds, predName+": "+typeName+"."+predName)
			}
		}
	}
}

// Set writes the model in the database
func Set[T proto.Message](ctx context.Context, m T) error {
	return upsert(ctx, m)
}

func SetBulk[T proto.Message](ctx context.Context, m []T) error {
	// TODO: implement
	return fmt.Errorf("not implemented")
}

// Query performs a query to retrieve models from the database
func Query[T any](ctx context.Context, params *QueryParams[T]) (T, error) {
	var nilT T

	dgraph, cErr, ctx := db.GetClient(ctx)

	if cErr != nil {
		return nilT, cErr
	}

	txn := db.GetTxn(ctx)
	extTxn := txn != nil

	if !extTxn {
		txn = dgraph.NewReadOnlyTxn()
		defer txn.Discard(ctx)
	}

	qResp, qErr := txn.QueryWithVars(ctx, params.Query, params.Vars)

	if qErr != nil {
		return nilT, qErr
	}

	if err := json.Unmarshal(qResp.Json, &params.Res); err != nil {
		return nilT, err
	}

	return params.Res, nil
}

// Increment performs an increment operation on a scalar predicate of a model
func Increment[T Number](ctx context.Context, model proto.Message, field string, delta T) error {
	modelName := string(model.ProtoReflect().Type().Descriptor().Name())

	modelId := ""

	model.ProtoReflect().Range(func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		//fmt.Printf("[DEBUG] model field: %s, value: %v\n", descriptor.Name(), value)
		if descriptor.Name() == ModelIdField {
			id, _ := value.Interface().(string)

			modelId = id

			return false
		}

		return true
	})

	//fmt.Printf("[DEBUG] Incrementing %s.%s (uid: %s)\n", modelName, field, modelId)

	if !db.ValidId(modelId) {
		return errors.New("model must have a valid id")
	}

	deltaStr := ""

	switch d := any(delta).(type) {
	case int, int8, int16, int32, int64:
		deltaStr = fmt.Sprintf("%d", d)
		break

	case uint, uint8, uint16, uint32, uint64:
		deltaStr = fmt.Sprintf("%d", d)
		break

	case float32, float64:
		deltaStr = fmt.Sprintf("%f", d)
		break
	}

	query := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(incrementQuery, "%model%", modelName), "%field%", field), "%delta%", deltaStr)

	queryVars := map[string]string{
		"$id": modelId,
	}

	nQuads := strings.ReplaceAll(strings.ReplaceAll(incrementSetNQuads, "%model%", modelName), "%field%", field)

	//fmt.Printf("[DEBUG] Increment query: %s\n", query)
	//fmt.Printf("[DEBUG] Increment nquads: %s\n", nQuads)

	mutations := make([]*api.Mutation, 0)

	mutations = append(mutations, &api.Mutation{
		SetNquads: []byte(nQuads),
		Cond:      "@if(gt(len(A), 0))",
	})

	_, rErr := rawQuery(ctx, query, queryVars, mutations)

	if rErr != nil {
		return rErr
	}

	return nil
}

func ToNQuads(i any, uidCount ...int) (string, string) {
	lines := []string{}
	uidCnt := 1

	if len(uidCount) > 0 && uidCount[0] > 1 {
		uidCnt = uidCount[0]
	}

	v := reflect.ValueOf(i)
	vk := v.Kind()

	if vk == reflect.Pointer || vk == reflect.Interface {
		// dereference the pointer or interface
		v = v.Elem()
		vk = v.Kind()
	}

	if vk != reflect.Struct {
		// cannot create nquads from this type
		return "", ""
	}

	typeName := v.Type().Name()

	uid := ""
	uidPlaceHolder := "$$uid$$"

	addRelationshipNquad := func(vf reflect.Value, predName string) {
		defer func() {
			if e := recover(); e != nil {
				goErr := goerrors.Wrap(e, 3)
				reset := string([]byte{27, 91, 48, 109})

				fmt.Printf("\n[Models ToNquad Add Relationship Recovery] panic recovered:\n\n%s\n\n%v\n\n%s\n\n%s%s\n", predName, vf, goErr.Error(), goErr.Stack(), reset)

				// continue to panic
				panic(e)
			}
		}()

		vfk := vf.Kind()

		if vfk == reflect.Pointer || vfk == reflect.Interface {
			// dereference the pointer or interface
			vf = vf.Elem()
			vfk = vf.Kind()
		}

		if vfk != reflect.Struct {
			// cannot create relationship from this type
			return
		}

		fVal, fvOk := vf.Addr().Interface().(Scalar)

		if fvOk {
			// this is a scalar value, not a relationship
			// - primarily used for in-built geo types like Point and Polygon
			lines = append(lines, fmt.Sprintf("%s <%s.%s> %s .", uidPlaceHolder, typeName, predName, fVal.NQuadValue()))

			return
		}

		vfId := vf.FieldByName("Id")

		if vfId.IsZero() {
			// try with uid
			vfId = vf.FieldByName("Uid")

			if vfId.IsNil() || vfId.IsZero() {
				// no id or uid fields. cannot create a relationship
				return
			}
		}

		if vfId.Kind() != reflect.String {
			// only string ids are supported
			return
		}

		lines = append(lines, fmt.Sprintf("%s <%s.%s> <%s> .", uidPlaceHolder, typeName, predName, vfId.String()))
	}

	for _, field := range reflect.VisibleFields(v.Type()) {
		if !field.IsExported() {
			// ignore unexported fields
			continue
		}

		fName := field.Name
		predName := str.FirstToLower(fName)
		tag := field.Tag.Get("nquad_type")

		vf := v.FieldByName(fName)
		vfk := vf.Kind()

		if vfk == reflect.Pointer || vfk == reflect.Interface {
			// dereference the pointer or interface
			// necessary for handling nullable fields as they are pointers in the model struct
			vf = vf.Elem()
			vfk = vf.Kind()
		}

		switch strings.ToLower(fName) {
		case "id", "uid":
			// this is the uid for the nquads
			//   it doesn't get its own nquad line, it will be the first item in all nquad lines
			if vfk != reflect.String {
				// uid must be a string
				continue
			}

			uid = vf.String()

			if uid == "" {
				// this is an insert, not an update
				// set a placeholder uid to be referenced by other nquad lines
				uid = fmt.Sprintf("_:%s%d", typeName, uidCnt)

				// make sure the inserted object has a type
				// add the type of the object
				lines = append(lines, fmt.Sprintf("%s <%s> \"%s\" .", uidPlaceHolder, DgraphTypePredicate, typeName))
			} else {
				uid = fmt.Sprintf("<%s>", uid)
			}
		default:
			switch vfk {
			case reflect.String:
				// it's possible that a protobuf string might actually be a uid
				//   the `nquad_type` struct field tag is used to indicate this case
				if tag == "uid" {
					id := vf.String()

					if id == "" {
						// no need to add anything for a relationship that is not specified
						continue
					}

					lines = append(lines, fmt.Sprintf("%s <%s.%s> <%s> .", uidPlaceHolder, typeName, predName, id))

					continue
				}

				lines = append(lines, fmt.Sprintf("%s <%s.%s> \"%s\"^^<xs:string> .", uidPlaceHolder, typeName, predName, vf.String()))

				continue
			case reflect.Slice, reflect.Array:
				// this is a list of relationships
				for i := 0; i < vf.Len(); i++ {
					vfR := vf.Index(i)
					addRelationshipNquad(vfR, predName)
				}

				continue
			case reflect.Struct:
				// this is a single relationship
				addRelationshipNquad(vf, predName)

				continue
			case reflect.Bool:
				fVal := "false"

				if vf.Bool() {
					fVal = "true"
				}

				lines = append(lines, fmt.Sprintf("%s <%s.%s> \"%s\"^^<xs:boolean> .", uidPlaceHolder, typeName, predName, fVal))

				continue
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fVal := strconv.Itoa(int(vf.Int()))

				lines = append(lines, fmt.Sprintf("%s <%s.%s> \"%s\"^^<xs:int> .", uidPlaceHolder, typeName, predName, fVal))

				continue
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fVal := strconv.FormatUint(vf.Uint(), 10)

				lines = append(lines, fmt.Sprintf("%s <%s.%s> \"%s\"^^<xs:int> .", uidPlaceHolder, typeName, predName, fVal))

				continue
			case reflect.Float32, reflect.Float64:
				fVal := strconv.FormatFloat(vf.Float(), 'f', -1, 64)

				lines = append(lines, fmt.Sprintf("%s <%s.%s> \"%s\"^^<xs:float> .", uidPlaceHolder, typeName, predName, fVal))

				continue
			case reflect.Invalid:
				// invalid type. ignore it
				continue
			default:
				// do nothing
			}

			if vf.CanConvert(reflect.TypeOf("s")) {
				fVal := vf.Convert(reflect.TypeOf("s")).String()

				lines = append(lines, fmt.Sprintf("%s <%s.%s> \"%s\"^^<xs:string> .", uidPlaceHolder, typeName, predName, fVal))
			}
		}
	}

	return strings.ReplaceAll(strings.Join(lines, "\n"), uidPlaceHolder, uid), uid
}

func rawQuery(ctx context.Context, query string, queryVars map[string]string, mutations []*api.Mutation) (*api.Response, error) {
	dgraph, cErr, ctx := db.GetClient(ctx)

	if cErr != nil {
		return nil, cErr
	}

	txn := db.GetTxn(ctx)
	extTxn := txn != nil

	if !extTxn {
		txn = dgraph.NewTxn()
		defer txn.Discard(ctx)
	}

	for i, _ := range mutations {
		mutations[i].CommitNow = !extTxn
	}

	req := &api.Request{
		Query:     query,
		Vars:      queryVars,
		Mutations: mutations,
		CommitNow: !extTxn,
	}

	return txn.Do(ctx, req)
}

func mutate(ctx context.Context, nquads string) (*api.Response, error) {
	dgraph, cErr, ctx := db.GetClient(ctx)

	if cErr != nil {
		return nil, cErr
	}

	txn := db.GetTxn(ctx)
	extTxn := txn != nil

	if !extTxn {
		txn = dgraph.NewTxn()
		defer txn.Discard(ctx)
	}

	mut := &api.Mutation{
		SetNquads: []byte(nquads),
		CommitNow: !extTxn,
	}

	return txn.Mutate(ctx, mut)
}

func upsert[T proto.Message](ctx context.Context, d T) (ret error) {
	defer func() {
		if e := recover(); e != nil {
			goErr := goerrors.Wrap(e, 3)
			reset := string([]byte{27, 91, 48, 109})

			fmt.Printf("\n[Upsert Recovery] panic recovered:\n\n%v\n\n%s\n\n%s%s\n", d, goErr.Error(), goErr.Stack(), reset)

			ret = goErr
		}
	}()

	nquads, uid := ToNQuads(d)

	resp, rErr := mutate(ctx, nquads)

	if rErr != nil {
		return rErr
	}

	if strings.HasPrefix(uid, "_:") {
		// a new uid should have been created. grab it and add it to the object
		if id, ok := resp.Uids[strings.TrimPrefix(uid, "_:")]; ok {
			v := reflect.ValueOf(d)
			vk := v.Kind()

			if vk == reflect.Pointer || vk == reflect.Interface {
				// dereference the pointer or interface
				v = v.Elem()
				vk = v.Kind()
			}

			vId := v.FieldByName("Id")

			if vId.CanSet() {
				vId.SetString(id)
			}
		}
	}

	return nil
}
