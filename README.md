# kafka


# docker-compose down
# docker-compose up --build
    MySQL (kafka-mysql-1) on port 3306.
    ZooKeeper (kafka-zookeeper-1) on port 2181.
    Kafka (kafka-kafka-1) on port 9092.
    Go app (kafka-app-1) on port 8080.

# create user logic topic
    docker exec -it kafka-kafka-1 /opt/bitnami/kafka/bin/kafka-topics.sh --create --topic user-logs --bootstrap-server kafka:9092 --partitions 4 --replication-factor 1

# verify topic
    docker exec -it kafka-kafka-1 /opt/bitnami/kafka/bin/kafka-topics.sh --list --bootstrap-server kafka:9092

# create user
    curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"email":"alice@example.com","name":"Alice Smith"}'

# get user
    curl http://localhost:8080/users/1

# update user
    curl -X PUT http://localhost:8080/users/1 -H "Content-Type: application/json" -d '{"email":"alice.new@example.com","name":"Alice Johnson"}'

# delete user
    curl -X DELETE http://localhost:8080/users/1

# verify logs in mysql
    docker exec -it kafka-mysql-1 mysql -uroot -proot -e "SELECT * FROM users.user_logs;"

