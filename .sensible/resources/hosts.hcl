# Every host must be a part of a group. The group name must be unique.

# ============= WEB SERVERS ============== #
group "Web Servers" {
host "server1" {
  address = "${address}"
  username = "root"
  password = "root"
  timeout = 1
  }

host "server2" {
  address = "${address}"
  username = "root"
  password = "root"
  timeout = 1
  }
}

# ============= DB SERVERS ============== #
group "DB Servers" {
host "db_server_1" {
  address = "${address}"
  username = "root"
  password = "root"
  timeout = 1
  }
}
