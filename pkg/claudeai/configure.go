package claudeai

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var cfgPath string

func init() {
	d, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	cfgPath = path.Join(d, "/.gochat", "claudeai")
}

var configureCmd = &cobra.Command{
	Use: "configure",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("claudeai configure called")
		log.Printf("cfg path %s\n", cfgPath)

		_, err := os.Stat(cfgPath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if os.IsNotExist(err) {
			log.Println("file exists")
			prompt := promptui.Prompt{
				Label: fmt.Sprintf("override %s? (Y/n)", cfgPath),
			}
			val, err := prompt.Run()
			if err != nil {
				return err
			}
			if val != "Y" {
				return nil
			}
		}

		prompt := promptui.Prompt{
			Label:       "Enter your API key",
			Mask:        '*',
			HideEntered: true,
		}
		val, err := prompt.Run()
		if err != nil {
			return err
		}

		err = os.MkdirAll(path.Dir(cfgPath), 0770)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(cfgPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}

		_, err = f.WriteString(val)
		if err != nil {
			return err
		}

		return err
	},
}
