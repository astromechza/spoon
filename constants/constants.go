package constants

// ValidPathPartRegex is the format of a part of a metric path
const ValidPathPartRegex = "[a-zA-Z0-9\\-\\_]+"

// ValidAgentPathRegex is a pattern matching the path of a metric as reported by
// an agent. It must contain one or more path segments and may begin with a dot
// to indicate a relative path.
// examples:
// - .something
// - .something.else.entirely
// - another.thing
// - singularity
const ValidAgentPathRegex = "\\.?" + ValidPathPartRegex + "(?:\\." + ValidPathPartRegex + ")*"

// ValidBasePathRegex is a stricter form of the ValidAgentPathRegex which must
// contain one or more path segments but may not start with a dot
// examples:
// - something
// - something.else
const ValidBasePathRegex = ValidPathPartRegex + "(?:\\." + ValidPathPartRegex + ")*"

// ValidAgentPathRegexStrict is ValidAgentPathRegex with string bounderies on
// either side
const ValidAgentPathRegexStrict = "^" + ValidAgentPathRegex + "$"

// ValidBasePathRegexStrict is ValidBasePathRegex with string bounderies on
// either side
const ValidBasePathRegexStrict = "^" + ValidBasePathRegex + "$"
