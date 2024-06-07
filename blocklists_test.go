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
