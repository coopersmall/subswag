package secret

//tygo:emit
var _ = `import { Metadata } from "./domain.generated.ts";

export type SecretType =
  | 'environment'
  | 'stored'

export type Secret = 
  | EnvironmentSecret
  | StoredSecret
`
