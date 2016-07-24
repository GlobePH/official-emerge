var gulp = require('gulp');


var plugins = require('gulp-load-plugins')({
  pattern:  ['gulp-*', 'gulp.', 'main-bower-files'],
  replaceString: /\bgulp[\-.]/
});

// destination
var destination = 'js';


// JS task
// gulp.task('js', function() {
//   var jsFiles = ['js/*'];
// 
//   console.log(gulp.src(plugins.mainBowerFiles()).pipe(plugins.filter('*.js')));
//   gulp.src(plugins.mainBowerFiles())
//     .pipe(plugins.filter('*.js'))
//     .pipe(plugins.concat('main.js'))
//     .pipe(gulp.dest('./js'))
//     .pipe(plugins.rename('main.min.js'))
//     .pipe(plugins.uglify())
//     .pipe(gulp.dest("./js"));
// });

gulp.task('default', ['clean', 'js', 'css']);

gulp.task('clean', function() {
  var cleanThese = ['js/', 'css/'];
  cleanThese.forEach(function(element, index, array) {
    gulp.src(element, {read: false})
      .pipe(plugins.clean());
  });
});

gulp.task('js', function() {
  var jsFiles = ['src/js/*.js'];
   gulp.src(plugins.mainBowerFiles().concat(jsFiles)
    .concat("http://maps.googleapis.com/maps/api/js?sensor=false&language=en")) 
    .pipe(plugins.filter('**/*.js'))
    .pipe(plugins.concat('main.min.js'))
    .pipe(gulp.dest('./js'));
});

gulp.task('css', function() {
  var cssFiles = ['src/css/*.css'];
  return gulp.src(plugins.mainBowerFiles().concat(cssFiles))
          .pipe(plugins.filter('**/*.css'))
          .pipe(plugins.concat('main.min.css'))
          .pipe(gulp.dest('./css'));
});
