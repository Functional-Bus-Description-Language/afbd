package json

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/Functional-Bus-Description-Language/afbd/internal/args"

	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/fn"
	"github.com/Functional-Bus-Description-Language/go-fbdl/pkg/fbdl/pkg"
)

func Generate(bus *fn.Block, pkgsConsts map[string]*pkg.Package) {
	err := os.MkdirAll(args.Json.Path, os.FileMode(int(0775)))
	if err != nil {
		log.Fatalf("generate reg json: %v", err)
	}

	regFile, err := os.Create(path.Join(args.Json.Path, "reg.json"))
	if err != nil {
		log.Fatalf("generate reg json: %v", err)
	}

	byteArray, err := json.MarshalIndent(bus, "", "\t")
	if err != nil {
		log.Fatalf("generate reg json: %v", err)
	}

	_, err = regFile.Write(byteArray)
	if err != nil {
		log.Fatalf("generate reg json: %v", err)
	}

	err = regFile.Close()
	if err != nil {
		log.Fatalf("generate reg json: %v", err)
	}

	constsFile, err := os.Create(path.Join(args.Json.Path, "const.json"))
	if err != nil {
		log.Fatalf("generate constants json: %v", err)
	}

	byteArray, err = json.MarshalIndent(pkgsConsts, "", "\t")
	if err != nil {
		log.Fatalf("generate constants json: %v", err)
	}

	_, err = constsFile.Write(byteArray)
	if err != nil {
		log.Fatalf("generate constants json: %v", err)
	}

	err = constsFile.Close()
	if err != nil {
		log.Fatalf("generate constants json: %v", err)
	}
}
