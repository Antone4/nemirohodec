package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type Context struct {
	Request    *http.Request
	Data       string
	Title      string
	Tovars     *[]Tovar
	Email      string
	Authorized bool
}

var titles map[string]string = map[string]string{
	"about":   "О нас",
	"catalog": "Каталог",
	"contact": "Контакты",
	"index":   "Не Мироходец",
}

var templateNames = []string{"about", "catalog", "contact"}

func init() {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			content := titles["index"]
			tmpl, err := template.ParseFiles("templates/index.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			templateData, err := os.ReadFile("templates/defines.html")
			if err != nil {
				panic(err)
			}
			thisTovars := getTovarsFromPGX(conn)
			tmpl = template.Must(tmpl.Parse(string(templateData)))

			session, _ := store.Get(r, "authdata")
			interfaceEmail, authorized := session.Values["email"]
			email := ""
			if authorized {
				email = interfaceEmail.(string)
			}

			err = tmpl.Execute(w, &Context{r, "hi", content, thisTovars, email, authorized})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			Printfln("%v", r.FormValue("psw"))
		} else {
			http.Redirect(w, r, "/PageNotFound", http.StatusMovedPermanently)
		}
	}).Methods("GET")
	router.HandleFunc("/", RegisterLoginHandler).Methods("POST")

	HandleItems()

	for _, name := range templateNames {
		newName := name
		router.HandleFunc("/"+newName, func(w http.ResponseWriter, r *http.Request) {
			content := titles[newName]
			tmpl, err := template.ParseFiles("templates/" + newName + ".html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			templateData, err := os.ReadFile("templates/defines.html")
			if err != nil {
				panic(err)
			}
			tmpl = template.Must(tmpl.Parse(string(templateData)))
			session, _ := store.Get(r, "authdata")
			interfaceEmail, authorized := session.Values["email"]
			email := ""
			if authorized {
				email = interfaceEmail.(string)
			}

			err = tmpl.Execute(w, &Context{r, "hi", content, &tovars, email, authorized})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}).Methods("GET")
		router.HandleFunc("/"+newName, RegisterLoginHandler).Methods("POST")
	}

	router.HandleFunc("/PageNotFound", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile("static/PageNotFound.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/PageNotFound", http.StatusFound)
	})

	HandleTest()
}

