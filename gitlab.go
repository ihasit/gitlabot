package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

var GitEmojiMap = map[string]string{
	":bulb:":                      "💡",
	":heavy_minus_sign:":          "➖",
	":bug:":                       "🐛",
	":art:":                       "🎨",
	":hammer:":                    "🔨",
	":sparkles:":                  "✨",
	":building_construction:":     "🏗️",
	":wrench:":                    "🔧",
	":triangular_flag_on_post:":   "🚩",
	":arrow_down:":                "⬇️",
	":label:":                     "🏷️",
	":dizzy:":                     "💫",
	":white_check_mark:":          "✅",
	":mag:":                       "🔍️",
	":bento:":                     "🍱",
	":chart_with_upwards_trend:":  "📈",
	":beers:":                     "🍻",
	":boom:":                      "💥",
	":bookmark:":                  "🔖",
	":monocle_face:":              "🧐",
	":recycle:":                   "♻️",
	":card_file_box:":             "🗃️",
	":globe_with_meridians:":      "🌐",
	":adhesive_bandage:":          "🩹",
	":pushpin:":                   "📌",
	":iphone:":                    "📱",
	":test_tube:":                 "🧪",
	":page_facing_up:":            "📄",
	":alien:":                     "👽️",
	":children_crossing:":         "🚸",
	":poop:":                      "💩",
	":heavy_plus_sign:":           "➕",
	":necktie:":                   "👔",
	":rotating_light:":            "🚨",
	":memo:":                      "📝",
	":loud_sound:":                "🔊",
	":construction:":              "🚧",
	":fire:":                      "🔥",
	":zap:":                       "⚡️",
	":stethoscope:":               "🩺",
	":package:":                   "📦️",
	":camera_flash:":              "📸",
	":lipstick:":                  "💄",
	":mute:":                      "🔇",
	":rocket:":                    "🚀",
	":lock:":                      "🔒️",
	":ambulance:":                 "🚑️",
	":pencil2:":                   "✏️",
	":arrow_up:":                  "⬆️",
	":clown_face:":                "🤡",
	":truck:":                     "🚚",
	":goal_net:":                  "🥅",
	":egg:":                       "🥚",
	":speech_balloon:":            "💬",
	":construction_worker:":       "👷",
	":passport_control:":          "🛂",
	":rewind:":                    "⏪️",
	":wheelchair:":                "♿️",
	":alembic:":                   "⚗️",
	":seedling:":                  "🌱",
	":green_heart:":               "💚",
	":tada:":                      "🎉",
	":busts_in_silhouette:":       "👥",
	":twisted_rightwards_arrows:": "🔀",
	":wastebasket:":               "🗑️",
	":coffin:":                    "⚰️",
	":see_no_evil:":               "🙈",
}

func trans2Emoji(content string) string {
	for k, v := range GitEmojiMap {
		content = strings.ReplaceAll(content, k, v)
	}
	return content
}

func NewClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

type WxResp struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// Push events
type PushBody struct {
	ObjectKind string     `json:"object_kind"`
	Ref        string     `json:"ref"`
	Commits    []Commit   `json:commits`
	Repository Repository `json:"repository"`
	After      string     `json:"after"`
	UserName   string     `json:"user_name"`
}

// TagPushBody Tag events
type TagPushBody struct {
	ObjectKind string     `json:"object_kind"`
	EventName  string     `json:"event_name"`
	Before     string     `json:"before"`
	After      string     `json:"after"`
	Ref        string     `json:"ref"`
	Commits    []Commit   `json:"commits"`
	Repository Repository `json:"repository"`
	UserName   string     `json:"user_name"`
}

// IssuePushBody Issues events
type IssuePushBody struct {
	User             IssueUser   `json:"user"`
	Repository       Repository  `json:"repository"`
	ObjectAttributes IssueObject `json:"object_attributes"`
}

// CommentPushBody comment
type CommentPushBody struct {
	User             IssueUser     `json:"user"`
	Repository       Repository    `json:"repository"`
	ObjectAttributes CommentObject `json:"object_attributes"`
}

type CommentObject struct {
	Id        int64  `json:"id"`
	Note      string `json:"note"`
	UpdatedAt string `json:"updated_at"`
	Url       string `json:"url"`
}

// MRPushBody
type MRPushBody struct {
	User             IssueUser  `json:"user"`
	Repository       Repository `json:"repository"`
	ObjectAttributes MRObjects  `json:"object_attributes"`
}

// PipelineBody
type PipelineBody struct {
	ObjectAttributes PipelineObject `json:"object_attributes"`
	User             IssueUser      `json:"user"`
	Project          Project        `json:"project"`
}

