package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"unicode/utf8"
)

var encodeCmdExample = `
  uc encode hello world -x
  output: 
  68 65 6C 6C 6F 
  77 6F 72 6C 64 
`

var (
	removeSpaceEncodeCmd    bool
	outputHexEncodeCmd      bool
	prefixEncodeCmd         string
	inputCodePointEncodeCmd bool
	encodeCmd               = &cobra.Command{
		Use:     "encode [<args>]",
		Short:   "Encode string using UTF-8",
		Long:    "Convert a string to a sequence of UTF-8 encoded values",
		Example: encodeCmdExample,
		Args:    cobra.MinimumNArgs(1),
		RunE:    runEncodeCmd,
	}
)

func init() {
	encodeCmd.Flags().BoolVarP(&removeSpaceEncodeCmd, "remove-space", "", false, "removes space between each two hex digits in output")
	encodeCmd.Flags().BoolVarP(&inputCodePointEncodeCmd, "unicode", "u", false, "input is a sequence of Unicode code points")
	encodeCmd.Flags().BoolVarP(&outputHexEncodeCmd, "hex", "x", false, "output hexadecimal numbers")
	encodeCmd.Flags().StringVarP(&prefixEncodeCmd, "prefix", "", "", "add prefix to every two hex digits")
	rootCmd.AddCommand(encodeCmd)
}

func runEncodeCmd(cmd *cobra.Command, args []string) error {
	result := bytes.NewBuffer([]byte{})
	// utf8 uses up to 4 bytes
	result.Grow(len(args) * 4)
	buf := [4]byte{}
	if inputCodePointEncodeCmd {
		// input is sequence of Unicode code points
		for _, str := range args {
			if str[:2] == "U+" {
				str = str[2:]
			}
			codepoint, err := strconv.ParseUint(str, 16, 32)
			if err != nil {
				return err
			}
			n := utf8.EncodeRune(buf[:], rune(codepoint))
			result.Write(buf[:n])
		}
		printBytes(result.Bytes())
	} else {
		// input is a string
		for _, str := range args {
			for len(str) > 0 {
				r, size := utf8.DecodeRuneInString(str)
				n := utf8.EncodeRune(buf[:], r)
				result.Write(buf[:n])
				str = str[size:]
			}
			printBytes(result.Bytes())
			result.Reset()
		}
	}
	return nil
}

func printBytes(buff []byte) {
	space := " "
	if removeSpaceEncodeCmd {
		space = ""
	}
	for _, b := range buff {
		if outputHexEncodeCmd {
			fmt.Printf("%s%X%s", prefixEncodeCmd, b, space)
		} else {
			fmt.Printf("%d%s", b, space)
		}
	}
	fmt.Printf("\n")
}
