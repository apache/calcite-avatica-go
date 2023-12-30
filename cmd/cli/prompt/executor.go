/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package prompt

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/apache/calcite-avatica-go/v5"
	"github.com/c-bata/go-prompt"
	calcitesql "github.com/apache/calcite-avatica-go/cmd/cli/calcitesql"
)

var isMultiline bool

func CreateAndRunPrompt(db *sql.DB) {
	fmt.Println("Welcome! Use SQL to query Apache Calcite.\nUse Ctrl+D, type \"exit\" or \"quit\" to exit.")
	fmt.Println()

	p := prompt.New(
		executeQueryWrapper(db),
		CustomCompleter,
		prompt.OptionLivePrefix(LivePrefix),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSuggestionBGColor(prompt.White),
		prompt.OptionSuggestionTextColor(prompt.Black),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionCompletionOnDown(),
		prompt.OptionTitle("Calcite CLI Prompt"),                 // Set a title for the prompt
		prompt.OptionInputTextColor(prompt.Fuchsia),              // Customize input text color
		prompt.OptionDescriptionTextColor(prompt.Black),          // Customize description text color
		prompt.OptionSelectedSuggestionTextColor(prompt.White),   // Customize selected suggestion text color
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray), // Customize selected suggestion background color
		prompt.OptionPrefix("calcite \U0001F48E:sql> "),          // Set a custom prefix for the prompt
	)

	p.Run()
}

func LivePrefix() (prefix string, useLivePrefix bool) {
	if isMultiline {
		prefix = "... "
		useLivePrefix = true
	} else {
		prefix = "calcite \U0001F48E:sql> "
		useLivePrefix = !isMultiline
	}
	return prefix, useLivePrefix
}

func executeQueryWrapper(db *sql.DB) func(string) {
	var multiLineQuery strings.Builder

	return func(query string) {
		// Check for exit command
		if strings.ToLower(query) == "exit" || strings.ToLower(query) == "quit" {
			fmt.Println("Exiting calcite CLI Prompt...")
			os.Exit(0)
		}

		trimmedQuery := strings.TrimSpace(query)

		// Check if it is a multiline query
		if strings.HasSuffix(trimmedQuery, ";") {
			multiLineQuery.WriteString(trimmedQuery)
			calcitesql.ExecuteQuery(db, multiLineQuery.String())
			multiLineQuery.Reset()
			isMultiline = false
		} else {
			if !isMultiline {
				multiLineQuery.Reset()
				isMultiline = true
			}
			multiLineQuery.WriteString(trimmedQuery)
			multiLineQuery.WriteString(" ")
		}
	}
}
