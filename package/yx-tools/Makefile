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

# 在 Build/Compile 阶段优先选择仓库中为目标架构准备的预构建二进制：
# 支持的文件名（按优先级）：
#   files/yx-<arch>           (例如 files/yx-aarch64_generic)
#   files/yx-<target_arch>    (尝试常见的 OpenWrt 变量名)
#   files/yx                  (通用二进制)
# 若找到则复制到 $(PKG_BUILD_DIR)/yx 以便 install 阶段使用。
define Build/Compile
	mkdir -p $(PKG_BUILD_DIR)
	# 尝试一系列可能的文件名
	FOUND=""
	for f in \
		./files/yx-$(ARCH) \
		./files/yx-$(TARGET_ARCH) \
		./files/yx-$(PKG_ARCH) \
		./files/yx-$(BOARD) \
		./files/yx; do \
		if [ -f $$f ]; then \
			cp $$f $(PKG_BUILD_DIR)/yx; chmod 755 $(PKG_BUILD_DIR)/yx; FOUND=1; break; \
		fi; \
	done
	if [ -z "$$FOUND" ]; then \
		echo "No prebuilt binary found in ./files for yx; package will be built empty unless you provide a binary.\nLooking for: ./files/yx-$(ARCH), ./files/yx-$(TARGET_ARCH), ./files/yx-$(PKG_ARCH), ./files/yx-$(BOARD), ./files/yx"; \
	fi
endef

define Build/Install
	true
endef

define Package/yx-tools/install
	$(INSTALL_DIR) $(1)/usr/bin
	# 从 $(PKG_BUILD_DIR)/yx 安装（Build/Compile 会把仓库 ./files 中的预构建二进制复制到此处）
	$(if $(wildcard $(PKG_BUILD_DIR)/yx),$(INSTALL_BIN) $(PKG_BUILD_DIR)/yx $(1)/usr/bin/yx,$(error "yx binary not found in $(PKG_BUILD_DIR)/yx - please place built binary in ./files/yx-<arch> or ./files/yx"))
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) ./files/etc/init.d/yx-tools $(1)/etc/init.d/yx-tools
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) ./files/etc/config/yx-tools $(1)/etc/config/yx-tools
	$(INSTALL_DIR) $(1)/usr/share/yx-tools
	$(CP) ./files/www/* $(1)/usr/share/yx-tools/ || true
	$(INSTALL_DIR) $(1)/www
	$(LN) yx-tools $(1)/www/yx-tools
endef

$(eval $(call BuildPackage,yx-tools))
