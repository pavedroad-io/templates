{{define "testAll.sh"}}
#!/bin/bash

#curl -X GET -v http://localhost:8080/api/v1/namespace/pavedroad.io/{{.Name}}LIST
#curl -X GET -v http://localhost:8080/api/v1/namespace/pavedroad.io/{{.Name}}/12345
#curl -X PUT -d "{}" -v http://localhost:8080/api/v1/namespace/pavedroad.io/{{.Name}}/12345
#curl -X DELETE -d "{}" -v http://localhost:8080/api/v1/namespace/pavedroad.io/{{.Name}}/12345
{{end}}
