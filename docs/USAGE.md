# ntfy-hub-mcp USAGE Guide

このドキュメントでは、`ntfy-hub-mcp` を Gemini CLI や他の AI エージェントに統合し、利用する方法について説明します。

## 1. `ntfy-hub-mcp` のセットアップ

### 1.1. ビルド
`ntfy-hub-mcp` ディレクトリで以下のコマンドを実行して実行ファイルをビルドします。

```bash
task build
```
これにより、`ntfy-hub-mcp/ntfy-hub-mcp.exe` (Windows の場合) が生成されます。

### 1.2. 環境変数の設定
`ntfy-hub-mcp` は以下の環境変数を使用して動作を設定します。MCP サーバーを起動する前にこれらの環境変数を設定してください。

| 環境変数         | 説明                                                              | デフォルト値     |
| :--------------- | :---------------------------------------------------------------- | :--------------- |
| `NTFY_URL`       | ntfy.sh サーバーのベース URL (自前サーバーの URL を指定)         | `https://ntfy.sh` |
| `NTFY_TOPIC_OUT` | エージェントから人間への通知に使用されるデフォルトのトピック名    | `agent-output`   |
| `NTFY_TOPIC_IN`  | 人間からエージェントへの指示に使用されるデフォルトのトピック名    | `agent-input`    |

**例 (PowerShell):**
```powershell
$env:NTFY_URL="https://ntfy.your-domain.com"
$env:NTFY_TOPIC_OUT="my-ai-notifications"
$env:NTFY_TOPIC_IN="my-ai-commands"
```

## 2. Gemini CLI への統合

`ntfy-hub-mcp` を Gemini CLI に統合するには、`gemini_cli_config.json` ファイルを編集して MCP サーバーとして登録します。

### 2.1. `gemini_cli_config.json` の設定例

以下の設定を `gemini_cli_config.json` に追加します。`"command"` のパスは、ビルドされた `ntfy-hub-mcp.exe` へのフルパスに置き換えてください。

```json
{
  "mcpServers": {
    "ntfy": {
      "command": "C:\workspace
tfy-hub-mcp
tfy-hub-mcp.exe",
      "env": {
        "NTFY_URL": "https://ntfy.your-domain.com",
        "NTFY_TOPIC_OUT": "my-ai-notifications",
        "NTFY_TOPIC_IN": "my-ai-commands"
      },
      "description": "MCP server for sending and receiving notifications via ntfy.sh"
    }
  }
}
```
**注意**: `env` フィールドに指定した値は、`ntfy-hub-mcp.exe` 起動時に直接環境変数として渡されます。上記の例では、`NTFY_URL` を直接指定していますが、Powershell などで事前に設定した環境変数を `ntfy-hub-mcp.exe` に継承させたい場合は、`gemini_cli_config.json` の `env` フィールドから該当する変数を削除してください。

### 2.2. エージェントからの利用方法

`ntfy-hub-mcp` が登録されると、Gemini CLI は `ntfy_publish` と `ntfy_wait_for_reply` の2つのツールを認識します。

#### 通知の送信 (`ntfy_publish`)
エージェントがユーザーに何かを通知したい場合に使用します。

**例:**
```json
{
  "tool_code": "ntfy_publish",
  "tool_name": "ntfy_publish",
  "tool_args": {
    "message": "The current task has been completed.",
    "topic": "project-status",
    "title": "Task Update"
  }
}
```
*   `message`: 必須。通知メッセージ。
*   `topic`: オプション。通知を送信するトピック。指定しない場合、`NTFY_TOPIC_OUT` が使用されます。
*   `title`: オプション。通知のタイトル。

#### 返信の待機 (`ntfy_wait_for_reply`)
エージェントがユーザーからの承認や追加の指示を必要とする場合に使用します。

**例:**
```json
{
  "tool_code": "ntfy_wait_for_reply",
  "tool_name": "ntfy_wait_for_reply",
  "tool_args": {
    "topic": "my-ai-commands",
    "timeout_seconds": 300,
    "prompt": "I need your approval to proceed. Please reply 'continue' or 'cancel'."
  }
}
```
*   `topic`: オプション。返信を待機するトピック。指定しない場合、`NTFY_TOPIC_IN` が使用されます。
*   `timeout_seconds`: オプション。返信を待機する秒数。デフォルトは 60 秒です。
*   `prompt`: オプション。返信を待つ前に、人間への指示として送信されるメッセージ。これは `NTFY_TOPIC_OUT` に公開されます。

## 3. ntfy.sh クライアントの設定

あなたのスマートフォンやデスクトップに ntfy.sh クライアントをインストールし、`NTFY_TOPIC_OUT` および `NTFY_TOPIC_IN` に設定したトピックを購読してください。
人間からの指示を送る場合は、`NTFY_TOPIC_IN` に設定したトピックに対してメッセージを送信します。
（例: `curl -d "continue" https://ntfy.your-domain.com/my-ai-commands`）

## 4. トラブルシューティング
*   `./docs/TROUBLESHOOTING.md` を参照してください。