func HandleItems() {
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM tovars").Scan(&count)
	if err != nil {
		Printfln("Failed to get count")
		return
	}
	for i := 0; i < count; i++ {
		j := i + 1
		var mu sync.Mutex
		clients := make(map[*websocket.Conn]bool)
		router.HandleFunc("/item_"+strconv.Itoa(j), func(w http.ResponseWriter, r *http.Request) {
			var likesCount int
			tovar := Tovar{}
			//Получаю товар из БД
			row := conn.QueryRow(fmt.Sprintf("SELECT * FROM tovars WHERE id = %v", j))
			//Printfln("%v", err)
			err = row.Scan(&tovar.Id, &tovar.ImgRef, &tovar.Name, &tovar.Description, &tovar.Price)
			if err != nil {
				Printfln("%v FFFF", err.Error())
			}
			// connection.Close()
			row = conn.QueryRow(fmt.Sprintf("SELECT COALESCE(SUM(is_like), 0) AS total_likes FROM likes WHERE tovar_id = %v", j))
			err = row.Scan(&likesCount)
			if err != nil {
				Printfln("%v QQQQQ", err.Error())
			}

			fmt.Println(tovar)
			//Создаю шаблон на основе товара из БД

			tmpl, err := template.ParseFiles("templates/infoCard.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			templateData, err := os.ReadFile("templates/defines.html")
			if err != nil {
				panic(err)
			}
			tmpl = template.Must(tmpl.Parse(string(templateData)))
			// err = tmpl.Execute(w, &Context{r, "hi", "Товар", &tovars})

			session, _ := store.Get(r, "authdata")
			interfaceEmail, authorized := session.Values["email"]
			email := ""
			if authorized {
				email = interfaceEmail.(string)
			}
			//Printfln("Total likes: %v", likesCount)
			err = tmpl.Execute(w, &struct {
				Request    *http.Request
				Title      string
				ThisTovar  *Tovar
				Email      string
				Authorized bool
				LikesCount int
				ItemId     int
				ThisPage   string
			}{
				Request:    r,
				Title:      "Товар",
				ThisTovar:  &tovar,
				Email:      email,
				Authorized: authorized,
				LikesCount: likesCount,
				ItemId:     j,
				ThisPage:   thisPage,
			})

			if err != nil {
				Printfln("%v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		}).Methods("GET")
		router.HandleFunc("/item_"+strconv.Itoa(j), RegisterLoginHandler).Methods("POST")

		router.HandleFunc("/ws_item_"+strconv.Itoa(j), func(w http.ResponseWriter, r *http.Request) {
			//Printfln("Enter to вебсоккет")
			connWebSocket, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				Printfln("Error1: %v", err.Error())
				return
			}
			mu.Lock()
			clients[connWebSocket] = true
			mu.Unlock()

			defer func() {
				mu.Lock()
				delete(clients, connWebSocket)
				mu.Unlock()
				connWebSocket.Close()
			}()
			session, _ := store.Get(r, "authdata")
			userId := session.Values["id"].(int)

			setLike := func(num int) {
				// var likesCount int
				// row, err := conn.Query(fmt.Sprintf("SELECT SUM(is_like) AS total_likes FROM likes WHERE tovar_id = %v", j))
				// row.Next()
				// err = row.Scan(&likesCount)
				// if err != nil {
				// 	Printfln("%v", err.Error())
				// }
				// likesCount += num
				// return

				//Если существует установить значение на num
				//Printfln("Вошли в SetLike")
				var exists bool
				row := conn.QueryRow(fmt.Sprintf("SELECT EXISTS ( SELECT 1 FROM likes WHERE tovar_id = %v AND user_id = %v)", j, userId))
				err = row.Scan(&exists)
				if err != nil {
					Printfln("%v ERRRR!!!", err.Error())
				}
				if exists {
					_, err = conn.Exec(fmt.Sprintf("UPDATE likes SET is_like = %v WHERE tovar_id = %v AND user_id = %v", num, j, userId))
					if err != nil {
						Printfln("Error 1z: %v", err.Error())
						return
					}
				} else {
					Printfln("Запись не существует")
					_, err = conn.Exec(fmt.Sprintf("INSERT INTO likes (tovar_id, user_id, is_like) VALUES (%v, %v, %v)", j, userId, num))
					if err != nil {
						Printfln("Error 2f: %v", err.Error())
						return
					}
				}
			}

			broadcastMessage := func(message string) {
				mu.Lock()
				defer mu.Unlock()
				for client := range clients {
					if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
						Printfln("Write error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}

			for {
				//Printfln("Бесконечный цикл")
				_, message, err := connWebSocket.ReadMessage()
				//Printfln("ReadMessage")
				if err != nil {
					Printfln("Read error: %v", err)
					break
				}
				//Printfln("%v", string(message))
				mu.Lock()
				if string(message)[0] == 'l' {
					setLike(1)
				} else if string(message)[0] == 'd' {
					setLike(-1)
				} else {
					Printfln("ELSE")
					setLike(100)
				}
				mu.Unlock()
				var currentNum int
				row := conn.QueryRow(fmt.Sprintf("SELECT COALESCE(SUM(is_like), 0) AS total_likes FROM likes WHERE tovar_id = %v", j))
				err = row.Scan(&currentNum)
				broadcastMessage(strconv.Itoa(currentNum))
			}
		}).Methods("GET")
	}
}

func RegisterLoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("emailEnter")
	if email != "" {
		loginHandler(w, r)
		return
	}
	email = r.FormValue("emailReg")
	if email != "" {
		registerHandler(w, r)
		return
	}
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("emailEnter")
	password := r.FormValue("pswEnter")
	var (
		dbEmail    string
		dbPassword string
		dbId       int
	)
	row := conn.QueryRow("SELECT username, password, id FROM users WHERE username = $1", email)
	err := row.Scan(&dbEmail, &dbPassword, &dbId)
	if err != nil {
		Printfln("Error: %v", err)
		http.Error(w, "Invalid username or password 1", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		Printfln("%v", err.Error())
		http.Error(w, "Invalid username or password 2", http.StatusUnauthorized)
		return
	}
	session, _ := store.Get(r, "authdata")
	session.Values["id"] = dbId
	session.Values["authorized"] = true
	session.Values["email"] = email

	session.Save(r, w)
	Printfln("SUCCES LOGGING!!")

	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("emailReg")
	password := r.FormValue("pswReg")
	passwordRepeat := r.FormValue("pswRepeatReg")
	var dbId int
	defer http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)

	var count int
	row := conn.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", email)
	err := row.Scan(&count) //Использую dbId чтобы не создавать лишнюю переменную
	if err != nil {
		Printfln("Error: %v", err)
		http.Error(w, "Register error", http.StatusUnauthorized)
		return
	}
	if count != 0 {
		Printfln("Account already exists")
		http.Error(w, "Account already exists", http.StatusUnauthorized)
		return
	}

	if isValidEmail(email) {
		fmt.Printf("%s is a valid email\n", email)
	} else {
		fmt.Printf("%s is not a valid email\n", email)
	}

	if 6 > len(password) && len(password) > 64 {
		return
	}
	if password != passwordRepeat {
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	Printfln("HashedPassword: %v", string(hashedPassword))
	if err != nil {
		return
	}
	_, err = conn.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", email, string(hashedPassword))
	//_, err = conn.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", 2, email, string(hashedPassword))
	if err != nil {
		Printfln("Error: %v", err.Error())
		return
	}

	row = conn.QueryRow("SELECT id FROM users WHERE username = $1", email)
	err = row.Scan(&dbId)
	if err != nil {
		http.Error(w, "Invalid username or password 1", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "authdata")
	session.Values["id"] = dbId
	session.Values["authorized"] = true
	session.Values["email"] = email
	session.Save(r, w)
}

func isValidEmail(email string) bool {
	// Паттерн для проверки email
	emailPattern := `^[a-zA-Z0
	-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Компиляция регулярного выражения
	regExp := regexp.MustCompile(emailPattern)

	// Проверка на соответствие паттерну
	return regExp.MatchString(email)
}
