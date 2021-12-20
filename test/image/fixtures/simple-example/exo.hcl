exo = "0.1"

components {
  component "echo" {
    type = "process"
    spec = jsonencode({
      program   = "python3"
      arguments = ["./listen.py"]
      environment = {
        "PORT": "44222"
      }
    })
  }

  process "echo-short" {
    program   = "python3"
    arguments = ["./listen.py"]
    environment = {
      "PORT": "44223"
    }
  }

  container "web" {
    build = "."
    ports = [ "44224:44224" ]
    environment = {
      PORT = "44224"
    }
    env_file = "./env"
  }

  volume "logvolume01" {}
}
