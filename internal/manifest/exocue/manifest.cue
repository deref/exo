#Manifest: {
  exo: {
    format: "exo.deref.io/cue-v1"
  }

  environment: #EnvironmentExpression
  components: #ComponentsByName
}

#EnvironmentExpression: {
  [string]: #EnvironmentVariableExpression
}

#EnvironmentVariableExpression: #EnvironmentVariableValue | null | #EnvironmentVariableSecret
#EnvironmentVariableValue: string
#EnvironmentVariableSecret: _|_ // TODO: Secret environment variables.

#Component: {
  type: string
  name: string
  spec: #Model
  run: bool | *true
  environment: #Environment
}

#Model: { [string]: _ }

#ComponentsByName: {
  [Name=string]: #Component & {
    name: Name
  }
}
