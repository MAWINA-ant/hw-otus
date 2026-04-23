package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("ErrOffsetExceedsFileSize", func(t *testing.T) {
		err := Copy("./testdata/input.txt", "./test_out.txt", 100000, 0)
		require.Equal(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("ErrUnsupportedFile", func(t *testing.T) {
		err := Copy("/dev/urandom", "./test_out.txt", 0, 0)
		require.Equal(t, err, ErrUnsupportedFile)
	})

	tests := []struct {
		name     string
		from     string
		to       string
		offset   int64
		limit    int64
		standard string
	}{
		{
			name: "no limit and no offset", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 0, limit: 0, standard: "./testdata/out_offset0_limit0.txt",
		},
		{
			name: "limit 10 and no offset", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 0, limit: 10, standard: "./testdata/out_offset0_limit10.txt",
		},
		{
			name: "limit 1000 and no offset", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 0, limit: 1000, standard: "./testdata/out_offset0_limit1000.txt",
		},
		{
			name: "limit 10000 and no offset", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 0, limit: 10000, standard: "./testdata/out_offset0_limit10000.txt",
		},
		{
			name: "limit 1000 and offset 100", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 100, limit: 1000, standard: "./testdata/out_offset100_limit1000.txt",
		},
		{
			name: "limit 1000 and offset 6000", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 6000, limit: 1000, standard: "./testdata/out_offset6000_limit1000.txt",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.NoError(t, err)
			outFile, _ := os.Open("test_out.txt")
			standardFile, _ := os.Open(tc.standard)
			outFileStat, _ := outFile.Stat()
			standardFileStat, _ := standardFile.Stat()
			require.Equal(t, outFileStat.Size(), standardFileStat.Size())
		})
	}
	os.Remove("./test_out.txt")
}

/*
	{name: "no limit and no offset", from: "./testdata/input.txt",
		to: "./test_out.txt", offset: 0, limit: 0, standard: "./testdata/out_offset0_limit0.txt"},
		{name: "limit 10 and no offset", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 0, limit: 10, standard: "./testdata/out_offset0_limit10.txt"},
		{name: "limit 1000 and no offset", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 0, limit: 1000, standard: "./testdata/out_offset0_limit1000.txt"},
		{name: "limit 10000 and no offset", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 0, limit: 10000, standard: "./testdata/out_offset0_limit10000.txt"},
		{name: "limit 1000 and offset 100", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 100, limit: 1000, standard: "./testdata/out_offset100_limit1000.txt"},
		{name: "limit 1000 and offset 6000", from: "./testdata/input.txt",
			to: "./test_out.txt", offset: 6000, limit: 1000, standard: "./testdata/out_offset6000_limit1000.txt"},
*/
