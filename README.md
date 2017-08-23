# DB (MySQL) query logger

How to use.

```console
$ go get -u github.com/koron/dbquerylog
$ sudo tcpdump -s 0 -l -w - "port 3306" | dbquerylog
2017-08-23T02:55:59Z    1503456959619049265     10.0.2.2:42342  10.0.2.15:3306  vagrant vagrant 108     1       2       12662   SHOW DATABASES
2017-08-23T02:55:59Z    1503456959619164414     10.0.2.2:42342  10.0.2.15:3306  vagrant vagrant 7       0       0       4271    CREATE TABLE IF NOT EXISTS users (\n\t\tid INT PRIMARY KEY AUTO_INCREMENT,\n\t\tname VARCHAR(255) UNIQUE,\n\t\tpassword VARCHAR(255)\n\t)
[WARN] ERROR: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near '' at line 1 (1064)
2017-08-23T02:55:59Z    1503456959619361930     10.0.2.2:42342  10.0.2.15:3306  vagrant vagrant 7       0       0       8546    INSERT INTO users (name, password) VALUES (?, ?)    "foo", "pass1234"
2017-08-23T02:55:59Z    1503456959619420401     10.0.2.2:42342  10.0.2.15:3306  vagrant vagrant 7       0       0       7069    INSERT INTO users (name, password) VALUES (?, ?)    "baz", "pass1234"
2017-08-23T02:55:59Z    1503456959619510382     10.0.2.2:42342  10.0.2.15:3306  vagrant vagrant 7       0       0       6352    INSERT INTO users (name, password) VALUES (?, ?)    "bar", "pass1234"
2017-08-23T02:55:59Z    1503456959619549421     10.0.2.2:42342  10.0.2.15:3306  vagrant vagrant 7       0       0       5649    INSERT INTO users (name, password) VALUES (?, ?)    "user001", "pass1234"
2017-08-23T02:55:59Z    1503456959619626550     10.0.2.2:42342  10.0.2.15:3306  vagrant vagrant 7       0       0       6033    INSERT INTO users (name, password) VALUES (?, ?)    "user002", "pass1234"
2017-08-23T02:55:59Z    1503456959619663378     10.0.2.2:42342  10.0.2.15:3306  vagrant vagrant 7       0       0       5425    INSERT INTO users (name, password) VALUES (?, ?)    "user003", "pass1234"
2017-08-23T02:55:59Z    1503456959619786417     10.0.2.2:42342  10.0.2.15:3306  vagrant vagrant 7       0       0       4739    DROP TABLE users
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
