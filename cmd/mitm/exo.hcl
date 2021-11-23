exo = "0.1"
components {
  container "b1" {
    build = "."
    ports = [ "44224:44224" ]
    environment = {
      PORT = "44224"
      RESPONSE = "a"
    }
  }
  #container "b2" {
    #build = "."
    #ports = [ "44225:44225" ]
    #environment = {
      #PORT = "44225"
      #RESPONSE = "b"
    #}
  #}
}
