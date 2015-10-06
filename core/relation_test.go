package core

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"testing"
)

var testStaffSchema = Schema{Attr{"name", reflect.String}, Attr{"age", reflect.Int64}, Attr{"job", reflect.String}}
var testStaff = &Relation{
	Name:  "staff",
	Attrs: testStaffSchema,
	Data: []Tuple{
		Tuple{Attrs: testStaffSchema, Data: map[string]Value{"name": "清水", "age": int64(17), "job": "エンジニア"}},
		Tuple{Attrs: testStaffSchema, Data: map[string]Value{"name": "田中", "age": int64(34), "job": "デザイナー"}},
		Tuple{Attrs: testStaffSchema, Data: map[string]Value{"name": "佐藤", "age": int64(21), "job": "マネージャー"}},
	},
}
var testStaff3 *Relation
var testRank3 *Relation

var testRankSchema = Schema{Attr{"name", reflect.String}, Attr{"rank", reflect.Int64}}
var testRank = &Relation{
	Name:  "rank",
	Attrs: testRankSchema,
	Data: []Tuple{
		Tuple{Attrs: testRankSchema, Data: map[string]Value{"name": "清水", "rank": int64(78)}},
		Tuple{Attrs: testRankSchema, Data: map[string]Value{"name": "田中", "rank": int64(46)}},
		Tuple{Attrs: testRankSchema, Data: map[string]Value{"name": "佐藤", "rank": int64(33)}},
	},
}

var testData1Schema = Schema{Attr{"id", reflect.String}, Attr{"job", reflect.String}}
var testRelation1 = &Relation{
	Name:  "testData1",
	Attrs: testData1Schema,
	Data: []Tuple{
		Tuple{Attrs: testData1Schema, Data: map[string]Value{"id": 1, "job": "abc"}},
		Tuple{Attrs: testData1Schema, Data: map[string]Value{"id": 2, "job": "def"}},
		Tuple{Attrs: testData1Schema, Data: map[string]Value{"id": 3, "job": "ghi"}},
		Tuple{Attrs: testData1Schema, Data: map[string]Value{"id": 4, "job": "jkl"}},
		Tuple{Attrs: testData1Schema, Data: map[string]Value{"id": 5, "job": "mno"}},
		Tuple{Attrs: testData1Schema, Data: map[string]Value{"id": 6, "job": "pqr"}},
	},
}

var testData2Schema = Schema{Attr{"id", reflect.String}, Attr{"job", reflect.String}}
var testRelation2 = &Relation{
	Name:  "testData2",
	Attrs: testData2Schema,
	Data: []Tuple{
		Tuple{Attrs: testData2Schema, Data: map[string]Value{"id": 1, "job": "zxy"}},
		Tuple{Attrs: testData2Schema, Data: map[string]Value{"id": 2, "job": "zxy"}},
		Tuple{Attrs: testData2Schema, Data: map[string]Value{"id": 3, "job": "cjk"}},
		Tuple{Attrs: testData2Schema, Data: map[string]Value{"id": 4, "job": "hij"}},
		Tuple{Attrs: testData2Schema, Data: map[string]Value{"id": 5, "job": "cjk"}},
		Tuple{Attrs: testData2Schema, Data: map[string]Value{"id": 6, "job": "zxy"}},
	},
}

var testData1 = &Node{
	Name:     "root",
	FullPath: Relations{},
	Nodes: Nodes{
		"test": &Node{
			FullPath: Relations{},
			Name:     "test",
			Relations: Relations{
				"staff1": testStaff,
				"rank1":  testRank,
			},
		},
		"20150926": &Node{
			FullPath: Relations{},
			Nodes: Nodes{
				"data": &Node{
					FullPath: Relations{},
					Name:     "data",
					Relations: Relations{
						"staff2": testStaff,
						"rank2":  testRank,
					},
				},
			},
			Relations: Relations{
				"staff3": testStaff,
				"rank3":  testRank,
			},
		},
	},
}

