set dir=D:\book\guaji_server\server\src\server\gameServer\game\config\csv\.
del /s /q "%dir%"
xcopy D:\book\guaji_art\data\csv\. D:\book\guaji_server\server\src\server\gameServer\game\config\csv\ /s/e

set dir=D:\book\guaji_server\server\src\server\gameServer\game\config\json\.
del /s /q "%dir%"
xcopy D:\book\guaji_art\templates\Jason文件及导出工具\result\. D:\book\guaji_server\server\src\server\gameServer\game\config\json /s/e

