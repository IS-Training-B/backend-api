#/bin/bash

# APIにPOSTリクエストする
# sh $0 <json> <localhost:3000/以降のエンドポイント>

### 引数一覧
# $1 : エンドポイント
# $2 : リクエストボディ(JSON)

### 変数定義

# エンドポイント
ENDPOINT="localhost:3000/${1}"

# リクエストボディ(JSON形式で受け取る)
JSON="${2}"

# リクエストを出す
RESPONSE=$(curl -X POST -H "Content-Type: application/json" -d "${JSON}" "${ENDPOINT}")

echo "Response: $RESPONSE"

