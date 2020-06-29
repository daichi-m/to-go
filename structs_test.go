package togo_test

/*
import (
	"reflect"
	"testing"
)

var fieldCloner func(f Field) Field
var gsCloner func(gs GoStruct) GoStruct

func init() {
	fieldCloner = func(f Field) Field {
		tmp, _ := f.Clone()
		return *(tmp.(*Field))
	}

	gsCloner = func(gs GoStruct) GoStruct {
		tmp, _ := gs.Clone()
		return *(tmp.(*GoStruct))
	}
}

func TestFieldDT_primitive(t *testing.T) {
	tests := []struct {
		tc   string
		f    FieldDT
		want bool
	}{
		{"Initial", 0, false},
		{"Bool", 1, true},
		{"Int", 2, true},
		{"Int64", 3, true},
		{"Float64", 4, true},
		{"String", 5, true},
		{"Slice", 6, false},
		{"Map", 7, false},
	}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			if got := tt.f.primitive(); got != tt.want {
				t.Errorf("FieldDT.primitive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldDT_str(t *testing.T) {
	tests := []struct {
		tc   string
		f    FieldDT
		want string
	}{
		{"Initial", 0, "Initial"},
		{"Bool", 1, "Bool"},
		{"Int", 2, "Int"},
		{"Int64", 3, "Int64"},
		{"Float64", 4, "Float64"},
		{"String", 5, "String"},
		{"Slice", 6, "Slice"},
		{"Map", 7, "Map"},
		{"Unknown", 10, "Unknown Type"},
	}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("FieldDT.str() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toFieldDT(t *testing.T) {
	tests := []struct {
		tc  string
		k   reflect.Kind
		res FieldDT
		ok  bool
	}{
		{"Bool", reflect.Bool, Bool, true},
		{"Int", reflect.Int, Int, true},
		{"Int64", reflect.Int64, Int64, true},
		{"Float32", reflect.Float32, Float64, true},
		{"Float64", reflect.Float64, Float64, true},
		{"String", reflect.String, String, true},
		{"Slice", reflect.Slice, Slice, true},
		{"Map", reflect.Map, Map, true},
		{"Unknown", reflect.Chan, Initial, false},
	}
	for _, tc := range tests {
		t.Run(tc.tc, func(t *testing.T) {
			res, ok := fieldDT(tc.k)
			if res != tc.res || ok != tc.ok {
				t.Errorf("TC: %s: Unexpected result. Expected (%v, %v) got (%v, %v)",
					tc.tc, tc.res, tc.ok, res, ok)
			}
		})
	}
}

func TestField_Equals(t *testing.T) {

	field := createField(Map, "FooType", 2)

	tests := []struct {
		tc     string
		field  Field
		equals bool
	}{
		{
			tc: "Name Not Equal",
			field: Field{
				"Random", field.annotation, field.dataType, field.dtStruct, field.sliceNesting,
			},
			equals: false,
		},
		{
			tc: "Type Not Equal",
			field: Field{
				field.name, field.annotation, Slice, field.dtStruct, field.sliceNesting,
			},
			equals: false,
		},
		{
			tc: "DTStruct Not Equal",
			field: Field{
				field.name, field.annotation, Map, "BarType", field.sliceNesting,
			},
			equals: false,
		},
		{
			tc: "Nesting Not Equal",
			field: Field{
				field.name, field.annotation, Map, field.dtStruct, 1,
			},
			equals: false,
		},
		{
			tc:     "Equal",
			field:  field,
			equals: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			got := field.Equals(&tt.field)
			if got != tt.equals {
				t.Errorf("TC: %s, Expacted equality to be %v but got %v", tt.tc, tt.equals, got)
			}
		})
	}
}

func TestField_Grow(t *testing.T) {
	field := createStringField()

	verify := func(f Field, a string) bool {
		return f.annotation == a
	}
	tests := []struct {
		tc               string
		fieldSupplier    func() Field
		annotation       string
		resultAnnotation string
	}{
		{
			tc: "New Annotation",
			fieldSupplier: func() Field {
				f := fieldCloner(field)
				f.SetAnnotation("")
				return f
			},
			annotation:       "json:" + field.name,
			resultAnnotation: "json:" + field.name,
		},
		{
			tc: "Exisiting Annotation",
			fieldSupplier: func() Field {
				return fieldCloner(field)
			},
			annotation:       "json:" + field.name,
			resultAnnotation: "json:" + field.name,
		},
		{
			tc: "Add to Annotation",
			fieldSupplier: func() Field {
				return fieldCloner(field)
			},
			annotation:       "omitempty",
			resultAnnotation: "json:" + field.name + ",omitempty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			f := tt.fieldSupplier()
			f2 := fieldCloner(field)
			f2.SetAnnotation(tt.annotation)
			f.Grow(&f2)
			if !verify(f, tt.resultAnnotation) {
				t.Errorf("TC: %s: Expected annotation to be %v but found %v",
					tt.tc, tt.resultAnnotation, f.annotation)
			}
		})
	}
}

func Test_MakeIField(t *testing.T) {

	verify := func(exp *Field, actual *Field) bool {
		if exp.name != actual.name ||
			exp.annotation != actual.annotation ||
			exp.dataType != actual.dataType {

			return false
		}
		if len(actual.dtStruct) != 0 || actual.sliceNesting != 0 {
			return false
		}
		return true
	}
	tests := []struct {
		tc     string
		val    interface{}
		name   string
		result Field
		resErr bool
	}{
		{
			tc:   "Int Field",
			val:  5,
			name: "IntField",
			result: Field{
				name:       "IntField",
				annotation: "",
				dataType:   Int,
			},
			resErr: false,
		},
		{
			tc:   "String Field",
			val:  "foo",
			name: "StrField",
			result: Field{
				name:       "StrField",
				annotation: "",
				dataType:   String,
			},
			resErr: false,
		},
		{
			tc:   "Map Field",
			val:  make(map[string]string),
			name: "MapField",
			result: Field{
				name:       "MapField",
				annotation: "",
				dataType:   Map,
			},
			resErr: false,
		},
		{
			tc:   "Slice Field",
			val:  make([]string, 10),
			name: "SlcField",
			result: Field{
				name:       "SlcField",
				annotation: "",
				dataType:   Slice,
			},
			resErr: false,
		},
		{
			tc:   "Channel Field",
			val:  make(chan string, 10),
			name: "ChanField",
			result: Field{
				name:       "ChanField",
				annotation: "",
				dataType:   Initial,
			},
			resErr: true,
		},
	}
	fm := FieldMaker{}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			fld, err := fm.MakeIField(tt.name, tt.val)
			if tt.resErr && err == nil {
				t.Errorf("TC: %s: Expected to get an error, but error is nil", tt.tc)
			}
			if fld != nil && !verify(&tt.result, fld.(*Field)) {
				t.Errorf("TC: %s: Verification failed, expcted %#v but received %#v",
					tt.tc, tt.result, fld)
			}
		})
	}
}

func TestGoStruct_AddField(t *testing.T) {

	gs := createSimpleGoStruct()
	field := createIntegerField()
	field.name = "Foo"
	gs.fields["Foo"] = &field

	tests := []struct {
		tc            string
		fieldSupplier func() *Field
		verify        func(gs *GoStruct) bool
		expError      bool
	}{
		{
			tc: "New Field",
			fieldSupplier: func() *Field {
				fld := fieldCloner(field)
				fld.name = "NewField"
				return &fld
			},
			verify: func(gs *GoStruct) bool {
				fld, ok := gs.fields["NewField"].(*Field)
				if !ok {
					return false
				}
				if fld.name != "NewField" ||
					fld.dataType != field.dataType ||
					fld.annotation != field.annotation {
					return false
				}
				return true
			},
			expError: false,
		},
		{
			tc: "Existing Field",
			fieldSupplier: func() *Field {
				fld, _ := gs.fields["Foo"].(*Field)
				return fld
			},
			verify: func(gs *GoStruct) bool {
				fld := gs.fields["Foo"].(*Field)
				if fld.name != field.name ||
					fld.annotation != field.annotation ||
					fld.dataType != field.dataType {
					return false
				}
				return true
			},
			expError: false,
		},
		{
			tc: "Existing Name Different Field",
			fieldSupplier: func() *Field {
				fld := createStringField()
				fld.name = "Foo"
				return &fld
			},
			verify: func(gs *GoStruct) bool {
				return false
			},
			expError: true,
		},
		{
			tc: "Nil Field",
			fieldSupplier: func() *Field {
				return nil
			},
			verify: func(gs *GoStruct) bool {
				return false
			},
			expError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			f := tt.fieldSupplier()
			_, err := gs.AddField(f)
			if tt.expError && err == nil {
				t.Errorf("TC: %s: Expected error, but did not see error", tt.tc)
			}
			if err == nil && !tt.verify(&gs) {
				t.Errorf("TC: %s: Verfication failed, GoStruct creted was %#v", tt.tc, gs)
			}
		})
	}
}

func TestGoStruct_Equals(t *testing.T) {

	gs := createSimpleGoStruct()
	tests := []struct {
		tc         string
		gsSupplier func() *GoStruct
		equal      bool
	}{
		{
			tc: "Equality",
			gsSupplier: func() *GoStruct {
				return &gs
			},
			equal: true,
		},
		{
			tc: "Not Equal",
			gsSupplier: func() *GoStruct {
				gsl := createSimpleGoStruct()
				if gs.name == gsl.name {
					gsl.name = "Random"
				}
				return &gsl
			},
			equal: false,
		},
		{
			tc: "Nil Struct",
			gsSupplier: func() *GoStruct {
				return nil
			},
			equal: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			eq := gs.Equals(tt.gsSupplier())
			if eq != tt.equal {
				t.Errorf("TC: %s: Expected equality %v but got %v instead", tt.tc, tt.equal, eq)
			}
		})
	}
}

func TestGoStruct_Grow(t *testing.T) {
	tests := []struct {
		tc            string
		gsSupplier    func() *GoStruct
		argGsSupplier func(gs GoStruct) *GoStruct
		verify        func(gs *GoStruct) bool
		expectErr     bool
	}{
		{
			tc: "GS not equal",
			gsSupplier: func() *GoStruct {
				gs := (createNamedGoStruct("FooStruct", "FooField", "BarField", "BazField"))
				return &gs
			},
			argGsSupplier: func(gs GoStruct) *GoStruct {
				argGs := (createNamedGoStruct("BarStruct", "DoField", "DarField", "DazField"))
				return &argGs
			},
			verify: func(gs *GoStruct) bool {
				return false
			},
			expectErr: true,
		},
		{
			tc: "No New Field",
			gsSupplier: func() *GoStruct {
				gs := (createNamedGoStruct("FooStruct", "FooField", "BarField", "BazField"))
				return &gs
			},
			argGsSupplier: func(gs GoStruct) *GoStruct {
				return &gs
			},
			verify: func(gs *GoStruct) bool {
				names := []string{"FooField", "BarField", "BazField"}
				for _, n := range names {
					fld, ok := gs.fields[n].(*Field)
					if !ok {
						return false
					}
					if fld.name != n {
						return false
					}
				}
				return true
			},
			expectErr: false,
		},
		{
			tc: "New Field",
			gsSupplier: func() *GoStruct {
				gs := (createNamedGoStruct("FooStruct", "FooField", "BarField", "BazField"))
				return &gs
			},
			argGsSupplier: func(gs GoStruct) *GoStruct {
				fld := createNamedField("QuxField", String, "", 0)
				ngs := gsCloner(gs)
				ngs.AddField(&fld)
				return &ngs
			},
			verify: func(gs *GoStruct) bool {
				names := []string{"FooField", "BarField", "BazField", "QuxField"}
				for _, n := range names {
					fld, ok := gs.fields[n].(*Field)
					if !ok {
						return false
					}
					if fld.name != n {
						return false
					}
				}
				return true
			},
			expectErr: false,
		},
		{
			tc: "Existing Field with different type",
			gsSupplier: func() *GoStruct {
				gs := (createNamedGoStruct("FooStruct", "FooField", "BarField", "BazField"))
				return &gs
			},
			argGsSupplier: func(gs GoStruct) *GoStruct {
				fld := createNamedField("FooField", String, "", 0)
				ngs := gsCloner(gs)
				delete(ngs.fields, "FooField")
				ngs.AddField(&fld)
				return &ngs
			},
			verify: func(gs *GoStruct) bool {
				return false
			},
			expectErr: true,
		},
		{
			tc: "Empty Struct",
			gsSupplier: func() *GoStruct {
				gs := new(GoStruct)
				gs.name = "FooStruct"
				return gs
			},
			argGsSupplier: func(gs GoStruct) *GoStruct {
				ngs := createNamedGoStruct("FooStruct", "FooField", "BarField")
				return &ngs
			},
			verify: func(gs *GoStruct) bool {
				names := []string{"FooField", "BarField"}
				for _, n := range names {
					fld, ok := gs.fields[n].(*Field)
					if !ok {
						return false
					}
					if fld.name != n {
						return false
					}
				}
				return true
			},
			expectErr: false,
		},
		{
			tc: "Empty Arg Struct",
			gsSupplier: func() *GoStruct {
				gs := createNamedGoStruct("FooStruct", "FooField", "BarField")
				return &gs
			},
			argGsSupplier: func(gs GoStruct) *GoStruct {
				ngs := new(GoStruct)
				ngs.name = "FooStruct"
				return ngs
			},
			verify: func(gs *GoStruct) bool {
				names := []string{"FooField", "BarField"}
				for _, n := range names {
					fld, ok := gs.fields[n].(*Field)
					if !ok {
						return false
					}
					if fld.name != n {
						return false
					}
				}
				return true
			},
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {
			gs := tt.gsSupplier()
			argGs := tt.argGsSupplier(*gs)
			_, err := gs.Grow(argGs)
			if tt.expectErr && err == nil {
				t.Errorf("TC: %s: Expected error but did not get any error", tt.tc)
			} else if !tt.expectErr && err != nil {
				t.Errorf("TC: %s: Did not expect error but got %#v", tt.tc, err)
			}

			if err == nil && !tt.verify(gs) {
				t.Errorf("TC: %s: Verification failed for GoStruct: %#v", tt.tc, gs)
			}
		})
	}
}

/*
func TestGoStruct_IsEmpty(t *testing.T) {
	tests := []struct {
		tc    string
		gs    GoStruct
		equal bool
	}{
		{"Non Empty", createSimpleGoStruct(), false},
		{"Empty", GoStruct{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.tc, func(t *testing.T) {

			if eq := tt.gs.IsEmpty(); eq != tt.equal {
				t.Errorf("TC: %s: Emptyness check failed. Expected: %v got %v",
					tt.tc, tt.equal, eq)
			}
		})
	}
}*/
