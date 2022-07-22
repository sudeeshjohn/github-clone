package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v45/github"
	"strings"
)

var allIssues []string

func createLabels(ctx context.Context, client *github.Client, sourceOrg string, sourceRepoSit string, targetOrg string, targetRepoSit string) ([]string, string, error) {
	var message string
	var allLabels []string
	var pageNumber int
	pageNumber = 1

	lstopt := &github.ListOptions{
		Page:    pageNumber,
		PerPage: 100,
	}
	opts := &github.IssueListByRepoOptions{ListOptions: *lstopt}
	for {

		labels, resp, err := client.Issues.ListLabels(ctx, sourceOrg, sourceRepoSit, lstopt)
		if err != nil {
			fmt.Errorf("Error: %s", err)
		}
		for _, issu := range labels {

			createLabel(ctx, client, targetOrg, targetRepoSit, issu)
		}
		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	return allLabels, message, nil
}

func createMilestones(ctx context.Context, client *github.Client, sourceOrg string, sourceRepoSit string, targetOrg string, targetRepoSit string) ([]string, string, error) {
	var message string
	var allLabels []string
	var pageNumber int
	var milestoneList []*github.Milestone
	pageNumber = 1

	lstopt := &github.ListOptions{
		Page:    pageNumber,
		PerPage: 100,
	}
	opts := &github.MilestoneListOptions{ListOptions: *lstopt, State: "all"}
	for {
		milestones, resp, err := client.Issues.ListMilestones(ctx, sourceOrg, sourceRepoSit, opts)
		if err != nil {
			fmt.Errorf("Error: %s", err)
		}
		if len(milestones) > 0 {
			for _, miles := range milestones {
				milestoneList = append(milestoneList, miles)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	previosMilestoneNumber := 0
	for i, _ := range milestoneList {
		miles := milestoneList[len(milestoneList)-i-1]
		fmt.Printf("Milestone: %s, number: %d\n", miles.GetTitle(), miles.GetNumber())
		for j := miles.GetNumber(); previosMilestoneNumber+1 < miles.GetNumber(); previosMilestoneNumber++ {
			title := "Deleted MileStone"
			state := "closed"
			description := "THis is deleted milestone in the original repo"
			mile := &github.Milestone{
				State:       &state,
				Title:       &title,
				Description: &description,
			}
			createMilestone(ctx, client, targetOrg, targetRepoSit, mile)
			j++
		}
		mile := &github.Milestone{
			State:       miles.State,
			Title:       miles.Title,
			Description: miles.Description,
		}
		createMilestone(ctx, client, targetOrg, targetRepoSit, mile)
		previosMilestoneNumber = miles.GetNumber()
	}
	return allLabels, message, nil
}

func createIssues(ctx context.Context, client *github.Client, sourceOrg string, sourceRepoSit string, targetOrg string, targetRepoSit string) ([]string, string, error) {
	var message string
	var allIssues []string
	var issueList []*github.Issue
	issueReq := &github.IssueRequest{}
	var pageNumber int
	pageNumber = 1
	lstopt := &github.ListOptions{
		Page:    pageNumber,
		PerPage: 100,
	}
	opts := &github.IssueListByRepoOptions{ListOptions: *lstopt, State: "all"}
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, sourceOrg, sourceRepoSit, opts)
		if err != nil {
			fmt.Errorf("Error: %s\n", err)
		}
		if len(issues) > 0 {
			for _, issu := range issues {
				issueList = append(issueList, issu)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.ListOptions.Page = resp.NextPage
	}
	previosIssueNumber := 0
	for i, _ := range issueList {
		issu := issueList[len(issueList)-i-1]
		fmt.Printf("Issue: %s, number: %d\n", issu.GetTitle(), issu.GetNumber())
		var labls *[]string
		var lab []string
		var assignees *[]string
		var assign []string
		if len(issu.Labels) > 0 {
			for _, la := range issu.Labels {
				if la != nil {
					lab = append(lab, *la.Name)
				}
			}
			labls = &lab
		} else {
			labls = nil
		}
		if len(issu.Assignees) != 0 {
			for _, ass := range issu.Assignees {
				if ass != nil {
					assign = append(assign, ass.GetLogin())
				}
			}
			assignees = &assign
		} else {
			labls = nil
		}

		for j := issu.GetNumber(); previosIssueNumber+1 < issu.GetNumber(); previosIssueNumber++ {
			title := "Deleted Issue"
			body := "This issue has been deleted in the original repo"
			issueReq = &github.IssueRequest{
				Title: &title,
				Body:  &body,
			}
			as, stat, _ := createIssue(ctx, client, targetOrg, targetRepoSit, issueReq, issu.GetNumber())
			if stat {
				state := "closed"
				issueReq.State = &state
				client.Issues.Edit(ctx, targetOrg, targetRepoSit, *as.Number, issueReq)
				j++
			}
		}
		issueReq = &github.IssueRequest{
			Title:     issu.Title,
			Body:      issu.Body,
			Labels:    nil,
			Assignee:  nil,
			Milestone: nil,
			Assignees: nil,
		}
		if len(issu.Labels) > 0 {
			issueReq.Labels = labls
		}
		if len(issu.Assignees) != 0 {
			fmt.Printf("Assignees:%s\n", *assignees)
			issueReq.Assignees = assignees
		}
		if issu.Milestone.GetNumber() != 0 {
			issueReq.Milestone = issu.Milestone.Number
			fmt.Printf("Milestone Number:%d, Milestone Title: %s\n", issu.Milestone.GetNumber(), issu.Milestone.GetTitle())
		}
		as, stat, err := createIssue(ctx, client, targetOrg, targetRepoSit, issueReq, issu.GetNumber())
		if stat {
			if *issu.State == "closed" {
				issueReq.State = issu.State
				client.Issues.Edit(ctx, targetOrg, targetRepoSit, *as.Number, issueReq)
			}
		} else {
			if strings.Contains(err.Error(), "alreay there is an issue in the same number") {
				fmt.Printf("In target repository, %s alreay there is an issue in the same number %d, with different title %s \n", targetRepoSit, as.GetNumber(), as.GetTitle())
				continue
			}
		}
		previosIssueNumber = issu.GetNumber()
	}
	return allIssues, message, nil
}

func createLabel(ctx context.Context, client *github.Client, targetOrg string, targetRepoSit string, label *github.Label) (bool, error) {
	_, rsp, err := client.Issues.CreateLabel(ctx, targetOrg, targetRepoSit, label)
	if err != nil {
		fmt.Errorf("error from create label: %s\n", err)
	}
	if rsp.StatusCode == 201 {
		fmt.Printf("Label %s Created\n", label.GetName())
		return true, nil
	}
	if rsp.StatusCode == 422 {
		if strings.Contains(err.Error(), "already_exists") {
			fmt.Printf("Label %s already exists\n", label.GetName())
			return true, nil
		}
	}
	return false, err
}
func createMilestone(ctx context.Context, client *github.Client, targetOrg string, targetRepoSit string, mile *github.Milestone) (bool, error) {
	_, rsp, err := client.Issues.CreateMilestone(ctx, targetOrg, targetRepoSit, mile)
	if err != nil {
		fmt.Errorf("error from create milestone: %s\n", err)
	}
	if rsp.StatusCode == 201 {
		fmt.Printf("Milestone %s Created\n", mile.GetTitle())
		return true, nil
	}
	if rsp.StatusCode == 422 {
		if strings.Contains(err.Error(), "already_exists\n") {
			fmt.Printf("milestone %s already exists\n", mile.GetTitle())
			return true, nil
		}
	}
	return false, err
}
func createIssue(ctx context.Context, client *github.Client, targetOrg string, targetRepoSit string, req *github.IssueRequest, num int) (*github.Issue, bool, error) {
	iss, _, err := client.Issues.Get(ctx, targetOrg, targetRepoSit, num)
	if err != nil {
		fmt.Errorf("Error from issue get: %s\n", err)
	}
	if len(iss.GetTitle()) != 0 {
		if iss.GetTitle() == req.GetTitle() {
			return iss, true, nil
		} else if iss.GetTitle() != req.GetTitle() {
			return iss, false, fmt.Errorf("In target repository, %s alreay there is an issue with the same number %d, with different title %s \n", targetRepoSit, iss.GetNumber(), iss.GetTitle())
		}
	}
	issue, _, err := client.Issues.Create(ctx, targetOrg, targetRepoSit, req)
	if err != nil {
		return nil, false, fmt.Errorf("Error from issue create: %s\n", err)
	}
	return issue, true, nil
}
