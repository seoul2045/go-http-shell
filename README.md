# go-http-shell
A basic Golang HTTP REST API for executing shell commands using the standard library.  
API accepts a shell command via a POST request to `localhost:8080/api/cmd` and JSON body: `{"request": "pwd"}`    
Errors are handled gracefully.  
If a successful command has no output (e.g. touch, rm), then a SUCCESS message with exit status 0 no error will be returned.

![image](https://github.com/seoul2045/go-http-shell/assets/13395406/c99a4355-42f5-41a1-916e-69214139a804)

![image](https://github.com/seoul2045/go-http-shell/assets/13395406/fc691afa-c3b1-4479-bb88-795e3516a236)

![image](https://github.com/seoul2045/go-http-shell/assets/13395406/2d302478-6d39-4c63-b2f8-7496814d2fcd)

![image](https://github.com/seoul2045/go-http-shell/assets/13395406/ec92a3d2-ace1-4a75-bad5-fcba2516ebbd)

![image](https://github.com/seoul2045/go-http-shell/assets/13395406/7789a9ed-1b7d-41d4-ac61-b14496538f84)
