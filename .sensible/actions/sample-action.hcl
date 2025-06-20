action "Sample Action" {

  groups = ["Web Servers"]

  shell "Test shell" {
    command = "echo Hello world!"
  }

  installer "Test installer" {
    preferred = "apt"
    packages = ["cowsay", "git"]
  }

  cron "Test crom" {
    expression = "* * * * *" 
    job = "echo This is a cron job >> ~/random-stuffs/sensible/cron-test.log"
    type = "add"
  }

}
