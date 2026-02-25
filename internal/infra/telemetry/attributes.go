package telemetry

import "go.opentelemetry.io/otel/attribute"

const (
	AttrServiceLayer   = attribute.Key("service.layer")
	AttrServiceMethod  = attribute.Key("service.method")
	AttrCommandType    = attribute.Key("command.type")
	AttrCommandName    = attribute.Key("command.name")
	AttrHandlerName    = attribute.Key("handler.name")
	AttrRepositoryName = attribute.Key("repository.name")
	AttrUserID         = attribute.Key("user.id")
	AttrSessionID      = attribute.Key("session.id")
	AttrEntityType     = attribute.Key("entity.type")
	AttrEntityID       = attribute.Key("entity.id")
	AttrOperationType  = attribute.Key("operation.type")
	AttrOperationName  = attribute.Key("operation.name")
	AttrDBSystem       = attribute.Key("db.system")
	AttrDBOperation    = attribute.Key("db.operation")
	AttrDBTable        = attribute.Key("db.sql.table")
	AttrResultCount    = attribute.Key("result.count")
	AttrResultError    = attribute.Key("result.error")
)

const (
	LayerHTTP       = "http"
	LayerService    = "service"
	LayerRepository = "repository"
	LayerCommand    = "command"
)

const (
	OpCreate = "create"
	OpRead   = "read"
	OpUpdate = "update"
	OpDelete = "delete"
	OpList   = "list"
	OpQuery  = "query"
)

func ServiceLayer(layer string) attribute.KeyValue {
	return AttrServiceLayer.String(layer)
}

func ServiceMethod(method string) attribute.KeyValue {
	return AttrServiceMethod.String(method)
}

func CommandType(cmdType string) attribute.KeyValue {
	return AttrCommandType.String(cmdType)
}

func CommandName(name string) attribute.KeyValue {
	return AttrCommandName.String(name)
}

func UserID(id string) attribute.KeyValue {
	return AttrUserID.String(id)
}

func SessionID(id string) attribute.KeyValue {
	return AttrSessionID.String(id)
}

func EntityType(entityType string) attribute.KeyValue {
	return AttrEntityType.String(entityType)
}

func EntityID(id string) attribute.KeyValue {
	return AttrEntityID.String(id)
}

func OperationType(opType string) attribute.KeyValue {
	return AttrOperationType.String(opType)
}

func OperationName(name string) attribute.KeyValue {
	return AttrOperationName.String(name)
}

func DBSystem(system string) attribute.KeyValue {
	return AttrDBSystem.String(system)
}

func DBOperation(op string) attribute.KeyValue {
	return AttrDBOperation.String(op)
}

func DBTable(table string) attribute.KeyValue {
	return AttrDBTable.String(table)
}

func ResultCount(count int) attribute.KeyValue {
	return AttrResultCount.Int(count)
}

func ResultError(hasError bool) attribute.KeyValue {
	return AttrResultError.Bool(hasError)
}
