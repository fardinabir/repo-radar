package finders

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	googlesearch "github.com/rocketlaunchr/google-search"
	"io"
	"net/http"
	"os"
	"star-seeker/models"
	"strings"
)

func FindRepos(keywords []string) {
	ctx := context.Background()
	name := strings.Join(keywords, " ")
	site := "github.com"

	result, _ := googlesearch.Search(ctx, name+" site:"+site)
	repoUrls := processSearchResult(&result)
	prepareRow(repoUrls)

}

func processSearchResult(result *[]googlesearch.Result) []string {
	var repoUrls []string
	uniqueUrls := make(map[string]bool)
	for i := 0; i < len(*result); i++ {
		url := (*result)[i].URL
		if uniqueUrls[url] {
			continue
		}
		repoUserSuffix := trimRepoUser(url[19:len(url)])
		//fmt.Println(url)
		newApiUrl := "https://api.github.com/repos/" + repoUserSuffix
		//fmt.Println(newApiUrl)
		repoUrls = append(repoUrls, newApiUrl)
		uniqueUrls[url] = true
	}
	return repoUrls
}

func prepareRow(repoUrls []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Repo URL", "Description", "Language", "Updated At", "Stars", "Forks", "Subscribers"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	//table.SetAutoFormatHeaders(false)
	//table.SetAutoWrapText(false)
	for _, newApiUrl := range repoUrls {
		resp, err := http.Get(newApiUrl)
		if err != nil {
			fmt.Println("Error in Getting Repo")
		}
		body, err := io.ReadAll(resp.Body)
		//fmt.Println(resp.StatusCode, "---------------------------------------------------------------------")
		if resp.StatusCode == 200 {
			var respRepo models.RepoDetails
			err = json.Unmarshal(body, &respRepo)
			if err != nil {
				fmt.Println("Error while unmarshalling")
			}
			if respRepo.FullName != "" {
				shortRepoDetails := respRepo.GetShortDetails()
				table.Append(shortRepoDetails)
				//fmt.Println(" https://github.com/"+respRepo.FullName, "\t>>>>>>>>>>>>>>>\t", respRepo.StargazersCount)
			}
		}
	}
	table.Render()
}

func formatRow(repo models.RepoDetailsShort) {

}
