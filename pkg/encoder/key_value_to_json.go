package encoder

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	inputTextFile   string
	outputJsonFile1 string
	printOutput     bool
)

// textToJsonCmd represents the aa command
var keyValueToJSONCmd = &cobra.Command{
	Use:   "kvtj [flags]",
	Short: "Converts Key-Value (text) to JSON.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Read the input file
		content, err := os.ReadFile(inputTextFile)
		if err != nil {
			return fmt.Errorf("read key-value file: %w", err)
		}

		// Check if the input file is empty
		if len(content) == 0 {
			return fmt.Errorf("input file is empty")
		}

		// Convert the input file to JSON
		entries := strings.Split(string(content), "\n")
		m := make(map[string]string)
		for _, e := range entries {

			// Skip empty lines, comments and lines without "="
			if e == "" {
				continue
			}

			// Skip comments
			if strings.HasPrefix(e, "#") || strings.HasPrefix(e, "//") {
				continue
			}

			// Skip lines without "="
			if !strings.Contains(e, "=") {
				continue
			}

			// Split the line by "="
			parts := strings.Split(e, "=")

			var value string

			// If the value contains "=" then join the parts otherwise take the second part
			if len(parts) > 2 {
				value = strings.Join(parts[1:], "=")
			} else {
				value = parts[1]
			}

			// Remove " and ' from the value string
			value = strings.Trim(value, "\"")
			value = strings.Trim(value, "'")

			// Add the key-value pair to the map
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(value)
		}
		jsonString, _ := json.MarshalIndent(m, "", "  ")

		// Print the output to the console
		if printOutput {
			fmt.Println(string(jsonString))
			return nil
		}

		if outputJsonFile1 == "" {
			outputJsonFile1 = "output.json"
		}

		// Write the output file
		if err := os.WriteFile(outputJsonFile1, jsonString, 0644); err != nil {
			return fmt.Errorf("write JSON output: %w", err)
		}

		fmt.Println("Operation completed successfully. Check the", outputJsonFile1, "file.")
		return nil
	},
}

func init() {

	// Flags for the TTJ command
	keyValueToJSONCmd.Flags().StringVarP(&inputTextFile, "file", "f", "", "Input the text file name. Eg: keys.txt or .env")
	_ = keyValueToJSONCmd.MarkFlagRequired("file")

	keyValueToJSONCmd.Flags().StringVarP(&outputJsonFile1, "output", "o", "", "Output JSON file name (default is output.json)")
	keyValueToJSONCmd.Flags().BoolVarP(&printOutput, "print", "p", false, "Print the output to the console")
}

// KeyValueToJSONCommand exposes the converter as an encode subcommand.
func KeyValueToJSONCommand() *cobra.Command {
	return keyValueToJSONCmd
}
