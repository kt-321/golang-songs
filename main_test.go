package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUpHandler(t *testing.T) {
	url := "http://localhost:8081/api/signup"

	// テスト用の JSON ボディ作成
	b, err := json.Marshal(Form{Email: "i@i", Password: "iiiiii"})
	if err != nil {
		t.Fatal(err)
	}

	// テスト用のリクエスト作成
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(b))

	////headerをセット
	//req.Header.Set("Content-Type", "application/json")
	////httpクライアント
	//client := &http.Client{}

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()
	// ハンドラーの実行
	SignUpHandler(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	//// レスポンスの JSON ボディのテスト
	//resp := JsonResponse{}
	//if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
	//	t.Errorf("errpr: %#v, res: %#v", err, res)
	//}
	//if resp.Message != "hello world" {
	//	t.Errorf("invalid response: %#v", resp)
	//}
	//
	//t.Logf("%#v", resp)

	////実行
	//resp, err := client.Do(req)
	//if err != nil {
	//	//panic(err)
	//	//t.Error("リクエスト実行エラー")
	//	t.Error("リクエスト実行エラー", err)
	//}

	////body使わなくなったら閉じる
	//defer resp.Body.Close()

	//エラー検証
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	//ステータスコード確認
	//if resp.StatusCode != 200 {
	//	t.Error(resp.StatusCode)
	//	return
	//}
}

func TestGetUserHandler(t *testing.T) {
	//var j = []byte(`
	//{
	//  "id": "test"
	//}`)

	url := "http://localhost:8081/api/user/1"
	//url := "/api/user/1"

	//postrequest作成
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req, err := http.NewRequest("GET", url, nil)

	//headerをセット
	req.Header.Set("Content-Type", "application/json")

	//httpクライアント
	client := &http.Client{}

	//実行
	resp, err := client.Do(req)
	if err != nil {
		//panic(err)
		//t.Error("リクエスト実行エラー")
		t.Error("リクエスト実行エラー", err)
	}

	//body使わなくなったら閉じる
	defer resp.Body.Close()

	//エラー検証
	if err != nil {
		t.Error(err)
		return
	}

	//ステータスコード確認
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
		return
	}

	//レスポンスBODY取得
	body, _ := ioutil.ReadAll(resp.Body)

	actual := string(body)
	expected := `{"id":1, "createdAt": "2020-05-23T19:17:21+09:00", "updatedAt": "2020-05-23T19:21:04+09:00", "deletedAt": null, "name": "a", "email": "a@a", "age": 20, "gender": 1, "imageUrl": "", favoriteMusicAge": 1980, "favoriteArtist":"椎名林檎", "comment": テストテスト", "followings": [], "bookmarkings": []}`

	if actual != expected {
		t.Error("response error")
	}
	//return

	//t.Logf("result is %d", actual)

	//result := GetUserHandler("山澤さん")
	//expext := "こんにちは、山澤さん！！"
	//if result != expext {
	//	t.Error("\n実際： ", result, "\n理想： ", expext)
	//}
	//
	//t.Log("TestHello終了")

}

func TestLoginHandler(t *testing.T) {
	url := "http://localhost:8081/api/login"

	// テスト用の JSON ボディ作成
	b, err := json.Marshal(Form{Email: "a@a", Password: "aaaaaa"})
	if err != nil {
		t.Fatal(err)
	}

	// テスト用のリクエスト作成
	//req := httptest.NewRequest("GET", "http://example.com/?name=world", nil)
	req := httptest.NewRequest("POST", url, bytes.NewBuffer(b))

	////headerをセット
	//req.Header.Set("Content-Type", "application/json")
	////httpクライアント
	//client := &http.Client{}

	// テスト用のレスポンス作成
	res := httptest.NewRecorder()
	// ハンドラーの実行
	LoginHandler(res, req)

	// レスポンスのステータスコードのテスト
	if res.Code != http.StatusOK {
		t.Errorf("invalid code: %d", res.Code)
	}

	//// レスポンスの JSON ボディのテスト
	//resp := JsonResponse{}
	//if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
	//	t.Errorf("errpr: %#v, res: %#v", err, res)
	//}
	//if resp.Message != "hello world" {
	//	t.Errorf("invalid response: %#v", resp)
	//}
	//
	//t.Logf("%#v", resp)

	////実行
	//resp, err := client.Do(req)
	//if err != nil {
	//	//panic(err)
	//	//t.Error("リクエスト実行エラー")
	//	t.Error("リクエスト実行エラー", err)
	//}

	////body使わなくなったら閉じる
	//defer resp.Body.Close()

	//エラー検証
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	//ステータスコード確認
	//if resp.StatusCode != 200 {
	//	t.Error(resp.StatusCode)
	//	return
	//}
}
