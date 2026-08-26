# Cairn Sonar

Cairn Sonar 是一个面向石窟岩壁声学巡检的 Go 后端服务。它记录勘测批次、声脉冲与回波剖面，支持异常复核、巡检路线规划、仪器校准和归档。数据落在本地 SQLite 文件中，服务重启后仍可恢复。

## 启动

```bash
go run . --db ./cairn-sonar.db
```

## 主要接口

- `GET /healthz`
- `POST /surveys` / `GET /surveys`
- `GET /survey/{id}` / `POST /survey/{id}/start`
- `POST /route/{surveyID}` / `GET /route/{routeID}`
- `GET /export/{surveyID}`

操作页位于 `web/index.html`，仅作为后端服务的薄客户端。
