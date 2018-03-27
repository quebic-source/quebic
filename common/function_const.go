package common

//Runtime function runtime
type Runtime string

//RuntimeJava java
const RuntimeJava = "java"

//RuntimeNodeJS nodejs
const RuntimeNodeJS = "nodejs"

//RuntimeValidate runtime validate
func RuntimeValidate(runtime Runtime) bool {

	runtimesAviable := [2]string{RuntimeJava, RuntimeNodeJS}

	for _, runtimeAviable := range runtimesAviable {

		if Runtime(runtimeAviable) == runtime {
			return true
		}

	}

	return false

}
