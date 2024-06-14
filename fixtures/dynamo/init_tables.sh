#!/bin/bash

: ${PROJECT_NAME:=yy-go-backend-template}
: ${ENVCODE:=local}
: ${DYNAMO_ENDPOINT=http://localhost:8000/}

DIR=$(dirname $0)
UPPER_FIRST_ENVCODE=$(echo $ENVCODE | awk '{ print toupper(substr($0, 1, 1)) substr($0, 2, length($0) - 1) }')

# NOTE テーブル名はファイル名と同名
TABLES=$(find ${DIR} -maxdepth 2 -name '*.json' -exec sh -c "basename {} | cut -d. -f1" \;)

for TABLE in ${TABLES}
do
  TABLE_NAME=${TABLE}-${UPPER_FIRST_ENVCODE}
  echo $TABLE_NAME

  if aws dynamodb list-tables \
    --region ap-northeast-1 \
    --endpoint-url ${DYNAMO_ENDPOINT} \
    | grep ${TABLE_NAME} 1> /dev/null \
  ; then
    echo Delete ${TABLE_NAME}
    aws dynamodb delete-table \
      --table-name ${TABLE_NAME} \
      --endpoint-url ${DYNAMO_ENDPOINT} \
      1> /dev/null
  fi

  # テーブルの作成
  echo Create ${TABLE_NAME}
  aws dynamodb create-table \
    --table-name ${TABLE_NAME} \
    --cli-input-json file://fixtures/dynamo/structure/${TABLE}.json \
    --endpoint-url ${DYNAMO_ENDPOINT} \
    1> /dev/null

  # TTLの設定
  echo Update TTL: ${TABLE_NAME}
  aws dynamodb update-time-to-live \
    --table-name ${TABLE_NAME} \
    --time-to-live-specification AttributeName=TTL,Enabled=true \
    --endpoint-url ${DYNAMO_ENDPOINT} \
    1> /dev/null

  # テストデータの書き込み
  TEST_DATA=${DIR}/data/${ENVCODE}/${TABLE_NAME}.json
  if [[ -e $TEST_DATA ]]; then
    echo Test Data: ${ENVCODE}/${TABLE_NAME}
    aws dynamodb transact-write-items \
      --transact-items file://${TEST_DATA} \
      --endpoint-url ${DYNAMO_ENDPOINT} \
      1> /dev/null
  fi
done
