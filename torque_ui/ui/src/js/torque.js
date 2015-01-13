/*jshint multistr: true */
/*jshint expr: true */
$(document).ready(function(event) {
    var workoutTagCount = 1,
        $workoutTagSection = $("#wkt-tag-group"),
        $addWktTag = $("#add-workout-tag"),
		$exercises = $('#exercises');
		$addExercise = $('.exercises-header i'),
		tag_template = '\
<div class="row"> \
  <div class="form-group col-sm-4"> \
    <label class="sr-only" for="wkt-tag-key-{{ i }}">Workout Tag Key {{i}}</label> \
    <input type="text" class="form-control" id="wkt-tag-key-{{ i }}" placeholder="Key"> \
  </div> \
  <div class="form-group col-sm-4"> \
    <label class="sr-only" for="wkt-tag-value-{{ i }}">Workout Tag Value {{i}}</label> \
    <input type="text" class="form-control" id="wkt-tag-value-{{ i }}" placeholder="Value"> \
  </div> \
</div>',
		setTemplate = '\
<div class="col-sm-6">\
  <label class="sr-only" for="ex{{exID}}-set{{setID}}-weight">Weight</label> \
  <input type="numeric" id="ex{{exID}}-set{{setID}}-weight" \
		class="form-control" min="0" max="1000" step="5" value="100">\
</div>\
<div class="col-sm-6">\
  <label class="sr-only" for="ex{{exID}}-set{{setID}}-reps">Repetitions</label> \
  <input type="numeric" id="ex{{exID}}-set{{setID}}-reps" class="form-control" min="1" max="100" step="1" value="10">\
</div>',
        swagger = InitializeSwaggerClient();

    $addWktTag.click(function(e) {
        e.preventDefault();
		var tag_html = Mustache.to_html(tag_template, { 'i': workoutTagCount++});
        $workoutTagSection.append(tag_html);
    });

    $addExercise.click(function(e) {
        e.preventDefault();
		// Grab the last exercise div and clone
		$dupeExercise = $exercises.children('div.exercise').last().clone().appendTo($exercises);
    });

    /**
     * Form submission is hijacked to convert the form fields into a JSON
     * object for AJAX POST'ing to the server
     */
    $('form button').click(function(e) {
        workout = assembleWorkout($('form.workout-form'));
        swagger.wkt.toWktFormat(
            {body: workout},
            function(success){
                console.debug(success.obj.Wkt);
                $('.wkt-display').text(success.obj.Wkt);
                $('.wkt-display').show();
            },
            function(error){
               console.log("Failed to convert to .wkt: " + error.statusText);
            }
    );
    });

    // Disable normal form submission; hijack for JSON conversion
    $('form.workout-form').submit(false);
});

function assembleWorkout($form) {
    var now = moment().toISOString();
    workout = {
        workout_id: -1, // Will be ignored/overridden server-side
        last_modified: now,
        user_id: -1, // Once the app is user-aware, this will be a meaningful value
        exercises: [],
        tags: $form.find("input[name='workout-tags']").val()
    };
    // Clean and add each listed exercise
    $form.find('.exercise').each(function(i){
        var ex = {
            exercise_id: i,
            workout_id: -1,
            last_modified: now,
            // Input with substring 'name' as value for name attribute
            movement: $(this).find("input[name*='name']").val(),
            modifiers: "", // Ignore these for now
            sets: $(this).find("input[name*='sets']").val(),
            tags: $(this).find("input[name*='tags']").val()
        };
        workout.exercises.push(ex);
    });
    console.debug(prettyJSON(workout));
    return workout;
}

function prettyJSON(obj) {
    return JSON.stringify(obj, undefined, 2);
}

function InitializeSwaggerClient() {
    client = new SwaggerClient({
        url: "http://localhost:18001/spec/swagger.yaml",
        success: function() {
            this.setHost("localhost:18001");
            console.debug("Swagger client created");
        },
        useJQuery: true
    });
    return client;
}
