compose:
	echo "Starting docker compose"
	docker-compose -f docker-compose.yml up --build

create-topic:
	docker exec -it kafka1 kafka-topics -bootstrap-server kafka1:19091 --create --topic message_topic --partitions 3 --replication-factor 2
