package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

const testEnv = "PF_TMP_TEST"

func Test_pathfinder(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.Setenv(testEnv, tmpDir)
	if err != nil {
		t.Fatalf("Error setting test environment key : %v\n", err)
	}

	defer t.Cleanup(func() {
		err = os.Unsetenv(testEnv)
		if err != nil {
			HandleError(err)
		}
	})

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Error changing directories : %v\n", err)
	}

	var dirs [5]string
	for i := 1; i <= 4; i++ {
		dirs[i] = filepath.Join(tmpDir, fmt.Sprintf("dir%d", i))
		err = os.MkdirAll(dirs[i], 0777)
		if err != nil {
			t.Fatalf("Error creating sub directories : %v\n", err)
		}
	}

	dirs[0] = filepath.Join(tmpDir, "uniquedirthatdoesntexistsanywhereelse")
	// This is done to make sure test2 passes without any external conflicts
	err = os.MkdirAll(dirs[0], 0777)
	if err != nil {
		t.Fatalf("Error creating sub directories : %v\n", err)
	}

	type args struct {
		w         bytes.Buffer
		c         *Cache
		ignoreDir bool
		path      string
	}
	testCache := InitCache()
	tests := []struct {
		name    string
		args    args
		wantBuf string // Contents of bytes.Buffer
		wantRet int    // Return code
	}{
		{
			name: "Folder not found", args: args{c: testCache, ignoreDir: false, path: "notfounddir"}, wantBuf: "notfounddir", wantRet: EXIT_FOLDERNOTFOUND,
		},
		{
			name: "Folder exists but ignored (-i)", args: args{c: testCache, ignoreDir: true, path: filepath.Base(dirs[0])}, wantBuf: filepath.Base(dirs[0]), wantRet: EXIT_FOLDERNOTFOUND,
		},
		{
			name: "Folder found without ignore", args: args{c: testCache, ignoreDir: false, path: filepath.Base(dirs[3])}, wantBuf: dirs[3], wantRet: EXIT_SUCCESS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.c.cleanCache()
			gotRet := pathfinder(&tt.args.w, tt.args.c, tt.args.ignoreDir, tt.args.path)
			gotBuf := tt.args.w.String()
			if gotBuf != tt.wantBuf {
				t.Errorf("pathfinder = %v, wantBuf = %v", gotBuf, tt.wantBuf)
			}
			if gotRet != tt.wantRet {
				t.Errorf("pathfinder = %v, wantRet = %v", gotRet, tt.wantRet)
			}
		})
	}
}
