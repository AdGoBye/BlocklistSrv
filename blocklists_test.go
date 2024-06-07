package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_fetchBlocklist(t *testing.T) {
	type args struct {
		location string
	}
	tests := []struct {
		name    string
		args    args
		want    Blocklist
		wantErr bool
	}{
		{
			"invalid schema", args{location: "bogus://zombo.com"}, Blocklist{}, true,
		},
		{
			"invalid url", args{location: "notaurlreallytrustme"}, Blocklist{}, true,
		},
		{
			"invalid file", args{location: "file://notafile"}, Blocklist{}, true,
		},
		{
			"valid url but invalid location", args{location: "https://realcountry.realcountry"}, Blocklist{}, true,
		},
		{
			"valid web location but not blocklist",
			args{location: "https://example.com"}, Blocklist{}, true,
		},
		{
			name: "valid blocklist from web",
			args: args{
				location: "https://raw.githubusercontent.com/AdGoBye/AdGoBye-Blocklists/5646b6d5aecf00d336184bd70fd4c090b3a25f86/AGBCommunity.toml"},
			want: Blocklist{Blocks: []Block{
				{FriendlyName: "Default Home", WorldId: "wrld_4432ea9b-729c-46e3-8eaf-846aa0a37fdd",
					GameObjects: []Gameobject{{Name: "posterlight (8)"}},
				},
				{FriendlyName: "Movie & Chill", WorldId: "wrld_791ebf58-54ce-4d3a-a0a0-39f10e1b20b2",
					GameObjects: []Gameobject{{Name: "Label (2)"}},
				},
				{
					FriendlyName: "The Black Cat",
					WorldId:      "wrld_4cf554b4-430c-4f8f-b53e-1f294eed230b",
					GameObjects: []Gameobject{{Name: "cork medium",
						Position: Pointer(GameobjectPosition{X: 26.773, Y: 3.244, Z: -13.982})},
					},
				},
				{
					FriendlyName: "Furry Hideout", WorldId: "wrld_4b341546-65ff-4607-9d38-5b7f8f405132",
					GameObjects: []Gameobject{
						{Name: "PPSUI (2)"},
						{Name: "Cube (5)", Position: Pointer(GameobjectPosition{X: -29.597, Y: 44.894, Z: 6.501})},
					},
				},
				{
					FriendlyName: "Furry Talk and Chill", WorldId: "wrld_e76f0ce1-8b2f-4fd7-a6ac-84443d6f26f1",
					GameObjects: []Gameobject{{Name: "Bottom Tex"}},
				},
				{FriendlyName: "Murder 4", WorldId: "wrld_858dfdfc-1b48-4e1e-8a43-f0edc611e5fe",
					GameObjects: []Gameobject{{Name: "Link (2)"}},
				},
				{FriendlyName: "Prison Escape!", WorldId: "wrld_14750dd6-26a1-4edb-ae67-cac5bcd9ed6a",
					GameObjects: []Gameobject{
						{Name: "Group Sign"},
						{
							Name: "Image (1)", Position: Pointer(GameobjectPosition{X: -78.4, Y: -95, Z: 0}),
							Parent: Pointer(Gameobject{Name: "Panel"}),
						},
						{
							Name: "Image (2)", Position: Pointer(GameobjectPosition{X: -140, Y: -135, Z: 0}),
							Parent: Pointer(Gameobject{Name: "Panel"}),
						},
						{
							Name: "Image (3)", Position: Pointer(GameobjectPosition{X: -175, Y: -175, Z: 0}),
							Parent: Pointer(Gameobject{Name: "Panel"}),
						},
						{
							Name: "Text (TMP)", Position: Pointer(GameobjectPosition{X: -2.5, Y: -95, Z: 0}),
							Parent: Pointer(Gameobject{Name: "Panel"}),
						},
						{
							Name: "Text (TMP) (1)", Position: Pointer(GameobjectPosition{X: 14, Y: -135, Z: 0}),
							Parent: Pointer(Gameobject{Name: "Panel"}),
						},
						{
							Name: "Text (TMP) (2)", Position: Pointer(GameobjectPosition{X: 18, Y: -175, Z: 0}),
							Parent: Pointer(Gameobject{Name: "Panel"}),
						},
					},
				},
				{FriendlyName: "Just B Club 3", WorldId: "wrld_e6569266-21cd-4275-8aef-47fcb7458931",
					GameObjects: []Gameobject{
						{Name: "Discord TV Ad (1)"},
						{Name: "TV Prefab UNIQUE", Position: Pointer(GameobjectPosition{X: -4.366071, Y: 3.072498, Z: -54.21133})},
						{Name: "Poster (9)", Parent: Pointer(Gameobject{Name: "Poster (9)"})},
					},
				},
				{FriendlyName: "The room of the rain", WorldId: "wrld_fae3fa95-bc18-46f0-af57-f0c97c0ca90a",
					GameObjects: []Gameobject{
						{Name: "Neverphone"},
						{Name: "Patreon ui", Parent: Pointer(Gameobject{Name: "Patreon Things"})},
						{Name: "Patreon panel", Parent: Pointer(Gameobject{Name: "Patreon Things"})},
						{Name: "Patreon texture Changer", Parent: Pointer(Gameobject{Name: "Patreon Things"})},
						{Name: "Patreon texture Changer (1)", Parent: Pointer(Gameobject{Name: "Patreon Things"})},
					},
				},
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchBlocklist(tt.args.location)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchBlocklist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.EqualValuesf(t, tt.want, got, "fetchBlocklist() got = %v, want %v", got, tt.want)
		})
	}
}

func TestWorldObjectIndex_GetWorldById(t *testing.T) {
	type fields struct {
		Index map[string]WorldObject
	}
	type args struct {
		HashedWorldId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *WorldObject
	}{
		{
			name: "exists in index",
			fields: fields{Index: map[string]WorldObject{
				"ZDZkNWExN2IzMGEyYWZmNmZiYzlhOGZlYzhiZDQ4MGRiZTdiYzEzMTVlYWQ0NzhjOTRmYjdiMDUxOTQ4MDI2Ng==": {FriendlyName: "Test"},
			}},
			args: args{HashedWorldId: "ZDZkNWExN2IzMGEyYWZmNmZiYzlhOGZlYzhiZDQ4MGRiZTdiYzEzMTVlYWQ0NzhjOTRmYjdiMDUxOTQ4MDI2Ng=="},
			want: Pointer(WorldObject{FriendlyName: "Test", GameObjectMapping: nil}),
		},
		{
			name:   "doesn't exist in index",
			fields: fields{Index: map[string]WorldObject{}},
			args:   args{HashedWorldId: "ZDZkNWExN2IzMGEyYWZmNmZiYzlhOGZlYzhiZDQ4MGRiZTdiYzEzMTVlYWQ0NzhjOTRmYjdiMDUxOTQ4MDI2Ng=="},
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexObj := WorldObjectIndex{
				Index: tt.fields.Index,
			}
			assert.Equalf(t, tt.want, indexObj.GetWorldById(tt.args.HashedWorldId), "GetWorldById(%v)", tt.args.HashedWorldId)
		})
	}
}