var testData2 = NewNode("root")

func init() {
	err := testData2.SetRelations("test", Relations{"staff1": testStaff, "rank1": testRank})
	if err != nil {
		log.Fatalln(err)
	}
	err = testData2.SetRelations("20150926/data", Relations{"staff2": testStaff, "rank2": testRank})
	if err != nil {
		log.Fatalln(err)
	}
	testRank3 = testRank.Clone()
	testRank3.Name = "rank3"
	testStaff3 = testStaff.Clone()
	testStaff3.Name = "staff3"
	err = testData2.SetRelation("20150926", testRank3)
	if err != nil {
		log.Fatalln(err)
	}
	err = testData2.SetRelation("20150926", testStaff3)
	if err != nil {
		log.Fatalln(err)
	}
}

func TestGetRelation1(t *testing.T) {

	r, err := testData1.GetRelation("/test/rank1")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(r, testRank) {
		t.Errorf("Does not match %# v, want:%# v", r, testRank)
	}
	r, err = testData2.GetRelation("/test/rank1")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(r, testRank) {
		t.Errorf("Does not match %# v, want:%# v", r, testRank)
	}
}

func TestGetRelation2(t *testing.T) {

	r, err := testData1.GetRelation("20150926/data/staff2")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(r, testStaff) {
		t.Errorf("Does not match %# v, want:%# v", r, testStaff)
	}
	r, err = testData2.GetRelation("20150926/data/staff2")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(r, testStaff) {
		t.Errorf("Does not match %# v, want:%# v", r, testStaff)
	}
}

func TestGetRelation3(t *testing.T) {

	r, err := testData2.GetRelation("20150926/staff3")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(r, testStaff3) {
		t.Errorf("Does not match %# v, want:%# v", r, testStaff3)
	}
}

func TestJsonSelectionStream(t *testing.T) {
	schema := Schema{Attr{"name", reflect.String}, Attr{"age", reflect.Int64}, Attr{"job", reflect.String}}
	var want = &Relation{
		Attrs: schema,
		Data: []Tuple{
			Tuple{Attrs: schema, Data: map[string]Value{"name": "田中", "age": int64(34), "job": "デザイナー"}},
			Tuple{Attrs: schema, Data: map[string]Value{"name": "佐藤", "age": int64(21), "job": "マネージャー"}},
		},
	}
	const jsonStream = `{ "selection": {
			"input": { 
				"relation": {"name":"test/staff1"}},
				"attr": "age", "selector": ">", "arg": 20
	}}`
	m := Stream{}
	if err := json.NewDecoder(strings.NewReader(jsonStream)).Decode(&m); err != nil {
		log.Fatal(err)
	}
	result, _ := StreamToRelation(m, testData2)
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Does not match 'SELECT * FROM Staff WHERE age > 20'\nresult:% #v,\n want:% #v", result, want)
	}
}

func TestJsonProjectionStream(t *testing.T) {
	schema := Schema{Attr{"age", reflect.Int64}, Attr{"job", reflect.String}}
	var want = &Relation{
		Attrs: schema,
		Data: []Tuple{
			Tuple{Attrs: schema, Data: map[string]Value{"age": int64(17), "job": "エンジニア"}},
			Tuple{Attrs: schema, Data: map[string]Value{"age": int64(34), "job": "デザイナー"}},
			Tuple{Attrs: schema, Data: map[string]Value{"age": int64(21), "job": "マネージャー"}},
		},
	}
	const jsonStream = `{ "projection": {
			"input": { "relation": {"name":"test/staff1"}},
			"attrs": [ "age","job" ]
	}}`
	m := Stream{}
	if err := json.NewDecoder(strings.NewReader(jsonStream)).Decode(&m); err != nil {
		log.Fatal(err)
	}
	result, _ := StreamToRelation(m, testData2)
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Does not match 'SELECT age,job FROM Staff'\nresult:% #v,\n want:% #v", result, want)
	}
}
