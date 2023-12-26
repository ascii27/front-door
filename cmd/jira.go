/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var jql = make(map[string]string)

// jiraCmd represents the jira command
var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {

		projects := viper.GetStringSlice("jira.projects")

		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %w", err))
		}

		tp := jira.BasicAuthTransport{
			Username: viper.GetString("jira.username"),
			Password: viper.GetString("jira.apitoken"),
		}

		jiraClient, err := jira.NewClient(tp.Client(), viper.GetString("jira.host"))
		if err != nil {
			panic(err)
		}

		query := fmt.Sprintf(jql["myTeamCurrentSprint"], strings.Join(projects, ","))

		issues := new([]jira.Issue)
		*issues, _, err = jiraClient.Issue.Search(
			query,
			nil)
		if err != nil {
			panic(err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Key", "Assignee", "Priority", "Status", "Summary"})
		t.SetColumnConfigs([]table.ColumnConfig{
			{
				Name:     "Summary",
				WidthMax: 120,
			},
		})

		currentProject := ""

		for _, issue := range *issues {

			assignee := "Unassigned"
			if issue.Fields.Assignee != nil {
				assignee = issue.Fields.Assignee.DisplayName
			}

			if currentProject == "" {
				currentProject = issue.Fields.Project.Name
			}

			if currentProject != issue.Fields.Project.Name {
				t.AppendSeparator()
				currentProject = issue.Fields.Project.Name
			}

			t.AppendRows([]table.Row{{
				issue.Key,
				assignee,
				issue.Fields.Status.Name,
				issue.Fields.Priority.Name,
				issue.Fields.Summary}})
		}
		t.Render()
	},
}

func init() {

	jql["myIssues"] = "assignee=currentUser() ORDER BY priority"
	jql["myTeamCurrentSprint"] = "project in (%s) AND Sprint in openSprints() AND Sprint not in futureSprints() ORDER BY project, status, created DESC"

	listCmd.AddCommand(jiraCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jiraCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jiraCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
