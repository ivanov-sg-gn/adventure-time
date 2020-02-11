# Console
go run parser https://ivanov-host.ru/


# DB
./connectionDB/connectionDB.go
	dbLogin = "golang"
	dbPassword = "password"
	dbName = "golang"
	
	
# DB tables
./structure.sql


# Clear tables
DELETE FROM `Links` WHERE 1; ALTER TABLE `Links` AUTO_INCREMENT=1;
DELETE FROM `PhotoLinks` WHERE 1; ALTER TABLE `PhotoLinks` AUTO_INCREMENT=1;
