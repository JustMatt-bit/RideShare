run:
	docker container start rideshare-mysql-container
	go run ./

reset:
	docker container stop rideshare-mysql-container
	docker container rm rideshare-mysql-container
	docker rmi rideshare-mysql
	docker build -t rideshare-mysql ../
	docker run -d -p 3306:3306 --name rideshare-mysql-container rideshare-mysql

up:
	docker build -t rideshare-mysql ../
	docker run -d -p 3306:3306 --name rideshare-mysql-container rideshare-mysql