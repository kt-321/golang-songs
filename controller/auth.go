package controller


//func gormConnect() *gorm.DB {
//	db, err := gorm.Open("mysql", "root:@/laravel6?charset=utf8&parseTime=True&loc=Local")
//	if err != nil {
//		panic(err.Error())
//	}
//	fmt.Println("接続成功")
//	return db
//}
//
//// レスポンスにエラーを突っ込んで、返却するメソッド
//func errorInResponse(w http.ResponseWriter, status int, error model.Error) {
//	w.WriteHeader(status)
//	json.NewEncoder(w).Encode(error)
//	return
//}

//func SignUpHandler(w http.ResponseWriter, r *http.Request) {
//	user := model.User{}
//	error := model.Error{}
//
//	email := r.FormValue("email")
//	password := r.FormValue("password")
//
//	fmt.Println("signup中")
//	fmt.Println(email)
//	fmt.Println(password)
//
//	//使えるようにしたい
//	//fmt.Println(r.Body)
//
//	//一旦断念するが後でしたいJson.NewEncoder, json.NewDecoder で、エンコード(構造体から文字列)、デコード(文字列から構造体)の処理を行なっている。
//	//json.NewDecoder(r.Body).Decode(&user)
//
//	if email == "" {
//		//if user.Email == "" {
//		error.Message = "Emailは必須です。"
//		errorInResponse(w, http.StatusBadRequest, error)
//		return
//	}
//
//	if password == "" {
//		//if user.Password == "" {
//		error.Message = "パスワードは必須です。"
//		errorInResponse(w, http.StatusBadRequest, error)
//		return
//	}
//
//	fmt.Println(user)
//
//	// dump も出せる
//	fmt.Println("---------------------")
//	spew.Dump(user)
//
//	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
//	//hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("パスワード: ", password)
//	//fmt.Println("パスワード: ", user.Password)
//	fmt.Println("ハッシュ化されたパスワード", hash)
//
//	user.Email = email
//	user.Password = string(hash)
//	//ほんとはuserに入れて処理していきたい
//	password = string(hash)
//	//fmt.Println("コンバート後のパスワード: ", user.Password)
//	fmt.Println("コンバート後のパスワード: ", password)
//
//	db := gormConnect()
//	defer db.Close()
//	db.Create(&model.User{Email: email, Password: password})
//	//TOD　エラーハンドリング
//	//err = db.Create(&User{Email: email, Password: password})
//	//if err != nil {
//	//	error.Message = "サーバーエラー"
//	//	errorInResponse(w, http.StatusInternalServerError, error)
//	//	return
//	//}
//	// DB に登録できたらパスワードをからにしておく
//	user.Password = ""
//	w.Header().Set("Content-Type", "application/json")
//	// 使えなかった JSON 形式で結果を返却
//	//responseByJSON(w, user)
//
//	v, err := json.Marshal(user)
//	if err != nil {
//		println(string(v))
//	}
//	w.Write(v)
//}