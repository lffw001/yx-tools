include $(TOPDIR)/rules.mk

LUCI_TITLE:=luci-app-yx-tools
LUCI_DEPENDS:=+yx-tools +luci

PKG_NAME:=luci-app-yx-tools
PKG_VERSION:=1
PKG_RELEASE:=1

PKG_MAINTAINER:=byJoey <https://github.com/byJoey>
PKG_LICENSE:=MIT

include $(INCLUDE_DIR)/package.mk
include $(INCLUDE_DIR)/luci.mk

define Package/$(PKG_NAME)
  SECTION:=luci
  CATEGORY:=LuCI
  SUBMENU:=3. Applications
  TITLE:=LuCI support for yx-tools
  DEPENDS:=$(LUCI_DEPENDS)
endef

define Package/$(PKG_NAME)/description
 Web UI for yx-tools (start/stop and basic configuration).
endef

define Build/Compile
  true
endef

define Package/$(PKG_NAME)/install
	$(INSTALL_DIR) $(1)/usr/lib/lua/luci/controller
	# Install CBI model into the root of model/cbi (NOT a subdirectory).
	# LuCI's cbi("yx_tools") expects /usr/lib/lua/luci/model/cbi/yx_tools.lua
	$(INSTALL_DIR) $(1)/usr/lib/lua/luci/model/cbi
	$(INSTALL_DATA) ./luasrc/controller/*.lua $(1)/usr/lib/lua/luci/controller/
	$(INSTALL_DATA) ./luasrc/model/cbi/yx_tools.lua $(1)/usr/lib/lua/luci/model/cbi/yx_tools.lua
	$(INSTALL_DIR) $(1)/www/luci-static/resources/yx-tools
	$(CP) ./htdocs/* $(1)/www/luci-static/resources/yx-tools/ || true
endef

$(eval $(call BuildPackage,$(PKG_NAME)))
