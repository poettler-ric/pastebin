package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	// read config
	viper.SetConfigName("pastebin")
	viper.AddConfigPath("$HOME/.pastebin")
	viper.ReadInConfig()

	pflag.StringP("token", "t", "", "github api token to use")
	pflag.BoolP("public", "p", false, "create public gist")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// create client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("token")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// create gist
	gist := new(github.Gist)
	gist.Files = make(map[github.GistFilename]github.GistFile)
	public := viper.GetBool("public")
	gist.Public = &public

	// fill gist
	if files := pflag.Args(); len(files) > 0 {
		// post files
		for _, f := range files {
			if strings.ContainsRune(f, '/') {
				log.Fatalln("file must not be in a subdirectory")
			}
			bytes, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatalf("error while reading %s: %v", f, err)
			}
			content := string(bytes)
			gist.Files[github.GistFilename(f)] = github.GistFile{
				Content: &content,
			}
		}
	} else {
		// post stdin
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error while reading stdin: %v", err)
		}

		content := string(bytes)
		gist.Files["stdin"] = github.GistFile{
			Content: &content,
		}
	}

	// post gist
	gist, _, err := client.Gists.Create(ctx, gist)
	if err != nil {
		log.Fatalf("error while creating gist: %v", err)
	}
	fmt.Println(*gist.HTMLURL)
}
