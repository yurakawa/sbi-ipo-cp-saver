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

#### 初期設定

1. 以下の環境変数を設定します。
   - `SBI_USERNAME`
   - `SBI_PASSWORD`
   - `SBI_TORIHIKI_PASSWORD`
   - `GCP_PROJECT_ID`
   - `REGION`

2. 初期設定を行い、スケジュールを登録して定期実行します。また、デプロイまで行います。

   ```bash
   make initializing
   ```

#### 初回以降
1. イメージやパスワードの更新を反映したい場合

   ```bash
   make update
   ```
