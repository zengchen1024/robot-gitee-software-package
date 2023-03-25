package messageserver

import (
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/app"
	"github.com/opensourceways/robot-gitee-software-package/softwarepkg/domain"
)

type msgToHandleNewPkg struct {
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

func (msg *msgToHandleNewPkg) toCmd() app.CmdToNewPkg {
	return app.CmdToNewPkg{
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
