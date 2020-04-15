package controllers

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/monsterry/go-openvpn/client/config"
	"github.com/monsterry/openvpn-web-ui/lib"
	"github.com/monsterry/openvpn-web-ui/models"
)

type NewCertParams struct {
	Name string `form:"Name" valid:"Required;"`
}

type CertificatesController struct {
	BaseController
}

func (c *CertificatesController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Certificates",
	}
}

// @router /certificates/single-config/:key [get]
func (c *CertificatesController) DownloadSingleConfig() {
	name := c.GetString(":key")
	filename := fmt.Sprintf("%s.ovpn", name)

	c.Ctx.Output.Header("Content-Type", "text/plain")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	if cfgPath, err := saveClientSingleConfig(name); err == nil {
		c.Ctx.Output.Download(cfgPath, filename)
	}

}

// @router /certificates/:key [get]
func (c *CertificatesController) Download() {
	name := c.GetString(":key")
	filename := fmt.Sprintf("%s.zip", name)

	c.Ctx.Output.Header("Content-Type", "application/zip")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	zw := zip.NewWriter(c.Controller.Ctx.ResponseWriter)

	if cfgPath, err := saveClientConfig(name); err == nil {
		addFileToZip(zw, cfgPath)
	}
	keysPath := models.GlobalCfg.OVPkiPath + "/pki"
	addFileToZip(zw, keysPath+"/ca.crt")
	addFileToZip(zw, keysPath+"/issued/"+name+".crt")
	addFileToZip(zw, keysPath+"/private/"+name+".key")

	files := strings.Split(strings.Trim(string(models.GlobalCfg.OVExtraFiles), " "), " ")
	for _, file := range files {
		addFileToZip(zw, keysPath+"/"+file)
	}

	if err := zw.Close(); err != nil {
		beego.Error(err)
	}
}

func addFileToZip(zw *zip.Writer, path string) error {
	header := &zip.FileHeader{
		Name:         filepath.Base(path),
		Method:       zip.Store,
		ModifiedTime: uint16(time.Now().UnixNano()),
		ModifiedDate: uint16(time.Now().UnixNano()),
	}
	fi, err := os.Open(path)
	if err != nil {
		beego.Error(err)
		return err
	}

	fw, err := zw.CreateHeader(header)
	if err != nil {
		beego.Error(err)
		return err
	}

	if _, err = io.Copy(fw, fi); err != nil {
		beego.Error(err)
		return err
	}

	return fi.Close()
}

// @router /certificates [get]
func (c *CertificatesController) Get() {
	c.TplName = "certificates.html"
	c.showCerts()
}

func (c *CertificatesController) showCerts() {
	path := models.GlobalCfg.OVPkiPath + "/pki/index.txt"
	certs, err := lib.ReadCerts(path)
	if err != nil {
		beego.Error(err)
	}
	lib.Dump(certs)
	c.Data["certificates"] = &certs
	lib.Dump(models.GlobalCfg.ServerName)
	c.Data["serverName"] = models.GlobalCfg.ServerName
}

// @router /certificates [post]
func (c *CertificatesController) Post() {
	c.TplName = "certificates.html"
	flash := beego.NewFlash()

	cParams := NewCertParams{}
	if err := c.ParseForm(&cParams); err != nil {
		beego.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		if vMap := validateCertParams(cParams); vMap != nil {
			c.Data["validation"] = vMap
		} else {
			if err := lib.CreateCertificate(cParams.Name); err != nil {
				beego.Error(err)
				flash.Error(err.Error())
				flash.Store(&c.Controller)
			}
		}
	}
	c.showCerts()
}

// @router /certificates/renew [post]
func (c *CertificatesController) RenewCertificate() {
	c.TplName = "certificates.html"
	flash := beego.NewFlash()

	cParams := NewCertParams{}
	if err := c.ParseForm(&cParams); err != nil {
		beego.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		if vMap := validateCertParams(cParams); vMap != nil {
			c.Data["validation"] = vMap
		} else {
			if err := lib.RenewCertificate(cParams.Name); err != nil {
				beego.Error(err)
				flash.Error(err.Error())
				flash.Store(&c.Controller)
			}
		}
	}
	c.Redirect(beego.URLFor("CertificatesController.Get"), 303)
}

// @router /certificates/revoke [post]
func (c *CertificatesController) RevokeCertificate() {
	c.TplName = "certificates.html"
	flash := beego.NewFlash()

	cParams := NewCertParams{}
	if err := c.ParseForm(&cParams); err != nil {
		beego.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		if vMap := validateCertParams(cParams); vMap != nil {
			c.Data["validation"] = vMap
		} else {
			if err := lib.RevokeCertificate(cParams.Name, "unspecified"); err != nil {
				beego.Error(err)
				flash.Error(err.Error())
				flash.Store(&c.Controller)
			}
		}
	}
	c.Redirect(beego.URLFor("CertificatesController.Get"), 303)
}

func validateCertParams(cert NewCertParams) map[string]map[string]string {
	valid := validation.Validation{}
	b, err := valid.Valid(&cert)
	if err != nil {
		beego.Error(err)
		return nil
	}
	if !b {
		return lib.CreateValidationMap(valid)
	}
	return nil
}

func saveClientConfig(name string) (string, error) {
	cfg := config.New()
	cfg.ServerAddress = models.GlobalCfg.ServerAddress
	cfg.Cert = name + ".crt"
	cfg.Key = name + ".key"
	serverConfig := models.OVConfig{Profile: "default"}
	serverConfig.Read("Profile")
	cfg.Port = serverConfig.Port
	cfg.Proto = serverConfig.Proto
	cfg.Auth = serverConfig.Auth
	cfg.Cipher = serverConfig.Cipher
	cfg.Keysize = serverConfig.Keysize

	destPath := models.GlobalCfg.OVPkiPath + "/" + name + ".ovpn"
	if err := config.SaveToFile("conf/openvpn-client-config.tpl",
		cfg, destPath); err != nil {
		beego.Error(err)
		return "", err
	}

	return destPath, nil
}

func saveClientSingleConfig(name string) (string, error) {
	cfg := config.New()
	cfg.ServerAddress = models.GlobalCfg.ServerAddress
	pathString := models.GlobalCfg.OVPkiPath + "/pki"
	cfg.Cert = readCert(pathString + name + ".crt")
	cfg.Key = readCert(pathString + name + ".key")
	cfg.Ca = readCert(pathString + "ca.crt")
	serverConfig := models.OVConfig{Profile: "default"}
	serverConfig.Read("Profile")
	cfg.Port = serverConfig.Port
	cfg.Proto = serverConfig.Proto
	cfg.Auth = serverConfig.Auth
	cfg.Cipher = serverConfig.Cipher
	cfg.Keysize = serverConfig.Keysize

	destPath := models.GlobalCfg.OVPkiPath + name + ".ovpn"
	if err := config.SaveToFile("conf/openvpn-client-config.ovpn.tpl",
		cfg, destPath); err != nil {
		beego.Error(err)
		return "", err
	}

	return destPath, nil
}

func readCert(path string) string {
	buff, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		beego.Error(err)
		return ""
	}
	return string(buff)
}
