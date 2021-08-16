build:
	docker-compose build
	# docker-compose build --no-cache

up:
	docker-compose up -d web db

task:
	docker-compose run --rm worker /go/main

down:
	docker-compose down

bash:
	docker-compose run --rm web bash