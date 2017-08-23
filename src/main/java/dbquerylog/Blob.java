package dbquerylog;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.Random;
import javax.sql.rowset.serial.SerialBlob;

public class Blob {
    public static void test0(Statement st) throws SQLException {
        st.execute(
                "CREATE TABLE IF NOT EXISTS items(" +
                "  id INT PRIMARY KEY AUTO_INCREMENT," +
                "  name VARCHAR(255) UNIQUE," +
                "  image MEDIUMBLOB" +
                ")");
    }

    static Random rand = new Random();

    public static SerialBlob newRandomBlob() throws SQLException {
        byte[] b = new byte[16];
        rand.nextBytes(b);
        return new SerialBlob(b);
    }

    public static void test1(Connection c, Statement st) throws SQLException {
        try (PreparedStatement ins = c.prepareStatement("INSERT INTO items (name, image) VALUES (?, ?)")) {
            ins.setString(1, "foo");
            ins.setBlob(2, newRandomBlob());
            ins.execute();
            ins.clearParameters();
            /*
            ins.setString(1, "bar");
            ins.setBlob(2, newRandomBlob());
            ins.execute();
            ins.clearParameters();
            ins.setString(1, "baz");
            ins.setBlob(2, newRandomBlob());
            ins.execute();
            ins.clearParameters();
            */
        }
    }

    public static void test1a(Connection c, Statement st) throws SQLException {
        try (PreparedStatement ins = c.prepareStatement("INSERT INTO items (name) VALUES (?)")) {
            ins.setString(1, "PRE1");
            ins.execute();
            ins.clearParameters();
            ins.setString(1, "PRE2");
            ins.execute();
            ins.clearParameters();
            ins.setString(1, "PRE3");
            ins.execute();
            ins.clearParameters();
        }
    }

    public static void test1b(Connection c, Statement st) throws SQLException {
        try (PreparedStatement ins = c.prepareStatement("INSERT INTO items (name) VALUES (?)")) {
            ins.setString(1, "POST1");
            ins.execute();
            ins.clearParameters();
            ins.setString(1, "POST2");
            ins.execute();
            ins.clearParameters();
            ins.setString(1, "POST3");
            ins.execute();
            ins.clearParameters();
        }
    }

    public static void test99(Statement st) throws SQLException {
        st.execute("DROP TABLE items");
    }

    public static void main(String[] args) {
        try (
                Connection c = DriverManager.getConnection("jdbc:mysql://127.0.0.1/vagrant?useSSL=false", "vagrant", "db1234");
                Statement st = c.createStatement();
            ) {
            test0(st);
            test1a(c, st);
            test1(c, st);
            test1b(c, st);
            test99(st);
        } catch (SQLException e) {
            e.printStackTrace();
        }
    }
}
