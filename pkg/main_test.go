package main

import "testing"

func Test_parseArgs(t *testing.T) {
	type args struct {
		configFile *string
		inputPath  *string
		scriptFile *string
		scriptPath *string
		outputPath *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parseArgs(tt.args.configFile, tt.args.inputPath, tt.args.scriptFile, tt.args.scriptPath, tt.args.outputPath); (err != nil) != tt.wantErr {
				t.Errorf("parseArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
