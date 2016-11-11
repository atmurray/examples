package main

import "strings"
import "gopkg.in/distil.v1"
import "os"

func main() {
	// Use default connection params, this makes the resulting executable
	// portable to different installations
	ds := distil.NewDISTIL(distil.FromEnvVars())

	path := os.Getenv("REF_PMU_PATH")

	if ds.StreamFromPath(path+"/L1MAG") != nil {
		instance1 := &SymmetricalComponentsDistiller{basefreq: 120}
		registration1 := &distil.Registration{
			Instance:    instance1,
			UniqueName:  "pq1_" + strings.Replace(path, "/", "_", -1),
			InputPaths:  []string{path + "/L1MAG", path + "/L1ANG", path + "/L2MAG", path + "/L2ANG", path + "/L3MAG", path + "/L3ANG"},
			OutputPaths: []string{path + "/LPMAG", path + "/LPANG", path + "/LNMAG", path + "/LNANG", path + "/LZMAG", path + "/LZANG"},
		}
		ds.RegisterDistillate(registration1)
	}

	ds.StartEngine()
}
