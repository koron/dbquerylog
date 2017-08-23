package dbquerylog;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.Random;
import javax.sql.rowset.serial.SerialBlob;

public class DateTest {
    public static void test0(Statement st) throws SQLException {
        st.execute(
                "CREATE TABLE IF NOT EXISTS records(" +
                "  id INT PRIMARY KEY AUTO_INCREMENT," +
                "  name VARCHAR(255) UNIQUE," +
                "  ctime DATETIME" +
                ")");
    }

    public static void test1(Connection c, Statement st) throws SQLException {
        try (PreparedStatement ins = c.prepareStatement("INSERT INTO records (name, ctime) VALUES (?, ?)")) {
            ins.setString(1, "foo");
            ins.setTimestamp(2,
                    new java.sql.Timestamp(new java.util.Date().getTime()));
            ins.execute();
            ins.clearParameters();
        }
    }

    public static void test99(Statement st) throws SQLException {
        st.execute("DROP TABLE records");
    }

    public static void main(String[] args) {
        try (
                Connection c = DriverManager.getConnection("jdbc:mysql://127.0.0.1/vagrant?useSSL=false", "vagrant", "db1234");
                Statement st = c.createStatement();
            ) {
            test0(st);
            test1(c, st);
            test99(st);
        } catch (SQLException e) {
            e.printStackTrace();
        }
    }
}
