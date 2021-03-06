var gulp = require('gulp');
// requires node version >= 12
var child = require('child_process');
const babel = require('gulp-babel');
const browserify = require('gulp-browserify');
var run = require('gulp-run');
var chalk = require('chalk')
var util = require('gulp-util')
var rename = require('gulp-rename')
const rewrite = require('gulp-rewrite-css')
const concatCss = require('gulp-concat-css')
var readline = require('readline')

var gulpProcess;
var verbose = true;
var short = true
var rl = readline.createInterface({input: process.stdin, output: process.stdin});

gulp.task('test', function() {
  args = ["test"]
  if (verbose) { args.push("-v");}
  if(short) {args.push("-test.short");}

  test = child.spawnSync("go", args)
  if(test.status == 0) {
    util.log(chalk.white.bgGreen.bold(' Go Test Successful'));
    if (verbose) {
      var lines = test.stdout.toString().split("\n");
      for (var l in lines) {
        util.log(lines[l]);
      }
    }
  } else {
    util.log(chalk.white.bgRed.bold(" GO Test Failed "))
    var lines = test.stdout.toString().split("\n");
    for (var l in lines) {
      util.log(chalk.red(lines[l]));
    }
    var errLines = test.stderr.toString().split("\n");
    for (var l in errLines) {
      util.log(chalk.black.bold(errLines[l]));
    }
  }
  return test
})

gulp.task('build', function() {
  build = child.spawnSync("go", ["install"])
  if(build.status == 0) {
    util.log(chalk.white.bgBlue.bold(' Go Install Successful '));
  } else {
    util.log(chalk.white.bgRed.bold(" GO Install Failed "))
    var lines = build.stderr.toString().split('\n');
    for (var l in lines)
      util.log(chalk.red(lines[l]));
  }
  return build;
});

const assetsSrc = 'app/assets/'
const assetsBuild = 'app/assets/build/'
const assetsFinal = 'app/'
const cssFinal = assetsFinal + "styles/"
const jsFinal = assetsFinal + "js/"

gulp.task('babel', function() {
  return gulp.src(assetsSrc + 'js/**/*.js')
    .pipe(babel({presets: ['react', 'es2015']}))
    .pipe(gulp.dest(assetsBuild + 'js'))
});

// Get his from the build abvove.
gulp.task('browserify', function() {
  return gulp.src(assetsBuild + 'js/app.js')
  .pipe(browserify({
    insertGlobals: true,
    debug: true
  }))
  .pipe(rename('bundle.js'))
  .pipe(gulp.dest(assetsFinal + 'js'))
});

gulp.task('css', function() {
  var dest = cssFinal
  return gulp.src(assetsSrc + '/css/**/*.css')
  .pipe(concatCss("bundle.css")) // This will also rebase the URLS by default.
  .pipe(gulp.dest(dest))
});

function doCommand(command) {
  retVal = true
  commands = command.split(" ")
  switch (commands[0]) {
    case 'verbose': 
      if (verbose) {
        verbose = false;
      } else {
        verbose = true;
      }
      util.log("Verbose is now " + verbose.toString())
      break;
    case 'long':
      short = false;
      util.log("Short test is now " + short.toString())
      break;
    case 'short':
      short = true;
      util.log("Short test is now " + short.toString())
      break;
    case 'quit':
      process.exit();
    break;
    case '': // just eat returns.
      break
    default:
      util.log("Unknown command: ", command)
      retVal = false;
  }
  return retVal;
}

// This is just user interface hackery to get a command line to appear reasonablly.
// really only works if you run: gulp watch --silent
gulp.task('newline', function() {process.stdout.write("\n"); return true;});
gulp.task('prompt', function() {rl.prompt(); return true;});

function doCommandPrompt(answer) {
  doCommand(answer);
  rl.prompt();
}

function commandPromptLoop() {
  rl.setPrompt("command: ");
  rl.on('line', doCommandPrompt);
  rl.prompt();
}

gulp.task('watch', function(){
  gulp.watch('**/*.go', ['newline', 'build', 'test', 'prompt']);
  gulp.watch(assetsSrc + 'js/**/*.js', ['newline', 'babel', 'browserify', 'prompt']);
  gulp.watch(assetsSrc + 'css/**/*.css', ['newline', 'css', 'prompt'])
  commandPromptLoop()
})

// Variety of ways explored to get gulp to reload on gulpfile edit.
// In the end went with gulper: npm install -g gulper.
// I can't figure out how to turn off notificationn, but it's a small
// price.
