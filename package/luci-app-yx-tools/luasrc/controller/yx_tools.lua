module("luci.controller.yx_tools", package.seeall)

function index()
  entry({"admin", "services", "yx_tools"}, cbi("yx_tools"), _("YX Tools"), 60).dependent = true
end
