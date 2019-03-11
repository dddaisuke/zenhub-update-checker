package main

import (
	"appengine"
	"appengine/memcache"
	// "appengine/taskqueue"
	"appengine/urlfetch"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
  HipchatRoomId     = "HipChat Room Id"
  HipchatApiKey     = "HipChat API Key"
  ZenhubAccessToken = "zenhub oauth2 access token"
	TimeFormat        = "2006-01-02T15:04:05.000Z"
	MLastCheckTime    = "last_check_time"
	MIssue            = "issue/"
)

type board struct {
	Id string `json:"_id"`
	CreatedBy string
	V uint64 `json:"__v"`
	LastModified string
	DefaultPrPipelineId string
	CreatedAt string
	Pipelines []pipeline
}

type pipeline struct {
	Name string
	Id string `json:"_id"`
	Issues []issue
}

type issue struct {
	RepoId uint64 `json:"repo_id"`
	IssueNumber uint64 `json:"issue_number"`
	Id string `json:"_id"`
}

func sendToHipchat(client *http.Client, roomId string, from string, message string, color string) {
	log.Printf("send to hipchat : %s", message)

	data := url.Values{
		"room_id":        {roomId},
		"from":           {from},
		"message":        {message},
		"message_format": {"html"},
		"notify":         {"1"},
		"color":          {color},
		"format":         {"json"},
		"auth_token":     {HipchatApiKey},
	}

	client.PostForm(
		"https://api.hipchat.com/v1/rooms/message",
		data,
	)
}

func fetchZenhubBoard(client *http.Client) board {
	url := fmt.Sprintf("https://api.zenhub.io/v3/repos/11913693/board")
	req, _ := http.NewRequest("GET", url, nil)
	// cookie := http.Cookie{"connect.sid", "s%3AFGcOZXVzXJQuG7fLzYbfpW8l.DX4%2BS6j9Imd9%2F%2FT0TcQNplMR%2BqRLxcVRQTMtGstcdEU", nil, nil, nil, nil, nil, true, true, "connect.sid=s%3AFGcOZXVzXJQuG7fLzYbfpW8l.DX4%2BS6j9Imd9%2F%2FT0TcQNplMR%2BqRLxcVRQTMtGstcdEU", []string{"connect.sid=s%3AFGcOZXVzXJQuG7fLzYbfpW8l.DX4%2BS6j9Imd9%2F%2FT0TcQNplMR%2BqRLxcVRQTMtGstcdEU"}}
  // req.AddCookie(&cookie)
	req.Header.Set("x-authentication-token", ZenhubAccessToken)
	// req.Header.Set("connect.sid", "s%3AFGcOZXVzXJQuG7fLzYbfpW8l.DX4%2BS6j9Imd9%2F%2FT0TcQNplMR%2BqRLxcVRQTMtGstcdEU")

	resp, _ := client.Do(req)
	byteArray, _ := ioutil.ReadAll(resp.Body)
	dec := json.NewDecoder(bytes.NewReader(byteArray))

	var d board
	dec.Decode(&d)
	log.Print(d)
	return d
}

func fetchZenhubBoard(client *http.Client) board {
	url := fmt.Sprintf("https://api.zenhub.io/v3/repos/11913693/board")
	req, _ := http.NewRequest("GET", url, nil)
	// cookie := http.Cookie{"connect.sid", "s%3AFGcOZXVzXJQuG7fLzYbfpW8l.DX4%2BS6j9Imd9%2F%2FT0TcQNplMR%2BqRLxcVRQTMtGstcdEU", nil, nil, nil, nil, nil, true, true, "connect.sid=s%3AFGcOZXVzXJQuG7fLzYbfpW8l.DX4%2BS6j9Imd9%2F%2FT0TcQNplMR%2BqRLxcVRQTMtGstcdEU", []string{"connect.sid=s%3AFGcOZXVzXJQuG7fLzYbfpW8l.DX4%2BS6j9Imd9%2F%2FT0TcQNplMR%2BqRLxcVRQTMtGstcdEU"}}
  // req.AddCookie(&cookie)
	req.Header.Set("x-authentication-token", ZenhubAccessToken)
	// req.Header.Set("connect.sid", "s%3AFGcOZXVzXJQuG7fLzYbfpW8l.DX4%2BS6j9Imd9%2F%2FT0TcQNplMR%2BqRLxcVRQTMtGstcdEU")

	resp, _ := client.Do(req)
	byteArray, _ := ioutil.ReadAll(resp.Body)
	dec := json.NewDecoder(bytes.NewReader(byteArray))

	var d board
	dec.Decode(&d)
	log.Print(d)
	return d
}



