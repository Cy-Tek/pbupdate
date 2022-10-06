package pbupdate

import (
	"reflect"
	"testing"
)

type address struct {
	State  string  `json:"state"`
	Zip    string  `json:"zip"`
	Unit   *string `json:"unit"`
	Street string  `json:"street"`
}

type user struct {
	Connections    int     `json:"connections"`
	Location       address `json:"location"`
	Name           string  `json:"name"`
	IntegrationKey *string `json:"integrationKey"`
}

type requestUser struct {
	metadata map[string]string

	Connections    int     `json:"connections"`
	Location       address `json:"location"`
	Name           string  `json:"name"`
	IntegrationKey *string `json:"integrationKey"`
	AdditionalInfo string  `json:"additionalInfo"`
}

func Test_objectToMap(t *testing.T) {
	type args struct {
		object any
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := objectToMap(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readJsonPath(t *testing.T) {
	type args struct {
		path   string
		object any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "Should read a top level initialized field",
			args: args{
				path: "connections",
				object: user{
					Connections: 0,
				},
			},
			// Json is automatically returned as a float64 when not being specifically unmarshaled into a specific type
			want:    0.0,
			wantErr: false,
		},
		{
			name: "Should read a top level nil field",
			args: args{
				path:   "integrationKey",
				object: user{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Should read a sub-object field",
			args: args{
				path: "location.state",
				object: user{
					Location: address{
						State: "NC",
					},
				},
			},
			want:    "NC",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			object, err := objectToMap(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectToMap() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := readJsonPath(tt.args.path, object)
			if (err != nil) != tt.wantErr {
				t.Errorf("readJsonPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readJsonPath() got = %v with type %T, want %v with Type %T", got, got, tt.want, tt.want)
			}
		})
	}
}

func Test_splitDotPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitDotPath(tt.args.path)
			if got != tt.want {
				t.Errorf("splitDotPath() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitDotPath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_updateJsonValue(t *testing.T) {
	type args struct {
		path     string
		object   any
		newValue any
	}
	tests := []struct {
		name    string
		args    args
		after   user
		wantErr bool
	}{
		{
			name: "Should update a top level initialized field",
			args: args{
				path: "connections",
				object: user{
					Connections: 0,
				},
				newValue: 20.0,
			},
			after: user{
				Connections: 20,
			},
			wantErr: false,
		},
		{
			name: "Should update a sub-object's field",
			args: args{
				path: "location.state",
				object: user{
					Location: address{
						State: "SC",
					},
				},
				newValue: "TN",
			},
			after: user{
				Location: address{
					State: "TN",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			object, err := objectToMap(tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("objectToMap() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := updateJsonValue(tt.args.path, object, tt.args.newValue); (err != nil) != tt.wantErr {
				t.Errorf("updateJsonValue() error = %v, wantErr %v", err, tt.wantErr)
			}

			after, _ := objectToMap(tt.after)

			if !reflect.DeepEqual(object, after) {
				t.Errorf("readJsonPath() got = %v, want %v", object, after)
			}
		})
	}
}

func TestCopyValuesFromPaths(t *testing.T) {
	unit := "Unit 202"
	newUnit := "Unit 3000"
	integrationKey := "C0UFJ-ASDFE"
	newIntegrationKey := "C0UFJ-TEST"

	type args struct {
		paths  []string
		source requestUser
		dest   user
	}
	tests := []struct {
		name    string
		args    args
		want    user
		wantErr bool
	}{
		{
			name: "Should handle a single path",
			args: args{
				paths: []string{"location.state"},
				source: requestUser{
					metadata:    nil,
					Connections: 20,
					Location: address{
						State: "NC",
					},
					Name:           "",
					IntegrationKey: nil,
					AdditionalInfo: "",
				},
				dest: user{
					Connections: 10,
					Location: address{
						State:  "WV",
						Zip:    "11111",
						Unit:   &unit,
						Street: "2020 Vision Dr",
					},
					Name:           "Test User",
					IntegrationKey: &integrationKey,
				},
			},
			want: user{
				Connections: 10,
				Location: address{
					State:  "NC",
					Zip:    "11111",
					Unit:   &unit,
					Street: "2020 Vision Dr",
				},
				Name:           "Test User",
				IntegrationKey: &integrationKey,
			},
			wantErr: false,
		},
		{
			name: "Should handle every possible path",
			args: args{
				paths: []string{"location.state", "location.zip", "location.unit", "location.street", "connections", "name", "integrationKey"},
				source: requestUser{
					metadata:    nil,
					Connections: 20,
					Location: address{
						State:  "NC",
						Zip:    "22222",
						Unit:   &newUnit,
						Street: "202 Baker Street",
					},
					Name:           "Test User 2",
					IntegrationKey: &newIntegrationKey,
					AdditionalInfo: "This object can be copied over now",
				},
				dest: user{
					Connections: 10,
					Location: address{
						State:  "WV",
						Zip:    "11111",
						Unit:   &unit,
						Street: "2020 Vision Dr",
					},
					Name:           "Test User",
					IntegrationKey: &integrationKey,
				},
			},
			want: user{
				Connections: 20,
				Location: address{
					State:  "NC",
					Zip:    "22222",
					Unit:   &newUnit,
					Street: "202 Baker Street",
				},
				Name:           "Test User 2",
				IntegrationKey: &newIntegrationKey,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CopyValuesFromPaths(tt.args.paths, tt.args.source, tt.args.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("CopyValuesFromPaths() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("CopyValuesFromPaths() got = %v, want %v", got, tt.want)
			}
		})
	}
}
