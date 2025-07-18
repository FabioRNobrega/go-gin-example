up:
	docker-compose up --build -d

down:
	docker-compose down

shell:
	docker exec -it go1.23 bash

logs:
	docker-compose logs -f
