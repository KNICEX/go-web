package web

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	f, err := os.Open("file_test.txt")
	require.NoError(t, err)
	data := make([]byte, 1024)
	n, err := f.Read(data)
	require.NoError(t, err)
	fmt.Println(n)

	n, err = f.WriteString("hello")
	fmt.Println(n)
	require.Error(t, err)
	// access denied
	fmt.Println(err)
	_ = f.Close()

	f, err = os.OpenFile("file_test.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	require.NoError(t, err)

	n, err = f.WriteString("hello")
	fmt.Println(n)
	require.NoError(t, err)
	_ = f.Close()
}
