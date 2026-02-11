# AI Agent Interaction System with ntfy

ntfyを利用して、AIエージェントと人間がスマホやPCを通じて対話できる仕組みをGoで構築します。

## 1. アーキテクチャ

エージェントは `ntfy.sh` の特定のトピックを購読（Subscribe）し、人間からのメッセージを待ち受けます。人間からの入力を検知すると、エージェントは応答を生成し、同じ（または別）トピックにメッセージを公開（Publish）します。

### 通信フロー
1.  **人間**: ntfyアプリまたは `curl` を使ってメッセージを送信。
    - `POST https://ntfy.sh/<topic_human_to_agent>`
2.  **エージェント (Go)**: SSE (Server-Sent Events) で上記トピックを監視。
3.  **エージェント (Go)**: 応答を生成。
4.  **エージェント (Go)**: 応答を送信。
    - `POST https://ntfy.sh/<topic_agent_to_human>`

## 2. 実装状況

- [x] **Project Setup**: `go mod init` 完了。
- [x] **Client Implementation**:
    - `ntfy.Client` を実装。Publish/Subscribe (SSE) に対応。
- [x] **Main Logic**: 基本的なメッセージループを実装。
- [ ] **Advanced Features**:
    - [ ] LLM (OpenAI API等) との連携。
    - [ ] ntfy Actions (ボタン) の活用。
    - [ ] 設定ファイルの外部化 (YAML/JSON)。

## 3. 使い方

### 実行
```bash
go run main.go
```

### 人間側からの送信テスト
```bash
curl -d "Hello" https://ntfy.sh/my-ai-agent-input
```

### エージェントの応答確認
ブラウザ等で `https://ntfy.sh/my-ai-agent-output` を開いておくと、エージェントからの返信を確認できます。