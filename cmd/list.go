/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	jira "github.com/andygrunwald/go-jira"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	// Jira token:

	Run: func(cmd *cobra.Command, args []string) {


		tp := jira.BasicAuthTransport{
			Username: "",
			Password: "",
		}

		//jiraClient, err := jira.NewClient(tp.Client(), "https://nitebritestudio.atlassian.net")
		jiraClient, err := jira.NewClient(tp.Client(), "https://hashicorp.atlassian.net")
		if err != nil {
			panic(err)
		}

		issues := new([]jira.Issue)
		*issues, _, err = jiraClient.Issue.Search("assignee=currentUser() ORDER BY priority", nil)
		if err != nil {
			panic(err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Key", "Priority", "Summary"})
		t.SetColumnConfigs([]table.ColumnConfig{
			{
				Name:     "Summary",
				WidthMax: 120,
			},
		})
		for _, issue := range *issues {

			t.AppendRows([]table.Row{{issue.Key, issue.Fields.Priority.Name, issue.Fields.Summary}})
		}
		t.Render()

		/*
			issue, _, err := jiraClient.Issue.Get("PLAT-1372", nil)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)
			fmt.Printf("Type: %s\n", issue.Fields.Type.Name)
			fmt.Printf("Priority: %s\n", issue.Fields.Priority.Name)
		*/
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
