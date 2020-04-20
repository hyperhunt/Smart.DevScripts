package osm

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestOSM_Append(t *testing.T) {
	o := &OSM{}

	o.Append(&Node{ID: 1})
	if n := o.Nodes[0]; n.ID != 1 {
		t.Errorf("incorrect node: %v", n)
	}

	o.Append(&Way{ID: 2})
	if w := o.Ways[0]; w.ID != 2 {
		t.Errorf("incorrect way: %v", w)
	}

	o.Append(&Relation{ID: 3})
	if r := o.Relations[0]; r.ID != 3 {
		t.Errorf("incorrect relation: %v", r)
	}

	o.Append(&Changeset{ID: 4})
	if cs := o.Changesets[0]; cs.ID != 4 {
		t.Errorf("incorrect changeset: %v", cs)
	}

	o.Append(&User{ID: 5})
	if u := o.Users[0]; u.ID != 5 {
		t.Errorf("incorrect user: %v", u)
	}

	o.Append(&Note{ID: 6})
	if n := o.Notes[0]; n.ID != 6 {
		t.Errorf("incorrect note: %v", n)
	}
}

func TestOSM_Elements(t *testing.T) {
	var no *OSM
	no.Elements() // should not panic if OSM is nil

	es := Elements{
		&Node{ID: 1},
		&Way{ID: 2},
		&Relation{ID: 3},
	}

	o := &OSM{}
	o.Append(es[0])
	o.Append(es[1])
	o.Append(es[2])

	elements := o.Elements()
	if !reflect.DeepEqual(elements, es) {
		t.Errorf("incorrect elements: %v", elements)
	}
}

func TestOSM_Objects(t *testing.T) {
	var no *OSM
	no.Objects() // should not panic if OSM is nil

	os := Objects{
		&Node{ID: 1},
		&Way{ID: 2},
		&Relation{ID: 3},
		&Changeset{ID: 4},
		&User{ID: 5},
		&Note{ID: 6},
	}

	o := &OSM{}
	for _, obj := range os {
		o.Append(obj)
	}

	objects := o.Objects()
	if !reflect.DeepEqual(objects, os) {
		t.Errorf("incorrect objects: %v", objects)
	}
}

func TestOSM_FeatureIDs(t *testing.T) {
	var no *OSM
	no.FeatureIDs() // should not panic if OSM is nil

	es := Elements{
		&Node{ID: 1},
		&Way{ID: 2},
		&Relation{ID: 3},
	}

	o := &OSM{}
	o.Append(es[0])
	o.Append(es[1])
	o.Append(es[2])

	expected := FeatureIDs{
		NodeID(1).FeatureID(),
		WayID(2).FeatureID(),
		RelationID(3).FeatureID(),
	}

	if ids := o.FeatureIDs(); !reflect.DeepEqual(ids, expected) {
		t.Errorf("incorrect ids: %v", ids)
	}
}

func TestOSM_ElementIDs(t *testing.T) {
	var no *OSM
	no.ElementIDs() // should not panic if OSM is nil

	es := Elements{
		&Node{ID: 1, Version: 4},
		&Way{ID: 2, Version: 5},
		&Relation{ID: 3, Version: 6},
	}

	o := &OSM{}
	o.Append(es[0])
	o.Append(es[1])
	o.Append(es[2])

	expected := ElementIDs{
		NodeID(1).ElementID(4),
		WayID(2).ElementID(5),
		RelationID(3).ElementID(6),
	}

	if ids := o.ElementIDs(); !reflect.DeepEqual(ids, expected) {
		t.Errorf("incorrect ids: %v", ids)
	}
}

