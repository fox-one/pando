export PANDO_DATABASE_DIALECT="mysql"
export PANDO_DATABASE_DATASOURCE="root@tcp(localhost:3306)/test?parseTime=true"

docker run                                \
    -p 3306:3306                          \
    --env MYSQL_DATABASE=test             \
    --env MYSQL_ALLOW_EMPTY_PASSWORD=yes  \
    --name mysqltest                      \
    --detach                              \
    --rm                                  \
    mysql:5.7                             \
    --character-set-server=utf8mb4        \
    --collation-server=utf8mb4_unicode_ci \

# wait mysql to launch
sleep 10

go test -v github.com/fox-one/pando/store/vault
docker kill mysqltest
