#!/bin/bash
RED='\033[0;31m'
GRN='\033[0;32m'
BLU='\033[0;34m'
NC='\033[0m'
echo ""
echo -e "监管机一站式自动化工具"
echo ""
PS3='Please enter your choice: '
options=("自动绕过（恢复模式）" "屏蔽通知（桌面模式）" "屏蔽通知（恢复模式）" "查看监管状态" "退出")
select opt in "${options[@]}"; do
	case $opt in
	"自动绕过（恢复模式）")
		echo -e "${GRN}自动绕过（恢复模式）"
		if [ -d "/Volumes/Data" ]; then
   			diskutil rename "Data" "Data"
		fi
		echo -e "${GRN}创建新用户"
        echo -e "${BLU}请按下回车继续"
  		echo -e "请输入用户名 请使用纯英文且不包含特殊符号（默认:MacBook）"
		read realName
  		realName="${realName:=MacBook}"
    	echo -e "${BLUE}请输入用户名 ${RED}不允许字母间有空格和特殊符号，且保持与刚才输入的用户名一致！ ${GRN}（默认:MacBook）"
      	read username
		username="${username:=MacBook}"
  		echo -e "${BLUE}请输入密码（默认: 123456）"
    	read passw
      	passw="${passw:=123456}"
		dscl_path='/Volumes/Data/private/var/db/dslocal/nodes/Default' 
        echo -e "${GREEN}创建用户中"
  		# Create user
    	dscl -f "$dscl_path" localhost -create "/Local/Default/Users/$username"
      	dscl -f "$dscl_path" localhost -create "/Local/Default/Users/$username" UserShell "/bin/zsh"
	    dscl -f "$dscl_path" localhost -create "/Local/Default/Users/$username" RealName "$realName"
	 	dscl -f "$dscl_path" localhost -create "/Local/Default/Users/$username" RealName "$realName"
	    dscl -f "$dscl_path" localhost -create "/Local/Default/Users/$username" UniqueID "501"
	    dscl -f "$dscl_path" localhost -create "/Local/Default/Users/$username" PrimaryGroupID "20"
		mkdir "/Volumes/Data/Users/$username"
	    dscl -f "$dscl_path" localhost -create "/Local/Default/Users/$username" NFSHomeDirectory "/Users/$username"
	    dscl -f "$dscl_path" localhost -passwd "/Local/Default/Users/$username" "$passw"
	    dscl -f "$dscl_path" localhost -append "/Local/Default/Groups/admin" GroupMembership $username
		echo "0.0.0.0 deviceenrollment.apple.com" >>/Volumes/Macintosh\ HD/etc/hosts
		echo "0.0.0.0 mdmenrollment.apple.com" >>/Volumes/Macintosh\ HD/etc/hosts
		echo "0.0.0.0 iprofiles.apple.com" >>/Volumes/Macintosh\ HD/etc/hosts
         echo "0.0.0.0 acmdm.apple.com" >>/Volumes/Macintosh\ HD/etc/hosts
         echo "0.0.0.0 albert.apple.com" >>/Volumes/Macintosh\ HD/etc/hosts
        echo -e "${GREEN}成功隔绝本机与MDM服务器的通信${NC}"
		# echo "Remove config profile"
  	touch /Volumes/Data/private/var/db/.AppleSetupDone
        rm -rf /Volumes/Macintosh\ HD/var/db/ConfigurationProfiles/Settings/.cloudConfigHasActivationRecord
	rm -rf /Volumes/Macintosh\ HD/var/db/ConfigurationProfiles/Settings/.cloudConfigRecordFound
	touch /Volumes/Macintosh\ HD/var/db/ConfigurationProfiles/Settings/.cloudConfigProfileInstalled
	touch /Volumes/Macintosh\ HD/var/db/ConfigurationProfiles/Settings/.cloudConfigRecordNotFound
     launchctl disable system/com.apple.ManagedClient.enroll
         echo "---------------------------"
         echo -e "${CYAN}------ 成功自动绕过 ------${NC}"
         echo -e "${CYAN}------ 在终端手动输入reboot进行重启！ ------${NC}"
		break
		;;
    "屏蔽通知（桌面模式）")
    	echo -e "${RED}请直接输入您的密码以继续（此处输入不会显示状态，请直接键入并按下回车）${NC}"
        sudo profiles remove -all
        sudo rm /var/db/ConfigurationProfiles/Settings/.cloudConfigHasActivationRecord
        sudo rm /var/db/ConfigurationProfiles/Settings/.cloudConfigRecordFound
        sudo touch /var/db/ConfigurationProfiles/Settings/.cloudConfigProfileInstalled
        sudo touch /var/db/ConfigurationProfiles/Settings/.cloudConfigRecordNotFound
        launchctl disable system/com.apple.ManagedClient.enroll
         echo "------ 上述如果报错为正常现象，并代表屏蔽成功 ------"
         echo -e "${CYAN}------ 成功屏蔽通知 ------${NC}"
         echo -e "${CYAN}------ 重启电脑即可正常使用，享受吧！ ------${NC}"
        break
        ;;
    "屏蔽通知（恢复模式）")
         echo "0.0.0.0 deviceenrollment.apple.com" >>/Volumes/Macintosh\ HD/etc/hosts
		echo "0.0.0.0 mdmenrollment.apple.com" >>/Volumes/Macintosh\ HD/etc/hosts
		echo "0.0.0.0 iprofiles.apple.com" >>/Volumes/Macintosh\ HD/etc/hosts
         echo "0.0.0.0 acmdm.apple.com" >>/Volumes/Macintosh\ HD/etc/hosts
         echo "0.0.0.0 albert.apple.com" >>/Volumes/Macintosh\ HD/etc/hosts
        rm -rf /Volumes/Macintosh\ HD/var/db/ConfigurationProfiles/Settings/.cloudConfigHasActivationRecord
	rm -rf /Volumes/Macintosh\ HD/var/db/ConfigurationProfiles/Settings/.cloudConfigRecordFound
	touch /Volumes/Macintosh\ HD/var/db/ConfigurationProfiles/Settings/.cloudConfigProfileInstalled
	touch /Volumes/Macintosh\ HD/var/db/ConfigurationProfiles/Settings/.cloudConfigRecordNotFound
     launchctl disable system/com.apple.ManagedClient.enroll
         echo "------ 上述如果报错为正常现象，并代表屏蔽成功 ------"
         echo -e "${CYAN}------ 成功屏蔽通知，输入reboot并回车电脑便会自动重启 ------${NC}"
         echo -e "${CYAN}------ 进入桌面以后，打开终端，再次索要代码运行“屏蔽通知（桌面模式）”！ ------${NC}"
        break
        ;;
	"查看监管状态")
		echo ""
		echo -e "${GRN}查看监管状态，报错则代表屏蔽成功${NC}"
		echo ""
		echo -e "${RED}请输入您的密码以继续（此处输入不会显示状态，请直接键入并按下回车）${NC}"
		echo ""
		sudo profiles show -type enrollment
		break
		;;
	"退出")
		break
		;;
	*) echo "Invalid option $REPLY" ;;
	esac
done
