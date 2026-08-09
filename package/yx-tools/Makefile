# SPDX-License-Identifier: GPL-2.0-only
#
# Copyright (C) 2024 byJoey
#

include $(TOPDIR)/rules.mk

PKG_NAME:=yx-tools
PKG_VERSION:=3.0.0
PKG_RELEASE:=1

PKG_MAINTAINER:=byJoey <https://github.com/byJoey>
PKG_LICENSE:=MIT
PKG_LICENSE_FILES:=LICENSE

include $(INCLUDE_DIR)/package.mk

define Package/yx-tools
  SECTION:=net
  CATEGORY:=Network
  SUBMENU:=Web Servers/Proxies
  TITLE:=Cloudflare优选IP测速工具
  URL:=https://github.com/byJoey/yx-tools
  DEPENDS:=+curl
endef

define Package/yx-tools/description
 Cloudflare优选IP测速工具，支持CLI和Web界面，
 可测速Cloudflare各数据中心延迟和下载速度，
 支持一键上报到cfnew面板或GitHub。
endef

define Package/yx-tools/conffiles
/etc/config/yx-tools
endef

define Build/Prepare
true
endef

define Build/Configure
true
endef

define Build/Compile
true
endef

define Package/yx-tools/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/yx $(1)/usr/bin/yx
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) ./files/etc/init.d/yx-tools $(1)/etc/init.d/yx-tools
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) ./files/etc/config/yx-tools $(1)/etc/config/yx-tools
	$(INSTALL_DIR) $(1)/usr/share/yx-tools
	$(CP) ./files/www/* $(1)/usr/share/yx-tools/
	$(INSTALL_DIR) $(1)/www
	$(LN) yx-tools $(1)/www/yx-tools
endef

$(eval $(call BuildPackage,yx-tools))
