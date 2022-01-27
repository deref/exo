$Manifest: {
  //exo: string
  environment: $Environment
  $environment: environment
  components: $Components
}

$Environment: { [string]: string }

$Components: {
  [Name=string]: $Component & {
    name: Name
  }
}

$Component: {
  type: string
  name: string
  spec: {}
}

$Stack: {
  environment: $Environment
  components: $Components
}

// TODO: Move to os package, prefix type string, etc.
$Daemon: $Component & {
  type: "daemon"
  spec: {
    program: string
    arguments: [...string] | *[]
    environment: $stack.environment
  }
}

$stack: $Stack
