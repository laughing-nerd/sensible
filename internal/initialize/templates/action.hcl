// for more information on the syntax and the list of available components,
// you can refer to the documentation at https://sensible.dev/docs/action

action "Sample Action" {
  groups = ["Web Server 1"] // Omit this line if you want the action to run locally on your machine

  shell "Test shell" {
    command = "echo Hello world!"
  }

  installer "Test installer" {
    preferred = "apt"
    packages = ["cowsay", "git"]
  }

  cron "Test cron" {
    expression = "* * * * *" 
    job = "echo This is a cron job >> ~/random-stuffs/sensible/cron-test.log"
    type = "add" 
  }

}
