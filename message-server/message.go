package messageserver

import (
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
	PRNum int `json:"pr_num"`
}

func (msg *messageOfApprovedPkg) toCmd() app.CmdToMergePR {
	return app.CmdToMergePR{
		PRNum: msg.PRNum,
	}
}

type messageOfRejectedPkg struct {
	PRNum  int    `json:"pr_num"`
	Reason string `json:"reason"`
}

func (msg *messageOfRejectedPkg) toCmd() app.CmdToClosePR {
	return app.CmdToClosePR{
		PRNum:  msg.PRNum,
		Reason: msg.Reason,
	}
}
