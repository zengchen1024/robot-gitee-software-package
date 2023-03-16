package messageserver

import (
	"errors"
	"strconv"
	"strings"

	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

type messageOfNewPkg struct {
	Importer          string `json:"importer"`
	ImporterEmail     string `json:"importer_email"`
	PkgId             string `json:"pkg_id"`
	PkgName           string `json:"pkg_name"`
	PkgDesc           string `json:"pkg_desc"`
	SourceCodeURL     string `json:"source_code_url"`
	SourceCodeLicense string `json:"source_code_license"`
	ImportingPkgSig   string `json:"sig"`
	ReasonToImportPkg string `json:"reason_to_import"`
}

func (msg *messageOfNewPkg) toCmd() app.CmdToCreatePR {
	return app.CmdToCreatePR{
		SoftwarePkgBasic: domain.SoftwarePkgBasic{
			Id:   msg.PkgId,
			Name: msg.PkgName,
		},
		ImporterName:  msg.Importer,
		ImporterEmail: msg.ImporterEmail,
		Application: domain.SoftwarePkgApplication{
			SourceCode: domain.SoftwarePkgSourceCode{
				Address: msg.SourceCodeURL,
				License: msg.SourceCodeLicense,
			},
			PackageDesc:       msg.PkgDesc,
			ImportingPkgSig:   msg.ImportingPkgSig,
			ReasonToImportPkg: msg.ReasonToImportPkg,
		},
	}
}

type messageOfApprovedPkg struct {
	PkgId      string `json:"pkg_id"`
	PkgName    string `json:"pkg_name"`
	RelevantPR string `json:"pr"`
}

func (msg *messageOfApprovedPkg) toCmd() (cmd app.CmdToMergePR, err error) {
	sp := strings.Split(strings.TrimSuffix(msg.RelevantPR, "/"), "/")
	if len(sp) == 0 {
		err = errors.New("relevant pr is empty")

		return
	}

	prNumInt, err := strconv.Atoi(sp[len(sp)-1])
	if err == nil {
		cmd.PRNum = prNumInt
	}

	return
}

type messageOfRejectedPkg = messageOfApprovedPkg
