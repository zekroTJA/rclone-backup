package list

import (
	"reflect"
	"strings"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type testStruct struct {
		name    string
		in      string
		exp     []Entry
		wantErr bool
	}

	tests := []testStruct{
		{
			name: "paths",
			in: "/foo/bar/baz\n" +
				"  /hello/world  \n" +
				"\n" +
				"C:\\hello\\world\n",
			exp: []Entry{
				{Path: "/foo/bar/baz", Args: []string{}},
				{Path: "/hello/world", Args: []string{}},
				{Path: "C:\\hello\\world", Args: []string{}},
			},
		},
		{
			name: "paths-with-args",
			in: "/home/user/importantdata,**/exclude1/**\n" +
				"/mnt/data/photos,{{.*\\.ARW.*}}\n" +
				"D:\\archiv,**.mp4,max-size=15G\n",
			exp: []Entry{
				{Path: "/home/user/importantdata",
					Args: []string{"--exclude=**/exclude1/**"}},
				{Path: "/mnt/data/photos",
					Args: []string{"--exclude={{.*\\.ARW.*}}"}},
				{Path: "D:\\archiv",
					Args: []string{"--exclude=**.mp4", "--max-size=15G"}},
			},
		},
		{
			name: "comments",
			in: "  # This is also a comment!\n" +
				"/foo\n" +
				"# this is a comment!\n" +
				"/bar/baz,**.mp4",
			exp: []Entry{
				{Path: "/foo", Args: []string{}},
				{Path: "/bar/baz",
					Args: []string{"--exclude=**.mp4"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := Unmarshal(strings.NewReader(tt.in))
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.exp) {
				t.Errorf("Unmarshal() gotRes = %v, want %v", gotRes, tt.exp)
			}
		})
	}
}
