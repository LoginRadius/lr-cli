package cmdutil

import (
	"errors"
	"encoding/binary"
	"github.com/loginradius/lr-cli/config"
	"regexp"
	"strings"
	"net"
)

var ThemeMap = map[string]string{
	"0": "Template_1", // Handled fallback logic to Template_1.
	"1": "Template_2",
	"2": "Template_3",
	"3": "Template_4",
	"4": "Template_5",
}

var conf = config.GetInstance()



var DomainValidation = regexp.MustCompile(`^((([\S]+:\/\/?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]*|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-]*)?\??(?:[-\+=&;%@.\w]*)#?(?:[\w]*))?)`)

var AccessRestrictionDomain = regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\")){0,}@{0,1}((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)

var ValidateEmail = regexp.MustCompile(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)

var ValidateIP = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[0-9]{1,2})\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[0-9]{1,2})$`)


func isNotLeadingZeroIp (ip string) bool {
	octets := strings.Split(ip, ".")
  for _,val := range octets {
    if (len(val) > 1 && strings.HasPrefix(val, "0") ) {
      return false
    }
  }
  return true
}

func CompareTwoIPs (ip1 string, ip2 string) bool {
	newip1 := net.ParseIP(ip1)
	newip2 := net.ParseIP(ip2)
	ip41 := newip1.To4()
	ip42 := newip2.To4()
  ip1Int := binary.BigEndian.Uint32(ip41)
  ip2Int := binary.BigEndian.Uint32(ip42)
  if ip1Int < ip2Int {
    return true
  } else {
    return false
  }
}

func ValidateIPorIPRange (ip string) (bool, error) {
	if strings.Contains(ip, "-") {
		ips := strings.Split(ip, "-")
		ip1 := strings.TrimSpace(ips[0])
		ip2 := strings.TrimSpace(ips[1])
		isLowerValidIP := ValidateIP.MatchString(ip1)
      	isHighetValidIP := ValidateIP.MatchString(ip2)
      if isNotLeadingZeroIp(ip1) && isNotLeadingZeroIp(ip2) {
        	if isLowerValidIP && isHighetValidIP && len(ips) == 2 {
        	  if CompareTwoIPs(ip1, ip2) {
        	    return true, nil
        	  } else {
				  return false,
				  errors.New("Start IP should be lower than the end IP")
        	  }
        	} else {
				return false,
				errors.New("IP address or IP range must be in a valid format, example - 192.168.0.1 or 192.168.0.1-192.168.0.255")
        	}
		} else {
			 
			return false, errors.New("IP address or IP range must be in a valid format")
		}
		

	} else {
		if isNotLeadingZeroIp(ip) {
			if ValidateIP.MatchString(ip){
				return true, nil
			} else {
				return false,
				errors.New("IP address or IP range must be in a valid format, example - 192.168.0.1 or 192.168.0.1-192.168.0.255")
			}

		} else {
			return false,
			errors.New("IP address or IP range must be in a valid format")
		}
	}
}



type ThemeType struct {
	PageType     string        `json:"PageType"`
	CustomCss    []string      `json:"CustomCss"`
	HeadTags     []interface{} `json:"HeadTags"`
	FavIcon      string        `json:"FavIcon"`
	HtmlBody     string        `json:"HtmlBody"`
	EndScript    string        `json:"EndScript"`
	BeforeScript string        `json:"BeforeScript"`
	CustomJS     []string 	   `json:"CustomJS"`
	IsActive     bool          `json:"IsActive"`
	MainScript   string        `json:"MainScript"`
	Status       string        `json:"Status"`
}

type SmtpProviderSchema struct {
	Name 		string 		
	Display 	string 	
	SmtpHost 	string 		
	SmtpPort 	string 		
	EnableSSL 	bool 		
}

  
var SmtpOptionNames = map[string]string {
	"FromName": 	"From Name",
	"FromEmailId": 	"From Email Id",
	"UserName": 	"SMTP User Name",
	"Password": 	"SMTP Password",
	"SmtpHost": 	"SMTP Host",
	"SmtpPort": 	"SMTP Port",
	"IsSsl": 		"Enable SSL",
	"Provider": 	"SMTP Providers",
	"Key": 			"Key",
	"Secret": 		"Secret",
}

var PermissionCommands = map[string]string {
"lr_demo":		 				"API_ViewConfiguration",
"lr_reset_secret":	 			"API_EditCredentials",
"lr_verify": 					"UserManagement_Admin",
"lr_get_config":				"API_EditCredentials",
"lr_get_domain":				"API_ViewConfiguration",
"lr_get_server-info": 			"API_ViewConfiguration",
"lr_get_theme": 				"API_ViewConfiguration",
"lr_get_social": 				"API_ViewThirdPartyCredentials",
"lr_get_sott":					"API_ViewConfiguration",
"lr_get_hooks":					"ThirdPartyIntegration_View",
"lr_get_schema":				"API_ViewConfiguration",
// "lr_get_site":					"API_ViewConfiguration",
"lr_get_login-method":			"API_ViewConfiguration",
"lr_get_account":				"UserManagement_View",
"lr_get_profile":				"UserManagement_View",
"lr_get_smtp-configuration":	"API_ViewThirdPartyCredentials",
"lr_get_access-restriction":	"SecurityPolicy_View",
"lr_add_domain":				"API_AdminConfiguration",
"lr_add_social": 				"API_EditThirdPartyCredentials",
"lr_add_sott":					"API_AdminConfiguration",
"lr_add_custom-field":			"API_AdminConfiguration",
"lr_add_account":				"UserManagement_Admin",        
"lr_add_hooks":					"ThirdPartyIntegration_Admin",
"lr_add_smtp-configuration":	"API_EditThirdPartyCredentials",
"lr_add_access-restriction":	"SecurityPolicy_Admin",
"lr_set_domain":				"API_EditConfiguration",
"lr_set_social": 				"API_EditThirdPartyCredentials",
"lr_set_theme":					"API_EditConfiguration",
"lr_set_schema":				"API_EditConfiguration",
"lr_set_smtp-configuration":	"API_EditThirdPartyCredentials",
"lr_set_access-restriction":	"SecurityPolicy_Edit",
"lr_delete_domain":				"API_AdminConfiguration",
"lr_delete_social": 			"API_EditThirdPartyCredentials",
"lr_delete_sott":				"API_AdminConfiguration",
"lr_delete_custom-field":		"API_AdminConfiguration",
"lr_delete_account":			"UserManagement_Admin",        
"lr_delete_hooks":				"ThirdPartyIntegration_Admin",
"lr_delete_smtp-configuration":	"API_EditThirdPartyCredentials",
"lr_delete_access-restriction":	"SecurityPolicy_Admin",
"lr_verify_resend":				"UserManagement_Admin",
}

func updatePath(themechildkey string) string {
	var hubDomain = conf.HubPageDomain
    var CdnIDXPath string;
	var CdnPath string;
	if strings.Contains(hubDomain ,"//devhub."){
		CdnIDXPath = "hosted-pages-dev.lrinternal.com"
	CdnPath = "https://cdn-dev.lrinternal.com"
	} else if strings.Contains(hubDomain ,"//staginghub."){
	CdnIDXPath = "hosted-pages-stag.lrinternal.com"
	CdnPath = "https://cdn-stag.lrinternal.com"
	} else {
    CdnIDXPath = "hosted-pages.lrcontent.com";
	CdnPath = "https://cdn.loginradius.com";
	}
    if 
      !strings.Contains(themechildkey,"class=") &&
      !strings.Contains(themechildkey,"https://") {
      if !strings.Contains(themechildkey,"/Themes/") {
        var middlePath = "/";
        if strings.Contains(CdnPath,".lrinternal.com") {
          middlePath += "hosted-page";
        } else {
          middlePath = "/hub/prod/v1";
        }
        themechildkey = CdnPath + middlePath + themechildkey;
      } else {
        themechildkey = "https://" + CdnIDXPath + themechildkey;
      }
    }

    return themechildkey;

}

var Theme1Profile = ThemeType{

	PageType:     "Profile",
	CustomCss:    []string{0: updatePath("/css/hosted-auth-default.css")},
	FavIcon:      updatePath(("/images/favicon.ico")),
	HtmlBody:     "<div class='grid lr-hostr-container lr-hostr-logged-in'>\n<div id='lr-raas-message' class='loginradius-raas-success-message'></div>\n<div class='grid lr-hostr-frame cf'>\n<div class='lr-profile-frame lr-social-login-frame lr-frames lr-sample-background-enabled cf'>\n<div class='lr-profile-image'><img alt=''></div>\n<h1 class='lr-profile-name'></h1>\n<div class='lr-profile-info'><p></p></div>\n<div class='lr-link-social-container'>\n<div class='lr-linked-social-frame' id='lr-linked-social'>\n<h5 class='lr-heading'>Linked social accounts</h5>\n</div>\n<script type=\"text/html\" id=\"linkedAccountsTemplate\"><# if(isLinked) { #> <div class=\"lr-social-account\"><span class=\"lr-social-icon lr-flat-<#= Name.toLowerCase() #> button-shade lr-sl-icon lr-sl-icon-<#= Name.toLowerCase() #>\"></span><span class=\"lr-social-info\"><#= Name #></span><span class=\"lr-social-unlink\"><a onclick='return window[\"loginradiusv1\"]? unLinkAccount(\"<#= Name.toLowerCase() #>\",\"<#= providerId #>\") : LRObject.util.unLinkAccount(\"<#= Name.toLowerCase() #>\",\"<#= providerId #>\")'>Unlink</a></span></div><# } #></script><div class='lr-not-linked-social-frame' id='lr-not-linked-social'>\n<h5 class='lr-heading'>Link more social accounts</h5>\n</div>\n<script type=\"text/html\" id=\"notLinkedAccountsTemplate\"><# if(!isLinked) { #> <span class=\"lr-social-icon lr-flat-<#= Name.toLowerCase() #> button-shade lr-sl-icon lr-sl-icon-<#= Name.toLowerCase() #>\" onclick='LRObject.util.openWindow(\"<#= Endpoint #>\");'></span><# } #></script></div>\n<div class='lr-menu lr-account-menu'>\n<div class='lr-menu-button'></div>\n<div class='lr-menu-list-frame'>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-edit-profile'>Edit Profile</a>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-change-password'>Change Password</a>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-account-settings'>Account Settings</a>\n<a class='lr-logout lr-menu-list' href='auth.aspx?action=logout'>Logout</a>\n</div>\n</div>\n</div>\n<div class='lr-frames lr-more-info-container'>\n<div class='lr-more-info-frame'>\n<div class='lr-more-info-heading'>\n<h2>My Profile</h2>\n<div class='lr-edit-profile lr-button outline-grey lr-show-settings' data-query='lr-edit-profile'>Edit</div>\n</div>\n<div class='lr-content-section cf' id='profile-viewer'>\n</div>\n</div>\n<script type=\"text/html\" id=\"profileViewTemplate\"><# if(typeof value !=\"undefined\" ) { #> <div class=\"lr-content-group\"><h6 class=\"lr-label\"><#= display #></h6> <div class=\"lr-data\"><#= value #></div></div><# } #></script><div class='lr-more-menu-contents'>\n<div id='lr-edit-profile' class='lr-more-menu-frame lr-edit-profile lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Edit My Profile</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-editable-fields-frame' id='profile-editor-container'>\n</div>\n</div>\n<div id='lr-account-settings' class='lr-more-menu-frame lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Account Settings</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-account-settings-frame'>\n<h5 class='lr-setting-label'>Delete Account</h5>\n<div class='lr-action-box'>\n<a class='lr-button white button-shade' href='javascript:void(0)' onclick='deleteAccount();'>Delete my account</a>\n</div>\n</div>\n</div>\n<div id='lr-change-password' class='lr-more-menu-frame lr-field-editor lr-change-password lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Change Password</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-editable-fields-frame' id='change-password'>\n</div>\n<div class='lr-editable-fields-frame' id='set-password'>\n</div>\n</div>\n</div>\n</div>\n</div>\n</div>\n<div class='lr_fade lr-loading-screen-overlay' id='loading-spinner'>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n</div>",
	EndScript:    "",
	BeforeScript: updatePath("/js/default-profile-before-script.js"),
	IsActive:     true,
	MainScript:   "",
	Status:       "0",
}

var Theme2Profile = ThemeType{

	PageType:     "Profile",
	CustomCss:    []string{0: updatePath("/Themes/Theme-Default/profile/css/profile-style.css")},
	FavIcon:      updatePath("/images/favicon.ico"),
	HtmlBody:     "<div class=\"row content\" id=\"lr-showifjsenabled\">\n    <div class=\"conform-ovelay\" style=\"display: none\">\n        <div id=\"confirm\">\n            <h2 class=\"confirmationtitle\" style=\"color: rgb(2, 11, 19);\"></h2>\n            <div class=\"message\"></div>\n            <button class=\"yes_btn\"></button>\n            <button class=\"no_btn\"></button>\n        </div>\n    </div>\n    <div class=\"col\">\n        <div id=\"lr-raas-message\" class=\"loginradius-raas-success-message\"></div>\n        <div class=\"top-container\">\n            <div class=\"top-left-container\">\n                <img class=\"lr-profile-image\"\n                    src=\"data:image/svg+xml,%3Csvg version='1.1' id='Layer_1' xmlns='http://www.w3.org/2000/svg' xmlns:xlink='http://www.w3.org/1999/xlink' x='0px' y='0px' viewBox='0 0 250 250' style='enable-background:new 0 0 250 250;' xml:space='preserve'%3E%3Cstyle type='text/css'%3E .st0%7Bfill:%23757A7E;%7D .st1%7Bfill:%23FFFFFF;%7D%0A%3C/style%3E%3Cg%3E%3Crect class='st0' width='250' height='250'/%3E%3C/g%3E%3Cpath class='st1' d='M148.1,167.3c0,2.2-1.7,3.9-3.9,3.9h-54c-8.6,0-15.6-7-15.6-15.6v-44.8c0-3.3,1-6.4,2.9-9.1 c1.3-1.7,3.7-2.1,5.4-0.9c1.7,1.3,2.1,3.7,0.9,5.4c-1,1.3-1.5,2.9-1.5,4.5v44.8c0,4.3,3.5,7.8,7.8,7.8h54 C146.4,163.4,148.1,165.2,148.1,167.3z M105.2,123.9c-2.1-0.4-4.1,1.1-4.5,3.2c-0.2,1.4-0.4,2.8-0.4,4.2c0,13.4,10.9,24.4,24.4,24.4 c1.6,0,3.2-0.2,4.8-0.5c2.1-0.4,3.5-2.5,3.1-4.6c-0.4-2.1-2.5-3.5-4.6-3.1c-1.1,0.2-2.1,0.3-3.2,0.3c-9.1,0-16.6-7.4-16.6-16.6 c0-1,0.1-1.9,0.2-2.8C108.8,126.3,107.3,124.3,105.2,123.9z M173.5,174.8c-0.8,0.8-1.8,1.1-2.8,1.1s-2-0.4-2.8-1.1l-92-92 c-1.5-1.5-1.5-4,0-5.5c1.5-1.5,4-1.5,5.5,0l17.1,17.1c0.5-0.4,0.9-1,1.2-1.7l1.3-3.7c1.7-4.6,6.1-7.7,11-7.7H137 c5,0,9.5,3.2,11.1,7.9l1.2,3.5c0.5,1.6,2,2.6,3.7,2.6h5.9c8.6,0,15.6,7,15.6,15.6v25.3c0,2.2-1.7,3.9-3.9,3.9 c-2.2,0-3.9-1.7-3.9-3.9v-25.3c0-4.3-3.5-7.8-7.8-7.8h-5.9c-5,0-9.5-3.2-11.1-7.9l-1.2-3.5c-0.5-1.6-2-2.6-3.7-2.6h-24.9 c-1.6,0-3.1,1-3.7,2.6l-1.3,3.7c-0.6,1.8-1.7,3.3-3,4.5l9.7,9.7c3.4-1.7,7.1-2.6,11-2.6c13.4,0,24.4,10.9,24.4,24.4 c0,3.8-0.9,7.6-2.6,11l18.4,18.4c1.1-1.4,1.8-3.1,1.8-5c0-2.2,1.7-3.9,3.9-3.9c2.2,0,3.9,1.7,3.9,3.9c0,4-1.5,7.7-4.1,10.5l3.1,3.1 C175,170.8,175,173.2,173.5,174.8z M119.7,115.5l20.8,20.8c0.5-1.6,0.8-3.3,0.8-5c0-9.1-7.4-16.6-16.6-16.6 C123,114.7,121.3,115,119.7,115.5z'/%3E%3C/svg%3E\"\n                    alt=\"Profile-picture\" />\n            </div>\n            <div class=\"top-right-container\">\n                <h2 class=\"lr-profile-name\"></h2>\n                <button type=\"button\" onclick=\"window.location.href = 'auth.aspx?action=logout';\">\n                    LOGOUT\n                </button>\n            </div>\n        </div>\n        <ul class=\"nav\">\n            <li class=\"nav-item\">\n                <a class=\"nav-link active profile\" onclick=\"navigation('profile')\">Profile</a>\n            </li>\n            <li class=\"nav-item\">\n                <a class=\"nav-link edit\" onclick=\"navigation('edit')\">Edit</a>\n            </li>\n            <li class=\"nav-item\">\n                <a class=\"nav-link security\" onclick=\"navigation('security')\">Security</a>\n            </li>\n        </ul>\n        <div class=\"row\">\n            <div class=\"col\">\n                <div class=\"row\" id=\"navupdatepassword\">\n                    <div class=\"col\">\n                        <h2 class=\"text-primary\">Change Password</h2>\n                        <div class=\"row\">\n                            <div class=\"col\" id=\"change-password\"></div>\n                        </div>\n                    </div>\n                </div>\n\n                <div class=\"row\" id=\"navdownload\">\n                    <div class=\"col\">\n                        <h2 class=\"text-primary\">Download My Profile Data</h2>\n                        <p>\n                            This option is provided to download the data of your account\n                            it will be downloaded as an text file.\n                        </p>\n                        <button type=\"button\" class=\"download-btn\">\n                            DOWNLOAD PROFILE\n                        </button>\n                    </div>\n                </div>\n                <div class=\"row\" id=\"navdeleteacc\">\n                    <div class=\"col\">\n                        <h2 class=\"text-primary\">Delete My Account</h2>\n                        <p>\n                            Kindly be aware of deleting account\n                            <strong>It cannot be undo after delete</strong>.\n                        </p>\n                        <button type=\"button\"\n                            onclick=\"confirmationd('Do you want to delete your account?','You are about to delete your account. All your data will be permanently removed. This action cannot be undone.','deleteaccount');\">\n                            PERMENENTLY DELETE MY ACCOUNT\n                        </button>\n                    </div>\n                </div>\n\n                <div class=\" row\" id=\"navinfo\">\n                    <div class=\"col\">\n                        <h2 class=\"text-primary\">Personal information</h2>\n                        <div class=\"row\">\n                            <div class=\"col profile-viewer\" id=\"profile-viewer\"></div>\n                            <div class=\"hidden-div\">Hidden Div</div>\n                            <script type=\"text/html\"\n                                id=\"profileViewTemplate\"><# if(typeof value !=\"undefined\" ) { #> <div class=\"lr-content-group\"><h6 class=\"lr-label\"><#= display #></h6> <div class=\"lr-data\"><#= value #></div></div><# } #></script>\n                        </div>\n                    </div>\n                </div>\n                <div class=\"row\" id=\"navedit\">\n                    <div class=\"col\">\n                        <h2 class=\"text-primary\">Update Profile</h2>\n                        <div class=\"row\">\n                            <div class=\"col\" id=\"profile-editor-container\">\n                                <p id=\"lr-profile-name\"></p>\n                            </div>\n                        </div>\n                    </div>\n                    <div class=\"col\">\n                        <h2 class=\"text-primary\">Add/Remove Email</h2>\n                        <div class=\"row\">\n                            <table class=\"emailList\" id=\"emailList\" cellspacing=\"10\" cellpadding=\"10\"></table>\n                            <button class=\"addEmail\" onclick=\"addEmail()\">Add Email</button>\n                        </div>\n                    </div>\n                </div>\n\n            </div>\n\n            <div class=\"dialog-ovelay1\">\n                <div class=\"col dialog1\">\n                    <h2 class=\"text-primary-dialog\">Add New Email</h2>\n                    <hr class=\"rounded\">\n                    <div class=\"row\">\n                        <div class=\"col\" id=\"addemail-col\">\n                            <div id=\"addemail-container\"></div>\n                        </div>\n                        <hr>\n                        <div class=\"align-right\">\n                            <button class=\"cance-overlay1 no_btn\" onclick=\"cancle()\">Cancel</button>\n                        </div>\n                    </div>\n                </div>\n            </div>\n        </div>\n\n        <div class=\"row\" id=\"navsocial\">\n            <div class=\"col\">\n\n                <h2 class=\"text-primary\">Social Accounts</h2>\n                <div class=\"row\">\n                    <div class=\"col\" id=\"lr-linked-social\">\n                        <script type=\"text/html\" id=\"linkedAccountsTemplate\">\n                <# if(isLinked) { #>\n                <div class=\"lr-social-account\">\n                  <a\n                    onclick='confirmationd(\"Unlink My Account\",\"Are you sure? You want to unlink <#= Name.toLowerCase() #> account?\",\"unlinksocial\",\"<#= Name.toLowerCase() #>\",\"<#= providerId #>\");'\n                    title=\"Unlink <#= Name #>\" alt=\"Unlink the <#=Name#> account\">\n                    <span\n                      class=\"lr-social-icon lr-flat-<#= Name.toLowerCase() #> button-shade lr-sl-icon lr-sl-icon-<#= Name.toLowerCase() #>\"\n                    ></span\n                  ></a>\n                </div>\n                <# } #>\n              </script>\n                        <span class=\"lr-social-icon button-shade\" onclick=\"addSocialAccount()\">\n                            <img alt=\"Add Social Account\"\n                                src=\"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAYAAACqaXHeAAAABmJLR0QA/wD/AP+gvaeTAAAEO0lEQVR4nO3aS4gcRRjA8d/sSyGKekk0PlDBB4ga4yMnNQjGGFldNYIvFPTkUQwxRxW86BoTRPHi4xCiqCeNR1+niFk2BoyK4kqym002BoVEzUbNroeaYbrbnp2ZnumZTOw/FPRUV3/1ffX4quqroaCgoKCgoKCgoKCg4H9IqUt1LsflWFzOO4jvsBPzXdCpI5yNTdgnGJmW9uFlLOmSjrlQwnr8rrbhyXQE63RghOZdwal4C/cn8g/gc0yXdTgHK4VREuUdPIbZfNXMh5JgQLRnx7AKfSnl+3AbxhPfbNUdX9Uy68UN2Yj+Br4bEHxF9NunctIxNxbjsKoBL2WQsVncJySnxwlNtAfHNNbzSQbEp8PGtmmXMyXxpW5VC7JWR+RM6RFfcJ2q0tPSHV6j9AkrRkXe8pa1S6mg3Vwaef4Ccy3ImivLSJPdFvJogKWR56k2yJuMPJ/bBnkx8miA6F6+HXM2KqPt54Q8GmA68tyOHjuvhuwTlmtVndZ+rTvBmYi8a1rWrgMkl8HVLchaE5EzqUeWQcKRtqL4uLCpaZYBfK213WTXSG6FN2WQ8Yr4VrjnYgTrxA80mzU2EgbEjZ/HkznpmCsl4SgbNWRc8Am1jsNrxIf9PLbIce53IiDyJh5I5M+oBkQIy+VK1Rhhha14XI8GRCqUhPP8EY2HxA4Lw75nvH4jLBGOtFNqGz4pePvkSMiNboXFlwlh8UqQ4wC+V53/BQUFnaHTPmAIV+MinImzyvmzwvXYBL4VVoyThivxPHbgmPpL4D/YhVHcqEeXwkHhRucbja/9tdLreSqa5ZRWj4fxLC5OeTePH7Abh/BbOX8I5+MSXJHQ66Y69Q3iXoxgherSOoMv8SE+wF9N2tE0S7FN+kXnFgzjjAbkLMLteAPbcc8CZUfwY0qdyfST0Ei5cSt+TVT6C57GaTnUV8JzQtS4men0omyXNAvyqDC8KpUcFwIieRhe4RnNGR5No+1U5AnxXtgjeO48uUPzPZ9M97VDkWFhyaoI3Sl+J5AH/cI+oZZhu4Qo8mVCZ9QqNyE43sxcJf6Pj+04vRWBDTJi4Z7dECmbjCol04NkC1kPCEGOReXfE7hLZ3Zvd9Z5/3fk+VidssNk2wdsEGL/8KcQxjqYQU4WbkjJ242PhV79KpL/mTBF+7EWFya+W5FFgT78oTqMOh2sPCQ+jI9rLHiyzH+nwFGUmp0Cc0LgAj4VIr2dJE3frGeFEtmmwM1Ci45p7eo7C/tVT5CEBvlE2OoeFTplR/ndLbhecM4Ppciq/O+gp3jbwp49+meq0Tpl3yOf2+E8+ajO+8HI8yl1ym5rUZeuMCCcJhfaCF0gBFz3LlBuj3Bn0ZPcrbVt8Lx0n9BTvCC78a92Qd+206++k0tLr4n7iZ5nrbAVr2f4XifBsK/FkHCweRc/C/uBWcHRvY9H9LDDy51/AXIizctR/fJVAAAAAElFTkSuQmCC\">\n                        </span>\n                    </div>\n                </div>\n            </div>\n        </div>\n        <div id=\"add-social-account \" class=\"dialog-ovelay\">\n            <div class='dialog'>\n                <h2 class=\"text-primary-dialog\">Add Social Account</h2>\n                <hr class=\"rounded\">\n                <div class=\"col\" id=\"lr-not-linked-social\">\n\n                    <script type=\"text/html\"\n                        id=\"notLinkedAccountsTemplate\"><# if(!isLinked) { #> <span class=\"lr-social-icon lr-flat-<#= Name.toLowerCase() #> button-shade lr-sl-icon lr-sl-icon-<#= Name.toLowerCase() #>\" onclick='LRObject.util.openWindow(\"<#= Endpoint #>\");'title=\"Sign up with <#= Name #>\" alt=\"Sign in with <#=Name#>\"></span><# } #></script>\n                </div>\n                <div class=\"align-right\">\n                    <button class=\"cance-overlay no_btn\">Cancel</button>\n                </div>\n            </div>\n        </div>\n    </div>\n</div>\n<div class=\"lr_fade lr-loading-screen-overlay\" id=\"loading-spinner\">\n    <div id=\"loader\" class=\"lr-ls-page-loadwheel\"></div>\n</div>\n",
	EndScript:    "",
	BeforeScript: updatePath("/Themes/Theme-Default/profile/js/profile-before-script.js"),
	IsActive:     true,
	MainScript:   "",
	Status:       "1",
}

var Theme3Profile = ThemeType{

	PageType:     "Profile",
	CustomCss:    []string{0: updatePath("/Themes/profile/css/profile.css")},
	FavIcon:      updatePath("/images/favicon.ico"),
	HtmlBody:     "<div class='grid lr-hostr-container lr-hostr-logged-in'>\n<div id='lr-raas-message' class='loginradius-raas-success-message'></div>\n<div class='grid lr-hostr-frame cf'>\n<div class='lr-profile-frame lr-social-login-frame lr-frames lr-sample-background-enabled cf'>\n<div class='lr-profile-image'><img alt=''></div>\n<h1 class='lr-profile-name'></h1>\n<div class='lr-profile-info'><p></p></div>\n<div class='lr-link-social-container'>\n<div class='lr-linked-social-frame' id='lr-linked-social'>\n<h5 class='lr-heading'>Linked social accounts</h5>\n</div>\n<script type=\"text/html\" id=\"linkedAccountsTemplate\"><# if(isLinked) { #> <div class=\"lr-social-account\"><span class=\"lr-social-icon lr-flat-<#= Name.toLowerCase() #> button-shade lr-sl-icon lr-sl-icon-<#= Name.toLowerCase() #>\"></span><span class=\"lr-social-info\"><#= Name #></span><span class=\"lr-social-unlink\"><a onclick='return window[\"loginradiusv1\"]? unLinkAccount(\"<#= Name.toLowerCase() #>\",\"<#= providerId #>\") : LRObject.util.unLinkAccount(\"<#= Name.toLowerCase() #>\",\"<#= providerId #>\")'>Unlink</a></span></div><# } #></script><div class='lr-not-linked-social-frame' id='lr-not-linked-social'>\n<h5 class='lr-heading'>Link more social accounts</h5>\n</div>\n<script type=\"text/html\" id=\"notLinkedAccountsTemplate\"><# if(!isLinked) { #> <span class=\"lr-social-icon lr-flat-<#= Name.toLowerCase() #> button-shade lr-sl-icon lr-sl-icon-<#= Name.toLowerCase() #>\" onclick='LRObject.util.openWindow(\"<#= Endpoint #>\");'></span><# } #></script></div>\n<div class='lr-menu lr-account-menu'>\n<div class='lr-menu-button'></div>\n<div class='lr-menu-list-frame'>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-edit-profile'>Edit Profile</a>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-change-password'>Change Password</a>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-account-settings'>Account Settings</a>\n<a class='lr-logout lr-menu-list' href='auth.aspx?action=logout'>Logout</a>\n</div>\n</div>\n</div>\n<div class='lr-frames lr-more-info-container'>\n<div class='lr-more-info-frame'>\n<div class='lr-more-info-heading'>\n<h2>My Profile</h2>\n<div class='lr-edit-profile lr-button outline-grey lr-show-settings' data-query='lr-edit-profile'>Edit</div>\n</div>\n<div class='lr-content-section cf' id='profile-viewer'>\n</div>\n</div>\n<script type=\"text/html\" id=\"profileViewTemplate\"><# if(typeof value !=\"undefined\" ) { #> <div class=\"lr-content-group\"><h6 class=\"lr-label\"><#= display #></h6> <div class=\"lr-data\"><#= value #></div></div><# } #></script><div class='lr-more-menu-contents'>\n<div id='lr-edit-profile' class='lr-more-menu-frame lr-edit-profile lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Edit My Profile</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-editable-fields-frame' id='profile-editor-container'>\n</div>\n</div>\n<div id='lr-account-settings' class='lr-more-menu-frame lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Account Settings</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-account-settings-frame'>\n<h5 class='lr-setting-label'>Delete Account</h5>\n<div class='lr-action-box'>\n<a class='lr-button white button-shade' href='javascript:void(0)' onclick='deleteAccount();'>Delete my account</a>\n</div>\n</div>\n</div>\n<div id='lr-change-password' class='lr-more-menu-frame lr-field-editor lr-change-password lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Change Password</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-editable-fields-frame' id='change-password'>\n</div>\n<div class='lr-editable-fields-frame' id='set-password'>\n</div>\n</div>\n</div>\n</div>\n</div>\n</div>\n<div class='lr_fade lr-loading-screen-overlay' id='loading-spinner'>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n</div>",
	EndScript:    "",
	BeforeScript: updatePath("/Themes/profile/js/profile.js"),
	IsActive:     true,
	MainScript:   "",
	Status:       "2",
}

var Theme4Profile = ThemeType{
	PageType:     "Profile",
	CustomCss:    []string{ updatePath("/Themes/profile/css/profile.css")},
	FavIcon:      updatePath("/images/favicon.ico"),
	HtmlBody:     "<div class='grid lr-hostr-container lr-hostr-logged-in'>\n<div id='lr-raas-message' class='loginradius-raas-success-message'></div>\n<div class='grid lr-hostr-frame cf'>\n<div class='lr-profile-frame lr-social-login-frame lr-frames lr-sample-background-enabled cf'>\n<div class='lr-profile-image'><img alt=''></div>\n<h1 class='lr-profile-name'></h1>\n<div class='lr-profile-info'><p></p></div>\n<div class='lr-link-social-container'>\n<div class='lr-linked-social-frame' id='lr-linked-social'>\n<h5 class='lr-heading'>Linked social accounts</h5>\n</div>\n<script type=\"text/html\" id=\"linkedAccountsTemplate\"><# if(isLinked) { #> <div class=\"lr-social-account\"><span class=\"lr-social-icon lr-flat-<#= Name.toLowerCase() #> button-shade lr-sl-icon lr-sl-icon-<#= Name.toLowerCase() #>\"></span><span class=\"lr-social-info\"><#= Name #></span><span class=\"lr-social-unlink\"><a onclick='return window[\"loginradiusv1\"]? unLinkAccount(\"<#= Name.toLowerCase() #>\",\"<#= providerId #>\") : LRObject.util.unLinkAccount(\"<#= Name.toLowerCase() #>\",\"<#= providerId #>\")'>Unlink</a></span></div><# } #></script><div class='lr-not-linked-social-frame' id='lr-not-linked-social'>\n<h5 class='lr-heading'>Link more social accounts</h5>\n</div>\n<script type=\"text/html\" id=\"notLinkedAccountsTemplate\"><# if(!isLinked) { #> <span class=\"lr-social-icon lr-flat-<#= Name.toLowerCase() #> button-shade lr-sl-icon lr-sl-icon-<#= Name.toLowerCase() #>\" onclick='LRObject.util.openWindow(\"<#= Endpoint #>\");'></span><# } #></script></div>\n<div class='lr-menu lr-account-menu'>\n<div class='lr-menu-button'></div>\n<div class='lr-menu-list-frame'>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-edit-profile'>Edit Profile</a>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-change-password'>Change Password</a>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-account-settings'>Account Settings</a>\n<a class='lr-logout lr-menu-list' href='auth.aspx?action=logout'>Logout</a>\n</div>\n</div>\n</div>\n<div class='lr-frames lr-more-info-container'>\n<div class='lr-more-info-frame'>\n<div class='lr-more-info-heading'>\n<h2>My Profile</h2>\n<div class='lr-edit-profile lr-button outline-grey lr-show-settings' data-query='lr-edit-profile'>Edit</div>\n</div>\n<div class='lr-content-section cf' id='profile-viewer'>\n</div>\n</div>\n<script type=\"text/html\" id=\"profileViewTemplate\"><# if(typeof value !=\"undefined\" ) { #> <div class=\"lr-content-group\"><h6 class=\"lr-label\"><#= display #></h6> <div class=\"lr-data\"><#= value #></div></div><# } #></script><div class='lr-more-menu-contents'>\n<div id='lr-edit-profile' class='lr-more-menu-frame lr-edit-profile lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Edit My Profile</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-editable-fields-frame' id='profile-editor-container'>\n</div>\n</div>\n<div id='lr-account-settings' class='lr-more-menu-frame lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Account Settings</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-account-settings-frame'>\n<h5 class='lr-setting-label'>Delete Account</h5>\n<div class='lr-action-box'>\n<a class='lr-button white button-shade' href='javascript:void(0)' onclick='deleteAccount();'>Delete my account</a>\n</div>\n</div>\n</div>\n<div id='lr-change-password' class='lr-more-menu-frame lr-field-editor lr-change-password lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Change Password</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-editable-fields-frame' id='change-password'>\n</div>\n<div class='lr-editable-fields-frame' id='set-password'>\n</div>\n</div>\n</div>\n</div>\n</div>\n</div>\n<div class='lr_fade lr-loading-screen-overlay' id='loading-spinner'>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n</div>",
	EndScript:    "",
	BeforeScript: updatePath("/Themes/profile/js/profile.js"),
	IsActive:     true,
	MainScript:   "",
	Status:       "3",
}

var Theme5Profile = ThemeType{
	PageType:     "Profile",
	CustomCss:    []string{updatePath("/Themes/profile/css/profile.css")},
	FavIcon:      updatePath("/images/favicon.ico"),
	HtmlBody:     "<div class='grid lr-hostr-container lr-hostr-logged-in'>\n<div id='lr-raas-message' class='loginradius-raas-success-message'></div>\n<div class='grid lr-hostr-frame cf'>\n<div class='lr-profile-frame lr-social-login-frame lr-frames lr-sample-background-enabled cf'>\n<div class='lr-profile-image'><img alt=''></div>\n<h1 class='lr-profile-name'></h1>\n<div class='lr-profile-info'><p></p></div>\n<div class='lr-link-social-container'>\n<div class='lr-linked-social-frame' id='lr-linked-social'>\n<h5 class='lr-heading'>Linked social accounts</h5>\n</div>\n<script type=\"text/html\" id=\"linkedAccountsTemplate\"><# if(isLinked) { #> <div class=\"lr-social-account\"><span class=\"lr-social-icon lr-flat-<#= Name.toLowerCase() #> button-shade lr-sl-icon lr-sl-icon-<#= Name.toLowerCase() #>\"></span><span class=\"lr-social-info\"><#= Name #></span><span class=\"lr-social-unlink\"><a onclick='return window[\"loginradiusv1\"]? unLinkAccount(\"<#= Name.toLowerCase() #>\",\"<#= providerId #>\") : LRObject.util.unLinkAccount(\"<#= Name.toLowerCase() #>\",\"<#= providerId #>\")'>Unlink</a></span></div><# } #></script><div class='lr-not-linked-social-frame' id='lr-not-linked-social'>\n<h5 class='lr-heading'>Link more social accounts</h5>\n</div>\n<script type=\"text/html\" id=\"notLinkedAccountsTemplate\"><# if(!isLinked) { #> <span class=\"lr-social-icon lr-flat-<#= Name.toLowerCase() #> button-shade lr-sl-icon lr-sl-icon-<#= Name.toLowerCase() #>\" onclick='LRObject.util.openWindow(\"<#= Endpoint #>\");'></span><# } #></script></div>\n<div class='lr-menu lr-account-menu'>\n<div class='lr-menu-button'></div>\n<div class='lr-menu-list-frame'>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-edit-profile'>Edit Profile</a>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-change-password'>Change Password</a>\n<a class='lr-settings lr-menu-list lr-show-settings' data-query='lr-account-settings'>Account Settings</a>\n<a class='lr-logout lr-menu-list' href='auth.aspx?action=logout'>Logout</a>\n</div>\n</div>\n</div>\n<div class='lr-frames lr-more-info-container'>\n<div class='lr-more-info-frame'>\n<div class='lr-more-info-heading'>\n<h2>My Profile</h2>\n<div class='lr-edit-profile lr-button outline-grey lr-show-settings' data-query='lr-edit-profile'>Edit</div>\n</div>\n<div class='lr-content-section cf' id='profile-viewer'>\n</div>\n</div>\n<script type=\"text/html\" id=\"profileViewTemplate\"><# if(typeof value !=\"undefined\" ) { #> <div class=\"lr-content-group\"><h6 class=\"lr-label\"><#= display #></h6> <div class=\"lr-data\"><#= value #></div></div><# } #></script><div class='lr-more-menu-contents'>\n<div id='lr-edit-profile' class='lr-more-menu-frame lr-edit-profile lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Edit My Profile</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-editable-fields-frame' id='profile-editor-container'>\n</div>\n</div>\n<div id='lr-account-settings' class='lr-more-menu-frame lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Account Settings</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-account-settings-frame'>\n<h5 class='lr-setting-label'>Delete Account</h5>\n<div class='lr-action-box'>\n<a class='lr-button white button-shade' href='javascript:void(0)' onclick='deleteAccount();'>Delete my account</a>\n</div>\n</div>\n</div>\n<div id='lr-change-password' class='lr-more-menu-frame lr-field-editor lr-change-password lr-account-settings'>\n<div class='lr-more-info-heading'>\n<h2>Change Password</h2>\n<a id='lr-close' class='lr-close'>&times;</a>\n</div>\n<div class='lr-editable-fields-frame' id='change-password'>\n</div>\n<div class='lr-editable-fields-frame' id='set-password'>\n</div>\n</div>\n</div>\n</div>\n</div>\n</div>\n<div class='lr_fade lr-loading-screen-overlay' id='loading-spinner'>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n</div>",
	EndScript:    "",
	BeforeScript: updatePath("/Themes/profile/js/profile.js"),
	IsActive:     true,
	MainScript:   "",
	Status:       "4",
}



var Theme1Auth = ThemeType{
	PageType: "Auth",
	CustomCss: []string{
		updatePath("/css/hosted-auth-default.css"),
	},
	FavIcon:      updatePath("/images/favicon.ico"),
	HtmlBody:     "<div class='grid lr-hostr-container'>\n<div id='lr-raas-message' class='loginradius-raas-success-message'></div>\n<div class='grid lr-hostr-frame cf'>\n<div id='lr-social-login' class='lr-social-login-frame lr-frames lr-sample-background-enabled cf'>\n<h2 class='lr-social-login-message'>Login with your social account</h2>\n<div id='interfacecontainerdiv' class='lr-sl-shaded-brick-frame cf lr-widget-container'>\n</div>\n<script type=\"text/html\" id=\"loginradiuscustom_tmpl\"><span class=\"lr-provider-label lr-sl-shaded-brick-button lr-flat-<#=Name.toLowerCase()#>\" onclick=\" return LRObject.util.openWindow('<#= Endpoint #>'); \" title=\"Sign up with <#= Name #>\" alt=\"Sign in with <#= Name#>\"><span class=\"lr-sl-icon lr-sl-icon-<#= Name.toLowerCase()#>\"></span>Login with <#= Name#> </span> </script></div>\n<div class='lr-frames lr-forms-container'>\n<div id='lr-traditional-login' class='lr-form-frame'>\n<h2 class='lr-form-login-message'>...or login with your email</h2>\n<div id='login-container' class='lr-widget-container'>\n</div>\n<div class='lr-link-box'>\n<a class='lr-raas-forgot-password'>Forgot password?</a>\n<a class='lr-register-link'>Create Account</a>\n</div>\n</div>\n<div id='lr-raas-registartion' class='lr-form-frame'>\n<h2 class='lr-form-login-message'>...or create an account</h2>\n<div id='registration-container' class='lr-widget-container'>\n</div>\n<div class='lr-link-box'>\n<a class='lr-raas-forgot-password'>Forgot password?</a>\n<a class='lr-raas-login-link'>Login</a>\n</div>\n</div>\n<div id='lr-raas-forgotpassword' class='lr-form-frame'>\n<h2 class='lr-form-login-message'>Forgot Password</h2>\n<p class='lr-form-login-subnote'>We'll email you an instruction on how to reset your password.</p>\n<div id='forgotpassword-container' class='lr-widget-container'>\n</div>\n<div class='lr-link-box'>\n<a class='lr-raas-login-link'>Login</a>\n<a class='lr-register-link'>Create Account</a>\n</div>\n</div>\n<div id='lr-raas-sociallogin' class='lr-form-frame'>\n<h2 class='lr-form-login-message'>Complete your Profile</h2>\n<p class='lr-form-login-subnote'>Require to fill all mandatory fields.</p>\n<div id='sociallogin-container' class='lr-widget-container'>\n</div>\n</div>\n<div id='lr-raas-resetpassword' class='lr-form-frame'>\n<h2 class='lr-form-login-message'>Reset your Password</h2>\n<p class='lr-form-login-subnote'>Reset your password to get back access of your account</p>\n<div id='resetpassword-container' class='lr-widget-container'>\n</div>\n</div>\n</div>\n</div>\n</div>\n<div class='lr_fade lr-loading-screen-overlay' id='loading-spinner'>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n<div class='load-dot'></div>\n</div>",
	EndScript:    "",
	BeforeScript: updatePath("/js/default-auth-before-script.js"),
	IsActive:     true,
	MainScript:   "",
	Status:       "0",
}

var Theme2Auth = ThemeType{
	PageType: "Auth",
	CustomCss: []string{
		updatePath("/Themes/Theme-Default/auth/css/idx-style.css"),
	},
	CustomJS: []string{
		updatePath("/Themes/Theme-Default/auth/js/idx-selfhosted.1.0.0.js"),
	},
	FavIcon:      updatePath("/images/favicon.ico"),
	HtmlBody:     "<div id=\"idx-container\"></div>\n",
	EndScript:    "",
	BeforeScript: updatePath("/Themes/Theme-Default/auth/js/idx-options.js"),
	IsActive:     true,
	MainScript:   "",
	Status:       "1",
}

var Theme3Auth = ThemeType{
	PageType: "Auth",
	CustomCss: []string{
		updatePath("/Themes/Theme-1/auth/css/hosted-auth-default.css"),
		updatePath("/Themes/Theme-1/auth/css/jquery-ui.css"),
	},
	FavIcon:      updatePath("/images/favicon.ico"),
	HtmlBody:     "<div class=\"lr-hostr-main-container\">\n\n<div id=\"lr-showifjsenabled\" style=\"visibility: visible;\">\n<div class=\"grid lr-hostr-container\">\n<div id=\"lr-raas-message\" class=\"loginradius-raas-success-message\"></div>\n<div class=\"lr-logo-wrap\">\n<a href=\"#\">\n<div class=\"lr-logo\">\n<img id=\"logo-image\" class=\"lr-logo-size\" src=\"data:image/gif;base64,R0lGODlhAQABAAAAACwAAAAAAQABAAA=\"\n              alt=\"\" />\n\n\n\n<svg version=\"1.1\" id=\"lr-logo-svg\" class=\"lr-logo-size\"  xmlns=\"http://www.w3.org/2000/svg\"\n              xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"200\"\n              height=\"41\" x=\"0px\" y=\"0px\" viewBox=\"0 0 736 148.8\"\n              style=\"enable-background:new 0 0 736 148.8;\" xml:space=\"preserve\">\n<style type=\"text/css\">\n                .st0 {\n                  fill: #FFFFFF;\n                }\n\n                .st1 {\n                  fill: #E5E5E5;\n                }\n\n                .st2 {\n                  fill: #F9F9F9;\n                }\n\n                .st3 {\n                  fill: #D1D1D1;\n                }\n              </style>\n<g>\n<path class=\"st0\" d=\"M228,115.3c-0.3-0.3-0.4-0.6-0.4-1V83.5c0-0.3,0-0.6-0.1-0.7l-25-50.3c-0.2-0.3-0.2-0.6-0.2-0.7\n\t\tc0-0.6,0.4-1,1.3-1h15.3c0.8,0,1.4,0.4,1.7,1.1l15.2,32.3c0.2,0.5,0.5,0.5,0.7,0l15.2-32.3c0.3-0.7,0.9-1.1,1.7-1.1h15.5\n\t\tc0.6,0,1,0.1,1.2,0.4c0.2,0.3,0.2,0.7-0.1,1.3l-25.2,50.3c-0.1,0.2-0.1,0.4-0.1,0.7v30.7c0,0.4-0.1,0.7-0.4,1\n\t\tc-0.3,0.3-0.6,0.4-1,0.4h-14.1C228.7,115.7,228.3,115.6,228,115.3z\" />\n<path class=\"st0\" d=\"M280.4,111.5c-4.9-3.6-8.2-8.4-10-14.6c-1.1-3.8-1.7-7.9-1.7-12.4c0-4.8,0.6-9.1,1.7-12.9\n\t\tc1.9-6,5.2-10.7,10.1-14.1c4.9-3.4,10.7-5.1,17.5-5.1c6.6,0,12.3,1.7,17,5c4.7,3.4,8,8,10,14c1.3,4,1.9,8.3,1.9,12.7\n\t\tc0,4.4-0.6,8.5-1.7,12.3c-1.8,6.3-5.1,11.3-9.9,14.9c-4.8,3.6-10.6,5.4-17.4,5.4C291.1,116.8,285.3,115,280.4,111.5z M304.7,99.7\n\t\tc1.9-1.6,3.2-3.8,4-6.7c0.6-2.6,1-5.4,1-8.5c0-3.4-0.3-6.3-1-8.6c-0.9-2.8-2.3-4.9-4.1-6.4c-1.9-1.5-4.1-2.3-6.8-2.3\n\t\tc-2.8,0-5,0.8-6.9,2.3c-1.8,1.5-3.1,3.7-3.9,6.4c-0.6,1.9-1,4.8-1,8.6c0,3.6,0.3,6.5,0.8,8.5c0.8,2.8,2.2,5.1,4.1,6.7\n\t\tc1.9,1.6,4.2,2.4,7,2.4C300.6,102.1,302.8,101.3,304.7,99.7z\" />\n<path class=\"st0\" d=\"M374.8,53.9c0.3-0.3,0.6-0.4,1-0.4H390c0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v59.5c0,0.4-0.1,0.7-0.4,1\n\t\tc-0.3,0.3-0.6,0.4-1,0.4h-14.2c-0.4,0-0.7-0.1-1-0.4c-0.3-0.3-0.4-0.6-0.4-1v-4.1c0-0.2-0.1-0.4-0.2-0.4c-0.2,0-0.3,0.1-0.5,0.3\n\t\tc-3.2,4.4-8.3,6.6-15.1,6.6c-6.2,0-11.2-1.9-15.2-5.6c-4-3.7-5.9-8.9-5.9-15.7V54.9c0-0.4,0.1-0.7,0.4-1c0.3-0.3,0.6-0.4,1-0.4H353\n\t\tc0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v36.3c0,3.2,0.9,5.9,2.6,7.9c1.7,2,4.1,3,7.2,3c2.7,0,5-0.8,6.8-2.5\n\t\tc1.8-1.7,2.9-3.8,3.3-6.5V54.9C374.4,54.5,374.5,54.1,374.8,53.9z\" />\n<path class=\"st0\" d=\"M442.1,54.3c0.6,0.3,0.9,0.9,0.7,1.8l-2.5,13.8c-0.1,1-0.6,1.3-1.7,0.8c-1.2-0.4-2.6-0.6-4.2-0.6\n\t\tc-0.6,0-1.5,0.1-2.7,0.2c-2.9,0.2-5.4,1.3-7.4,3.2c-2,1.9-3,4.4-3,7.6v33.1c0,0.4-0.1,0.7-0.4,1c-0.3,0.3-0.6,0.4-1,0.4h-14.2\n\t\tc-0.4,0-0.7-0.1-1-0.4c-0.3-0.3-0.4-0.6-0.4-1V54.9c0-0.4,0.1-0.7,0.4-1c0.3-0.3,0.6-0.4,1-0.4h14.2c0.4,0,0.7,0.1,1,0.4\n\t\tc0.3,0.3,0.4,0.6,0.4,1v4.6c0,0.2,0.1,0.4,0.2,0.5c0.2,0.1,0.3,0,0.4-0.1c3.3-4.9,7.8-7.3,13.4-7.3\n\t\tC438.1,52.6,440.4,53.1,442.1,54.3z\" />\n<path class=\"st0\" d=\"M476.4,115.4c-0.3-0.3-0.4-0.6-0.4-1V32.3c0-0.4,0.1-0.7,0.4-1c0.3-0.3,0.6-0.4,1-0.4h14.2\n\t\tc0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v68.2c0,0.4,0.2,0.6,0.6,0.6h39.7c0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v11.8\n\t\tc0,0.4-0.1,0.7-0.4,1c-0.3,0.3-0.6,0.4-1,0.4h-56C477,115.8,476.7,115.7,476.4,115.4z\" />\n<path class=\"st0\" d=\"M554.4,111.5c-4.9-3.6-8.2-8.4-10-14.6c-1.1-3.8-1.7-7.9-1.7-12.4c0-4.8,0.6-9.1,1.7-12.9\n\t\tc1.9-6,5.2-10.7,10.1-14.1c4.9-3.4,10.7-5.1,17.5-5.1c6.6,0,12.3,1.7,17,5c4.7,3.4,8,8,10,14c1.3,4,1.9,8.3,1.9,12.7\n\t\tc0,4.4-0.6,8.5-1.7,12.3c-1.8,6.3-5.1,11.3-9.9,14.9c-4.8,3.6-10.6,5.4-17.4,5.4C565.1,116.8,559.2,115,554.4,111.5z M578.7,99.7\n\t\tc1.9-1.6,3.2-3.8,4-6.7c0.6-2.6,1-5.4,1-8.5c0-3.4-0.3-6.3-1-8.6c-0.9-2.8-2.3-4.9-4.1-6.4c-1.9-1.5-4.1-2.3-6.8-2.3\n\t\tc-2.8,0-5,0.8-6.9,2.3c-1.8,1.5-3.1,3.7-3.9,6.4c-0.6,1.9-1,4.8-1,8.6c0,3.6,0.3,6.5,0.8,8.5c0.8,2.8,2.2,5.1,4.1,6.7\n\t\tc1.9,1.6,4.2,2.4,7,2.4C574.5,102.1,576.8,101.3,578.7,99.7z\" />\n<path class=\"st0\" d=\"M650.5,53.9c0.3-0.3,0.6-0.4,1-0.4h14.2c0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v55.4\n\t\tc0,10.6-3.1,18.2-9.2,22.7c-6.1,4.5-14,6.8-23.6,6.8c-2.8,0-6-0.2-9.5-0.6c-0.8-0.1-1.2-0.6-1.2-1.6l0.5-12.5\n\t\tc0-1.1,0.6-1.5,1.7-1.3c2.9,0.5,5.6,0.7,8,0.7c5.2,0,9.2-1.1,12-3.4c2.8-2.3,4.2-5.9,4.2-10.9c0-0.2-0.1-0.4-0.2-0.4\n\t\tc-0.2,0-0.3,0-0.5,0.2c-3.1,3.3-7.6,5-13.5,5c-5.3,0-10.1-1.3-14.5-3.8c-4.4-2.5-7.5-6.7-9.4-12.3c-1.2-3.7-1.8-8.5-1.8-14.4\n\t\tc0-6.3,0.7-11.5,2.2-15.4c1.7-4.9,4.6-8.9,8.6-11.9c4-3,8.7-4.5,14.1-4.5c6.2,0,11,1.9,14.3,5.6c0.2,0.2,0.3,0.2,0.5,0.2\n\t\tc0.2,0,0.2-0.2,0.2-0.4v-2.9C650,54.5,650.2,54.1,650.5,53.9z M650,84.1c0-2.5-0.1-4.4-0.2-5.8c-0.2-1.4-0.4-2.7-0.8-3.9\n\t\tc-0.7-2.2-1.9-3.9-3.6-5.3c-1.7-1.3-3.8-2-6.3-2c-2.4,0-4.5,0.7-6.2,2c-1.7,1.3-3,3.1-3.8,5.3c-1.1,2.4-1.7,5.7-1.7,9.8\n\t\tc0,4.5,0.5,7.8,1.5,9.7c0.7,2.2,2,3.9,3.8,5.3c1.8,1.3,4,2,6.5,2c2.6,0,4.7-0.7,6.4-2c1.7-1.3,2.9-3.1,3.5-5.2\n\t\tC649.7,91.7,650,88.4,650,84.1z\" />\n<path class=\"st0\" d=\"M689.5,111.5c-4.9-3.6-8.2-8.4-10-14.6c-1.1-3.8-1.7-7.9-1.7-12.4c0-4.8,0.6-9.1,1.7-12.9\n\t\tc1.9-6,5.2-10.7,10.1-14.1c4.9-3.4,10.7-5.1,17.5-5.1c6.6,0,12.3,1.7,17,5c4.7,3.4,8,8,10,14c1.3,4,1.9,8.3,1.9,12.7\n\t\tc0,4.4-0.6,8.5-1.7,12.3c-1.8,6.3-5.1,11.3-9.9,14.9c-4.8,3.6-10.6,5.4-17.4,5.4C700.2,116.8,694.3,115,689.5,111.5z M713.8,99.7\n\t\tc1.9-1.6,3.2-3.8,4-6.7c0.6-2.6,1-5.4,1-8.5c0-3.4-0.3-6.3-1-8.6c-0.9-2.8-2.3-4.9-4.1-6.4c-1.9-1.5-4.1-2.3-6.8-2.3\n\t\tc-2.8,0-5,0.8-6.9,2.3c-1.8,1.5-3.1,3.7-3.9,6.4c-0.6,1.9-1,4.8-1,8.6c0,3.6,0.3,6.5,0.8,8.5c0.8,2.8,2.2,5.1,4.1,6.7\n\t\tc1.9,1.6,4.2,2.4,7,2.4C709.6,102.1,711.9,101.3,713.8,99.7z\" />\n</g>\n<g>\n<path class=\"st1\" d=\"M137.9,35.6c-3.7-0.6-7.6-0.9-11.5-0.9c-32.2,0-59.6,20.4-70,49.1C34.5,77.3,16.8,61.1,8.3,40.2\n\t\tC20.7,16.3,45.7,0,74.4,0C101.3,0,124.8,14.2,137.9,35.6z\" />\n<path class=\"st2\" d=\"M148.8,74.4c0,39.7-31.1,72.2-70.3,74.3c9.7-12.6,15.4-28.3,15.4-45.4c0-18.6-6.8-35.7-18.1-48.7l-0.3,0.3\n\t\tc0.1-0.1,0.2-0.2,0.3-0.3l0,0c0.1-0.1,0.3-0.3,0.5-0.4c0.1-0.1,0.3-0.3,0.5-0.4c0.7-0.6,1.4-1.2,2.1-1.8c0.1-0.1,0.2-0.2,0.3-0.2\n\t\tc0.6-0.5,1.3-1,1.9-1.5c0.7-0.6,1.5-1.1,2.3-1.7c0.8-0.5,1.5-1.1,2.3-1.6c0.8-0.5,1.5-1,2.3-1.4c0.4-0.3,0.9-0.5,1.3-0.8\n\t\tc0.4-0.2,0.8-0.5,1.2-0.7c0.4-0.2,0.8-0.5,1.2-0.7c1.7-0.9,3.4-1.7,5.1-2.5c0.4-0.2,0.9-0.4,1.3-0.5c0.9-0.4,1.8-0.7,2.6-1\n\t\tc0.4-0.2,0.9-0.3,1.3-0.5c0.4-0.2,0.9-0.3,1.4-0.5c0.4-0.1,0.9-0.3,1.4-0.4c1.7-0.5,3.4-1,5.1-1.3c0.5-0.1,1.1-0.2,1.6-0.3\n\t\tc0.3-0.1,0.7-0.1,1-0.2c0.7-0.1,1.4-0.3,2.2-0.4c0.5-0.1,1-0.2,1.4-0.2c0.5-0.1,1-0.1,1.5-0.2c0.5-0.1,1-0.1,1.5-0.2\n\t\tc0.6-0.1,1.1-0.1,1.7-0.2c0.6,0,1.3-0.1,1.9-0.1c0.6,0,1.3-0.1,1.9-0.1c0.6,0,1.3,0,1.9,0c0.8,0,1.6,0,2.4,0c0.3,0,0.6,0,0.9,0\n\t\tc0.8,0,1.7,0.1,2.5,0.1c1,0.1,1.9,0.2,2.8,0.3c0.5,0,0.9,0.1,1.4,0.2c0.5,0.1,0.9,0.1,1.4,0.2v0C144.8,46.9,148.8,60.2,148.8,74.4z\n\t\t\" />\n<path class=\"st3\" d=\"M93.9,103.3c0,17.1-5.8,32.8-15.4,45.4c-1.3,0.1-2.7,0.1-4.1,0.1C33.3,148.8,0,115.5,0,74.4\n\t\tc0-12.3,3-24,8.3-34.2c11,27.3,37.8,46.5,69,46.5c4.9,0,9.8-0.5,14.4-1.4C93.1,91.1,93.9,97.1,93.9,103.3z\" />\n</g>\n</svg>\n</div>\n</a>\n</div>\n<div class=\"grid lr-hostr-frame cf\">\n<div class=\"img-holder\">\n<div class=\"bg\"></div>\n<div class=\"info-holder\"></div>\n</div>\n<div class=\"lr-frames lr-forms-container\">\n<div class=\"lr-forms-container-inner\">\n<div class=\"lr-forms-items\">\n<div class=\"lr-forms-items-header\">\n<h2 class=\"lr-form-login-message\">\nLogin to your account or register yourself if you don't have an account yet.\n</h2>\n<p>\nIt's quick and easy!\n</p>\n</div>\n<div class=\"lr-tabs\">\n<a id=\"login-label\" class=\"active lr-raas-login-link\">Login</a>\n<a id=\"sign-up-label\" class=\"lr-register-link\">Register</a>\n</div>\n<div id=\"lr-traditional-login\" class=\"lr-form-frame\">\n<div id=\"login-container\" class=\"lr-widget-container\"></div>\n\n</div>\n<div id=\"lr-raas-registartion\" class=\"lr-form-frame\">\n<div id=\"registration-container\" class=\"lr-widget-container\"></div>\n\n</div>\n<div id=\"lr-raas-forgotpassword\" class=\"lr-form-frame\">\n<h2 id=\"forgot-pw-label\" class=\"lr-form-login-message\">Forgot Password</h2>\n<p id=\"forgot-pw-message\" class=\"lr-form-login-subnote\">\nWe'll email or sms you an instruction on how to reset your password.\n</p>\n<div id=\"forgotpassword-container\" class=\"lr-widget-container\"></div>\n<div class=\"lr-link-box\">\n<a class=\"lr-raas-login-link\">Login</a>\n<a class=\"lr-register-link\">Create Account</a>\n</div>\n</div>\n<div id=\"lr-raas-sociallogin\" class=\"lr-form-frame\">\n<h2 class=\"lr-form-login-message\">Complete your Profile</h2>\n<p class=\"lr-form-login-subnote\">\nRequire to fill all mandatory fields.\n</p>\n<div id=\"sociallogin-container\" class=\"lr-widget-container\"></div>\n</div>\n<div id=\"lr-raas-resetpassword\" class=\"lr-form-frame\">\n<h2 class=\"lr-form-login-message\">Reset your Password</h2>\n<p class=\"lr-form-login-subnote\">\nReset your password to get back access of your account\n</p>\n<div id=\"resetpassword-container\" class=\"lr-widget-container\"></div>\n</div>\n<div id=\"lr-social-login\" class=\"lr-social-login-frame lr-frames lr-sample-background-enabled cf\">\n<div id=\"social-block-label\" class=\"lr-social-login-message\">Or login with</div>\n<div id=\"interfacecontainerdiv\" class=\"lr-sl-shaded-brick-frame cf lr-widget-container\"></div>\n<script type=\"text/html\" id=\"loginradiuscustom_tmpl\">\n                  <div class=\"lr-provider-wrapper\">\n                    <span class=\"lr-provider-label lr-sl-shaded-brick-button lr-flat-<#=Name.toLowerCase()#>\" onclick=\" return LRObject.util.openWindow('<#= Endpoint #>'); \" title=\"Login with <#= Name #>\" alt=\"Login with <#= Name#>\">\n                      <span class=\"lr-sl-icon lr-sl-icon-<#= Name.toLowerCase()#>\"></span>\n                    </span>\n                    <span class=\"lr-provider-provider-name\" onclick=\" return LRObject.util.openWindow('<#= Endpoint #>'); \" title=\"Login with <#= Name #>\" alt=\"Login with <#= Name#>\"><#= Name#></span>\n                  </div>\n                </script>\n</div>\n</div>\n</div>\n</div>\n</div>\n</div>\n<div class=\"lr_fade lr-loading-screen-overlay\" id=\"loading-spinner\">\n<div class=\"load-dot\"></div>\n<div class=\"load-dot\"></div>\n<div class=\"load-dot\"></div>\n<div class=\"load-dot\"></div>\n</div>\n</div>\n\n</div>",
	EndScript:    "",
	BeforeScript: updatePath("/Themes/Theme-1/auth/js/before-script.js?v=123"),
	IsActive:     true,
	MainScript:   "",
	Status:       "2",
}

var Theme4Auth = ThemeType{
	PageType: "Auth",
	CustomCss: []string{
		updatePath("/Themes/Theme-2/auth/css/hosted-auth-default.css"),
		updatePath("/Themes/Theme-2/auth/css/jquery-ui.css"),
	},
	FavIcon:      updatePath("/images/favicon.ico"),
	HtmlBody:     "<div class=\"lr-hostr-main-container\">\n\n<div id=\"lr-showifjsenabled\" style=\"visibility: visible;\">\n<div class=\"grid lr-hostr-container\">\n<div id=\"lr-raas-message\" class=\"loginradius-raas-success-message\"></div>\n<div class=\"grid lr-hostr-frame cf\">\n<div class=\"lr-frames lr-forms-container\">\n<div class=\"lr-forms-container-inner\">\n<div class=\"lr-forms-items\">\n<div class=\"lr-logo-wrap\">\n<a href=\"#\">\n<div class=\"lr-logo\">\n<img id=\"logo-image\" class=\"lr-logo-size\"\n                      src=\"data:image/gif;base64,R0lGODlhAQABAAAAACwAAAAAAQABAAA=\" alt=\"\" />\n\n\n\n<svg version=\"1.1\" id=\"lr-logo-svg\" class=\"lr-logo-size\" xmlns=\"http://www.w3.org/2000/svg\"\n                      xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"200\"\n                      height=\"41\" x=\"0px\" y=\"0px\" viewBox=\"0 0 736 148.8\"\n                      style=\"enable-background:new 0 0 736 148.8;\" xml:space=\"preserve\">\n<style type=\"text/css\">\n                        .st0 {\n                          fill: #FFFFFF;\n                        }\n\n                        .st1 {\n                          fill: #E5E5E5;\n                        }\n\n                        .st2 {\n                          fill: #F9F9F9;\n                        }\n\n                        .st3 {\n                          fill: #D1D1D1;\n                        }\n                      </style>\n<g>\n<path class=\"st0\" d=\"M228,115.3c-0.3-0.3-0.4-0.6-0.4-1V83.5c0-0.3,0-0.6-0.1-0.7l-25-50.3c-0.2-0.3-0.2-0.6-0.2-0.7\n\t\tc0-0.6,0.4-1,1.3-1h15.3c0.8,0,1.4,0.4,1.7,1.1l15.2,32.3c0.2,0.5,0.5,0.5,0.7,0l15.2-32.3c0.3-0.7,0.9-1.1,1.7-1.1h15.5\n\t\tc0.6,0,1,0.1,1.2,0.4c0.2,0.3,0.2,0.7-0.1,1.3l-25.2,50.3c-0.1,0.2-0.1,0.4-0.1,0.7v30.7c0,0.4-0.1,0.7-0.4,1\n\t\tc-0.3,0.3-0.6,0.4-1,0.4h-14.1C228.7,115.7,228.3,115.6,228,115.3z\" />\n<path class=\"st0\" d=\"M280.4,111.5c-4.9-3.6-8.2-8.4-10-14.6c-1.1-3.8-1.7-7.9-1.7-12.4c0-4.8,0.6-9.1,1.7-12.9\n\t\tc1.9-6,5.2-10.7,10.1-14.1c4.9-3.4,10.7-5.1,17.5-5.1c6.6,0,12.3,1.7,17,5c4.7,3.4,8,8,10,14c1.3,4,1.9,8.3,1.9,12.7\n\t\tc0,4.4-0.6,8.5-1.7,12.3c-1.8,6.3-5.1,11.3-9.9,14.9c-4.8,3.6-10.6,5.4-17.4,5.4C291.1,116.8,285.3,115,280.4,111.5z M304.7,99.7\n\t\tc1.9-1.6,3.2-3.8,4-6.7c0.6-2.6,1-5.4,1-8.5c0-3.4-0.3-6.3-1-8.6c-0.9-2.8-2.3-4.9-4.1-6.4c-1.9-1.5-4.1-2.3-6.8-2.3\n\t\tc-2.8,0-5,0.8-6.9,2.3c-1.8,1.5-3.1,3.7-3.9,6.4c-0.6,1.9-1,4.8-1,8.6c0,3.6,0.3,6.5,0.8,8.5c0.8,2.8,2.2,5.1,4.1,6.7\n\t\tc1.9,1.6,4.2,2.4,7,2.4C300.6,102.1,302.8,101.3,304.7,99.7z\" />\n<path class=\"st0\" d=\"M374.8,53.9c0.3-0.3,0.6-0.4,1-0.4H390c0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v59.5c0,0.4-0.1,0.7-0.4,1\n\t\tc-0.3,0.3-0.6,0.4-1,0.4h-14.2c-0.4,0-0.7-0.1-1-0.4c-0.3-0.3-0.4-0.6-0.4-1v-4.1c0-0.2-0.1-0.4-0.2-0.4c-0.2,0-0.3,0.1-0.5,0.3\n\t\tc-3.2,4.4-8.3,6.6-15.1,6.6c-6.2,0-11.2-1.9-15.2-5.6c-4-3.7-5.9-8.9-5.9-15.7V54.9c0-0.4,0.1-0.7,0.4-1c0.3-0.3,0.6-0.4,1-0.4H353\n\t\tc0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v36.3c0,3.2,0.9,5.9,2.6,7.9c1.7,2,4.1,3,7.2,3c2.7,0,5-0.8,6.8-2.5\n\t\tc1.8-1.7,2.9-3.8,3.3-6.5V54.9C374.4,54.5,374.5,54.1,374.8,53.9z\" />\n<path class=\"st0\" d=\"M442.1,54.3c0.6,0.3,0.9,0.9,0.7,1.8l-2.5,13.8c-0.1,1-0.6,1.3-1.7,0.8c-1.2-0.4-2.6-0.6-4.2-0.6\n\t\tc-0.6,0-1.5,0.1-2.7,0.2c-2.9,0.2-5.4,1.3-7.4,3.2c-2,1.9-3,4.4-3,7.6v33.1c0,0.4-0.1,0.7-0.4,1c-0.3,0.3-0.6,0.4-1,0.4h-14.2\n\t\tc-0.4,0-0.7-0.1-1-0.4c-0.3-0.3-0.4-0.6-0.4-1V54.9c0-0.4,0.1-0.7,0.4-1c0.3-0.3,0.6-0.4,1-0.4h14.2c0.4,0,0.7,0.1,1,0.4\n\t\tc0.3,0.3,0.4,0.6,0.4,1v4.6c0,0.2,0.1,0.4,0.2,0.5c0.2,0.1,0.3,0,0.4-0.1c3.3-4.9,7.8-7.3,13.4-7.3\n\t\tC438.1,52.6,440.4,53.1,442.1,54.3z\" />\n<path class=\"st0\" d=\"M476.4,115.4c-0.3-0.3-0.4-0.6-0.4-1V32.3c0-0.4,0.1-0.7,0.4-1c0.3-0.3,0.6-0.4,1-0.4h14.2\n\t\tc0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v68.2c0,0.4,0.2,0.6,0.6,0.6h39.7c0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v11.8\n\t\tc0,0.4-0.1,0.7-0.4,1c-0.3,0.3-0.6,0.4-1,0.4h-56C477,115.8,476.7,115.7,476.4,115.4z\" />\n<path class=\"st0\" d=\"M554.4,111.5c-4.9-3.6-8.2-8.4-10-14.6c-1.1-3.8-1.7-7.9-1.7-12.4c0-4.8,0.6-9.1,1.7-12.9\n\t\tc1.9-6,5.2-10.7,10.1-14.1c4.9-3.4,10.7-5.1,17.5-5.1c6.6,0,12.3,1.7,17,5c4.7,3.4,8,8,10,14c1.3,4,1.9,8.3,1.9,12.7\n\t\tc0,4.4-0.6,8.5-1.7,12.3c-1.8,6.3-5.1,11.3-9.9,14.9c-4.8,3.6-10.6,5.4-17.4,5.4C565.1,116.8,559.2,115,554.4,111.5z M578.7,99.7\n\t\tc1.9-1.6,3.2-3.8,4-6.7c0.6-2.6,1-5.4,1-8.5c0-3.4-0.3-6.3-1-8.6c-0.9-2.8-2.3-4.9-4.1-6.4c-1.9-1.5-4.1-2.3-6.8-2.3\n\t\tc-2.8,0-5,0.8-6.9,2.3c-1.8,1.5-3.1,3.7-3.9,6.4c-0.6,1.9-1,4.8-1,8.6c0,3.6,0.3,6.5,0.8,8.5c0.8,2.8,2.2,5.1,4.1,6.7\n\t\tc1.9,1.6,4.2,2.4,7,2.4C574.5,102.1,576.8,101.3,578.7,99.7z\" />\n<path class=\"st0\" d=\"M650.5,53.9c0.3-0.3,0.6-0.4,1-0.4h14.2c0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v55.4\n\t\tc0,10.6-3.1,18.2-9.2,22.7c-6.1,4.5-14,6.8-23.6,6.8c-2.8,0-6-0.2-9.5-0.6c-0.8-0.1-1.2-0.6-1.2-1.6l0.5-12.5\n\t\tc0-1.1,0.6-1.5,1.7-1.3c2.9,0.5,5.6,0.7,8,0.7c5.2,0,9.2-1.1,12-3.4c2.8-2.3,4.2-5.9,4.2-10.9c0-0.2-0.1-0.4-0.2-0.4\n\t\tc-0.2,0-0.3,0-0.5,0.2c-3.1,3.3-7.6,5-13.5,5c-5.3,0-10.1-1.3-14.5-3.8c-4.4-2.5-7.5-6.7-9.4-12.3c-1.2-3.7-1.8-8.5-1.8-14.4\n\t\tc0-6.3,0.7-11.5,2.2-15.4c1.7-4.9,4.6-8.9,8.6-11.9c4-3,8.7-4.5,14.1-4.5c6.2,0,11,1.9,14.3,5.6c0.2,0.2,0.3,0.2,0.5,0.2\n\t\tc0.2,0,0.2-0.2,0.2-0.4v-2.9C650,54.5,650.2,54.1,650.5,53.9z M650,84.1c0-2.5-0.1-4.4-0.2-5.8c-0.2-1.4-0.4-2.7-0.8-3.9\n\t\tc-0.7-2.2-1.9-3.9-3.6-5.3c-1.7-1.3-3.8-2-6.3-2c-2.4,0-4.5,0.7-6.2,2c-1.7,1.3-3,3.1-3.8,5.3c-1.1,2.4-1.7,5.7-1.7,9.8\n\t\tc0,4.5,0.5,7.8,1.5,9.7c0.7,2.2,2,3.9,3.8,5.3c1.8,1.3,4,2,6.5,2c2.6,0,4.7-0.7,6.4-2c1.7-1.3,2.9-3.1,3.5-5.2\n\t\tC649.7,91.7,650,88.4,650,84.1z\" />\n<path class=\"st0\" d=\"M689.5,111.5c-4.9-3.6-8.2-8.4-10-14.6c-1.1-3.8-1.7-7.9-1.7-12.4c0-4.8,0.6-9.1,1.7-12.9\n\t\tc1.9-6,5.2-10.7,10.1-14.1c4.9-3.4,10.7-5.1,17.5-5.1c6.6,0,12.3,1.7,17,5c4.7,3.4,8,8,10,14c1.3,4,1.9,8.3,1.9,12.7\n\t\tc0,4.4-0.6,8.5-1.7,12.3c-1.8,6.3-5.1,11.3-9.9,14.9c-4.8,3.6-10.6,5.4-17.4,5.4C700.2,116.8,694.3,115,689.5,111.5z M713.8,99.7\n\t\tc1.9-1.6,3.2-3.8,4-6.7c0.6-2.6,1-5.4,1-8.5c0-3.4-0.3-6.3-1-8.6c-0.9-2.8-2.3-4.9-4.1-6.4c-1.9-1.5-4.1-2.3-6.8-2.3\n\t\tc-2.8,0-5,0.8-6.9,2.3c-1.8,1.5-3.1,3.7-3.9,6.4c-0.6,1.9-1,4.8-1,8.6c0,3.6,0.3,6.5,0.8,8.5c0.8,2.8,2.2,5.1,4.1,6.7\n\t\tc1.9,1.6,4.2,2.4,7,2.4C709.6,102.1,711.9,101.3,713.8,99.7z\" />\n</g>\n<g>\n<path class=\"st1\" d=\"M137.9,35.6c-3.7-0.6-7.6-0.9-11.5-0.9c-32.2,0-59.6,20.4-70,49.1C34.5,77.3,16.8,61.1,8.3,40.2\n\t\tC20.7,16.3,45.7,0,74.4,0C101.3,0,124.8,14.2,137.9,35.6z\" />\n<path class=\"st2\" d=\"M148.8,74.4c0,39.7-31.1,72.2-70.3,74.3c9.7-12.6,15.4-28.3,15.4-45.4c0-18.6-6.8-35.7-18.1-48.7l-0.3,0.3\n\t\tc0.1-0.1,0.2-0.2,0.3-0.3l0,0c0.1-0.1,0.3-0.3,0.5-0.4c0.1-0.1,0.3-0.3,0.5-0.4c0.7-0.6,1.4-1.2,2.1-1.8c0.1-0.1,0.2-0.2,0.3-0.2\n\t\tc0.6-0.5,1.3-1,1.9-1.5c0.7-0.6,1.5-1.1,2.3-1.7c0.8-0.5,1.5-1.1,2.3-1.6c0.8-0.5,1.5-1,2.3-1.4c0.4-0.3,0.9-0.5,1.3-0.8\n\t\tc0.4-0.2,0.8-0.5,1.2-0.7c0.4-0.2,0.8-0.5,1.2-0.7c1.7-0.9,3.4-1.7,5.1-2.5c0.4-0.2,0.9-0.4,1.3-0.5c0.9-0.4,1.8-0.7,2.6-1\n\t\tc0.4-0.2,0.9-0.3,1.3-0.5c0.4-0.2,0.9-0.3,1.4-0.5c0.4-0.1,0.9-0.3,1.4-0.4c1.7-0.5,3.4-1,5.1-1.3c0.5-0.1,1.1-0.2,1.6-0.3\n\t\tc0.3-0.1,0.7-0.1,1-0.2c0.7-0.1,1.4-0.3,2.2-0.4c0.5-0.1,1-0.2,1.4-0.2c0.5-0.1,1-0.1,1.5-0.2c0.5-0.1,1-0.1,1.5-0.2\n\t\tc0.6-0.1,1.1-0.1,1.7-0.2c0.6,0,1.3-0.1,1.9-0.1c0.6,0,1.3-0.1,1.9-0.1c0.6,0,1.3,0,1.9,0c0.8,0,1.6,0,2.4,0c0.3,0,0.6,0,0.9,0\n\t\tc0.8,0,1.7,0.1,2.5,0.1c1,0.1,1.9,0.2,2.8,0.3c0.5,0,0.9,0.1,1.4,0.2c0.5,0.1,0.9,0.1,1.4,0.2v0C144.8,46.9,148.8,60.2,148.8,74.4z\n\t\t\" />\n<path class=\"st3\" d=\"M93.9,103.3c0,17.1-5.8,32.8-15.4,45.4c-1.3,0.1-2.7,0.1-4.1,0.1C33.3,148.8,0,115.5,0,74.4\n\t\tc0-12.3,3-24,8.3-34.2c11,27.3,37.8,46.5,69,46.5c4.9,0,9.8-0.5,14.4-1.4C93.1,91.1,93.9,97.1,93.9,103.3z\" />\n</g>\n</svg>\n</div>\n</a>\n</div>\n<div class=\"lr-forms-items-header\">\n<h2 class=\"lr-form-login-message\">\nLogin to your account or register yourself if you don't have an account yet.\n</h2>\n<p>\nIt's quick and easy!\n</p>\n</div>\n<div class=\"lr-tabs\">\n<a id=\"login-label\" class=\"active lr-raas-login-link\">Login</a>\n<a id=\"sign-up-label\" class=\"lr-register-link\">Register</a>\n</div>\n<div id=\"lr-traditional-login\" class=\"lr-form-frame\">\n<div id=\"login-container\" class=\"lr-widget-container\"></div>\n\n</div>\n<div id=\"lr-raas-registartion\" class=\"lr-form-frame\">\n<div id=\"registration-container\" class=\"lr-widget-container\"></div>\n\n</div>\n<div id=\"lr-raas-forgotpassword\" class=\"lr-form-frame\">\n<h2 id=\"forgot-pw-label\" class=\"lr-form-login-message\">Forgot Password</h2>\n<p id=\"forgot-pw-message\" class=\"lr-form-login-subnote\">\nWe'll email or sms you an instruction on how to reset your password.\n</p>\n<div id=\"forgotpassword-container\" class=\"lr-widget-container\"></div>\n<div class=\"lr-link-box\">\n<a class=\"lr-raas-login-link\">Login</a>\n<a class=\"lr-register-link\">Create Account</a>\n</div>\n</div>\n<div id=\"lr-raas-sociallogin\" class=\"lr-form-frame\">\n<h2 class=\"lr-form-login-message\">Complete your Profile</h2>\n<p class=\"lr-form-login-subnote\">\nRequire to fill all mandatory fields.\n</p>\n<div id=\"sociallogin-container\" class=\"lr-widget-container\"></div>\n</div>\n<div id=\"lr-raas-resetpassword\" class=\"lr-form-frame\">\n<h2 class=\"lr-form-login-message\">Reset your Password</h2>\n<p class=\"lr-form-login-subnote\">\nReset your password to get back access of your account\n</p>\n<div id=\"resetpassword-container\" class=\"lr-widget-container\"></div>\n</div>\n<div id=\"lr-social-login\" class=\"lr-social-login-frame lr-frames lr-sample-background-enabled cf\">\n<span id=\"social-block-label\" class=\"lr-social-login-message\">Or login with</span>\n<div id=\"interfacecontainerdiv\" class=\"lr-sl-shaded-brick-frame cf lr-widget-container\"></div>\n<script type=\"text/html\" id=\"loginradiuscustom_tmpl\">\n                  <div class=\"lr-provider-wrapper\">\n                    <span class=\"lr-provider-label lr-sl-shaded-brick-button lr-flat-<#=Name.toLowerCase()#>\" onclick=\" return LRObject.util.openWindow('<#= Endpoint #>'); \" title=\"Login with <#= Name #>\" alt=\"Login with <#= Name#>\">\n                      <span class=\"lr-sl-icon lr-sl-icon-<#= Name.toLowerCase()#>\"></span>\n                    </span>\n                    <span class=\"lr-provider-provider-name\" onclick=\" return LRObject.util.openWindow('<#= Endpoint #>'); \" title=\"Login with <#= Name #>\" alt=\"Login with <#= Name#>\"><#= Name#></span>\n                  </div>\n                </script>\n</div>\n</div>\n</div>\n</div>\n</div>\n</div>\n<div class=\"lr_fade lr-loading-screen-overlay\" id=\"loading-spinner\">\n<div class=\"load-dot\"></div>\n<div class=\"load-dot\"></div>\n<div class=\"load-dot\"></div>\n<div class=\"load-dot\"></div>\n</div>\n</div>\n\n</div>",
	EndScript:    "",
	BeforeScript: updatePath("/Themes/Theme-2/auth/js/before-script.js?v=123"),
	IsActive:     true,
	MainScript:   "",
	Status:       "3",
}

var Theme5Auth = ThemeType{
	PageType: "Auth",
	CustomCss: []string{
		updatePath("/Themes/Theme-3/auth/css/hosted-auth-default.css"),
		updatePath("/Themes/Theme-3/auth/css/jquery-ui.css"),
	},
	FavIcon:      updatePath("/images/favicon.ico"),
	HtmlBody:     "<div class=\"lr-hostr-main-container\">\n\n<div id=\"lr-showifjsenabled\" style=\"visibility: visible;\">\n<div class=\"grid lr-hostr-container\">\n<div id=\"lr-raas-message\" class=\"loginradius-raas-success-message\"></div>\n<div class=\"lr-logo-wrap\">\n<a href=\"#\">\n<div class=\"lr-logo\">\n<img id=\"logo-image\" class=\"lr-logo-size\" src=\"data:image/gif;base64,R0lGODlhAQABAAAAACwAAAAAAQABAAA=\"\n              alt=\"\" />\n\n\n\n<svg version=\"1.1\" id=\"lr-logo-svg\" class=\"lr-logo-size\" xmlns=\"http://www.w3.org/2000/svg\"\n              xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"200\"\n              height=\"41\" x=\"0px\" y=\"0px\" viewBox=\"0 0 736 148.8\"\n              style=\"enable-background:new 0 0 736 148.8;\" xml:space=\"preserve\">\n<style type=\"text/css\">\n                .st0 {\n                  fill: #FFFFFF;\n                }\n\n                .st1 {\n                  fill: #E5E5E5;\n                }\n\n                .st2 {\n                  fill: #F9F9F9;\n                }\n\n                .st3 {\n                  fill: #D1D1D1;\n                }\n              </style>\n<g>\n<path class=\"st0\" d=\"M228,115.3c-0.3-0.3-0.4-0.6-0.4-1V83.5c0-0.3,0-0.6-0.1-0.7l-25-50.3c-0.2-0.3-0.2-0.6-0.2-0.7\n\t\tc0-0.6,0.4-1,1.3-1h15.3c0.8,0,1.4,0.4,1.7,1.1l15.2,32.3c0.2,0.5,0.5,0.5,0.7,0l15.2-32.3c0.3-0.7,0.9-1.1,1.7-1.1h15.5\n\t\tc0.6,0,1,0.1,1.2,0.4c0.2,0.3,0.2,0.7-0.1,1.3l-25.2,50.3c-0.1,0.2-0.1,0.4-0.1,0.7v30.7c0,0.4-0.1,0.7-0.4,1\n\t\tc-0.3,0.3-0.6,0.4-1,0.4h-14.1C228.7,115.7,228.3,115.6,228,115.3z\" />\n<path class=\"st0\" d=\"M280.4,111.5c-4.9-3.6-8.2-8.4-10-14.6c-1.1-3.8-1.7-7.9-1.7-12.4c0-4.8,0.6-9.1,1.7-12.9\n\t\tc1.9-6,5.2-10.7,10.1-14.1c4.9-3.4,10.7-5.1,17.5-5.1c6.6,0,12.3,1.7,17,5c4.7,3.4,8,8,10,14c1.3,4,1.9,8.3,1.9,12.7\n\t\tc0,4.4-0.6,8.5-1.7,12.3c-1.8,6.3-5.1,11.3-9.9,14.9c-4.8,3.6-10.6,5.4-17.4,5.4C291.1,116.8,285.3,115,280.4,111.5z M304.7,99.7\n\t\tc1.9-1.6,3.2-3.8,4-6.7c0.6-2.6,1-5.4,1-8.5c0-3.4-0.3-6.3-1-8.6c-0.9-2.8-2.3-4.9-4.1-6.4c-1.9-1.5-4.1-2.3-6.8-2.3\n\t\tc-2.8,0-5,0.8-6.9,2.3c-1.8,1.5-3.1,3.7-3.9,6.4c-0.6,1.9-1,4.8-1,8.6c0,3.6,0.3,6.5,0.8,8.5c0.8,2.8,2.2,5.1,4.1,6.7\n\t\tc1.9,1.6,4.2,2.4,7,2.4C300.6,102.1,302.8,101.3,304.7,99.7z\" />\n<path class=\"st0\" d=\"M374.8,53.9c0.3-0.3,0.6-0.4,1-0.4H390c0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v59.5c0,0.4-0.1,0.7-0.4,1\n\t\tc-0.3,0.3-0.6,0.4-1,0.4h-14.2c-0.4,0-0.7-0.1-1-0.4c-0.3-0.3-0.4-0.6-0.4-1v-4.1c0-0.2-0.1-0.4-0.2-0.4c-0.2,0-0.3,0.1-0.5,0.3\n\t\tc-3.2,4.4-8.3,6.6-15.1,6.6c-6.2,0-11.2-1.9-15.2-5.6c-4-3.7-5.9-8.9-5.9-15.7V54.9c0-0.4,0.1-0.7,0.4-1c0.3-0.3,0.6-0.4,1-0.4H353\n\t\tc0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v36.3c0,3.2,0.9,5.9,2.6,7.9c1.7,2,4.1,3,7.2,3c2.7,0,5-0.8,6.8-2.5\n\t\tc1.8-1.7,2.9-3.8,3.3-6.5V54.9C374.4,54.5,374.5,54.1,374.8,53.9z\" />\n<path class=\"st0\" d=\"M442.1,54.3c0.6,0.3,0.9,0.9,0.7,1.8l-2.5,13.8c-0.1,1-0.6,1.3-1.7,0.8c-1.2-0.4-2.6-0.6-4.2-0.6\n\t\tc-0.6,0-1.5,0.1-2.7,0.2c-2.9,0.2-5.4,1.3-7.4,3.2c-2,1.9-3,4.4-3,7.6v33.1c0,0.4-0.1,0.7-0.4,1c-0.3,0.3-0.6,0.4-1,0.4h-14.2\n\t\tc-0.4,0-0.7-0.1-1-0.4c-0.3-0.3-0.4-0.6-0.4-1V54.9c0-0.4,0.1-0.7,0.4-1c0.3-0.3,0.6-0.4,1-0.4h14.2c0.4,0,0.7,0.1,1,0.4\n\t\tc0.3,0.3,0.4,0.6,0.4,1v4.6c0,0.2,0.1,0.4,0.2,0.5c0.2,0.1,0.3,0,0.4-0.1c3.3-4.9,7.8-7.3,13.4-7.3\n\t\tC438.1,52.6,440.4,53.1,442.1,54.3z\" />\n<path class=\"st0\" d=\"M476.4,115.4c-0.3-0.3-0.4-0.6-0.4-1V32.3c0-0.4,0.1-0.7,0.4-1c0.3-0.3,0.6-0.4,1-0.4h14.2\n\t\tc0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v68.2c0,0.4,0.2,0.6,0.6,0.6h39.7c0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v11.8\n\t\tc0,0.4-0.1,0.7-0.4,1c-0.3,0.3-0.6,0.4-1,0.4h-56C477,115.8,476.7,115.7,476.4,115.4z\" />\n<path class=\"st0\" d=\"M554.4,111.5c-4.9-3.6-8.2-8.4-10-14.6c-1.1-3.8-1.7-7.9-1.7-12.4c0-4.8,0.6-9.1,1.7-12.9\n\t\tc1.9-6,5.2-10.7,10.1-14.1c4.9-3.4,10.7-5.1,17.5-5.1c6.6,0,12.3,1.7,17,5c4.7,3.4,8,8,10,14c1.3,4,1.9,8.3,1.9,12.7\n\t\tc0,4.4-0.6,8.5-1.7,12.3c-1.8,6.3-5.1,11.3-9.9,14.9c-4.8,3.6-10.6,5.4-17.4,5.4C565.1,116.8,559.2,115,554.4,111.5z M578.7,99.7\n\t\tc1.9-1.6,3.2-3.8,4-6.7c0.6-2.6,1-5.4,1-8.5c0-3.4-0.3-6.3-1-8.6c-0.9-2.8-2.3-4.9-4.1-6.4c-1.9-1.5-4.1-2.3-6.8-2.3\n\t\tc-2.8,0-5,0.8-6.9,2.3c-1.8,1.5-3.1,3.7-3.9,6.4c-0.6,1.9-1,4.8-1,8.6c0,3.6,0.3,6.5,0.8,8.5c0.8,2.8,2.2,5.1,4.1,6.7\n\t\tc1.9,1.6,4.2,2.4,7,2.4C574.5,102.1,576.8,101.3,578.7,99.7z\" />\n<path class=\"st0\" d=\"M650.5,53.9c0.3-0.3,0.6-0.4,1-0.4h14.2c0.4,0,0.7,0.1,1,0.4c0.3,0.3,0.4,0.6,0.4,1v55.4\n\t\tc0,10.6-3.1,18.2-9.2,22.7c-6.1,4.5-14,6.8-23.6,6.8c-2.8,0-6-0.2-9.5-0.6c-0.8-0.1-1.2-0.6-1.2-1.6l0.5-12.5\n\t\tc0-1.1,0.6-1.5,1.7-1.3c2.9,0.5,5.6,0.7,8,0.7c5.2,0,9.2-1.1,12-3.4c2.8-2.3,4.2-5.9,4.2-10.9c0-0.2-0.1-0.4-0.2-0.4\n\t\tc-0.2,0-0.3,0-0.5,0.2c-3.1,3.3-7.6,5-13.5,5c-5.3,0-10.1-1.3-14.5-3.8c-4.4-2.5-7.5-6.7-9.4-12.3c-1.2-3.7-1.8-8.5-1.8-14.4\n\t\tc0-6.3,0.7-11.5,2.2-15.4c1.7-4.9,4.6-8.9,8.6-11.9c4-3,8.7-4.5,14.1-4.5c6.2,0,11,1.9,14.3,5.6c0.2,0.2,0.3,0.2,0.5,0.2\n\t\tc0.2,0,0.2-0.2,0.2-0.4v-2.9C650,54.5,650.2,54.1,650.5,53.9z M650,84.1c0-2.5-0.1-4.4-0.2-5.8c-0.2-1.4-0.4-2.7-0.8-3.9\n\t\tc-0.7-2.2-1.9-3.9-3.6-5.3c-1.7-1.3-3.8-2-6.3-2c-2.4,0-4.5,0.7-6.2,2c-1.7,1.3-3,3.1-3.8,5.3c-1.1,2.4-1.7,5.7-1.7,9.8\n\t\tc0,4.5,0.5,7.8,1.5,9.7c0.7,2.2,2,3.9,3.8,5.3c1.8,1.3,4,2,6.5,2c2.6,0,4.7-0.7,6.4-2c1.7-1.3,2.9-3.1,3.5-5.2\n\t\tC649.7,91.7,650,88.4,650,84.1z\" />\n<path class=\"st0\" d=\"M689.5,111.5c-4.9-3.6-8.2-8.4-10-14.6c-1.1-3.8-1.7-7.9-1.7-12.4c0-4.8,0.6-9.1,1.7-12.9\n\t\tc1.9-6,5.2-10.7,10.1-14.1c4.9-3.4,10.7-5.1,17.5-5.1c6.6,0,12.3,1.7,17,5c4.7,3.4,8,8,10,14c1.3,4,1.9,8.3,1.9,12.7\n\t\tc0,4.4-0.6,8.5-1.7,12.3c-1.8,6.3-5.1,11.3-9.9,14.9c-4.8,3.6-10.6,5.4-17.4,5.4C700.2,116.8,694.3,115,689.5,111.5z M713.8,99.7\n\t\tc1.9-1.6,3.2-3.8,4-6.7c0.6-2.6,1-5.4,1-8.5c0-3.4-0.3-6.3-1-8.6c-0.9-2.8-2.3-4.9-4.1-6.4c-1.9-1.5-4.1-2.3-6.8-2.3\n\t\tc-2.8,0-5,0.8-6.9,2.3c-1.8,1.5-3.1,3.7-3.9,6.4c-0.6,1.9-1,4.8-1,8.6c0,3.6,0.3,6.5,0.8,8.5c0.8,2.8,2.2,5.1,4.1,6.7\n\t\tc1.9,1.6,4.2,2.4,7,2.4C709.6,102.1,711.9,101.3,713.8,99.7z\" />\n</g>\n<g>\n<path class=\"st1\" d=\"M137.9,35.6c-3.7-0.6-7.6-0.9-11.5-0.9c-32.2,0-59.6,20.4-70,49.1C34.5,77.3,16.8,61.1,8.3,40.2\n\t\tC20.7,16.3,45.7,0,74.4,0C101.3,0,124.8,14.2,137.9,35.6z\" />\n<path class=\"st2\" d=\"M148.8,74.4c0,39.7-31.1,72.2-70.3,74.3c9.7-12.6,15.4-28.3,15.4-45.4c0-18.6-6.8-35.7-18.1-48.7l-0.3,0.3\n\t\tc0.1-0.1,0.2-0.2,0.3-0.3l0,0c0.1-0.1,0.3-0.3,0.5-0.4c0.1-0.1,0.3-0.3,0.5-0.4c0.7-0.6,1.4-1.2,2.1-1.8c0.1-0.1,0.2-0.2,0.3-0.2\n\t\tc0.6-0.5,1.3-1,1.9-1.5c0.7-0.6,1.5-1.1,2.3-1.7c0.8-0.5,1.5-1.1,2.3-1.6c0.8-0.5,1.5-1,2.3-1.4c0.4-0.3,0.9-0.5,1.3-0.8\n\t\tc0.4-0.2,0.8-0.5,1.2-0.7c0.4-0.2,0.8-0.5,1.2-0.7c1.7-0.9,3.4-1.7,5.1-2.5c0.4-0.2,0.9-0.4,1.3-0.5c0.9-0.4,1.8-0.7,2.6-1\n\t\tc0.4-0.2,0.9-0.3,1.3-0.5c0.4-0.2,0.9-0.3,1.4-0.5c0.4-0.1,0.9-0.3,1.4-0.4c1.7-0.5,3.4-1,5.1-1.3c0.5-0.1,1.1-0.2,1.6-0.3\n\t\tc0.3-0.1,0.7-0.1,1-0.2c0.7-0.1,1.4-0.3,2.2-0.4c0.5-0.1,1-0.2,1.4-0.2c0.5-0.1,1-0.1,1.5-0.2c0.5-0.1,1-0.1,1.5-0.2\n\t\tc0.6-0.1,1.1-0.1,1.7-0.2c0.6,0,1.3-0.1,1.9-0.1c0.6,0,1.3-0.1,1.9-0.1c0.6,0,1.3,0,1.9,0c0.8,0,1.6,0,2.4,0c0.3,0,0.6,0,0.9,0\n\t\tc0.8,0,1.7,0.1,2.5,0.1c1,0.1,1.9,0.2,2.8,0.3c0.5,0,0.9,0.1,1.4,0.2c0.5,0.1,0.9,0.1,1.4,0.2v0C144.8,46.9,148.8,60.2,148.8,74.4z\n\t\t\" />\n<path class=\"st3\" d=\"M93.9,103.3c0,17.1-5.8,32.8-15.4,45.4c-1.3,0.1-2.7,0.1-4.1,0.1C33.3,148.8,0,115.5,0,74.4\n\t\tc0-12.3,3-24,8.3-34.2c11,27.3,37.8,46.5,69,46.5c4.9,0,9.8-0.5,14.4-1.4C93.1,91.1,93.9,97.1,93.9,103.3z\" />\n</g>\n</svg>\n</div>\n</a>\n</div>\n<div class=\"grid lr-hostr-frame cf\">\n<div class=\"lr-frames lr-forms-container\">\n<div class=\"lr-forms-container-inner\">\n<div class=\"img-holder\">\n<div class=\"info-holder\">\n<img src=\"https://cdn.loginradius.com/hub/prod/v1/hosted-page-default-images/lr-bg2.svg\" alt=\"\" />\n</div>\n</div>\n<div class=\"lr-forms-items\">\n<div class=\"lr-forms-items-header\">\n<h2>\nLogin to your account or register yourself if you don't have an account yet.\n</h2>\n<p>\nIt's quick and easy!\n</p>\n</div>\n<div id=\"lr-traditional-login\" class=\"lr-form-frame\">\n<div id=\"login-container\" class=\"lr-widget-container\"></div>\n\n</div>\n<div id=\"lr-raas-registartion\" class=\"lr-form-frame\">\n<div id=\"registration-container\" class=\"lr-widget-container\"></div>\n\n</div>\n<div id=\"lr-raas-forgotpassword\" class=\"lr-form-frame\">\n<h2 id=\"forgot-pw-label\" class=\"lr-form-login-message\">Forgot Password</h2>\n<p id=\"forgot-pw-message\" class=\"lr-form-login-subnote\">\nWe'll email or sms you an instruction on how to reset your password.\n</p>\n<div id=\"forgotpassword-container\" class=\"lr-widget-container\"></div>\n<div class=\"lr-link-box\">\n<a class=\"lr-raas-login-link\">Login</a>\n<a class=\"lr-register-link\">Register</a>\n</div>\n</div>\n<div id=\"lr-raas-sociallogin\" class=\"lr-form-frame\">\n<h2 class=\"lr-form-login-message\">Complete your Profile</h2>\n<p class=\"lr-form-login-subnote\">\nRequire to fill all mandatory fields.\n</p>\n<div id=\"sociallogin-container\" class=\"lr-widget-container\"></div>\n</div>\n<div id=\"lr-raas-resetpassword\" class=\"lr-form-frame\">\n<h2 class=\"lr-form-login-message\">Reset your Password</h2>\n<p class=\"lr-form-login-subnote\">\nReset your password to get back access of your account\n</p>\n<div id=\"resetpassword-container\" class=\"lr-widget-container\"></div>\n</div>\n<div id=\"lr-social-login\" class=\"lr-social-login-frame lr-frames lr-sample-background-enabled cf\">\n<span id=\"social-block-label\" class=\"lr-social-login-message\">Or login with</span>\n<div id=\"interfacecontainerdiv\" class=\"lr-sl-shaded-brick-frame cf lr-widget-container\"></div>\n<script type=\"text/html\" id=\"loginradiuscustom_tmpl\">\n                  <div class=\"lr-provider-wrapper\"><span class=\"lr-provider-label lr-sl-shaded-brick-button lr-flat-<#=Name.toLowerCase()#>\" onclick=\" return LRObject.util.openWindow('<#= Endpoint #>'); \" title=\"Login with <#= Name #>\" alt=\"Login with <#= Name#>\"><span class=\"lr-sl-icon lr-sl-icon-<#= Name.toLowerCase()#>\"></span> </span><span class=\"lr-provider-provider-name\" onclick=\" return LRObject.util.openWindow('<#= Endpoint #>'); \" title=\"Login with <#= Name #>\" alt=\"Login with <#= Name#>\"><#= Name#></span></div>\n                </script>\n</div>\n<div class=\"lr-tabs\">\n<a id=\"login-label\" class=\"lr-raas-login-link\">Login </a>\n<a id=\"sign-up-label\" class=\"lr-register-link\">Register</a>\n</div>\n</div>\n</div>\n</div>\n</div>\n</div>\n<div class=\"lr_fade lr-loading-screen-overlay\" id=\"loading-spinner\">\n<div class=\"load-dot\"></div>\n<div class=\"load-dot\"></div>\n<div class=\"load-dot\"></div>\n<div class=\"load-dot\"></div>\n</div>\n</div>\n\n</div>",
	EndScript:    "",
	BeforeScript: updatePath("/Themes/Theme-3/auth/js/before-script.js"),
	IsActive:     true,
	MainScript:   "",
	Status:       "4",
}

var SmtpProviders = map[int]SmtpProviderSchema {
    0: {
      Name: "Mailazy",
      Display: "Mailazy",
	  SmtpHost: "",
      SmtpPort: "",
      EnableSSL: false,
    },
    1: {
      Name: "Amazon SES (US East)",
      Display: "AmazonSES-USEast",
	  SmtpHost: "email-smtp.us-east-1.amazonaws.com",
      SmtpPort: "587",
      EnableSSL: true,
    },
    2: {
      Name: "Amazon SES (US West)",
      Display: "AmazonSES-USWest",
	  SmtpHost: "email-smtp.us-west-2.amazonaws.com",
      SmtpPort: "587",
      EnableSSL: true,
    },
    3: {
      Name: "AmazonSES(EU)",
      Display: "AmazonSES-EU",
	  SmtpHost: "email-smtp.eu-west-1.amazonaws.com",
      SmtpPort: "587",
      EnableSSL: true,
    },
    4: {
      Name: "Gmail",
      Display: "Gmail",
	  SmtpHost: "smtp.gmail.com",
      SmtpPort: "587",
      EnableSSL: true,
    },
    5: {
      Name: "Mandrill",
      Display: "Mandrill",
	  SmtpHost: "smtp.mandrillapp.com",
      SmtpPort: "587",
      EnableSSL: true,
    },
    6: {
      Name: "Rackspace-mailgun",
      Display: "Rackspace-mailgun",
	  SmtpHost: "smtp.mailgun.org",
      SmtpPort: "587",
      EnableSSL: true,
    },
    7: {
      Name: "SendGrid",
      Display: "SendGrid",
	  SmtpHost: "smtp.sendgrid.net",
      SmtpPort: "587",
      EnableSSL: true,
    },
    8: {
      Name: "Yahoo",
      Display: "Yahoo",
	  SmtpHost: "smtp.mail.yahoo.com",
      SmtpPort: "587",
      EnableSSL: true,
    },
    9: {
      Name: "Custom SMTP Providers",
      Display: "CustomSMTPProviders",
	  SmtpHost: "",
      SmtpPort: "",
      EnableSSL: false,
    },
}