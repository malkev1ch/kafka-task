compose:
	echo "Starting docker compose"
	docker-compose -f docker-compose.yml up --build

create-topic:
	echo "create topic message_topic"
	docker exec -it kafka1 kafka-topics -bootstrap-server kafka1:19091 --create --topic message_topic --partitions 3 --replication-factor 2

show-brokers:
	echo "show brokers"
	docker exec -it kafka1 zookeeper-shell zookeeper:2181 ls /brokers/ids

show-consumers:
	echo "show consumers"
	docker exec -it kafka1 kafka-consumer-groups --list --bootstrap-server kafka1:19091

describe-consumers:
	echo "describe consumers"
	docker exec -it kafka1 kafka-consumer-groups -describe --group consumers --bootstrap-server kafka1:19091

delete-containers:
	docker rm -f $(docker ps -aq)

delete-images:
	docker rmi -f $(docker images -aq)