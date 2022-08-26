package errs

func ExpHaveStackTrace(cause error) bool {
	return haveStackTrace(cause)
}

func ExpGetStackRecursive(cause error) StackTrace {
	return getStackRecursive(cause)
}
