package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
			{
				"id": "_pb_users_auth_",
				"created": "2024-08-12 08:16:27.245Z",
				"updated": "2024-09-03 03:52:40.840Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "2upz0oyd",
						"name": "has_changed_username",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id",
				"viewRule": "",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"onlyVerified": false,
					"requireEmail": false
				}
			},
			{
				"id": "4bz0rmd2bvmwuqf",
				"created": "2024-08-12 08:51:13.860Z",
				"updated": "2024-08-31 18:25:36.463Z",
				"name": "nodes",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "fxzbcyn2",
						"name": "text",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "fhloqlk2",
						"name": "previous_node",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "4bz0rmd2bvmwuqf",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "b0b1xn9u",
						"name": "author",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "0unm1cg2gqw6i4c",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "epujuzfu",
						"name": "likes",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "eflcuhgw",
						"name": "bookmarks",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "qjtqx1fq",
						"name": "x_position",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "bvkl8qye",
						"name": "y_position",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "n8ugqveb",
						"name": "state",
						"type": "select",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"moving",
								"pending",
								"editing",
								"completed"
							]
						}
					}
				],
				"indexes": [],
				"listRule": "",
				"viewRule": "",
				"createRule": "",
				"updateRule": "",
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "0unm1cg2gqw6i4c",
				"created": "2024-08-12 09:01:20.399Z",
				"updated": "2024-08-25 21:08:45.427Z",
				"name": "authors",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "qvxlbjvb",
						"name": "author_id",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "jcoukdwv",
						"name": "origin",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"human",
								"robot"
							]
						}
					}
				],
				"indexes": [],
				"listRule": "",
				"viewRule": "",
				"createRule": "@request.auth.id != \"\" && @request.auth.id = author_id && origin = \"human\"",
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "716s56dwn65vl48",
				"created": "2024-08-12 09:01:41.696Z",
				"updated": "2024-09-03 03:51:51.578Z",
				"name": "robots",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "iy9pjpgy",
						"name": "model",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "xuphmxtv2x0pacm",
				"created": "2024-09-03 04:03:57.646Z",
				"updated": "2024-10-17 19:36:36.539Z",
				"name": "bookmarks",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "d2uqckzl",
						"name": "node",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "4bz0rmd2bvmwuqf",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "nbzwdaos",
						"name": "user",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_IqQiAV0` + "`" + ` ON ` + "`" + `bookmarks` + "`" + ` (\n  ` + "`" + `node` + "`" + `,\n  ` + "`" + `user` + "`" + `\n)"
				],
				"listRule": "@request.auth.id = user.id",
				"viewRule": "@request.auth.id = user.id",
				"createRule": "@request.auth.id = user.id",
				"updateRule": null,
				"deleteRule": "@request.auth.id = user.id",
				"options": {}
			},
			{
				"id": "e4pxif1hl5ojgjs",
				"created": "2024-09-03 04:06:27.597Z",
				"updated": "2024-09-06 02:42:45.262Z",
				"name": "most_bookmarked",
				"type": "view",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "to8kqh8i",
						"name": "node",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "4bz0rmd2bvmwuqf",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "tzx3j0qe",
						"name": "bookmark_count",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					}
				],
				"indexes": [],
				"listRule": "",
				"viewRule": "",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {
					"query": "SELECT node.id, node.id as node, COUNT(bookmark.id) AS bookmark_count\nFROM nodes node\nLEFT JOIN bookmarks bookmark ON node.id = bookmark.node\nGROUP BY node.id;"
				}
			},
			{
				"id": "2w0q0ad3k5kzqux",
				"created": "2024-11-27 23:57:07.724Z",
				"updated": "2024-11-27 23:57:31.214Z",
				"name": "reports",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "kkpomcg3",
						"name": "node",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "4bz0rmd2bvmwuqf",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "l48ukr0s",
						"name": "user",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_lUNwRD8` + "`" + ` ON ` + "`" + `reports` + "`" + ` (\n  ` + "`" + `node` + "`" + `,\n  ` + "`" + `user` + "`" + `\n)"
				],
				"listRule": "@request.auth.id = user.id",
				"viewRule": null,
				"createRule": "@request.auth.id = user.id",
				"updateRule": null,
				"deleteRule": "@request.auth.id = user.id",
				"options": {}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