func Test_stringToHash(t *testing.T) {
	type args[inputs interface{ string | []byte }] struct {
		input inputs
	}
	type testCase[inputs interface{ string | []byte }] struct {
		name       string
		args       args[inputs]
		wantOutput []byte
	}
	stringTests := []testCase[string]{
		{
			name:       "Empty string",
			args:       args[string]{input: ""},
			wantOutput: []byte{0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55},
		},
		{
			name:       "Proper string",
			args:       args[string]{input: "wrld_00000000-0000-0000-0000-000000000000"},
			wantOutput: []byte{0xd6, 0xd5, 0xa1, 0x7b, 0x30, 0xa2, 0xaf, 0xf6, 0xfb, 0xc9, 0xa8, 0xfe, 0xc8, 0xbd, 0x48, 0xd, 0xbe, 0x7b, 0xc1, 0x31, 0x5e, 0xad, 0x47, 0x8c, 0x94, 0xfb, 0x7b, 0x5, 0x19, 0x48, 0x2, 0x66},
		},
	}
	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantOutput, stringToHash(tt.args.input), "stringToHash(%v)", tt.args.input)
		})
	}
}

func Test_generateObjectIndex(t *testing.T) {
	type args struct {
		blocklistsLocations []string
	}
	tests := []struct {
		name        string
		args        args
		wantMapping map[string]WorldObject
	}{
		{
			"index object from agbcommunity 5646b6d",
			args{blocklistsLocations: []string{"https://raw.githubusercontent.com/AdGoBye/AdGoBye-Blocklists/5646b6d5aecf00d336184bd70fd4c090b3a25f86/AGBCommunity.toml"}},
			map[string]WorldObject{"//72sZH9E1KefGiDV2vl6hS7tPFxRczk4f8tozmZT+A=": {FriendlyName: "The Black Cat", GameObjectMapping: map[string]Gameobject{"1rjfo3+mVV3ZAVcJI6voGUkdaE+MLvFyI0PMBzyG1XY=": {Name: "cork medium", Position: Pointer(GameobjectPosition{X: 26.773, Y: 3.244, Z: -13.982}), Parent: nil}}}, "UisKnWNb5njDLfcjjdHEql4PcSNaWkRudA2yXXDuAZQ=": {FriendlyName: "Furry Hideout", GameObjectMapping: map[string]Gameobject{"44+cpG4N7ycmV+zAvHftxlqjjIKtgIWwPBrRR6ECQfY=": {Name: "Cube (5)", Position: Pointer(GameobjectPosition{X: -29.597, Y: 44.894, Z: 6.501}), Parent: nil}, "eZpZL6VdV6MIwt5Zp85xa/bCb1uQvr7aNwMAa9Ci6r4=": {Name: "PPSUI (2)", Position: nil, Parent: nil}}}, "Wn6QIbBCeSTvZLenudccQBxdux4AbFIsCSRuiqFk0wc=": {FriendlyName: "Furry Talk and Chill", GameObjectMapping: map[string]Gameobject{"oZZfW96QQe/r97oBD+MJcFomNy/hpNTR3OeVj1QouWo=": {Name: "Bottom Tex", Position: nil, Parent: nil}}}, "ejki17vqRmaOoosLxZ7Btg9c0RP6Q4hl2vClhV6Mxjc=": {FriendlyName: "Murder 4", GameObjectMapping: map[string]Gameobject{"E+vtKYtVJUe5ZB6Exc9r6wtNhCTV7HoI1hPDVGByNr8=": {Name: "Link (2)", Position: nil, Parent: nil}}}, "ekogGjm7rMq9nDvpj5zvEnl6hu9kZFSbUTG61Uj+YzM=": {FriendlyName: "Prison Escape!", GameObjectMapping: map[string]Gameobject{"/JIRSiPxW8PbT3/6hiteuy19SvYXKRkiOF000vaz4IA=": {Name: "Image (3)", Position: Pointer(GameobjectPosition{X: -175, Y: -175, Z: 0}), Parent: Pointer(Gameobject{Name: "Panel"})}, "0H9wPe2SQozFHxSczgUCQX4HGA++4FiQ8NcXdh43mF8=": {Name: "Group Sign", Position: nil, Parent: nil}, "79fSWuUcG0xQRv3NVFxCZ08yq4Ir7+28e0BxhdBF+Lk=": {Name: "Text (TMP) (2)", Position: Pointer(GameobjectPosition{X: 18, Y: -175, Z: 0}), Parent: Pointer(Gameobject{Name: "Panel"})}, "NW0R88MKtatWNgpqkGRxsAPx+YYasT5FFh51sJPoO7w=": {Name: "Text (TMP) (1)", Position: Pointer(GameobjectPosition{X: 14, Y: -135, Z: 0}), Parent: Pointer(Gameobject{Name: "Panel"})}, "gVTW8UCy//rVmmcMU7GHTCvIl2rikJie1XDMoRZgIRU=": {Name: "Text (TMP)", Position: Pointer(GameobjectPosition{X: -2.5, Y: -95, Z: 0}), Parent: Pointer(Gameobject{Name: "Panel"})}, "hl4nmC0Va0NkjmiEd2yhKYBdwB0V95vxy1xNiQ+iRdo=": {Name: "Image (2)", Position: Pointer(GameobjectPosition{X: -140, Y: -135, Z: 0}), Parent: Pointer(Gameobject{Name: "Panel"})}, "mHCtrK1xyvq2zsj6JCg6GONmgzYE2hdQpp/vamFQURQ=": {Name: "Image (1)", Position: Pointer(GameobjectPosition{X: -78.4, Y: -95, Z: 0}), Parent: Pointer(Gameobject{Name: "Panel"})}}}, "jtLZiuvHwfHGwgUPKXv+znfckOxuWp7/PHoLQKlQcPs=": {FriendlyName: "Just B Club 3", GameObjectMapping: map[string]Gameobject{"2nK7NHwf3fLLVlNJkj3PsOm2St2M4U6Am0jNj1L9IXk=": {Name: "TV Prefab UNIQUE", Position: Pointer(GameobjectPosition{X: -4.366071, Y: 3.072498, Z: -54.21133}), Parent: nil}, "VaRJlJS9rr63cwr8+UirDj+/e7OqqAnwMykRccS9Xxs=": {Name: "Discord TV Ad (1)", Position: nil, Parent: nil}, "vm9Zs730Y6b3lkWASmTsoCxawEFZCUBOLADwjGtmFp0=": {Name: "Poster (9)", Position: nil, Parent: Pointer(Gameobject{Name: "Poster (9)"})}}}, "sYREJ3yvphTlFAF97DgXjRtoiLfWWoN1FcKOAelhPYs=": {FriendlyName: "The room of the rain", GameObjectMapping: map[string]Gameobject{"H3EGdEcNRGq2YGF2NnYj+XdjmJIFm7x50b1bgR3Er1Q=": {Name: "Patreon ui", Position: nil, Parent: Pointer(Gameobject{Name: "Patreon Things"})}, "H8sLfASKE2wV/UU/5/XYOOqxpSjKqtNbE2GLnDH1zys=": {Name: "Neverphone", Position: nil, Parent: nil}, "PSl0xoSem8oYbhJR2yH/GGCw3Q6xyN+WkX1GIiUjErA=": {Name: "Patreon texture Changer (1)", Position: nil, Parent: Pointer(Gameobject{Name: "Patreon Things"})}, "S8gfXVTsPW+pRlqKbf2bBKOS2AQ156bjRy8QHMxET7Y=": {Name: "Patreon panel", Position: nil, Parent: Pointer(Gameobject{Name: "Patreon Things"})}, "rRUqUuTSY4tekMftvMmKN/6ecwbED4yy4hGH/dITo3w=": {Name: "Patreon texture Changer", Position: nil, Parent: Pointer(Gameobject{Name: "Patreon Things"})}}}, "sbAfDcLEMfXdqt7ymd5y5wtGDnQHIa5oFqfjSnxSv+8=": {FriendlyName: "Default Home", GameObjectMapping: map[string]Gameobject{"tAyJF6Ht0JkOMLMlZI/5hacz365Y+DJSaNgayRDazkg=": {Name: "posterlight (8)", Position: nil, Parent: nil}}}, "zPjucIpmtcG2w2sUUVaq7j0tuRKLQINxrLRDEiJv+ZQ=": {FriendlyName: "Movie & Chill", GameObjectMapping: map[string]Gameobject{"NZk/vsuk5MZgw00AgfSzItvXCUMaCRgJTMVEA/pyXAc=": {Name: "Label (2)", Position: nil, Parent: nil}}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantMapping, generateObjectIndex(tt.args.blocklistsLocations), "generateObjectIndex(+%v)", tt.args.blocklistsLocations)
		})
	}
}
