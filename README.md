# XrayHelper
A unified helper for Android to control system proxy

# Config
XrayHelper use yml format configuration file, you can customize the path with the `-c` option.  
[Example of xrayhelper config](config.yml)

# Commands
## Control Core Service
`xrayhelper service start`, start core service  
`xrayhelper service stop`, stop core service  
`xrayhelper service restart`, restart core service  
`xrayhelper service status`, show core status  

## Control System Proxy
`xrayhelper proxy enable`, enable system proxy  
`xrayhelper proxy disable`, disable system proxy  
`xrayhelper proxy refresh`, refresh system proxy rule  
