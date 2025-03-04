//
//  Author:: Salim Afiune <afiune@chef.io>
//  Copyright:: Copyright 2017, Chef Software Inc.
//

package integration_test

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"

	gp "github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"

	"github.com/chef/automate/api/interservice/cfgmgmt/request"
	iBackend "github.com/chef/automate/components/ingest-service/backend"
	"github.com/chef/automate/lib/grpc/grpctest"
)

func TestSuggestionsEmptyRequestReturnsError(t *testing.T) {
	ctx := context.Background()
	req := request.Suggestion{}

	res, err := cfgmgmt.GetSuggestions(ctx, &req)
	grpctest.AssertCode(t, codes.InvalidArgument, err)
	assert.Nil(t, res)
}

func TestSuggestionsWithAnInvalidTypeReturnsError(t *testing.T) {
	ctx := context.Background()
	cases := []string{"fake", "my-platform", "PolicyName", "error", "something_else"}

	for _, cType := range cases {
		t.Run(fmt.Sprintf("with type '%v' it should throw error", cType), func(t *testing.T) {
			req := request.Suggestion{Type: cType}
			res, err := cfgmgmt.GetSuggestions(ctx, &req)
			grpctest.AssertCode(t, codes.InvalidArgument, err)
			assert.Nil(t, res)
		})
	}
}