// func createMessage(value data) string {
// 	message := ""
// 	switch value.TypeName {
// 	case "createPlus":
// 		message = fmt.Sprintf("<a href='https://github.com/%s'>%s</a> が <a href='https://github.com/%s/%s/issues/%.0f'>%s/%.0f</a> に <b>+1</b> しました。\n", value.Actor.Github.Username, value.Actor.Github.Username, value.Organization, value.Repository, value.Issue, value.Repository, value.Issue)
// 	case "transferIssue":
// 		if value.SrcPipelineName == value.DestPipelineName {
// 			message = fmt.Sprintf("<a href='https://github.com/%s'>%s</a> が %s 内で <a href='https://github.com/%s/%s/issues/%.0f'>%s/%.0f</a> の <b>優先順位を変更</b> しました。\n", value.Actor.Github.Username, value.Actor.Github.Username, value.SrcPipelineName, value.Organization, value.Repository, value.Issue, value.Repository, value.Issue)
// 		} else {
// 			message = fmt.Sprintf("<a href='https://github.com/%s'>%s</a> が <a href='https://github.com/%s/%s/issues/%.0f'>%s/%.0f</a> を <b>%s</b> から <b>%s</b> に移動しました。\n", value.Actor.Github.Username, value.Actor.Github.Username, value.Organization, value.Repository, value.Issue, value.Repository, value.Issue, value.SrcPipelineName, value.DestPipelineName)
// 		}
// 	case "createBoard":
// 		message = fmt.Sprintf("<a href='https://github.com/%s'>%s</a> が <a href='https://github.com/%s/%s/issues/%.0f'>%s/%.0f</a> の <b>ボードを作成</b> しました。\n", value.Actor.Github.Username, value.Actor.Github.Username, value.Organization, value.Repository, value.Issue, value.Repository, value.Issue)
// 	default:
// 		message = fmt.Sprintf("[Unknown] type = %s\n", value.TypeName)
// 	}
// 	return message
// }

func stringToTime(ttt string) time.Time {
	t, err := time.Parse(TimeFormat, ttt)
	if err != nil {
		panic(err)
	}
	return t
}

func getLastCheckTime(c appengine.Context) string {
	item, err := memcache.Get(c, MLastCheckTime)
	if err == memcache.ErrCacheMiss {
		c.Infof("item not in the cache")
		return ""
	} else if err != nil {
		c.Errorf("error getting item: %v", err)
		return ""
	} else {
		c.Infof("the lyric is %q", item.Value)
		return string(item.Value)
	}
}

func setLastCheckTime(c appengine.Context, time time.Time) {
	item := &memcache.Item{
		Key:   MLastCheckTime,
		Value: []byte(time.Format(TimeFormat)),
	}
	if err := memcache.Set(c, item); err == memcache.ErrNotStored {
		c.Infof("item with key %q already exists", item.Key)
	} else if err != nil {
		c.Errorf("error adding item: %v", err)
	}
}

