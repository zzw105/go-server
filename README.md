# Go Server

## 环境配置

项目使用环境变量来管理敏感信息（如数据库密码）。

### 首次设置

1. 复制环境变量模板文件：
   ```bash
   cp .env.example .env
   ```

2. 编辑 `.env` 文件，填入你的实际配置：
   ```
   DB_HOST=your_database_host
   DB_PORT=3306
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database_name
   DB_CHARSET=utf8mb4
   ```

3. 确保 `.env` 文件不会被提交到 Git（已在 `.gitignore` 中配置）

### 数据库迁移

执行以下 SQL 语句来更新数据库表结构：

```sql
-- 修改 classifications 表的字段类型为 BIGINT
ALTER TABLE classifications
  MODIFY COLUMN id BIGINT NOT NULL AUTO_INCREMENT,
  MODIFY COLUMN parent_id BIGINT NOT NULL DEFAULT 0;

-- 添加 sort 字段
ALTER TABLE classifications ADD COLUMN sort INT NOT NULL DEFAULT 0;
```

## 运行项目

```bash
go run main.go
```

## API 接口

### 分类管理

- `GET /classification` - 获取分类树
- `PUT /classification` - 更新整个分类树（全量替换）