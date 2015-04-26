package main

import "encoding/json"
import "io/ioutil"
import "net/http"
import "fmt"
import "net/url"

import "github.com/go-martini/martini"
import "github.com/martini-contrib/render"

type Student struct {
	ID 			string `json:"id"`
	Name 		string
	Email 		string
	Birthdate 	string
	Age 		string
}


const ApiRoot = "http://djholt.php.cs.dixie.edu/slim-example/index.php"
// const ApiRoot = "http://localhost:8888/slim_example/index.php"

func makeGetRequest(url string) []byte {
  // get url
  // make request ( http )
  // get the response
  // return the body
  
  res, err := http.Get(ApiRoot + url)
  if err != nil {
    panic(err)
  }

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    panic(err)
  }

  return body
}

func makeFormPostRequest(url string, data url.Values) []byte {
  res, err := http.PostForm(ApiRoot + url, data)
  
  if err != nil {
    panic(err)
  }

  body, err := ioutil.ReadAll(res.Body)

  if err != nil {
    panic(err)
  }

  return body
}

func createStudent(student Student) Student {
  formData := url.Values{
    "name": {student.Name},
    "email": {student.Email},
    "birthdate": {student.Birthdate},
    "age": {student.Age},
  }
  data := makeFormPostRequest("/students.json", formData)

  student = Student{}
  err:= json.Unmarshal(data, &student)
  if err != nil {
    panic(err)
  }

  return student
}

func getStudents() []Student {
  students := []Student{}
  data := makeGetRequest("/students.json")
  err := json.Unmarshal(data, &students)
  if err != nil {
    panic(err)
  }
  return students
}

func getStudent(id string) Student {
  student := Student{}
  data := makeGetRequest(fmt.Sprintf("/students/%s.json", id))
  err := json.Unmarshal(data, &student)
  if err != nil {
    panic(err)
  }
  return student
}
  // students = append(students, Student{
  //   ID: "1",
  //   Name: "Erik",
  //   Email: "ebylund@dmail.dixie.edu",
  //   Birthdate: "1990-11-20",
  //   Age: "27",
  //   })
  // students = append(students, Student{
  //   ID: "2",
  //   Name: "DJ",
  //   Email: "dbylund@dmail.dixie.edu",
  //   Birthdate: "1990-11-20",
  //   Age: "22",
  //   })

func main() {
  m := martini.Classic()
  m.Use(render.Renderer(render.Options{
  		Layout: "layout",
  	}))
  m.Get("/", func(r render.Render ) {
  	r.HTML(200, "hello", "Erik")
  })


//index
  m.Get("/students", func(r render.Render){
    students := getStudents()

  	r.HTML(200, "students/index", students)
  })



//show
  m.Get("/students/(?P<id>[0-9]+)", func(params martini.Params, r render.Render){
    student := getStudent(params["id"])
    r.HTML(200, "students/show", student)
  })

//new
  m.Get("/students/new", func(r render.Render){
    r.HTML(200, "students/new", nil)
  })

// create
  m.Post("/students", func(req *http.Request, r render.Render){

    student := Student{
      Name: req.FormValue("name"),
      Email: req.FormValue("email"),
      Birthdate: req.FormValue("birthdate"),
      Age: req.FormValue("age"),
    }

    createStudent(student)

    r.HTML(201, "students/create", student)
  })

  m.Get("/test", func() []byte {
    data := makeGetRequest("/students/json")
    return data
  })
  
  m.Run()
}
