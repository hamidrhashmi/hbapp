# hbapp
APIs for HBVideo Solution

Import the stored procedures in mysql 
```bash
 mysql -u root -p kamailio < /home/admin/kamailio_SPs.sql
```

export stored procedures with the following command 
```
mysqldump -u root -p -d -t -n -R --all-databases > SPs.sql
```
Create mysql User and grant permissions
```mysql
CREATE USER 'apiuser'@'%' IDENTIFIED BY 'password';
GRANT CREATE ROUTINE, ALTER ROUTINE, EXECUTE, SELECT, UPDATE, INSERT, DELETE on kamailio.* TO 'apiuser'@'%';
FLUSH PRIVILEGES;
```

