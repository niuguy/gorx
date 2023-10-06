package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

const openAIEndpoint = "https://api.openai.com/v1/completions"

type OpenAIRequest struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
}

type OpenAIResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func generateRegex(modeFlag *ModeFlag, inputStrings []string, apiKey string) (string, error) {

	var prompt string
	// fmt.Println(*mode)
	if modeFlag.pattern {
		prompt = "Generate a regular expression pattern (without explaination) that can generally match the common part of the following strings:\n" + formatInputStrings(inputStrings) + "Pattern: /"
	} else if modeFlag.context {
		prompt = "Generate a regular expression pattern (without explaination) that generally match the string:" + inputStrings[0] + " from the context string:" + inputStrings[1] + "Pattern: /"
	} else {
		prompt = "Generate a regular expression (without explaination) of " + formatInputStrings(inputStrings) + "Pattern:"

	}

	request := OpenAIRequest{
		Prompt:      prompt,
		Model:       "gpt-3.5-turbo-instruct",
		Temperature: 0,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", openAIEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	regexPattern := response.Choices[0].Text // Assuming the best choice is at index 0
	return regexPattern, nil
}

func formatInputStrings(inputStrings []string) string {
	var formattedStrings string
	for _, s := range inputStrings {
		formattedStrings += s + " "
	}

	return formattedStrings
}

type ModeFlag struct {
	pattern  bool
	context  bool
	semantic bool
}

func Execute() {

	var modelFlag ModeFlag

	var rgxCmd = &cobra.Command{
		Use:           "gorx [flags] [input strings]",
		Short:         "Generate a regular expression pattern using OpenAI's GPT-3 API ",
		SilenceUsage:  true,
		SilenceErrors: true,

		Run: func(cmd *cobra.Command, args []string) {
			apiKey := os.Getenv("OPENAI_API_KEY")

			if apiKey == "" {
				fmt.Println("OPENAI_API_KEY is not set.")
				fmt.Print("Please set the API_OPENAI_API_KEYKEY: ")
				var input string
				fmt.Scanln(&input)

				// Set the API_KEY in the environment
				os.Setenv("OPENAI_API_KEY", input)
				fmt.Println("OPENAI_API_KEY is now set.")
			}

			var inputStrings []string

			// take input strings from command line
			if len(os.Args) > 2 {
				inputStrings = os.Args[2:]
			} else {
				cmd.Help()
				return
			}

			regexPattern, err := generateRegex(&modelFlag, inputStrings, apiKey)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(regexPattern)
		},
	}

	rgxCmd.Flags().BoolVarP(&modelFlag.pattern, "pattern", "p", false, "Find common patterns of the given strings separated by empty space")
	rgxCmd.Flags().BoolVarP(&modelFlag.context, "context", "c", false, "Math the first input string from the second input string")
	rgxCmd.Flags().BoolVarP(&modelFlag.semantic, "semantic", "s", false, "Generate a regular expression based on semantic meaning, e.g. 'email address'")
	rgxCmd.Execute()
}
