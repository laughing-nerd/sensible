# Every host must be a part of a group. The group name must be unique.

# ============= WEB SERVERS ============== #
group "Web Servers" {
host "server1" {
  address = "${address}"
  username = "root"
  password = "root"
  }

host "server2" {
  address = "${address}"
  username = "root"
  password = "root"
  }
}

# ============= DB SERVERS ============== #
group "DB Servers" {
host "db_server_1" {
  address = "${address}"
  username = "root"
  password = "root"
  }
}
