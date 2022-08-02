package flags

import "flag"

var LunoKey = flag.String("lunoKey", "", "luno key")
var LunoSecret = flag.String("lunoSecret", "", "luno secret")
var LunoKeyFile = flag.String("lunoKeyFile", "", "luno KeyFile")

var ApplicationName = flag.String("AppName", "LunoApplication", "")
