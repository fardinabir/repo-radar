package finders

func trimRepoUser(suffix string) string {
	flagCut, i := 0, 0
	for i = 0; i < len(suffix); i++ {
		if suffix[i] == '/' {
			flagCut++
		}
		if flagCut == 2 {
			break
		}
	}
	//fmt.Println("trimFunc------=======----------", suffix[0:i])
	return suffix[0:i]
}
