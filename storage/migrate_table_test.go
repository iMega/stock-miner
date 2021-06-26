package storage

import (
	"testing"
)

func Test_hasColumn(t *testing.T) {
	type args struct {
		ti tableInfo
		c  col
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "column exitsts",
			args: args{
				ti: tableInfo{Columns: []col{{Name: "test"}}},
				c:  col{Name: "test"},
			},
			want: true,
		},
		{
			name: "column does'nt exist",
			args: args{
				ti: tableInfo{Columns: []col{{Name: "other"}}},
				c:  col{Name: "test"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasColumn(tt.args.ti, tt.args.c); got != tt.want {
				t.Errorf("hasColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_equalColumn(t *testing.T) {
	type args struct {
		ti tableInfo
		c  col
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "columns are equal",
			args: args{
				ti: tableInfo{Columns: []col{
					{
						CID:          1,
						Name:         "test",
						Type:         "VARCHAR(64)",
						NotNull:      1,
						DefaultValue: "value",
						PK:           1,
					},
				}},
				c: col{
					CID:          1,
					Name:         "test",
					Type:         "VARCHAR(64)",
					NotNull:      1,
					DefaultValue: "value",
					PK:           1,
				},
			},
			want: true,
		},
		{
			name: "columns are not equal. different PK",
			args: args{
				ti: tableInfo{Columns: []col{
					{
						CID:          1,
						Name:         "test",
						Type:         "VARCHAR(64)",
						NotNull:      1,
						DefaultValue: "value",
						PK:           0,
					},
				}},
				c: col{
					CID:          1,
					Name:         "test",
					Type:         "VARCHAR(64)",
					NotNull:      1,
					DefaultValue: "value",
					PK:           1,
				},
			},
			want: false,
		},
		{
			name: "columns are not equal. different names",
			args: args{
				ti: tableInfo{Columns: []col{
					{
						CID:          1,
						Name:         "other",
						Type:         "VARCHAR(64)",
						NotNull:      1,
						DefaultValue: "value",
						PK:           1,
					},
				}},
				c: col{
					CID:          1,
					Name:         "test",
					Type:         "VARCHAR(64)",
					NotNull:      1,
					DefaultValue: "value",
					PK:           1,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equalColumn(tt.args.ti, tt.args.c); got != tt.want {
				t.Errorf("equalColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}
