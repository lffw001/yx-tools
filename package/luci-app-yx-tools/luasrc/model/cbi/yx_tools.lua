m = Map("yx-tools", "YX Tools", "Configure and control yx-tools service")

s = m:section(TypedSection, "global", "Global Settings")
s.anonymous = true

listen = s:option(Value, "listen_addr", "Listen address")
listen.default = "0.0.0.0:8080"

startmode = s:option(ListValue, "start_mode", "Start mode")
startmode:value("web", "Web")
startmode:value("test", "Test")
startmode.default = "web"

start = s:option(Button, "_start", "Start service")
function start.write(self, section)
  luci.sys.call("/etc/init.d/yx-tools start >/dev/null 2>&1 &")
end

stop = s:option(Button, "_stop", "Stop service")
function stop.write(self, section)
  luci.sys.call("/etc/init.d/yx-tools stop >/dev/null 2>&1 &")
end

return m