func TestOSM_Marshal(t *testing.T) {
	c := loadChange(t, "testdata/changeset_38162206.osc")
	o1 := flattenOSM(c)
	o1.Bounds = &Bounds{1.1, 2.2, 3.3, 4.4}

	data, err := o1.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	o2, err := UnmarshalOSM(data)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Errorf("osm are not equal")
		t.Logf("%+v", o1)
		t.Logf("%+v", o2)
	}

	// second changeset
	c = loadChange(t, "testdata/changeset_38162210.osc")
	o1 = flattenOSM(c)

	data, err = o1.Marshal()
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	o2, err = UnmarshalOSM(data)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Errorf("osm are not equal")
		t.Logf("%+v", o1)
		t.Logf("%+v", o2)
	}
}

func TestOSM_MarshalJSON(t *testing.T) {
	o := &OSM{
		Version:   0.6,
		Generator: "osm-go",
		Nodes: Nodes{
			&Node{ID: 123},
		},
		Ways: Ways{
			&Way{ID: 456},
		},
		Relations: Relations{
			&Relation{ID: 789},
		},
		Changesets: Changesets{
			&Changeset{ID: 10},
		},
	}

	data, err := json.Marshal(o)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	if !bytes.Equal(data, []byte(`{"version":0.6,"generator":"osm-go","elements":[{"type":"node","id":123,"lat":0,"lon":0,"visible":false,"timestamp":"0001-01-01T00:00:00Z"},{"type":"way","id":456,"visible":false,"timestamp":"0001-01-01T00:00:00Z","nodes":[]},{"type":"relation","id":789,"visible":false,"timestamp":"0001-01-01T00:00:00Z","members":[]},{"type":"changeset","id":10,"created_at":"0001-01-01T00:00:00Z","closed_at":"0001-01-01T00:00:00Z","open":false}]}`)) {
		t.Errorf("incorrect json: %v", string(data))
	}
}

func TestOSM_MarshalXML(t *testing.T) {
	o := &OSM{
		Version:     0.7,
		Generator:   "osm-go-test",
		Copyright:   "copyright1",
		Attribution: "attribution1",
		License:     "license1",
		Nodes: Nodes{
			&Node{ID: 123},
		},
	}

	data, err := xml.Marshal(o)
	if err != nil {
		t.Fatalf("xml marshal error: %v", err)
	}

	expected := `<osm version="0.7" generator="osm-go-test" copyright="copyright1" attribution="attribution1" license="license1"><node id="123" lat="0" lon="0" user="" uid="0" visible="false" version="0" changeset="0" timestamp="0001-01-01T00:00:00Z"></node></osm>`

	if !bytes.Equal(data, []byte(expected)) {
		t.Errorf("incorrect marshal, got: %s", string(data))
	}

	// omit attributes if not defined
	o = &OSM{}
	data, err = xml.Marshal(o)
	if err != nil {
		t.Fatalf("xml marshal error: %v", err)
	}

	expected = `<osm></osm>`
	if !bytes.Equal(data, []byte(expected)) {
		t.Errorf("incorrect marshal, got: %s", string(data))
	}
}

func flattenOSM(c *Change) *OSM {
	o := c.Create
	if o == nil {
		o = &OSM{}
	}

	if c.Modify != nil {
		o.Nodes = append(o.Nodes, c.Modify.Nodes...)
		o.Ways = append(o.Ways, c.Modify.Ways...)
		o.Relations = append(o.Relations, c.Modify.Relations...)
	}

	if c.Delete != nil {
		o.Nodes = append(o.Nodes, c.Delete.Nodes...)
		o.Ways = append(o.Ways, c.Delete.Ways...)
		o.Relations = append(o.Relations, c.Delete.Relations...)
	}

	return o
}

func cleanXMLNameFromOSM(o *OSM) {
	for _, n := range o.Nodes {
		n.XMLName = xmlNameJSONTypeNode{}
	}

	for _, w := range o.Ways {
		w.XMLName = xmlNameJSONTypeWay{}
	}

	for _, r := range o.Relations {
		// r.XMLName = xml.Name{}
		r.XMLName = xmlNameJSONTypeRel{}
	}
}
