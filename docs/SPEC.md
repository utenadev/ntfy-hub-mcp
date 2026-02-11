# ntfy-hub-mcp SPECIFICATIONS

## 1. プロトコルバージョン
MCP プロトコルバージョン: `2025-06-18` (または最新の安定版)
ntfy プロトコルバージョン: 1

## 2. 環境変数

`ntfy-hub-mcp` は以下の環境変数をサポートします。

| 環境変数         | 説明                                                              | デフォルト値     |
| :--------------- | :---------------------------------------------------------------- | :--------------- |
| `NTFY_URL`       | ntfy.sh サーバーのベース URL (自前サーバーの URL を指定)         | `https://ntfy.sh` |
| `NTFY_TOPIC_OUT` | エージェントから人間への通知に使用されるデフォルトのトピック名    | `agent-output`   |
| `NTFY_TOPIC_IN`  | 人間からエージェントへの指示に使用されるデフォルトのトピック名    | `agent-input`    |

## 3. 提供される MCP ツール

### 3.1. `ntfy_publish`

指定されたトピックにメッセージを公開し、人間への通知を行います。

#### 説明 (Description)
Send a message to a ntfy topic (e.g., for notifications)

#### 引数 (Arguments)

| パラメータ | 型     | 必須 | 説明                                                                    |
| :--------- | :----- | :--- | :---------------------------------------------------------------------- |
| `message`  | string | はい  | 通知として送信するメッセージの内容                                      |
| `topic`    | string | いいえ | メッセージを公開するトピック名 (デフォルトは `NTFY_TOPIC_OUT` の値)     |
| `title`    | string | いいえ | 通知のタイトル (ntfy クライアントに表示される)                          |

#### 戻り値 (Result)
*   **成功**: string 型で、`"Message sent to topic 'TOPIC_NAME'"` の形式の成功メッセージ。
*   **失敗**: Error 型で、エラーメッセージ。`message` が空の場合や、ntfy.sh への publish に失敗した場合。

#### 使用例
```json
{
  "tool_code": "ntfy_publish",
  "tool_name": "ntfy_publish",
  "tool_args": {
    "message": "The build has completed successfully.",
    "topic": "my-project-alerts",
    "title": "Build Status"
  }
}
```

### 3.2. `ntfy_wait_for_reply`

特定のトピックで人間からの返信を待機します。人間の承認や入力が必要な場合に利用します。

#### 説明 (Description)
Wait for a reply from the human on a specific topic. Use this to get user input or approval.

#### 引数 (Arguments)

| パラメータ        | 型     | 必須 | 説明                                                                    |
| :---------------- | :----- | :--- | :---------------------------------------------------------------------- |
| `topic`           | string | いいえ | リッスンするトピック名 (デフォルトは `NTFY_TOPIC_IN` の値)              |
| `timeout_seconds` | number | いいえ | 返信を待機する時間 (秒単位)。デフォルトは `60` 秒。                     |
| `prompt`          | string | いいえ | 返信を待つ前に人間に送信するメッセージ (これは `NTFY_TOPIC_OUT` に公開されます)。 |

#### 戻り値 (Result)
*   **成功**: string 型で、受信したメッセージの内容。
*   **失敗**: Error 型で、エラーメッセージ。タイムアウトした場合や、ntfy.sh からの購読中にエラーが発生した場合。

#### 使用例
```json
{
  "tool_code": "ntfy_wait_for_reply",
  "tool_name": "ntfy_wait_for_reply",
  "tool_args": {
    "topic": "my-project-commands",
    "timeout_seconds": 120,
    "prompt": "Please confirm if I should proceed with deployment. Reply 'yes' or 'no'."
  }
}
```
