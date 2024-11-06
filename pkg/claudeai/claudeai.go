package claudeai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Body struct {
	Model     string     `json:"model"`
	MaxTokens int        `json:"max_tokens"`
	Messages  []*Message `json:"messages"`
}

// https://docs.anthropic.com/en/api/messages
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"` // suports other types but only string now
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"` // suports other types but only string now
}

type Response struct {
	ID string `json:"id"`

	Type    string     `json:"type"`
	Role    string     `json:"role"`
	Model   string     `json:"model"`
	Content []*Content `json:"content"`

	Usage Usage `json:"usage"`
}
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

var Cmd = &cobra.Command{
	Use: "claudeai",
	RunE: func(cmd *cobra.Command, args []string) error {
		keyb, err := os.ReadFile(cfgPath)
		if err != nil {
			return err
		}

		for {
			prompt := promptui.Prompt{
				Label:     "Ask>",
				AllowEdit: true,
			}
			val, err := prompt.Run()
			if err != nil {
				return err
			}

			bdy := &Body{
				Model:     "claude-3-5-sonnet-20241022",
				MaxTokens: 1024,
				Messages: []*Message{
					{
						Role:    "user",
						Content: val,
					},
				},
			}

			b, err := json.Marshal(bdy)
			if err != nil {
				return err
			}

			req, err := http.NewRequest(
				http.MethodPost,
				"https://api.anthropic.com/v1/messages",
				bytes.NewReader(b),
			)
			if err != nil {
				return err
			}

			req.Header.Add("x-api-key", string(keyb))
			req.Header.Add("anthropic-version", "2023-06-01")
			req.Header.Add("content-type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}

			rbdy, err := io.ReadAll(resp.Body)
			if resp.StatusCode != 200 {
				if err != nil {
					return err
				}
				return fmt.Errorf("received non-200 response %d. body %s", resp.StatusCode, string(rbdy))
			}

			rs := Response{}
			err = json.Unmarshal(rbdy, &rs)
			if err != nil {
				return err
			}

			for _, c := range rs.Content {
				fmt.Println(c.Text)
			}

			log.Printf("Usage %d %d\n", rs.Usage.InputTokens, rs.Usage.OutputTokens)
		}
	},
}

func init() {
	Cmd.AddCommand(configureCmd)
}