type PipelineObject struct {
	Id         int64  `json:"id"`
	Ref        string `json:"ref"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
	FinishedAt string `json:"finished_at"`
	Duration   int64  `json:"duration"`
	Tag        bool   `json:"tag"`
}

type MRObjects struct {
	Id           int64  `json:"id"`
	TargetBranch string `json:"target_branch"`
	SourceBranch string `json:"source_branch"`
	UpdatedAt    string `json:"updated_at"`
	Url          string `json:"url"`
	Action       string `json:"action"`
}

type IssueUser struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
}

type IssueObject struct {
	Id     int64  `json:"id"`
	Title  string `json:"title"`
	Url    string `jso:"url"`
	Action string `json:"action"`
}

type Commit struct {
	Id        string `json:"id"`
	Message   string `json:"message"`
	TimeStamp string `json:"timestamp"`
	Url       string `json:"url"`
	Author    Author `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Repository struct {
	Name      string `json:"name"`
	HomePage  string `json:"homepage"`
	GitSSHUrl string `json:"git_ssh_url"`
}

type Project struct {
	Name      string `json:"name"`
	WebUrl    string `json:"web_url"`
	GitSSHUrl string `json:"git_ssh_url"`
}

// ReleaseBody Release events
type ReleaseBody struct {
	ObjectKind string     `json:"object_kind"`
	Name       string     `json:"name"`
	Description string    `json:"description"`
	CreatedAt  string     `json:"created_at"`
	Url        string     `json:"url"`
	Assets     Assets     `json:"assets"`
	Project    Project    `json:"project"`
}

type Assets struct {
	Count int      `json:"count"`
	Links []Link   `json:"links"`
}

type Link struct {
	Id       int    `json:"id"`
	LinkType string `json:"link_type"`
	Name     string `json:"name"`
	Url      string `json:"url"`
}

func bindJson(ctx *gin.Context, m interface{}) error {
	err := ctx.BindJSON(m)
	if err != nil {
		ctx.JSON(400, WxResp{ErrCode: 400, ErrMsg: fmt.Sprintf("Parse gitlab requset body error: %s", err)})
		return err
	}
	return nil
}

func buildMsg(content string, markdown bool) string {
	if markdown {
		return fmt.Sprintf(`{"msgtype": "markdown", "markdown":{"content": "%s"}}`, content)
	}
	return fmt.Sprintf(`{"msgtype": "text", "text":{"content": "%s"}}`, content)
}

func TransmitRobot(ctx *gin.Context) {
	key := ctx.GetHeader("X-Gitlab-Token")
	if len(key) == 0 {
		ctx.Render(403, render.Data{ContentType: "application/json", Data: []byte("X-Gitlab-Token is empty")})
		return
	}
	requestUrl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", key)
	var resp *http.Response
	var wxErr error
	var content string
	pushEvent := ctx.GetHeader("X-Gitlab-Event")
	if pushEvent == "Push Hook" {
		pushBody := &PushBody{}
		if err := bindJson(ctx, pushBody); err != nil {
			return
		}
		if len(pushBody.Commits) == 0 && pushBody.After != "0000000000000000000000000000000000000000" {
			ctx.JSON(200, &WxResp{ErrCode: 0, ErrMsg: "no commit"})
			return
		}
		content = "# " + pushBody.Repository.Name + "\n"
		content += "### On branch `" + pushBody.Ref + "`\n"
		if len(pushBody.Commits) > 0 {
			v := pushBody.Commits[len(pushBody.Commits)-1]
			content += fmt.Sprintf("%s push a commit [%s](%s)  %s", v.Author.Name, strings.ReplaceAll(v.Message, "\n", ""), v.Url, v.TimeStamp) + "\n"
		}
		if pushBody.After == "0000000000000000000000000000000000000000" {
			content += fmt.Sprintf("%s `remove` it", pushBody.UserName)
		}
	} else if pushEvent == "Tag Push Hook" {
		tagPushBody := &TagPushBody{}
		if err := bindJson(ctx, tagPushBody); err != nil {
			return
		}
		content = "# " + tagPushBody.Repository.Name + "\n"
		content += fmt.Sprintf("%s push a tag: [%s](%s)\n", tagPushBody.UserName, tagPushBody.Ref, tagPushBody.Repository.HomePage+strings.Replace(tagPushBody.Ref, "refs/tags/", "tags/", -1))
		if len(tagPushBody.Commits) > 0 {
			commit := tagPushBody.Commits[len(tagPushBody.Commits)-1]
			content += fmt.Sprintf("Last Commit: [%s](%s) by %s\n", commit.Id, commit.Url, commit.Author.Name)
		}
	} else if pushEvent == "Issue Hook" {
		issueBody := &IssuePushBody{}
		if err := bindJson(ctx, issueBody); err != nil {
			return
		}
		content = "# " + issueBody.Repository.Name + "\n"
		content += fmt.Sprintf("%s %s a issue [%s](%s)", issueBody.User.Name, issueBody.ObjectAttributes.Action, issueBody.ObjectAttributes.Title, issueBody.ObjectAttributes.Url)
	} else if pushEvent == "Note Hook" {
		commentBody := &CommentPushBody{}
		if err := bindJson(ctx, commentBody); err != nil {
			return
		}
		content = "# " + commentBody.Repository.Name + "\n"
		content += fmt.Sprintf("%s leave a comment: %s  %s \n[Detail>>](%s)", commentBody.User.Name, commentBody.ObjectAttributes.Note, commentBody.ObjectAttributes.UpdatedAt, commentBody.ObjectAttributes.Url)
	} else if pushEvent == "Merge Request Hook" {
		mrBody := &MRPushBody{}
		if err := bindJson(ctx, mrBody); err != nil {
			return
		}
		content = "# " + mrBody.Repository.Name + "\n"
		content += fmt.Sprintf("%s `%s` a merge request from `%s` to `%s` \n[Detail>>](%s)", mrBody.User.Name, mrBody.ObjectAttributes.Action, mrBody.ObjectAttributes.SourceBranch, mrBody.ObjectAttributes.TargetBranch, mrBody.ObjectAttributes.Url)
	} else if pushEvent == "Pipeline Hook" {
		pipelineBody := &PipelineBody{}
		if err := bindJson(ctx, pipelineBody); err != nil {
			return
		}
		content = "# " + pipelineBody.Project.Name + "\n"
		branch := "branch"
		if pipelineBody.ObjectAttributes.Tag {
			branch = "tag"
		}
		content += fmt.Sprintf("### Pipeline on %s `%s`\n", branch, pipelineBody.ObjectAttributes.Ref)
		status := ""
		if pipelineBody.ObjectAttributes.Status == "failed" {
			status = "🐛"
		} else if pipelineBody.ObjectAttributes.Status == "running" {
			status = "🚀"
		} else if pipelineBody.ObjectAttributes.Status == "success" {
			status = "✅"
		} else if pipelineBody.ObjectAttributes.Status == "pending" {
			status = "🔒"
		}
		if len(status) == 0 {
			ctx.JSON(200, WxResp{ErrCode: 0, ErrMsg: "unknown status: " + pipelineBody.ObjectAttributes.Status})
			return
		}
		content += "`Status`: " + status + "\n"
		content += fmt.Sprintf("`Start at`: %s\n", pipelineBody.ObjectAttributes.CreatedAt)
		if len(pipelineBody.ObjectAttributes.FinishedAt) > 0 {
			content += fmt.Sprintf("`Finish at`: %s\n", pipelineBody.ObjectAttributes.FinishedAt)
		}
		if pipelineBody.ObjectAttributes.Duration > 0 {
			content += fmt.Sprintf("`Duration`: %ds", pipelineBody.ObjectAttributes.Duration)
		}
	} else if pushEvent == "Release Hook" {
		releaseBody := &ReleaseBody{}
		if err := bindJson(ctx, releaseBody); err != nil {
			return
		}
		content = "# " + releaseBody.Project.Name + "\n"
		content += fmt.Sprintf("Release: **%s**\n", releaseBody.Name)
		content += fmt.Sprintf("Description: %s\n", releaseBody.Description)
		content += fmt.Sprintf("Created at: %s\n", releaseBody.CreatedAt)
		content += fmt.Sprintf("URL: [Release Link](%s)\n", releaseBody.Url)
		for _, link := range releaseBody.Assets.Links {
			content += fmt.Sprintf("Asset: [%s](%s)\n", link.Name, link.Url)
		}
	}
	if len(content) == 0 {
		ctx.JSON(200, WxResp{ErrCode: 0, ErrMsg: "no content"})
		return
	}
	content = trans2Emoji(content)
	data := []byte(buildMsg(content, true))
	client := NewClient()
	resp, wxErr = client.Post(requestUrl, "application/json", bytes.NewBuffer(data))
	defer resp.Body.Close()
	if wxErr != nil {
		ctx.JSON(500, WxResp{ErrCode: 500, ErrMsg: fmt.Sprintf("Request wexin robot err: %s ", wxErr)})
		return
	}
	wxResp := &WxResp{}
	json.NewDecoder(resp.Body).Decode(wxResp)
	ctx.JSON(200, wxResp)
}
