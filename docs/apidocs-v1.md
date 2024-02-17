# v1

## user

### account


### **GET** _v1/user/_
> Display the Hello message
#### NEED AUTHENTICATION
False
#### PARAMETERS
| Name | Type | Description | Required |
|------|------|-------------|----------|

#### EXAMPLE RESPONSE
```json
200 OK
{
    "message": "Hello, World!"
}
```



### **POST** _v1/user/sendEmail_
> 发送验证邮件码到指定邮箱（目前直接显示在debug log）
#### NEED AUTHENTICATION
False
#### PARAMETERS
| Name | Type | Description | Required |
|------|------|-------------|----------|
email| `string` |邮箱地址（目前有指定后缀）| `True` 
#### EXAMPLE RESPONSE
```json
200 OK
{
    "msg": "Code send successfully!"
}
```

```json
400 Bad Request
{
    "msg": "Invalid email address!"
}
```



### **POST** _v1/user/verifyEmail_
> 验证邮箱
#### NEED AUTHENTICATION
False
#### PARAMETERS
| Name | Type | Description | Required |
|------|------|-------------|----------|
email| `string` |邮箱地址| `True` 
code| `string` |邮箱验证码验证码| `True` 
#### EXAMPLE RESPONSE
```json
200 OK
{
    "msg": "Email verified successfully!"
}
```

```json
400 Bad Request
{
    "msg": "Invalid email address or code!"
}
```



### **GET** _v1/user/isExistEmail_
> 判断邮箱是否已注册
#### NEED AUTHENTICATION
False
#### PARAMETERS
| Name | Type | Description | Required |
|------|------|-------------|----------|
email| `string` |邮箱地址| `True` 
#### EXAMPLE RESPONSE
```json
200 OK
{
    "msg": "Email is exist.",
    "data": true
}
```

```json
200 OK
{
    "msg": "Email is not exist.",
    "data": false
}
```

```json
400 Bad Request
{
    "msg": "Email is not valid."
}
```



### **POST** _v1/user/reg_
> 注册用户
#### NEED AUTHENTICATION
False
#### PARAMETERS
| Name | Type | Description | Required |
|------|------|-------------|----------|
username| `string` |用户名| `True` 
email| `string` |邮箱地址| `True` 
code| `string` |邮箱验证码| `True` 
password| `string` |密码 (sha-256加密后的密码)| `True` 
#### EXAMPLE RESPONSE
```json
200 OK
{
    "msg": "Register successfully!"
}
```

```json
400 Bad Request
{
    "msg": "Username is exist."
}
```



### **POST** _v1/user/login_
> 用户登录
#### NEED AUTHENTICATION
False
#### PARAMETERS
| Name | Type | Description | Required |
|------|------|-------------|----------|
username| `string` |用户名| `True` 
password| `string` |密码 (sha-256加密后的密码)| `True` 
#### EXAMPLE RESPONSE
```json
200 OK
{
    "msg": "Login Successfully!",
    "data": "this is token"
}
```

```json
400 Bad Request
{
    "msg": "Login Failed.",
    "data": false
}
```



### **GET** _v1/user/getUserInfo_
> 获取用户信息
#### NEED AUTHENTICATION
True
#### PARAMETERS
| Name | Type | Description | Required |
|------|------|-------------|----------|
user| `string` |指定的用户| `True` 
#### EXAMPLE RESPONSE
```json
200 OK
{
    "msg": "Get user info successfully!",
    "data": {
        "account": "<UserAccountObject>",
        "info": "<UserInfoObject>"
    }
}
```

```json
400 Bad Request
{
    "msg": "Get user info failed.",
    "data": false
}
```


### siteMessage(仅作为示例)

### attachment(仅作为示例)

