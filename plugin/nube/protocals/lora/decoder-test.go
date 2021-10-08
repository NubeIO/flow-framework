package main

// TODO: make this actual test file

import (
	"fmt"
	"encoding/json"

	"github.com/NubeDev/flow-framework/plugin/nube/protocals/lora/decoder"
)

func main() {
	strings := []string{
		"AAB296C4E5094228BA0000EC0000009A2D64",
		"CCB22E0BE8071D28C48E00EBA04F2FE03C64",
		"9BB2A166E9081C28373900EA21D0CEE03C61",
		"6FB214024D0AC627370000DA000000001B61",
		"6FB214025A0AC627B70000DA00000000385C",
		"A7AAB901000000C1FD005A00570000842261",
		"A7AAB901000000C2FD00670065000088255F",
		"11AA0203000000000D543E90000000855500",
		"16AB3241D9089B27B17200D7000000B73B00",
		"19ABAA51D3089A27B17700EA000000034800",
		"20ABBC90BA089327318700EC000000FB4F00",
	}

	for _, s := range strings {
		_, all := decoder.DecodePayload(s)
		j, _ := json.MarshalIndent(all, "", "  ")
        fmt.Println(s)
		fmt.Println(string(j))
        fmt.Println()
	}
}
