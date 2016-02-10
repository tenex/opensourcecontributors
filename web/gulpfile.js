var gulp = require('gulp');
var browserify = require('browserify');
var connect = require('connect');
var source = require('vinyl-source-stream');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var minifyCss = require('gulp-minify-css');
var imagemin = require('gulp-imagemin');
var sourcemaps = require('gulp-sourcemaps');
var templateCache = require('gulp-angular-templatecache');
var del = require('del');

var paths = {
  static: ['index.html', 'octicons/**/*'],
  fonts: ['fonts/*'],
  scripts: 'js/*.js',
  stylesheets: 'css/*.css',
  templates: 'templates/**/*.html',
  images: 'img/**/*'
};

gulp.task('browserify', function() {
  // Grabs the app.js file
  return browserify('js/app.js')
    .bundle()
    .pipe(source('main.js'))
    .pipe(gulp.dest('public/'));
});

gulp.task('stylesheets', function() {
  return gulp.src(paths.stylesheets)
    .pipe(sourcemaps.init())
    .pipe(minifyCss())
    .pipe(sourcemaps.write())
    .pipe(gulp.dest('public'));
});

gulp.task('static', function() {
  return gulp.src(paths.static)
    .pipe(gulp.dest('public'));
});

gulp.task('fonts', function() {
  return gulp.src(paths.fonts)
    .pipe(gulp.dest('public/fonts'));
});


gulp.task('clean', function() {
  return del(['public/**/*']);
});

gulp.task('templates', function() {
  return gulp.src(paths.templates)
    .pipe(templateCache('templates.js', {
      standalone: true
    }))
    .pipe(gulp.dest('js/'));
});

gulp.task('scripts', ['templates', 'browserify'], function() {
  // return gulp.src(paths.scripts)
  //     .pipe(sourcemaps.init())
  //     .pipe(concat('app.min.js'))
  //     .pipe(sourcemaps.write())
  //     .pipe(gulp.dest('public'));
});

gulp.task('watch', function() {
  gulp.watch([paths.templates, paths.scripts], ['scripts']);
  gulp.watch(paths.static, ['static']);
  gulp.watch(paths.stylesheets, ['stylesheets']);
});

gulp.task('connect', function () {
  connect.server({
    root: 'public',
    port: 4000
  });
});

gulp.task('default', [
  'scripts',
  'static',
  'stylesheets'
]);
