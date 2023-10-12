build:
	docker build -t blog_server_img .

run:
	docker run -d -p 8080:8080 --env-file ./.env --name blog_server_ctn blog_server_img

stop:
	docker stop blog_server_ctn

start:
	docker start blog_server_ctn

destroy:
	make stop
	docker rm blog_server_ctn
	docker rmi blog_server_img

create:
	make build
	make run