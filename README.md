# go-backend-template 
<a href="https://xcfile.dev"><img src="https://xcfile.dev/badge.svg" alt="xc compatible" /></a>
このテンプレートは、[@ras0q さんのテンプレート](https://github.com/ras0q/go-backend-template)を改変し、何もしなくてもNeoShowcase上でデプロイできるようにしたものです。
さらに、Swagger/OpenAPIドキュメントの自動生成にも対応しています。

> **Note**: シンプルにNeoShowcase対応だけしたバージョンは[neoshowcaseブランチ](https://github.com/pikachu0310/go-backend-template/tree/neoshowcase)をご利用ください。

以下の様に、何もしなくても正常に動きます(MariaDBも自動で環境変数を見て繋がります)。
![image](https://github.com/pikachu0310/go-backend-template/assets/17543997/dee159b2-598c-40ed-807a-9b5680f465a8)

ハッカソンなど短期間でWebアプリを開発する際のバックエンドのGo実装例です。
学習コストと開発コストを抑えることを目的としています。

## How to use

GitHubの `Use this template` ボタンからレポジトリを作成します。

[`gonew`](https://pkg.go.dev/golang.org/x/tools/cmd/gonew) コマンドからでも作成できます。`gonew` コマンドを使うと、モジュール名を予め変更した状態でプロジェクトを作成することができます。

```sh
gonew github.com/pikachu0310/go-backend-template {{ project_name }}
```

## Requirements

最低限[Docker](https://www.docker.com/)と[Docker Compose](https://docs.docker.com/compose/)が必要です。
[Compose Watch](https://docs.docker.com/compose/file-watch/)を使うため、Docker Composeのバージョンは2.22以上にしてください。

Linter, Formatterには[golangci-lint](https://golangci-lint.run/)を使っています。
VSCodeを使用する場合は`.vscode/settings.json`でLinterの設定を行ってください

```json
{
  "go.lintTool": "golangci-lint"
}
```

## Directory structure

[Organizing a Go module - The Go Programming Language](https://go.dev/doc/modules/layout#server-project) などを参考にしています。

```bash
$ tree | manual-explain
.
├── main.go # エントリーポイント
├── core # アプリケーション本体
│   ├── database # DBの初期化、マイグレーション
│   │   └── migrations # DBマイグレーションのスキーマ
│   ├── internal # ロジック (結合テストに公開する必要がないもの)
│   │   ├── handler # APIハンドラ
│   │   ├── repository # DBアクセス
│   │   └── services # 外部サービス, 複雑なビジネスロジック
│   └── config.go, deps.go, router.go # セットアップ (統合テスト用に公開)
├── frontend # フロントエンド
└── integration_tests # 結合テスト
```

特に重要なものは以下の通りです。

### `main.go`

アプリケーションのエントリーポイントを配置します。

**Tips**: 複数のエントリーポイントを実装する場合は、`cmd` ディレクトリを作成し、各エントリーポイントを `cmd/{app name}/main.go` に書くと見通しが良くなります。

### `core/internal/`

アプリケーション本体のロジックを配置します。
主に2つのパッケージに分かれています。

- `handler/`: ルーティング
  - 飛んできたリクエストを裁いてレスポンスを生成する
  - DBアクセスは`repository/`で実装したメソッドを呼び出す
  - **Tips**: リクエストのバリデーションがしたい場合は↓のどちらかを使うと良い
    - [go-playground/validator](https://github.com/go-playground/validator): タグベースのバリデーション
    - [go-ozzo/ozzo-validation](https://github.com/go-ozzo/ozzo-validation): コードベースのバリデーション
- `repository/`: ストレージ操作
  - DBや外部ストレージなどのストレージにアクセスする
    - 引数のバリデーションは`handler/`に任せる

**Tips**: `internal`パッケージは他モジュールから参照されません（参考: [Go 1.4 Release Notes](https://go.dev/doc/go1.4#internalpackages)）。
依存性注入や外部ライブラリの初期化のみを`core/`や`pkg/`で公開し、アプリケーションのロジックは`internal/`に閉じることで、後述の`integration_tests/go.mod`などの外部モジュールからの参照を最小限にすることができ、開発の効率を上げることができます。

### `core/database`

DBスキーマの定義、DBの初期化、マイグレーションを行っています。

マイグレーションツールは[pressly/goose](https://github.com/pressly/goose)を使っています。

### `integration_tests/`

結合テストを配置します。
APIエンドポイントに対してリクエストを送り、レスポンスを検証します。
短期開発段階では時間があれば書く程度で良いですが、長期開発に向けては書いておくと良いでしょう。

```go
package integration_tests

import (
  "testing"
  "gotest.tools/v3/assert"
)

func TestUser(t *testing.T) {
  t.Run("get users", func(t *testing.T) {
    t.Run("success", func(t *testing.T) {
      t.Parallel()
      rec := doRequest(t, "GET", "/api/v1/users", "")

      expectedStatus := `200 OK`
      expectedBody := `[{"id":"[UUID]","name":"test","email":"test@example.com"}]`
      assert.Equal(t, rec.Result().Status, expectedStatus)
      assert.Equal(t, escapeSnapshot(t, rec.Body.String()), expectedBody)
    })
  })
}
```

**Tips**: DBコンテナの立ち上げには[ory/dockertest](https://github.com/ory/dockertest)を使っています。

**Tips**: アサーションには[gotest.tools](https://github.com/gotestyourself/gotest.tools)を使っています。
`go test -update`を実行することで、`expectedXXX`のスナップショットを更新することができます（参考: [gotest.toolsを使う - 詩と創作・思索のひろば](https://motemen.hatenablog.com/entry/2022/03/gotest-tools)）。

外部サービス（traQ, Twitterなど）へのアクセスが発生する場合はTest Doublesでアクセスを置き換えると良いでしょう。

## API Documentation (Swagger/OpenAPI)

このテンプレートでは[swaggo/swag](https://github.com/swaggo/swag)を使用してSwagger/OpenAPIドキュメントを自動生成しています。

> **Note**: `docs/` ディレクトリは自動生成されますが、すぐに開発を始められるようにリポジトリにコミットしています。APIのアノテーションを変更した後は `swag init` を実行してドキュメントを更新し、変更をコミットしてください。

### アノテーションの書き方

各ハンドラー関数にコメント形式でSwaggerアノテーションを追加します：

```go
// GetUser godoc
// @Summary ユーザー情報取得
// @Description 指定したIDのユーザー情報を取得します
// @Tags users
// @Accept json
// @Produce json
// @Param userID path string true "User ID" format(uuid)
// @Success 200 {object} GetUserResponse "ユーザー情報"
// @Failure 400 {object} echo.HTTPError "Bad Request"
// @Router /users/{userID} [get]
func (h *Handler) GetUser(c echo.Context) error {
    // ...
}
```

詳しい書き方は[Swagドキュメント](https://github.com/swaggo/swag#declarative-comments-format)を参照してください。

### ドキュメントの生成と更新

`swag init` コマンドでドキュメントを生成します。生成されたファイルは `docs/` ディレクトリに配置され、アプリケーションに組み込まれます。

```bash
swag init
```

APIのアノテーションを変更した後は、このコマンドを実行してドキュメントを更新し、変更をコミットしてください。
開発時は `/swagger/index.html` にアクセスすることでSwagger UIからAPIをテストできます。

## Tasks

開発に用いるコマンド一覧

> [!TIP]
> `xc` を使うことでこれらのコマンドを簡単に実行できます。
> 詳細は以下のページをご覧ください。
>
> - [xc](https://xcfile.dev)
> - [MarkdownベースのGo製タスクランナー「xc」のススメ](https://zenn.dev/trap/articles/af32614c07214d)
>
> ```bash
> go install github.com/joerdav/xc/cmd/xc@latest
> ```

### Build-UI

フロントエンドをビルドします。

directory: ./frontend/app-ui

```sh
npm install
npm run build
```

### Build

アプリをビルドします。

requires: Build-UI

```sh

CMD=server
go mod download
go build -o ./bin/${CMD} ./main.go
```

### Generate-Swagger

Swagger/OpenAPIドキュメントを生成・更新します。
APIのアノテーションを変更した際に実行してください。

```sh
swag init
```

### Dev

ホットリロードの開発環境を構築します。

```sh
docker compose watch
```

API、DB、DB管理画面が起動します。
各コンテナが起動したら、以下のURLにアクセスすることができます。
Compose Watchにより、ソースコードの変更を検知して自動で再起動します。

- <http://localhost:8080/> (API)
- <http://localhost:8080/swagger/index.html> (Swagger UI)
- <http://localhost:8081/> (DBの管理画面)

### Test

全てのテストを実行します。

```sh
go test -v -cover -race -shuffle=on ./...
```

### Test-Unit

単体テストを実行します。

```sh
go test -v -cover -race -shuffle=on ./core/...
```

### Test-Integration

結合テストを実行します。

```sh
[ ! -e ./go.work ] && go work init . ./integration_tests
go test -v -cover -race -shuffle=on ./integration_tests/...
```

### Test-Integration:Update

結合テストのスナップショットを更新します。

```sh
[ ! -e ./go.work ] && go work init . ./integration_tests
go test -v -cover -race -shuffle=on ./integration_tests/... -update
```

### Lint

Linter (golangci-lint) を実行します。

```sh
golangci-lint run --timeout=5m --fix ./...
```

## Improvements

長期開発に向けた改善点をいくつか挙げておきます。

- ドメインを書く (`internal/domain/`など)
  - 現在は簡単のためにAPIスキーマとDBスキーマのみを書きこれらを直接やり取りしている
  - 本来はアプリの仕様や概念をドメインとして書き、スキーマの変換にはドメインを経由させるべき
- クライアントAPIスキーマを更に活用する
  - 現在はSwagger/OpenAPIでドキュメントを生成している
  - さらに [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) などでOpenAPIからGoコードを生成することで型安全性を高めることも可能
- 単体テスト・結合テストのカバレッジを上げる
  - カバレッジの可視化には[Codecov](https://codecov.io)(traPだと主流)や[Coveralls](https://coveralls.io)が便利
- ログの出力を整備する
  - ロギングライブラリは好みに合ったものを使うと良い

## Deploy

GitHub Actions でビルドした Docker image を Heroku や NeoShowcase のような PaaS にデプロイするためのJobを用意しています。

詳しくは [./.github/workflows/image.yaml](./.github/workflows/image.yaml) を参照してください。

## Troubleshooting

### Docker imageのビルドが失敗する

不要なファイルがDocker imageに混入するのを防ぐために `.dockerignore` をallowlist方式にしています。
ビルドに必要なファイルやディレクトリを追加した場合、 `.dockerignore` も編集してください。

### それでも解決しなければ？

Discussionに[Q&A](https://github.com/ras0q/go-backend-template/discussions/categories/q-a) を用意しています。

ここに @ras0q へのメンションを付けて質問を投げてくれれば答えます。

その他SNSなどの手段での質問ももちろん可能ですが、全体に公開されているDiscussionを推奨しています。
