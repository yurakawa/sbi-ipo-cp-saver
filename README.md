# sbi-ipo-cp-saver

sbi-ipo-cp-saverは、SBI証券のIPOに自動で応募し、チャレンジポイントを獲得します。

## インストール

1. リポジトリをクローンします。
   ```bash
   git clone https://github.com/yurakawa/sbi-ipo-cp-saver.git
   cd sbi-ipo-cp-saver
   ```

2. 必要なGoモジュールをインストールします。
   ```bash
   go mod tidy
   ```

## 使用方法

### ローカルでの実行

1. 以下の3つの環境変数を設定します。
   - `SBI_USERNAME`
   - `SBI_PASSWORD`
   - `SBI_TORIHIKI_PASSWORD`

2. プロジェクトをビルドして実行します。
   ```bash
   go build -o sbi-ipo-cp-saver
   ./sbi-ipo-cp-saver
   ```

3. または、Makefileを使用して実行します。
   ```bash
   make run
   ```

### CloudRunでの実行

#### 初回

1. 以下の環境変数を設定します。
   - `GCP_PROJECT_ID`
   - `REGION`

2. 初期設定を行い、スケジュールを登録して定期実行します。また、デプロイまで行います。

   ```bash
   make initializing
   ```

3. CloudRunでジョブを実行します。

   ```bash
   make exec
   ```

#### 初回以降
1. イメージやパスワードを更新する際に使用します。

   ```bash
   make update
   ```

2. CloudRunでジョブを実行します。

   ```bash
   make exec
   ```