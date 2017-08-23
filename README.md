# DB (MySQL) query logger

How to use.

```console
$ go get -u github.com/koron/dbquerylog
$ sudo tcpdump -s 0 -l -w - "port 3306" | dbquerylog
tcpdump: listening on lo, link-type EN10MB (Ethernet), capture size 262144 bytes
2017-08-01 12:50:11.700641772 +0000 UTC 10.0.2.2:41861  10.0.2.15:3306  vagrant 57      1       1       15167   SELECT @@max_allowed_packet
2017-08-01 12:50:11.700699346 +0000 UTC 10.0.2.2:41861  10.0.2.15:3306  vagrant 108     1       2       7495    SHOW DATABASES
2017-08-01 12:50:11.70072838 +0000 UTC  10.0.2.2:41861  10.0.2.15:3306  vagrant 7       0       0       7374    CREATE TABLE IF NOT EXISTS users (\n\t\tid INT PRIMARY KEY AUTO_INCREMENT,\n\t\tname VARCHAR(255) UNIQUE,\n\t\tpassword VARCHAR(255)\n\t)
[WARN] ERROR: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near '' at line 1 (1064)
2017-08-01 12:50:11.700885598 +0000 UTC 10.0.2.2:41861  10.0.2.15:3306  vagrant 7       0       0       12192   INSERT INTO users (name, password) VALUES (?, ?)    "foo", "pass1234"
2017-08-01 12:50:11.701014893 +0000 UTC 10.0.2.2:41861  10.0.2.15:3306  vagrant 7       0       0       12646   INSERT INTO users (name, password) VALUES (?, ?)    "baz", "pass1234"
2017-08-01 12:50:11.701141201 +0000 UTC 10.0.2.2:41861  10.0.2.15:3306  vagrant 7       0       0       10987   INSERT INTO users (name, password) VALUES (?, ?)    "bar", "pass1234"
2017-08-01 12:50:11.701184377 +0000 UTC 10.0.2.2:41861  10.0.2.15:3306  vagrant 7       0       0       5981    INSERT INTO users (name, password) VALUES (?, ?)    "user001", "pass1234"
2017-08-01 12:50:11.701219757 +0000 UTC 10.0.2.2:41861  10.0.2.15:3306  vagrant 7       0       0       36720   INSERT INTO users (name, password) VALUES (?, ?)    "user002", "pass1234"
2017-08-01 12:50:11.701287186 +0000 UTC 10.0.2.2:41861  10.0.2.15:3306  vagrant 7       0       0       5947    INSERT INTO users (name, password) VALUES (?, ?)    "user003", "pass1234"
2017-08-01 12:50:11.70131485 +0000 UTC  10.0.2.2:41861  10.0.2.15:3306  vagrant 220     3       3       11893   SELECT * FROM users WHERE name LIKE ?   "user%"
2017-08-01 12:50:11.70137674 +0000 UTC  10.0.2.2:41861  10.0.2.15:3306  vagrant 7       0       0       4744    DROP TABLE users
```

## Report format

Report format is based on TSV (tab separated values).
Each rows represent database queries.
Each columns are consisted by below:

*   start time (in RFC3339 format)
*   start time (in unix nanoseconds)
*   client address (IP address and port)
*   server address (IP address and port)
*   username
*   database name
*   total response size in byte
*   num of columns in response
*   num of rows which responded or updated
*   elapsed time (nanosecond)
*   query
*   parameters (available for prepared statement only)

## Options

*   `-select` include SELECT statemnets
*   `-debug` enable debug log
