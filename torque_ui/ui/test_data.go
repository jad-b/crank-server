package ui

import "time"

var (
	emptyWorkout = `{
  "tags": "",
  "exercises": [
    {
      "tags": "",
      "sets": "",
      "modifiers": "",
      "movement": "",
      "last_modified": "2015-12-31T00:58:51.935Z",
      "workout_id": -1,
      "exercise_id": 0
    }
  ],
  "user_id": -1,
  "last_modified": "2015-12-31T00:58:51.935Z",
  "workout_id": -1
}`
	then, _     = time.Parse("2015 Dec 26 @ 1712", WktTimeLayout)
	testWorkout = &Workout{
		ID:           -1,
		LastModified: then,
		Tags:         "unit= reps x kgs; this is a comment",
		Exercises: []Exercise{
			{
				ID:           0,
				WorkoutID:    -1,
				LastModified: then,
				Movement:     "Bench Press",
				Modifiers:    "",
				Sets:         "42 x 5, 52 x 5, 63 x 3, 77 x 5, 88 x 3, 100 x 4",
				Tags:         "Training Max=230 lbs;week=3",
			},
			{
				ID:           1,
				WorkoutID:    -1,
				LastModified: then,
				Movement:     "Row, BB",
				Modifiers:    "",
				Sets:         "49 x 9/4",
				Tags:         "prev=49 x 8/5",
			},
			{
				ID:           2,
				WorkoutID:    -1,
				LastModified: then,
				Movement:     "Bench Press",
				Modifiers:    "",
				Sets:         "77 x 9/5",
				Tags:         "prev=49 x 7/4",
			},
			{
				ID:           3,
				WorkoutID:    -1,
				LastModified: then,
				Movement:     "Pull-up, Ring",
				Modifiers:    "",
				Sets:         "8 x 8/4",
				Tags:         "prev=7 x 8/4",
			},
			{
				ID:           4,
				WorkoutID:    -1,
				LastModified: then,
				Movement:     "Press, BB, Standing, One-Arm, Thick",
				Modifiers:    "",
				Sets:         "20 x 7/6",
				Tags:         "prev=20 x 7/5",
			},
			{
				ID:           5,
				WorkoutID:    -1,
				LastModified: then,
				Movement:     "Tricep Extension, DB",
				Modifiers:    "",
				Sets:         "40 x 6/3",
				Tags:         "prev=35 x 10/4;unit=lbs x reps",
			},
		},
	}
)
