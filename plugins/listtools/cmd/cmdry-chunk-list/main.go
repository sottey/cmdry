package main

import (
	"fmt"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/listtools"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.chunk-list", Name: "Chunk List", Version: "0.1.0", Description: "Split a pasted newline-delimited list into fixed-size groups locally.", Category: "text", Icon: "list",
		SearchTerms: []string{"list", "chunk", "split", "batch", "groups", "lines"}, Pages: []cmdry.Page{{ID: "overview", Name: "Chunk", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New list", Method: "read"}, {ID: "chunk", Name: "Chunk list", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "chunk": chunk}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "Chunk List", Components: []cmdry.Component{{Type: "form", Title: "Split list into groups", Action: "chunk", Submit: "Chunk list", Description: "Uses one item per line and ignores blank lines. The output separates groups with one blank line.", Fields: []cmdry.Field{{Name: "input", Label: "List items", Type: "textarea", Required: true}, {Name: "size", Label: "Items per group", Type: "number", Value: "10", Min: "1", Max: "10000", Required: true}}}}}, nil
}

func chunk(request cmdry.Request) (cmdry.View, error) {
	size, err := strconv.Atoi(fmt.Sprint(request.Params["size"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("items per group must be a number")
	}
	chunks, err := listtools.ChunkLines(fmt.Sprint(request.Params["input"]), size)
	if err != nil {
		return cmdry.View{}, err
	}
	items := 0
	for _, group := range chunks {
		items += len(group)
	}
	return cmdry.View{Title: "Chunked list", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "List chunked", Message: fmt.Sprint(items) + " item(s) split into " + fmt.Sprint(len(chunks)) + " group(s)."}, {Type: "code", Title: "Chunked list", Text: listtools.FormatChunks(chunks)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Chunk another list"}}}}}, nil
}
