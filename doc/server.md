## Server API Docs

### 订阅事件 /sub

POST /sub HTTP/1.1

Content-Type: application/json

客户端轮询方式调用该接口向服务器注册，并且订阅消息事件。服务器端注册信息会在10秒后过期

请求参数：

|  参数   | 数据类型 | 说明  |
|  :----  |:---- |:----  |
| dev_uuid  | string |必填：设备UUID，客户端唯一标识，可使用UUID或GUID |
| describe  | string |描述信息 |
| bluetooth_id | string | 必填：蓝牙位置ID  |
| type | uint| 必填：设备类型: 1. PC  2. Mobile |

响应参数：

|  参数   | 数据类型 | 说明  |
|  :----  |:---- |:----  |
| message  | string |描述信息|
| status  | string |状态 |
| data | []string | 事件数据：例如图片下载列表  |

请求示例：

C#
```c#
var client = new RestClient("localhost:8000/sub");
client.Timeout = -1;
var request = new RestRequest(Method.POST);
request.AddHeader("Content-Type", "application/json");
request.AddParameter("application/json", "{\n    \"dev_uuid\": \"20190626092158-b2ccc519800e\",\n    \"type\": 2,\n    \"bluetooth_id\":\"bluetooth-1\"\n}",  ParameterType.RequestBody);
IRestResponse response = client.Execute(request);
Console.WriteLine(response.Content);

```

shell
``` shell script
curl -L -X POST 'localhost:8000/sub' \
-H 'Content-Type: application/json' \
-d '{
    "dev_uuid": "20190626092158-b2ccc519800e",
    "type": 2,
    "bluetooth_id":"bluetooth-1"
}'
```

```json
{
    "status": "OK",
    "message": "",
    "data": [
        "resource/opcom.png"
    ]
}
```

### 事件发布 /pub

POST /pub HTTP/1.1

Content-Type: application/json

PC 向对端设备发送消息，消息发出后，需要对端接收，当对端无法正常接收时，会报告异常信息

请求参数：

|  参数   | 数据类型 | 说明  |
|  :----  |:---- |:----  |
| from  | string |必填：发送方设备UUID，客户端唯一标识，可使用UUID或GUID |
| to  | string |必填：接收方设备UUID，客户端唯一标识，可使用UUID或GUID  |
| data | []string | 必填：消息数据，例如图片下载列表  |

响应参数：

|  参数   | 数据类型 | 说明  |
|  :----  |:---- |:----  |
| message  | string |描述信息|
| status  | string |状态 |

请求示例：

C#
```c#
var client = new RestClient("localhost:8000/pub");
client.Timeout = -1;
var request = new RestRequest(Method.POST);
request.AddHeader("Content-Type", "application/json");
request.AddParameter("application/json", "{\n\t\"from\":\"pc-20190626092158-b2ccc519800e\",\n\t\"to\":\"mobile-20190626092158-b2ccc519800e\",\n\t\"data\":[\"resource/opcom.png\"]\n}",  ParameterType.RequestBody);
IRestResponse response = client.Execute(request);
Console.WriteLine(response.Content);

```

shell
``` shell script
curl -L -X POST 'localhost:8000/pub' \
-H 'Content-Type: application/json' \
--data-raw '{
	"from":"pc-20190626092158-b2ccc519800e",
	"to":"mobile-20190626092158-b2ccc519800e",
	"data":["resource/opcom.png"]
}'

```

响应

```json
{
    "status": "Conflict",
    "message": "Consumer has a exception, message not be received"
}
```

### 设备列表 /devices

GET /devices HTTP/1.1

获取当前已注册的设备列表

请求参数：

|  参数   | 数据类型 | 说明  |
|  :----  |:---- |:----  |


响应参数：

|  参数   | 数据类型 | 说明  |
|  :----  |:---- |:----  |
| message  | string | 描述信息|
| status  | string | 状态 |
| count  | int | 设备数量|
| devices  | []object | 设备列表 |
| dev_uuid  | string |必填：设备UUID，客户端唯一标识，可使用UUID或GUID |
| describe  | string |描述信息 |
| bluetooth_id | string | 必填：蓝牙位置ID  |
| type | uint| 必填：设备类型: 1. PC  2. Mobile |

请求示例：

C#
```c#
var client = new RestClient("localhost:8000/devices");
client.Timeout = -1;
var request = new RestRequest(Method.GET);
IRestResponse response = client.Execute(request);
Console.WriteLine(response.Content);

```

shell
``` shell script
curl -L -X GET 'localhost:8000/devices'

```

响应

```json
{
    "status": "OK",
    "message": "",
    "count": 1,
    "devices": [
        {
            "dev_uuid": "1",
            "describe": "",
            "bluetooth_id": "a",
            "created_on": "2020-03-09T10:05:25.741984+08:00",
            "type": 2
        }
    ]
}
```

### 文件上传 /upload

POST /upload HTTP/1.1

Content-Type: multipart/form-data;

文件上传, 默认最大上传文件限制为8Mi

请求参数：

|  参数   | 数据类型 | 说明  |
|  :----  |:---- |:----  |
| file | boundry |
 

响应参数：

|  参数   | 数据类型 | 说明  |
|  :----  |:---- |:----  |
| message  | string | 描述信息|
| status  | string | 状态 |
| data | []string | 必填：消息数据，例如图片下载列表  |

请求示例：

c#
```c#
var client = new RestClient("https://localhost:8000");
RestRequest request = new RestRequest("/upload",Method.POST);
request.AddHeader("Content-Type", "multipart/form-data");
request.AddFile("file", fileStream.CopyTo, filename);
var response= client.Execute(request);
```

shell
``` shell script
curl -X POST http://localhost:8000/upload \
-F "file=@/Users/user/image/opcom.png" \
-H "Content-Type: multipart/form-data"

```

响应

```json
{
    "status": "OK",
    "message": "Uploaded",
    "data": [
        "resource/opcom.png"
    ]
}
```

### 文件下载 /resource

GET /resource HTTP/1.1

文件下载

请求参数：

|  参数   | 数据类型 | 说明  |
|  :----  |:---- |:----  |


响应参数：

|  参数   | 数据类型 | 说明  |
|  :----  |:---- |:----  |

请求示例：

shell
``` shell script
curl -L -X GET 'http://localhost:8000/resource/opcom.png'

```
