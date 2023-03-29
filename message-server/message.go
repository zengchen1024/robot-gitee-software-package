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
	SpecURL           string `json:"spec_url"`
	SrcRPMURL         string `json:"src_rpm_url"`
	ImportingPkgSig   string `json:"sig"`
	ReasonToImportPkg string `json:"reason_to_import"`
}

func (msg *msgToHandleNewPkg) toCmd() app.CmdToHandleNewPkg {
	return app.CmdToHandleNewPkg{
		SoftwarePkgBasic: domain.SoftwarePkgBasic{
			Id:   msg.PkgId,
			Name: msg.PkgName,
		},
		ImporterName:  msg.Importer,
		ImporterEmail: msg.ImporterEmail,
		Application: domain.SoftwarePkgApplication{
			SourceCode: domain.SoftwarePkgSourceCode{
				SpecURL:   msg.SpecURL,
				SrcRPMURL: msg.SrcRPMURL,
			},
			PackageDesc:       msg.PkgDesc,
			ImportingPkgSig:   msg.ImportingPkgSig,
			ReasonToImportPkg: msg.ReasonToImportPkg,
		},
	}
}
