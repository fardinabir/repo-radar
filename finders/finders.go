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
	"repo-radar/models"
	"strings"
	"sync"
	"time"
)

func FindRepos(keywords []string) {

	timeStart := time.Now()

	ctx := context.Background()
	name := strings.Join(keywords, " ")
	site := "github.com"

	result, _ := googlesearch.Search(ctx, name+" site:"+site)
	//fmt.Println("Google Search Ended...Total... ", len(result))
	repoUrls := processSearchResult(&result)
	prepareRow(repoUrls)

	fmt.Println("-------Total Exec Time : --------- : ", time.Now().Sub(timeStart))
}

func processSearchResult(result *[]googlesearch.Result) []string {
	var repoUrls []string
	uniqueUrls := make(map[string]bool)
	for i := 0; i < len(*result); i++ {
		url := (*result)[i].URL
		repoUserSuffix := trimRepoUser(url[19:len(url)])
		if uniqueUrls[repoUserSuffix] {
			continue
		}
		newApiUrl := "https://api.github.com/repos/" + repoUserSuffix
		repoUrls = append(repoUrls, newApiUrl)
		uniqueUrls[repoUserSuffix] = true
	}
	return repoUrls
}

func prepareRow(repoUrls []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Repo URL", "Description", "Language", "Updated At", "Stars", "Forks", "Subscribers"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	//table.SetAutoFormatHeaders(false)
	//table.SetAutoWrapText(false)
	//fmt.Println("Search Started, Preparing Table")
	var wg sync.WaitGroup
	for _, newApiUrl := range repoUrls {
		wg.Add(1)
		i := newApiUrl

		go func() {
			defer wg.Done()
			GetRepoDetails(i, table)
		}()
	}
	wg.Wait()

	//fmt.Println("Search Fininshed, Rendering Table")
	table.Render()
}

func GetRepoDetails(newApiUrl string, table *tablewriter.Table) {
	resp, err := http.Get(newApiUrl)
	if err != nil {
		fmt.Println("Error while Getting Repos Info")
		return
	}
	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		//fmt.Println(resp.StatusCode, "---------------GotUnexpectedResult-------------------", resp)
	}
	if resp.StatusCode == 200 {
		var respRepo models.RepoDetails
		err = json.Unmarshal(body, &respRepo)
		if err != nil {
			fmt.Println("Error while unmarshalling")
			return
		}
		if respRepo.FullName != "" {
			shortRepoDetails := respRepo.GetShortDetails()
			table.Append(shortRepoDetails)
			//fmt.Println(" https://github.com/"+respRepo.FullName, "\t>>>>>>>>>>>>>>>\t", respRepo.StargazersCount)
		}
	}
}