func isMovedPipeline(c appengine.Context, pipelineId string, issue issue) (bool, error) {
	item, err := memcache.Get(c, MIssue + issue.Id)
	if err == memcache.ErrCacheMiss {
		c.Infof("item not in the cache")
		setPipeline(c, pipelineId, issue)
		return false, err
	} else if err != nil {
		c.Errorf("error getting item: %v", err)
		return false, err
	} else {
		c.Infof("the lyric is %q", item.Value)
		isMoved := string(item.Value) != pipelineId
		return isMoved, nil
	}
}

func setPipeline(c appengine.Context, pipelineId string, issue issue) (bool, error) {
	item := &memcache.Item{
		Key:   MIssue + issue.Id,
		Value: []byte(pipelineId),
	}
	if err := memcache.Set(c, item); err == memcache.ErrNotStored {
		c.Infof("item with key %q already exists", item.Key)
		return false, err
	} else if err != nil {
		c.Errorf("error adding item: %v", err)
		return false, err
	}
	return true, nil
}

func checkZenhubMessage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	// var lastModifiedTime time.Time

	d := fetchZenhubBoard(client)
	log.Print("%d", d)
	// lastModifiedTimeString := d.LastModified
  // lastCheckTimeString := getLastCheckTime(c)
	// if lastModifiedTimeString == "" {
	// 	lastModifiedTime = time.Now().Add(-60 * time.Minute)
	// } else {
	// 	lastModifiedTime = stringToTime(lastModifiedTimeString)
	// }

	for _, pipeline := range d.Pipelines {
		if pipeline.Name == "In Stg-QA" {
			log.Print("-----")
			log.Print("%s : %s", pipeline.Name, pipeline.Id)
			for _, issue := range pipeline.Issues {
				isMoved, _ := isMovedPipeline(c, pipeline.Id, issue)
				if isMoved {
					log.Print("%s : %s", issue.RepoId, issue.IssueNumber)
					// message := fmt.Sprintf("<a href='https://github.com/%s'>%s</a> が <a href='https://github.com/%s/%s/issues/%.0f'>%s/%.0f</a> に <b>+1</b> しました。\n", value.Actor.Github.Username, value.Actor.Github.Username, value.Organization, value.Repository, value.Issue, value.Repository, value.Issue)
					message := fmt.Sprintf("%s", issue.Id)
					t := taskqueue.NewPOSTTask("/tasks/post_to_hipchat", map[string][]string{"message": {message}})
					t.Delay = count
					if _, err := taskqueue.Add(c, t, ""); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
			}
		}
	}

	// lastMessageTime := ""
	// var count time.Duration
	// count = 1
	// for _, value := range d {
	// 	log.Print("%s", value)
	// 	createdAt := value.CreatedAt
	// 	time := stringToTime(createdAt)
	// 	if !lastCheckTime.After(time) {
	// 		if time.Equal(lastCheckTime) {
	// 			continue
	// 		}
	// 		if lastMessageTime == "" {
	// 			lastMessageTime = value.CreatedAt
	// 		}
	// 		message := createMessage(value)

	// 		t := taskqueue.NewPOSTTask("/tasks/post_to_hipchat", map[string][]string{"message": {message}})
	// 		t.Delay = count
	// 		if _, err := taskqueue.Add(c, t, ""); err != nil {
	// 			http.Error(w, err.Error(), http.StatusInternalServerError)
	// 			return
	// 		}
	// 	}
	// 	count++
	// }
	// log.Printf("last_message_time: %s", lastMessageTime)
	// if lastMessageTime != "" {
	// 	setLastCheckTime(c, stringToTime(lastMessageTime))
	// }
}

func post_to_hipchat(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)

	message := r.FormValue("message")
	sendToHipchat(client, HipchatRoomId, "zenhub.io", message, "green")
}
func init() {
	http.HandleFunc("/zenhub", handler)
	http.HandleFunc("/tasks/post_to_hipchat", post_to_hipchat)
}

func handler(w http.ResponseWriter, r *http.Request) {
	checkZenhubMessage(w, r)
}
