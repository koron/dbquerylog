package dbquerylog;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;

public class Compress {
    public static void test0(Statement st) throws SQLException {
        try (ResultSet rs = st.executeQuery("SHOW DATABASES")) {
            while (rs.next()) {
                String name = rs.getString(1);
                System.out.println(String.format("test0: table:%s found", name));
            }
        }
    }

    public static void test1(Statement st) throws SQLException {
        st.execute(
                "CREATE TABLE IF NOT EXISTS users(" +
                "  id INT PRIMARY KEY AUTO_INCREMENT," +
                "  name VARCHAR(255) UNIQUE," +
                "  password VARCHAR(255)" +
                ")");
    }

    static PreparedStatement insert;
    static PreparedStatement select;
    static PreparedStatement update;
    static PreparedStatement delete;

    public static void test2(Connection c) throws SQLException {
        insert = c.prepareStatement("INSERT INTO users (name, password) VALUES (?, ?)");
        select = c.prepareStatement("SELECT name, password FROM users WHERE name LIKE ?");
        update = c.prepareStatement("UPDATE users SET name = ?, password = ? WHERE id = ?");
        delete = c.prepareStatement("DELETE FROM users WHERE id = ?");
    }

    public static void test3(Connection c) throws SQLException {
        c.prepareStatement("INSERT INTO users (name, password) VALUES (?,");
    }

    public static void insertUser(String name, String password) throws SQLException {
        try {
            insert.setString(1, name);
            insert.setString(2, password);
            insert.execute();
        } finally {
            insert.clearParameters();
        }
    }

    public static void test4() throws SQLException {
        insertUser("foo", "pass1234");
        insertUser("baz", "pass1234");
        insertUser("bar", "pass1234");
        insertUser("user001", "pass1234");
        insertUser("user002", "pass1234");
        insertUser("user003", "pass1234");
    }

    public static void test5() throws SQLException {
        select.setString(1, "user%");
        try (ResultSet rs = select.executeQuery()) {
            while (rs.next()) {
                String name = rs.getString(1);
                System.out.println("test5: " + name);
            }
        } finally {
            select.clearParameters();
        }
    }

    public static void test99(Statement st) throws SQLException {
        if (insert != null) {
            insert.close();
            insert = null;
        }
        if (select != null) {
            select.close();
            select = null;
        }
        if (update != null) {
            update.close();
            update = null;
        }
        if (delete != null) {
            delete.close();
            delete = null;
        }
        st.execute("DROP TABLE users");
    }

    public static void main(String[] args) {
        try (
                Connection c = DriverManager.getConnection("jdbc:mysql://127.0.0.1/vagrant?useSSL=false", "vagrant", "db1234");
                Statement st = c.createStatement();
            ) {
            test0(st);
            /*
            test1(st);
            test2(c);
            test3(c);
            test4();
            test5();
            test99(st);
            */
        } catch (SQLException e) {
            e.printStackTrace();
        }
    }
}
