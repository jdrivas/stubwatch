package server

import(
  "fmt"
  // "html/template"
  "net/http"
  // "strconv"
  "html/template"
  "path/filepath"
  // "sort"
  "strings"
  "github.com/gorilla/context"
  "github.com/Sirupsen/logrus"
)


// Render handles the return. Note, that if there is an error it's already written
// the erro to the response. So should be essentially the last thing you do
// in your controller.
func Render(w http.ResponseWriter, r *http.Request, renderValues interface{}) {
  err := doRender(w, r, renderValues)
  if err != nil {
    log.Error(nil, "Error rendering content.", err)
    http.Error(w, fmt.Sprintf("Error rendering content: %s", err), 
      http.StatusInternalServerError)
  }
}

func doRender(w http.ResponseWriter, r *http.Request, renderValues interface{}) (err error) {

  v, ok := context.GetOk(r, ViewsKey);
  if !ok {
    log.Error(nil, "Couldn't get the templates.", err)
    http.Error(w, fmt.Sprintf("Can't read templates. %s", err), http.StatusInternalServerError)
  }
  views := v.(*template.Template)



  t, ok := context.GetOk(r, ViewsListKey); 
  if !ok {
    log.Error(nil, "Can't get the list of view templates.", nil)
    http.Error(w, fmt.Sprintf("Error getting the view templates, Key not present.",), http.StatusInternalServerError)
  }
  templates := t.([]string)
  log.Debug(logrus.Fields{
    "templateName": views.Name(),
    "templates": templates,
    }, "Rendering.")

  // Render
  for _, templ := range templates {
    err = views.ExecuteTemplate(w, templ, renderValues)
    if err != nil {
      log.Error(logrus.Fields{"template": templ,}, "Couldn't write template.", err)
      http.Error(w, fmt.Sprintf("Can't render templates.%s", err), http.StatusInternalServerError)
      return
    }
  }
  views.Execute(w,nil)

  return err
}

func SetViewsHandler(next http.Handler) http.Handler {
  return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {

    log.Debug(nil,"Loading views ....")
    views, viewTemplateNames, err := getViews(w, r)
    if err != nil {
      http.Error(w, fmt.Sprintf("Error getting views: %s", err), http.StatusInternalServerError)
      return
    }
    log.Info(logrus.Fields{
      "viewTemplateNames": viewTemplateNames,
      }, "Loaded views and template names.")
    context.Set(r, ViewsKey, views)
    context.Set(r, ViewsListKey, viewTemplateNames)

    next.ServeHTTP(w,r)
  })
}

// TODO: This whole thing wants to be cached.
func getViews(w http.ResponseWriter, r *http.Request) (views *template.Template, templateNames []string, err error) {
  // THIS WANTS TO CACHED!
  layoutFile, controllerFiles, templateNames, err  := loadViews(r)
  if err != nil {
    return nil, nil, fmt.Errorf("Error loading views. files.", err)
  }
  f := logrus.Fields{
    "layout": layoutFile,
    "controllerFiles": controllerFiles,
    "templateNames": templateNames,
  }
  log.Debug(f,"Loaded view files names.")

  viewFiles := []string{layoutFile}
  for _, f := range controllerFiles {
    viewFiles = append(viewFiles, f)
  }
  views, err = template.ParseFiles(viewFiles...)
  if err != nil  || views == nil {
    log.Error(f, "Failed to parse view templates.", err)
    err = fmt.Errorf("Failed to parse view templates: %s", err)
  }

  return views, templateNames, err
}

func loadViews(r *http.Request) (layoutFile string, viewFiles, templateNames []string, err error) {

  controllerFiles, controllerName, err :=loadControllerViews(r)
  if err != nil {
    return layoutFile, viewFiles, templateNames, fmt.Errorf("Failed to load controller files: %s", err)
  }

  layoutFile, err = loadLayoutView()
  if err != nil {
    return layoutFile, viewFiles, templateNames, fmt.Errorf("Failed to load layout files: %s", err)
  }


  templateNames = getTemplateNames(viewFiles, controllerName, defaultLayoutTemplate)

  return layoutFile, controllerFiles, templateNames, err
}

const (
  viewsPath = "app/views/"
  layoutsPath  = viewsPath + "layouts/"
  defaultLayoutTemplate = "application"
  defaultLayout = defaultLayoutTemplate  + ".tmpl"
  rootController = "home"
  contentTemplateName = "content"
)

func loadControllerViews(r *http.Request) (controllerFiles []string, 
    controllerName string, err error) {

  controllerGlob := viewsPath + "*.tmpl"
  pathComponents := strings.Split(r.URL.Path, "/")
  controllerName = rootController
  if len(pathComponents) > 1 {
    controllerName = pathComponents[1]
    controllerGlob = viewsPath + controllerName + "/*.tmpl"
  }

  controllerFiles, err = filepath.Glob(controllerGlob)
  if err != nil {
    log.Error(logrus.Fields{
      "glob": controllerGlob,
      }, "Error loading controller view files %s", err)
    return nil, controllerName, fmt.Errorf("Error getting the controller files with glob: %s: %s", controllerGlob, err)
  } else {
    log.Debug(logrus.Fields{"files": controllerFiles, "glob": controllerGlob,},
     "Loaded controller views.")
  }
  return controllerFiles, controllerName, err
}

// TODO: have a think about multiple layouts.
func loadLayoutView() (layoutFiles string, err error) {
  return layoutsPath + defaultLayout, nil
  // layoutGlob := layoutsPath + "*.tmpl"
  // layoutFiles, err = filepath.Glob(layoutGlob)
  // if err != nil {
  //   log.Error(logrus.Fields{
  //     "glob": layoutGlob,
  //     }, "Can't get the shared view files.", err)
  //   return nil, fmt.Errorf("Error get controller files with glob: %s: %s", layoutGlob, err)
  // } else {
  //   log.Debug(logrus.Fields{"files": layoutFiles,"glob": layoutGlob}, "Loaded layout views.")
  // }
  // return layoutFiles, err
}

// The convention is that each filename (minus extension) is the name of the template.
// Except for the contorller file which must will be defined as "content".
// This allows an application layout to always include a content in the middle of it.
func getTemplateNames(files []string, controllerName, layoutName string) (ts []string)  {
  for _, f := range files {
    tn := filepath.Base(f)
    tn = strings.Split(tn, ".")[0] // making some assumptions here.
    ts = append(ts, tn)
  }
  ts = append(ts, layoutName)
  // ts = append(ts, contentTemplateName)
  return ts
}