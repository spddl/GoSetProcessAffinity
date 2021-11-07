## GoSetProcessAffinity
---------------

### input
```
-pid 1337 (running Process identifier)
-proc csgo.exe (running *.exe)
-exe "C:\Windows\notepad.exe" (new *.exe)
```

### configuration
```
-logging (enable logging in debug.log)
-boost (https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi)
-cpu "0,2,4" (ProcessAffinityMask)
-priorityClass (https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-setpriorityclass)
-ioPriority (cchanges the I/O priority)
-memoryPriority (Process Page Priority)
```

### examples

+ GoSetProcessAffinity.exe -exe csgo.exe -priorityClass 4 -ioPriority 3 -boost
+ GoSetProcessAffinity.exe -proc winlogon.exe -suspend
+ GoSetProcessAffinity.exe -proc winlogon.exe -resume
+ GoSetProcessAffinity.exe -logging -pid 8236,10640 -priorityClass 4 -ioPriority 3 -memoryPriority 5
+ GoSetProcessAffinity.exe -logging -pid 10640 -priorityClass 4 -ioPriority 3 -memoryPriority 5
+ GoSetProcessAffinity.exe -logging -proc vlc.exe -priorityClass 1 -ioPriority 3 -memoryPriority 5