func TestSuggestionsWithTableDriven(t *testing.T) {
	dataNodes := []struct {
		number     int
		node       iBackend.NodeInfo
		attributes []string
		exists     bool
		namePrefix string
	}{
		{10, iBackend.NodeInfo{
			Status: "success", Platform: "windows",
			PolicyGroup: "heros", PolicyRevision: "marvel", PolicyName: "comics",
			Cookbooks: []string{"avengers", "justice_league", "guardians_of_the_galaxy"},
			Recipes:   []string{"thor", "ironman", "ant-man", "spiderman", "dr_strange"},
			Roles:     []string{"heros", "villains"}, ResourceNames: []string{"super_powers"},
			Environment: "games", OrganizationName: "org2", ChefVersion: "1.0",
			ChefTags: []string{"chef"}},
			[]string{"can_fly", "magic", "god", "immortal", "complex_attr_lala"},
			true,
			"node-"},
		{10, iBackend.NodeInfo{
			Status: "success", Platform: "ubuntu",
			PolicyGroup: "RPGs", PolicyRevision: "fantasy", PolicyName: "boardgames",
			Cookbooks:     []string{"dungeons_n_dragons", "angels", "startwars", "starfinder"},
			Recipes:       []string{"bard", "elf", "human", "barbarian"},
			Roles:         []string{"wizard", "fighter", "necromancer"},
			ResourceNames: []string{"wand", "sword", "bow", "shield"},
			Environment:   "dev", OrganizationName: "org1", ChefVersion: "1.1",
			ChefTags: []string{"chef", "boop"}},
			[]string{"dexterity", "charisma", "strength", "constitution", "intelligence"},
			true,
			"node-"},
		{10, iBackend.NodeInfo{
			Status: "success", Platform: "ubuntu",
			PolicyGroup: "RPGs", PolicyRevision: "fantasy", PolicyName: "boardgames",
			Cookbooks:     []string{"dungeons_n_dragons", "angels", "startwars", "starfinder"},
			Recipes:       []string{"bard", "elf", "human", "barbarian"},
			Roles:         []string{"wizard", "fighter", "necromancer"},
			ResourceNames: []string{"wand", "sword", "bow", "shield"},
			Environment:   "dev", OrganizationName: "org1", ChefVersion: "2.0",
			ChefTags: []string{"cheese"}},
			[]string{"dexterity", "charisma", "strength", "constitution", "intelligence"},
			false,
			"deleted-"},
		{10, iBackend.NodeInfo{
			Status: "failure", Platform: "ubuntu",
			PolicyGroup: "nintendo", PolicyRevision: "zelda", PolicyName: "videogames",
			Cookbooks: []string{"twilight", "ocarina_of_time", "breath_of_the_wild"},
			Recipes:   []string{"zelda", "link", "ganon", "epona"},
			Roles:     []string{"heros", "villains"}, ResourceNames: []string{"triforce"},
			Environment: "prod", OrganizationName: "org3", ChefVersion: "3.0",
			ChefTags: []string{"boop"}},
			[]string{"sword", "shield", "horse", "bow", "finally_can_jump"},
			true,
			"node-"},
		{20, iBackend.NodeInfo{
			Status: "success", Platform: "centos",
			PolicyGroup: "tv_series", PolicyRevision: "friends", PolicyName: "time_to_be_funny",
			Cookbooks:     []string{"rachel", "monica", "chandler", "phoebe", "joey"},
			Recipes:       []string{"season1", "season2", "season3", "you_get_it"},
			ResourceNames: []string{"smart", "joke", "smile", "hilarious", "honest"},
			Environment:   "new-york", OrganizationName: "org1", ChefVersion: "1.1.0",
			ChefTags: []string{"chef"}},
			[]string{"funny", "friendship", "sarcasm"},
			true,
			"node-"},
		{10, iBackend.NodeInfo{
			Status: "failure", Platform: "arch",
			PolicyGroup: "no_tv", PolicyRevision: "no_gum", PolicyName: "no_hats",
			Cookbooks: []string{"breakfast_club"}, Recipes: []string{"students"},
			Roles:       []string{"claire", "john", "allison", "andrew", "brian"},
			Environment: "dev", OrganizationName: "org2", ChefVersion: "2.0",
			ChefTags: []string{"boop"}},
			[]string{"saturday", "detention", "school", "disparate"},
			false,
			"deleted-"},
		{10, iBackend.NodeInfo{
			Status: "failure", Platform: "oracle",
			PolicyGroup: "movies", PolicyRevision: "old_school", PolicyName: "time_to_be_funny",
			Cookbooks: []string{"breakfast_club"}, Recipes: []string{"students"},
			Roles:       []string{"claire", "john", "allison", "andrew", "brian"},
			Environment: "dev", OrganizationName: "org2", ChefVersion: "2.0",
			ChefTags: []string{"boop"}},
			[]string{"saturday", "detention", "school", "disparate"},
			true,
			"node-"},
		{20, iBackend.NodeInfo{
			Status: "missing", Platform: "redhat",
			PolicyGroup: "sports", PolicyRevision: "extream", PolicyName: "games",
			Cookbooks: []string{"ping_pong"}, Recipes: []string{"soccer"},
			Roles: []string{"defence"}, ResourceNames: []string{"attack"},
			Environment: "prod", OrganizationName: "org3", ChefVersion: "2.0",
			ChefTags: []string{"cheese"}},
			[]string{"i_ran_out_of_ideas", "please_forgive_my_typos"},
			true,
			"node-"},
		{20, iBackend.NodeInfo{
			Status: "missing", Platform: "redhat",
			PolicyGroup: "sports", PolicyRevision: "extream", PolicyName: "games",
			Cookbooks: []string{"ping_pong"}, Recipes: []string{"soccer"},
			Roles: []string{"defence"}, ResourceNames: []string{"attack"},
			Environment: "prod", OrganizationName: "org3", ChefVersion: "2.0",
			ChefTags: []string{"cheese"}},
			[]string{"i_ran_out_of_ideas", "please_forgive_my_typos"},
			false,
			"deleted-"},
		{20, iBackend.NodeInfo{
			Status: "missing", Platform: "solaris"},
			[]string{},
			true,
			"node-"},
	}

	var (
		nodes        = make([]iBackend.Node, 0)
		allNodeNames = make([]string, 0)
		index        = 0
	)

	for _, data := range dataNodes {
		for i := 0; i < data.number; i++ {
			data.node.EntityUuid = newUUID()
			data.node.NodeName = data.namePrefix + fmt.Sprintf("%03d", index)
			if data.exists {
				allNodeNames = append(allNodeNames, data.node.NodeName)
			}
			node := iBackend.Node{
				NodeInfo:   data.node,
				Attributes: data.attributes,
				Exists:     data.exists,
			}
			nodes = append(nodes, node)
			index++
		}
	}
	suite.IngestNodes(nodes)
	defer suite.DeleteAllDocuments()

	ctx := context.Background()
	cases := []struct {
		description string
		request     request.Suggestion
		expected    []string
	}{
		// Suggestions for Nodes
		{"should return all nodes suggestions",
			request.Suggestion{Type: "name"},
			allNodeNames},
		{"should return just the set of node names that match",
			request.Suggestion{Type: "name", Text: "node-00"},
			[]string{"node-000", "node-001", "node-002", "node-003", "node-004",
				"node-005", "node-006", "node-007", "node-008", "node-009"}},
		{"should return no matching names",
			request.Suggestion{Type: "name", Text: "deleted-"},
			[]string{}},

		// Suggestions for Environments
		{"should return all environment suggestions",
			request.Suggestion{Type: "environment"},
			[]string{"dev", "prod", "games", "new-york"}},
		{"should return zero environment suggestions",
			request.Suggestion{Type: "platform", Text: "lol"},
			[]string{}},
		{"should return results when starting typing 'new-york'",
			request.Suggestion{Type: "environment", Text: "n"}, // less than 2 characters we return all results?
			[]string{"dev", "prod", "games", "new-york"}},
		{"should return one environment suggestion 'new-york'",
			request.Suggestion{Type: "environment", Text: "ne"},
			[]string{"new-york"}},
		{"should return one environment suggestion 'new-york'",
			request.Suggestion{Type: "environment", Text: "new"},
			[]string{"new-york"}},

		// Suggestions for Platform
		{"should return all platform suggestions",
			request.Suggestion{Type: "platform"},
			[]string{"centos", "redhat", "ubuntu", "oracle", "windows", "solaris"}},
		{"should return zero platform suggestions",
			request.Suggestion{Type: "platform", Text: "lol"},
			[]string{}},
		{"should return results when starting typing 'oracle'",
			request.Suggestion{Type: "platform", Text: "o"}, // less than 2 characters we return all results?
			[]string{"centos", "redhat", "ubuntu", "oracle", "windows", "solaris"}},
		{"should return one platform suggestion 'oracle'",
			request.Suggestion{Type: "platform", Text: "or"},
			[]string{"oracle"}},
		{"should return one platform suggestion 'oracle'",
			request.Suggestion{Type: "platform", Text: "ora"},
			[]string{"oracle"}},

		// Suggestions for Policy Group
		{"should return all policy_group suggestions",
			request.Suggestion{Type: "policy_group"},
			[]string{"sports", "tv_series", "RPGs", "heros", "movies", "nintendo"}},
		{"should return zero policy_group suggestions",
			request.Suggestion{Type: "policy_group", Text: "lol"},
			[]string{}},
		{"should return results when starting typing 'nintendo'",
			request.Suggestion{Type: "policy_group", Text: "n"}, // less than 2 characters we return all results?
			[]string{"sports", "tv_series", "RPGs", "heros", "movies", "nintendo"}},
		{"should return one policy_group suggestion 'nintendo'",
			request.Suggestion{Type: "policy_group", Text: "ni"},
			[]string{"nintendo"}},
		{"should return one policy_group suggestion 'nintendo'",
			request.Suggestion{Type: "policy_group", Text: "nin"},
			[]string{"nintendo"}},

		// Suggestions for Policy Name
		{"should return all policy_name suggestions",
			request.Suggestion{Type: "policy_name"},
			[]string{"time_to_be_funny", "games", "boardgames", "comics", "videogames"}},
		{"should return zero policy_name suggestions",
			request.Suggestion{Type: "policy_name", Text: "lol"},
			[]string{}},
		{"should return results when starting typing 'games'",
			request.Suggestion{Type: "policy_name", Text: "g"}, // less than 2 characters we return all results?
			[]string{"time_to_be_funny", "games", "boardgames", "comics", "videogames"}},
		{"should return one policy_name suggestion 'games'",
			request.Suggestion{Type: "policy_name", Text: "ga"},
			[]string{"games"}},
		{"should return one policy_name suggestion 'games'",
			request.Suggestion{Type: "policy_name", Text: "gam"},
			[]string{"games"}},

		// Suggestions for Policy Revision
		{"should return all policy_revision suggestions",
			request.Suggestion{Type: "policy_revision"},
			[]string{"extream", "friends", "fantasy", "marvel", "old_school", "zelda"}},
		{"should return zero policy_revision suggestions",
			request.Suggestion{Type: "policy_revision", Text: "lol"},
			[]string{}},
		{"should return results when starting typing 'friends'",
			request.Suggestion{Type: "policy_revision", Text: "f"}, // less than 2 characters we return all results?
			[]string{"extream", "friends", "fantasy", "marvel", "old_school", "zelda"}},
		{"should return one policy_revision suggestion 'friends'",
			request.Suggestion{Type: "policy_revision", Text: "fr"},
			[]string{"friends"}},
		{"should return one policy_revision suggestion 'friends'",
			request.Suggestion{Type: "policy_revision", Text: "fri"},
			[]string{"friends"}},

		// Suggestions for Recipes
		{"should return all recipe suggestions",
			request.Suggestion{Type: "recipe"},
			[]string{"season1", "season2", "season3", "soccer", "you_get_it", "ant-man",
				"barbarian", "bard", "dr_strange", "elf", "epona", "ganon", "human", "ironman",
				"link", "spiderman", "students", "thor", "zelda"}},
		{"should return zero recipe suggestions",
			request.Suggestion{Type: "recipe", Text: "lol"},
			[]string{}},
		{"should return results when starting typing 'season2'",
			request.Suggestion{Type: "recipe", Text: "s"}, // only return the words that has 's'
			[]string{"season1", "season2", "season3", "soccer", "dr_strange", "spiderman", "students"}},
		{"should return one recipe suggestion 'season2'",
			request.Suggestion{Type: "recipe", Text: "se"},
			[]string{"season1", "season3", "season2"}},
		{"should return one recipe suggestion 'season2'",
			request.Suggestion{Type: "recipe", Text: "sea"},
			[]string{"season1", "season3", "season2"}},
		{"should return one recipe suggestion 'season2'",
			request.Suggestion{Type: "recipe", Text: "season2"},
			[]string{"season2"}},

		// Suggestions for Cookbooks
		{"should return all cookbook suggestions",
			request.Suggestion{Type: "cookbook"},
			[]string{"chandler", "joey", "monica", "phoebe", "ping_pong", "rachel", "angels", "avengers",
				"breakfast_club", "breath_of_the_wild", "dungeons_n_dragons", "guardians_of_the_galaxy",
				"justice_league", "ocarina_of_time", "starfinder", "startwars", "twilight"}},
		{"should return zero cookbook suggestions",
			request.Suggestion{Type: "cookbook", Text: "lol"},
			[]string{}},
		{"should return results when starting typing 'breakfast_club'",
			request.Suggestion{Type: "cookbook", Text: "b"}, // only return the words that has 'b'
			[]string{"phoebe", "breakfast_club", "breath_of_the_wild"}},
		{"should return one cookbook suggestion 'breakfast_club'",
			request.Suggestion{Type: "cookbook", Text: "br"},
			[]string{"breakfast_club", "breath_of_the_wild"}},
		{"should return one cookbook suggestion 'breakfast_club'",
			request.Suggestion{Type: "cookbook", Text: "bre"},
			[]string{"breakfast_club", "breath_of_the_wild"}},
		{"should return one cookbook suggestion 'breakfast_club'",
			request.Suggestion{Type: "cookbook", Text: "break"},
			[]string{"breakfast_club"}},

		// Suggestions for Resource Names
		{"should return all resource_name suggestions",
			request.Suggestion{Type: "resource_name"},
			[]string{"attack", "hilarious", "honest", "joke", "smart", "smile", "bow", "shield",
				"super_powers", "sword", "triforce", "wand"}},
		{"should return zero resource_name suggestions",
			request.Suggestion{Type: "resource_name", Text: "lol"},
			[]string{}},
		{"should return results when starting typing 'super_powers' or 'smart'",
			request.Suggestion{Type: "resource_name", Text: "s"}, // only return the words that has 's'
			[]string{"hilarious", "honest", "smart", "smile", "shield", "super_powers", "sword"}},
		{"should return one resource_name suggestion 'super_powers'",
			request.Suggestion{Type: "resource_name", Text: "su"},
			[]string{"super_powers"}},
		{"should return one resource_name suggestion 'super_powers'",
			request.Suggestion{Type: "resource_name", Text: "sup"},
			[]string{"super_powers"}},
		{"should return one resource_name suggestion 'smart'",
			request.Suggestion{Type: "resource_name", Text: "sm"},
			[]string{"smile", "smart"}},
		{"should return one resource_name suggestion 'smart'",
			request.Suggestion{Type: "resource_name", Text: "sma"},
			[]string{"smart"}},

		// Suggestions for Attributes
		{"should return all attribute suggestions",
			request.Suggestion{Type: "attribute"},
			[]string{"friendship", "funny", "i_ran_out_of_ideas", "please_forgive_my_typos", "sarcasm", "bow",
				"can_fly", "charisma", "complex_attr_lala", "constitution", "detention", "dexterity", "disparate",
				"finally_can_jump", "god", "horse", "immortal", "intelligence", "magic", "saturday", "school",
				"shield", "strength", "sword"}},
		{"should return zero attribute suggestions",
			request.Suggestion{Type: "attribute", Text: "lol"},
			[]string{}},
		{"should return results when starting typing 'dexterity'",
			request.Suggestion{Type: "attribute", Text: "d"}, // only return the words that has 'd'
			[]string{"friendship", "i_ran_out_of_ideas", "detention", "dexterity", "disparate",
				"god", "saturday", "shield", "sword"}},
		{"should return one attribute suggestion 'dexterity'",
			request.Suggestion{Type: "attribute", Text: "de"},
			[]string{"dexterity", "detention"}},
		{"should return one attribute suggestion 'dexterity'",
			request.Suggestion{Type: "attribute", Text: "dex"},
			[]string{"dexterity"}},

		// Suggestions for Roles
		{"should return all role suggestions",
			request.Suggestion{Type: "role"},
			[]string{"defence", "heros", "villains", "allison", "andrew", "brian", "claire",
				"fighter", "john", "necromancer", "wizard"}},
		{"should return zero role suggestions",
			request.Suggestion{Type: "role", Text: "lol"},
			[]string{}},
		{"should return results when starting typing 'necromancer'",
			request.Suggestion{Type: "role", Text: "n"}, // only return the words that has 'd'
			[]string{"defence", "villains", "allison", "andrew", "brian", "john", "necromancer"}},
		{"should return one role suggestion 'necromancer'",
			request.Suggestion{Type: "role", Text: "ne"},
			[]string{"necromancer"}},
		{"should return one role suggestion 'necromancer'",
			request.Suggestion{Type: "role", Text: "nec"},
			[]string{"necromancer"}},

		// Suggestions for chef version
		{"should return all chef_version suggestions",
			request.Suggestion{Type: "chef_version"},
			[]string{"1.0", "1.1", "1.1.0", "2.0", "3.0"}},
		{"should return zero chef_version suggestions",
			request.Suggestion{Type: "chef_version", Text: "lol"},
			[]string{}},
		{"should return chef_version results when starting typing '1'",
			request.Suggestion{Type: "chef_version", Text: "1."}, // only return the versions that have '1.'
			[]string{"1.0", "1.1", "1.1.0"}},
		{"should return one chef_version suggestion '2.0'",
			request.Suggestion{Type: "chef_version", Text: "2.0"},
			[]string{"2.0"}},

		// Suggestions for chef tags
		{"should return all chef_tags suggestions",
			request.Suggestion{Type: "chef_tags"},
			[]string{"chef", "cheese", "boop"}},
		{"should return zero chef_tags suggestions",
			request.Suggestion{Type: "chef_tags", Text: "lol"},
			[]string{}},
		{"should return chef_tags results when starting typing 'ch'",
			request.Suggestion{Type: "chef_tags", Text: "che"}, // only return the versions that have 'ch.'
			[]string{"chef", "cheese"}},
		{"should return one chef_tags suggestion 'boop'",
			request.Suggestion{Type: "chef_tags", Text: "boop"},
			[]string{"boop"}},
	}

	// Run all the cases!
	for _, test := range cases {
		t.Run(fmt.Sprintf("with request '%v' it %s", test.request, test.description),
			func(t *testing.T) {
				res, err := cfgmgmt.GetSuggestions(ctx, &test.request)
				assert.Nil(t, err)

				// We actually don't care about the scores since it is something
				// the UI uses to order the results, therefor we will just abstract
				// the text into an array and compare it
				actualSuggestionsArray := extractTextFromSuggestionsResponse(res, t)

				// Verify they both are the same length
				assert.Equal(t, len(test.expected), len(actualSuggestionsArray))

				// Verify that they contains all the fields
				// We don't do 'assert.Equal()' because that checks order
				assert.ElementsMatch(t, actualSuggestionsArray, test.expected)
			})
	}
}

// extractTextFromSuggestionsResponse will extract only the text from the suggestions
// response:
//
// values {
//   struct_value {
//     fields {
//       key: "score"
//       value {
//         number_value: 1
//       }
//     }
//     fields {
//       key: "text"
//       value {
//         string_value: "dummy"       <--- This are the values we are looking for! :smile:
//       }
//     }
//   }
// }
//
// The only problem with that is that if there are nodes that has empty fields we will find
// something similar to this response:
//
// values {
//   struct_value {
//     fields {
//       key: "score"
//       value {
//         number_value: 1
//       }
//     }
//   }
// }
//
// Where there is NO string value!
//
// TODO: (@afiune) Is this the normal behavior? If not lets fix it.
// for now the fixt in the tests will be to check if there is a "string"
// value or not.
func extractTextFromSuggestionsResponse(list *gp.ListValue, t *testing.T) []string {
	// We don't initialize the slice size since we might found empty Values
	textArray := make([]string, 0)

	if list != nil {
		for _, sugg := range list.Values {
			sugStruct := sugg.GetStructValue()
			if txt := sugStruct.Fields["text"].GetStringValue(); txt != "" {
				textArray = append(textArray, txt)
			}
		}
	}
	return textArray
}
