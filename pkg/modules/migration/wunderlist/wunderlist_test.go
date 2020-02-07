// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-2020 Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package wunderlist

import (
	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/files"
	"code.vikunja.io/api/pkg/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/d4l3k/messagediff.v1"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

func TestWunderlistParsing(t *testing.T) {

	config.InitConfig()

	time1, err := time.Parse(time.RFC3339Nano, "2013-08-30T08:29:46.203Z")
	assert.NoError(t, err)
	time2, err := time.Parse(time.RFC3339Nano, "2013-08-30T08:36:13.273Z")
	assert.NoError(t, err)
	time3, err := time.Parse(time.RFC3339Nano, "2013-09-05T08:36:13.273Z")
	assert.NoError(t, err)
	time4, err := time.Parse(time.RFC3339Nano, "2013-08-02T11:58:55Z")
	assert.NoError(t, err)

	exampleFile, err := ioutil.ReadFile(config.ServiceRootpath.GetString() + "/pkg/modules/migration/wunderlist/testimage.jpg")
	assert.NoError(t, err)

	createTestTask := func(id, listID int, done bool) *task {
		completedAt, err := time.Parse(time.RFC3339Nano, "1970-01-01T00:00:00Z")
		assert.NoError(t, err)
		if done {
			completedAt = time1
		}
		return &task{
			ID:          id,
			AssigneeID:  123,
			CreatedAt:   time1,
			DueDate:     "2013-09-05",
			ListID:      listID,
			Title:       "Ipsum" + strconv.Itoa(id),
			Completed:   done,
			CompletedAt: completedAt,
		}
	}

	createTestNote := func(id, taskID int) *note {
		return &note{
			ID:        id,
			TaskID:    taskID,
			Content:   "Lorem Ipsum dolor sit amet",
			CreatedAt: time3,
			UpdatedAt: time2,
		}
	}

	fixtures := &wunderlistContents{
		folders: []*folder{
			{
				ID:        123,
				Title:     "Lorem Ipsum",
				ListIds:   []int{1, 2, 3, 4},
				CreatedAt: time1,
				UpdatedAt: time2,
			},
		},
		lists: []*list{
			{
				ID:        1,
				CreatedAt: time1,
				Title:     "Lorem1",
			},
			{
				ID:        2,
				CreatedAt: time1,
				Title:     "Lorem2",
			},
			{
				ID:        3,
				CreatedAt: time1,
				Title:     "Lorem3",
			},
			{
				ID:        4,
				CreatedAt: time1,
				Title:     "Lorem4",
			},
			{
				ID:        5,
				CreatedAt: time4,
				Title:     "List without a namespace",
			},
		},
		tasks: []*task{
			createTestTask(1, 1, false),
			createTestTask(2, 1, false),
			createTestTask(3, 2, true),
			createTestTask(4, 2, false),
			createTestTask(5, 3, false),
			createTestTask(6, 3, true),
			createTestTask(7, 3, true),
			createTestTask(8, 3, false),
			createTestTask(9, 4, true),
			createTestTask(10, 4, true),
		},
		notes: []*note{
			createTestNote(1, 1),
			createTestNote(2, 2),
			createTestNote(3, 3),
		},
		files: []*file{
			{
				ID:          1,
				URL:         "https://vikunja.io/testimage.jpg", // Using an image which we are hosting, so it'll still be up
				TaskID:      1,
				ListID:      1,
				FileName:    "file.md",
				ContentType: "text/plain",
				FileSize:    12345,
				CreatedAt:   time2,
				UpdatedAt:   time4,
			},
			{
				ID:          2,
				URL:         "https://vikunja.io/testimage.jpg",
				TaskID:      3,
				ListID:      2,
				FileName:    "file2.md",
				ContentType: "text/plain",
				FileSize:    12345,
				CreatedAt:   time3,
				UpdatedAt:   time4,
			},
		},
		reminders: []*reminder{
			{
				ID:        1,
				Date:      time4,
				TaskID:    1,
				CreatedAt: time4,
				UpdatedAt: time4,
			},
			{
				ID:        2,
				Date:      time3,
				TaskID:    4,
				CreatedAt: time3,
				UpdatedAt: time3,
			},
		},
		subtasks: []*subtask{
			{
				ID:        1,
				TaskID:    2,
				CreatedAt: time4,
				Title:     "LoremSub1",
			},
			{
				ID:        2,
				TaskID:    2,
				CreatedAt: time4,
				Title:     "LoremSub2",
			},
			{
				ID:        3,
				TaskID:    4,
				CreatedAt: time4,
				Title:     "LoremSub3",
			},
		},
	}

	expectedHierachie := []*models.NamespaceWithLists{
		{
			Namespace: models.Namespace{
				Name:    "Lorem Ipsum",
				Created: time1.Unix(),
				Updated: time2.Unix(),
			},
			Lists: []*models.List{
				{
					Created: time1.Unix(),
					Title:   "Lorem1",
					Tasks: []*models.Task{
						{
							Text:        "Ipsum1",
							DueDateUnix: 1378339200,
							Created:     time1.Unix(),
							Description: "Lorem Ipsum dolor sit amet",
							Attachments: []*models.TaskAttachment{
								{
									File: &files.File{
										Name:        "file.md",
										Mime:        "text/plain",
										Size:        12345,
										Created:     time2,
										CreatedUnix: time2.Unix(),
										FileContent: exampleFile,
									},
									Created: time2.Unix(),
								},
							},
							RemindersUnix: []int64{time4.Unix()},
						},
						{
							Text:        "Ipsum2",
							DueDateUnix: 1378339200,
							Created:     time1.Unix(),
							Description: "Lorem Ipsum dolor sit amet",
							RelatedTasks: map[models.RelationKind][]*models.Task{
								models.RelationKindSubtask: {
									{
										Text: "LoremSub1",
									},
									{
										Text: "LoremSub2",
									},
								},
							},
						},
					},
				},
				{
					Created: time1.Unix(),
					Title:   "Lorem2",
					Tasks: []*models.Task{
						{
							Text:        "Ipsum3",
							Done:        true,
							DoneAtUnix:  time1.Unix(),
							DueDateUnix: 1378339200,
							Created:     time1.Unix(),
							Description: "Lorem Ipsum dolor sit amet",
							Attachments: []*models.TaskAttachment{
								{
									File: &files.File{
										Name:        "file2.md",
										Mime:        "text/plain",
										Size:        12345,
										Created:     time3,
										CreatedUnix: time3.Unix(),
										FileContent: exampleFile,
									},
									Created: time3.Unix(),
								},
							},
						},
						{
							Text:          "Ipsum4",
							DueDateUnix:   1378339200,
							Created:       time1.Unix(),
							RemindersUnix: []int64{time3.Unix()},
							RelatedTasks: map[models.RelationKind][]*models.Task{
								models.RelationKindSubtask: {
									{
										Text: "LoremSub3",
									},
								},
							},
						},
					},
				},
				{
					Created: time1.Unix(),
					Title:   "Lorem3",
					Tasks: []*models.Task{
						{
							Text:        "Ipsum5",
							DueDateUnix: 1378339200,
							Created:     time1.Unix(),
						},
						{
							Text:        "Ipsum6",
							DueDateUnix: 1378339200,
							Created:     time1.Unix(),
							Done:        true,
							DoneAtUnix:  time1.Unix(),
						},
						{
							Text:        "Ipsum7",
							DueDateUnix: 1378339200,
							Created:     time1.Unix(),
							Done:        true,
							DoneAtUnix:  time1.Unix(),
						},
						{
							Text:        "Ipsum8",
							DueDateUnix: 1378339200,
							Created:     time1.Unix(),
						},
					},
				},
				{
					Created: time1.Unix(),
					Title:   "Lorem4",
					Tasks: []*models.Task{
						{
							Text:        "Ipsum9",
							DueDateUnix: 1378339200,
							Created:     time1.Unix(),
							Done:        true,
							DoneAtUnix:  time1.Unix(),
						},
						{
							Text:        "Ipsum10",
							DueDateUnix: 1378339200,
							Created:     time1.Unix(),
							Done:        true,
							DoneAtUnix:  time1.Unix(),
						},
					},
				},
			},
		},
		{
			Namespace: models.Namespace{
				Name: "Migrated from wunderlist",
			},
			Lists: []*models.List{
				{
					Created: time4.Unix(),
					Title:   "List without a namespace",
				},
			},
		},
	}

	hierachie, err := convertWunderlistToVikunja(fixtures)
	assert.NoError(t, err)
	assert.NotNil(t, hierachie)
	if diff, equal := messagediff.PrettyDiff(hierachie, expectedHierachie); !equal {
		t.Errorf("ListUser.ReadAll() = %v, want %v, diff: %v", hierachie, expectedHierachie, diff)
	}
}