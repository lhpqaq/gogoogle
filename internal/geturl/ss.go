package geturl

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const ssURLFile = ".ss_config"
const salt = "1435660288"

var cmdHeader = "open"

type Flag struct {
	ShortArg  string
	Keyword   string
	Content   []string
	Debug     bool
	Translate bool
	ChatGPT   bool
	Cheat     bool
	Wait      bool
}

type SSWeb struct {
	Name     string `json:"name"`
	ShortArg string `json:"shortArg"`
	Search   string `json:"search"`
	Url      string `json:"url"`
	Delim    string `json:"delim"`
	Appid    string `json:"appid"`
	Key      string `json:"key"`
}

type TranslationResult struct {
	From        string `json:"from"`
	To          string `json:"to"`
	TransResult []struct {
		Src string `json:"src"`
		Dst string `json:"dst"`
	} `json:"trans_result"`
}

func md5V3(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

func getWebConfig() ([]SSWeb, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %v", err)
	}
	fileData, err := os.ReadFile(filepath.Join(currentUser.HomeDir, ssURLFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read ss_url file: %v", err)
	}

	var decodedWebs []SSWeb
	if err := json.Unmarshal(fileData, &decodedWebs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	return decodedWebs, nil
}

func GetWeb(f Flag) (SSWeb, []string, error) {
	decodedWebs, err := getWebConfig()
	if err != nil {
		return SSWeb{}, nil, err
	}

	for _, p := range decodedWebs {
		if len(f.ShortArg) > 0 && f.ShortArg == p.ShortArg {
			return p, append([]string{f.Keyword}, f.Content...), nil
		}
		if p.Name == f.Keyword {
			return p, f.Content, nil
		}
	}
	return decodedWebs[0], append([]string{f.Keyword}, f.Content...), nil
}

func ParseArgs(args []string) Flag {
	var f Flag
	for _, arg := range args[1:] {
		if arg[0] == '-' {
			if strings.Contains(arg, "d") {
				f.Debug = true
			} else if strings.Contains(arg, "t") {
				f.Translate = true
				f.ShortArg = "-t"
			} else if strings.Contains(arg, "c") {
				f.ChatGPT = true
				f.ShortArg = "-c"
			} else if strings.Contains(arg, "h") {
				f.Cheat = true
				f.ShortArg = "-h"
				f.Wait = true
				cmdHeader = "curl"
			} else {
				f.ShortArg = arg
			}
		} else {
			if len(f.Keyword) == 0 {
				f.Keyword = arg
			} else {
				f.Content = append(f.Content, arg)
			}
		}
	}

	return f
}

func ArrayToString(arr []string, delim string) string {
	var ret = ""
	if len(arr) > 0 {
		ret += arr[0]
	} else {
		return ret
	}
	for _, word := range arr[1:] {
		ret += delim
		ret += word
	}
	return ret
}

func getTranslateURL(web SSWeb, content string) string {
	string1 := fmt.Sprintf("%s%s%s%s", web.Appid, content, salt, web.Key)
	sign := md5V3(string1)
	encodedQuery := url.QueryEscape(content)
	return fmt.Sprintf(web.Url+web.Search, encodedQuery, "auto", "zh", web.Appid, salt, sign)
}

func GetURL(web SSWeb, content string) string {
	if web.ShortArg == "-t" {
		return getTranslateURL(web, content)
	}
	if len(web.Search) == 0 || len(content) == 0 {
		return web.Url
	}
	return fmt.Sprintf("%s%s", web.Url, fmt.Sprintf(web.Search, content))
}

func GetTranslate(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	var result TranslationResult
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
		fmt.Println("response body:", string(body))
		return
	}
	for _, trans := range result.TransResult {
		fmt.Printf("%s\n", trans.Dst)
	}
}

func Content2String(web SSWeb, content string) string {
	var ret string
	for _, s := range strings.Split(content, web.Delim) {
		fileInfo, err := os.Stat(s)
		isFile := false
		if err == nil && !fileInfo.IsDir() {
			file, err := os.Open(s)
			if err == nil {
				defer file.Close()
				data, err := io.ReadAll(file)
				if err == nil {
					ret += fmt.Sprintf("file name: %s\n", s)
					ret += "file content: \n```\n"
					ret += string(data)
					ret += "```\n"
					isFile = true
				}
			}
		}
		if !isFile {
			ret += fmt.Sprintf("%s\n", s)
		}
	}
	// fmt.Println(ret)
	return ret
}

func ChatGPT(web SSWeb, content string) {

	customConfig := openai.DefaultConfig(web.Key)
	customConfig.BaseURL = web.Url
	c := openai.NewClientWithConfig(customConfig)
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		// MaxTokens: 10,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: Content2String(web, content),
			},
		},
		Stream: true,
	}
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println()
			fmt.Println("-----------------Done-----------------")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}
		if len(response.Choices) > 0 {
			fmt.Printf(response.Choices[0].Delta.Content)
		}
	}
}
