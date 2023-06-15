docker container run --name test-mysql -v test:/var/lib/mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=hogehoge -d mysql:latest
