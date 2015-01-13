// grab our packages
var gulp   = require('gulp'),
    gulpFilter = require('gulp-filter'),
    mainBowerFiles = require('main-bower-files'),
    uglify = require('gulp-uglify');

var config = {
    jsPath: './ui/src/js/*.js',
    cssPath: './ui/src/css/*.css',
    assetsPath: './ui/assets'
};

// define the default task and add the watch task to it
gulp.task('default', ['watch']);


//// configure the jshint task
//gulp.task('jshint', function() {
  //return gulp.src(config.jsPath)
    //.pipe(jshint())
    //.pipe(jshint.reporter('jshint-stylish'));
//});

//// configure which files to watch and what tasks to use on file changes
gulp.task('watch', function() {
  gulp.watch(config.jsPath, ['jshint']);
});

// Minify JS
gulp.task('scripts', ['jshint'], function() {
    gulp.src(mainBowerFiles() + '/**/dist/*.min.js')
        .pipe(uglify())
        .pipe(gulp.dest(config.assetsPath + '/js/'));
});

// Move our custom CSS into place
gulp.task('css', function() {
    gulp.src(config.cssPath)
        .pipe(gulp.dest(config.assetsPath + '/css/'));
});
