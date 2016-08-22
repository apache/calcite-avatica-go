// Code generated by protoc-gen-go.
// source: requests.proto
// DO NOT EDIT!

package message

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Request for Meta#getCatalogs()
type CatalogsRequest struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *CatalogsRequest) Reset()                    { *m = CatalogsRequest{} }
func (m *CatalogsRequest) String() string            { return proto.CompactTextString(m) }
func (*CatalogsRequest) ProtoMessage()               {}
func (*CatalogsRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

// Request for Meta#getDatabaseProperties()
type DatabasePropertyRequest struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *DatabasePropertyRequest) Reset()                    { *m = DatabasePropertyRequest{} }
func (m *DatabasePropertyRequest) String() string            { return proto.CompactTextString(m) }
func (*DatabasePropertyRequest) ProtoMessage()               {}
func (*DatabasePropertyRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

// Request for Meta#getSchemas(String, org.apache.calcite.avatica.Meta.Pat)}
type SchemasRequest struct {
	Catalog       string `protobuf:"bytes,1,opt,name=catalog" json:"catalog,omitempty"`
	SchemaPattern string `protobuf:"bytes,2,opt,name=schema_pattern,json=schemaPattern" json:"schema_pattern,omitempty"`
	ConnectionId  string `protobuf:"bytes,3,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *SchemasRequest) Reset()                    { *m = SchemasRequest{} }
func (m *SchemasRequest) String() string            { return proto.CompactTextString(m) }
func (*SchemasRequest) ProtoMessage()               {}
func (*SchemasRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

// Request for Request for Meta#getTables(String, org.apache.calcite.avatica.Meta.Pat,
//   org.apache.calcite.avatica.Meta.Pat, java.util.List)
type TablesRequest struct {
	Catalog          string   `protobuf:"bytes,1,opt,name=catalog" json:"catalog,omitempty"`
	SchemaPattern    string   `protobuf:"bytes,2,opt,name=schema_pattern,json=schemaPattern" json:"schema_pattern,omitempty"`
	TableNamePattern string   `protobuf:"bytes,3,opt,name=table_name_pattern,json=tableNamePattern" json:"table_name_pattern,omitempty"`
	TypeList         []string `protobuf:"bytes,4,rep,name=type_list,json=typeList" json:"type_list,omitempty"`
	HasTypeList      bool     `protobuf:"varint,6,opt,name=has_type_list,json=hasTypeList" json:"has_type_list,omitempty"`
	ConnectionId     string   `protobuf:"bytes,7,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *TablesRequest) Reset()                    { *m = TablesRequest{} }
func (m *TablesRequest) String() string            { return proto.CompactTextString(m) }
func (*TablesRequest) ProtoMessage()               {}
func (*TablesRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

// Request for Meta#getTableTypes()
type TableTypesRequest struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *TableTypesRequest) Reset()                    { *m = TableTypesRequest{} }
func (m *TableTypesRequest) String() string            { return proto.CompactTextString(m) }
func (*TableTypesRequest) ProtoMessage()               {}
func (*TableTypesRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

// Request for Meta#getColumns(String, org.apache.calcite.avatica.Meta.Pat,
//   org.apache.calcite.avatica.Meta.Pat, org.apache.calcite.avatica.Meta.Pat).
type ColumnsRequest struct {
	Catalog           string `protobuf:"bytes,1,opt,name=catalog" json:"catalog,omitempty"`
	SchemaPattern     string `protobuf:"bytes,2,opt,name=schema_pattern,json=schemaPattern" json:"schema_pattern,omitempty"`
	TableNamePattern  string `protobuf:"bytes,3,opt,name=table_name_pattern,json=tableNamePattern" json:"table_name_pattern,omitempty"`
	ColumnNamePattern string `protobuf:"bytes,4,opt,name=column_name_pattern,json=columnNamePattern" json:"column_name_pattern,omitempty"`
	ConnectionId      string `protobuf:"bytes,5,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *ColumnsRequest) Reset()                    { *m = ColumnsRequest{} }
func (m *ColumnsRequest) String() string            { return proto.CompactTextString(m) }
func (*ColumnsRequest) ProtoMessage()               {}
func (*ColumnsRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

// Request for Meta#getTypeInfo()
type TypeInfoRequest struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *TypeInfoRequest) Reset()                    { *m = TypeInfoRequest{} }
func (m *TypeInfoRequest) String() string            { return proto.CompactTextString(m) }
func (*TypeInfoRequest) ProtoMessage()               {}
func (*TypeInfoRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{6} }

// Request for Meta#prepareAndExecute(Meta.StatementHandle, String, long, Meta.PrepareCallback)
type PrepareAndExecuteRequest struct {
	ConnectionId      string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
	Sql               string `protobuf:"bytes,2,opt,name=sql" json:"sql,omitempty"`
	MaxRowCount       uint64 `protobuf:"varint,3,opt,name=max_row_count,json=maxRowCount" json:"max_row_count,omitempty"`
	StatementId       uint32 `protobuf:"varint,4,opt,name=statement_id,json=statementId" json:"statement_id,omitempty"`
	MaxRowsTotal      int64  `protobuf:"varint,5,opt,name=max_rows_total,json=maxRowsTotal" json:"max_rows_total,omitempty"`
	FirstFrameMaxSize int32  `protobuf:"varint,6,opt,name=first_frame_max_size,json=firstFrameMaxSize" json:"first_frame_max_size,omitempty"`
}

func (m *PrepareAndExecuteRequest) Reset()                    { *m = PrepareAndExecuteRequest{} }
func (m *PrepareAndExecuteRequest) String() string            { return proto.CompactTextString(m) }
func (*PrepareAndExecuteRequest) ProtoMessage()               {}
func (*PrepareAndExecuteRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{7} }

// Request for Meta.prepare(Meta.ConnectionHandle, String, long)
type PrepareRequest struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
	Sql          string `protobuf:"bytes,2,opt,name=sql" json:"sql,omitempty"`
	MaxRowCount  uint64 `protobuf:"varint,3,opt,name=max_row_count,json=maxRowCount" json:"max_row_count,omitempty"`
	MaxRowsTotal int64  `protobuf:"varint,4,opt,name=max_rows_total,json=maxRowsTotal" json:"max_rows_total,omitempty"`
}

func (m *PrepareRequest) Reset()                    { *m = PrepareRequest{} }
func (m *PrepareRequest) String() string            { return proto.CompactTextString(m) }
func (*PrepareRequest) ProtoMessage()               {}
func (*PrepareRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{8} }

// Request for Meta#fetch(Meta.StatementHandle, List, long, int)
type FetchRequest struct {
	ConnectionId     string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
	StatementId      uint32 `protobuf:"varint,2,opt,name=statement_id,json=statementId" json:"statement_id,omitempty"`
	Offset           uint64 `protobuf:"varint,3,opt,name=offset" json:"offset,omitempty"`
	FetchMaxRowCount uint32 `protobuf:"varint,4,opt,name=fetch_max_row_count,json=fetchMaxRowCount" json:"fetch_max_row_count,omitempty"`
	FrameMaxSize     int32  `protobuf:"varint,5,opt,name=frame_max_size,json=frameMaxSize" json:"frame_max_size,omitempty"`
}

func (m *FetchRequest) Reset()                    { *m = FetchRequest{} }
func (m *FetchRequest) String() string            { return proto.CompactTextString(m) }
func (*FetchRequest) ProtoMessage()               {}
func (*FetchRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{9} }

// Request for Meta#createStatement(Meta.ConnectionHandle)
type CreateStatementRequest struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *CreateStatementRequest) Reset()                    { *m = CreateStatementRequest{} }
func (m *CreateStatementRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateStatementRequest) ProtoMessage()               {}
func (*CreateStatementRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{10} }

// Request for Meta#closeStatement(Meta.StatementHandle)
type CloseStatementRequest struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
	StatementId  uint32 `protobuf:"varint,2,opt,name=statement_id,json=statementId" json:"statement_id,omitempty"`
}

func (m *CloseStatementRequest) Reset()                    { *m = CloseStatementRequest{} }
func (m *CloseStatementRequest) String() string            { return proto.CompactTextString(m) }
func (*CloseStatementRequest) ProtoMessage()               {}
func (*CloseStatementRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{11} }

// Request for Meta#openConnection(Meta.ConnectionHandle, Map<String, String>)
type OpenConnectionRequest struct {
	ConnectionId string            `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
	Info         map[string]string `protobuf:"bytes,2,rep,name=info" json:"info,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *OpenConnectionRequest) Reset()                    { *m = OpenConnectionRequest{} }
func (m *OpenConnectionRequest) String() string            { return proto.CompactTextString(m) }
func (*OpenConnectionRequest) ProtoMessage()               {}
func (*OpenConnectionRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{12} }

func (m *OpenConnectionRequest) GetInfo() map[string]string {
	if m != nil {
		return m.Info
	}
	return nil
}

// Request for Meta#closeConnection(Meta.ConnectionHandle)
type CloseConnectionRequest struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *CloseConnectionRequest) Reset()                    { *m = CloseConnectionRequest{} }
func (m *CloseConnectionRequest) String() string            { return proto.CompactTextString(m) }
func (*CloseConnectionRequest) ProtoMessage()               {}
func (*CloseConnectionRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{13} }

type ConnectionSyncRequest struct {
	ConnectionId string                `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
	ConnProps    *ConnectionProperties `protobuf:"bytes,2,opt,name=conn_props,json=connProps" json:"conn_props,omitempty"`
}

func (m *ConnectionSyncRequest) Reset()                    { *m = ConnectionSyncRequest{} }
func (m *ConnectionSyncRequest) String() string            { return proto.CompactTextString(m) }
func (*ConnectionSyncRequest) ProtoMessage()               {}
func (*ConnectionSyncRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{14} }

func (m *ConnectionSyncRequest) GetConnProps() *ConnectionProperties {
	if m != nil {
		return m.ConnProps
	}
	return nil
}

// Request for Meta#execute(Meta.ConnectionHandle, list, long)
type ExecuteRequest struct {
	StatementHandle    *StatementHandle `protobuf:"bytes,1,opt,name=statementHandle" json:"statementHandle,omitempty"`
	ParameterValues    []*TypedValue    `protobuf:"bytes,2,rep,name=parameter_values,json=parameterValues" json:"parameter_values,omitempty"`
	FirstFrameMaxSize  uint64           `protobuf:"varint,3,opt,name=first_frame_max_size,json=firstFrameMaxSize" json:"first_frame_max_size,omitempty"`
	HasParameterValues bool             `protobuf:"varint,4,opt,name=has_parameter_values,json=hasParameterValues" json:"has_parameter_values,omitempty"`
}

func (m *ExecuteRequest) Reset()                    { *m = ExecuteRequest{} }
func (m *ExecuteRequest) String() string            { return proto.CompactTextString(m) }
func (*ExecuteRequest) ProtoMessage()               {}
func (*ExecuteRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{15} }

func (m *ExecuteRequest) GetStatementHandle() *StatementHandle {
	if m != nil {
		return m.StatementHandle
	}
	return nil
}

func (m *ExecuteRequest) GetParameterValues() []*TypedValue {
	if m != nil {
		return m.ParameterValues
	}
	return nil
}

type SyncResultsRequest struct {
	ConnectionId string      `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
	StatementId  uint32      `protobuf:"varint,2,opt,name=statement_id,json=statementId" json:"statement_id,omitempty"`
	State        *QueryState `protobuf:"bytes,3,opt,name=state" json:"state,omitempty"`
	Offset       uint64      `protobuf:"varint,4,opt,name=offset" json:"offset,omitempty"`
}

func (m *SyncResultsRequest) Reset()                    { *m = SyncResultsRequest{} }
func (m *SyncResultsRequest) String() string            { return proto.CompactTextString(m) }
func (*SyncResultsRequest) ProtoMessage()               {}
func (*SyncResultsRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{16} }

func (m *SyncResultsRequest) GetState() *QueryState {
	if m != nil {
		return m.State
	}
	return nil
}

// Request to invoke a commit on a Connection
type CommitRequest struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *CommitRequest) Reset()                    { *m = CommitRequest{} }
func (m *CommitRequest) String() string            { return proto.CompactTextString(m) }
func (*CommitRequest) ProtoMessage()               {}
func (*CommitRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{17} }

// Request to invoke rollback on a Connection
type RollbackRequest struct {
	ConnectionId string `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
}

func (m *RollbackRequest) Reset()                    { *m = RollbackRequest{} }
func (m *RollbackRequest) String() string            { return proto.CompactTextString(m) }
func (*RollbackRequest) ProtoMessage()               {}
func (*RollbackRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{18} }

// Request to prepare and execute a collection of sql statements.
type PrepareAndExecuteBatchRequest struct {
	ConnectionId string   `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
	StatementId  uint32   `protobuf:"varint,2,opt,name=statement_id,json=statementId" json:"statement_id,omitempty"`
	SqlCommands  []string `protobuf:"bytes,3,rep,name=sql_commands,json=sqlCommands" json:"sql_commands,omitempty"`
}

func (m *PrepareAndExecuteBatchRequest) Reset()                    { *m = PrepareAndExecuteBatchRequest{} }
func (m *PrepareAndExecuteBatchRequest) String() string            { return proto.CompactTextString(m) }
func (*PrepareAndExecuteBatchRequest) ProtoMessage()               {}
func (*PrepareAndExecuteBatchRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{19} }

// Each command is a list of TypedValues
type UpdateBatch struct {
	ParameterValues []*TypedValue `protobuf:"bytes,1,rep,name=parameter_values,json=parameterValues" json:"parameter_values,omitempty"`
}

func (m *UpdateBatch) Reset()                    { *m = UpdateBatch{} }
func (m *UpdateBatch) String() string            { return proto.CompactTextString(m) }
func (*UpdateBatch) ProtoMessage()               {}
func (*UpdateBatch) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{20} }

func (m *UpdateBatch) GetParameterValues() []*TypedValue {
	if m != nil {
		return m.ParameterValues
	}
	return nil
}

type ExecuteBatchRequest struct {
	ConnectionId string         `protobuf:"bytes,1,opt,name=connection_id,json=connectionId" json:"connection_id,omitempty"`
	StatementId  uint32         `protobuf:"varint,2,opt,name=statement_id,json=statementId" json:"statement_id,omitempty"`
	Updates      []*UpdateBatch `protobuf:"bytes,3,rep,name=updates" json:"updates,omitempty"`
}

func (m *ExecuteBatchRequest) Reset()                    { *m = ExecuteBatchRequest{} }
func (m *ExecuteBatchRequest) String() string            { return proto.CompactTextString(m) }
func (*ExecuteBatchRequest) ProtoMessage()               {}
func (*ExecuteBatchRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{21} }

func (m *ExecuteBatchRequest) GetUpdates() []*UpdateBatch {
	if m != nil {
		return m.Updates
	}
	return nil
}

func init() {
	proto.RegisterType((*CatalogsRequest)(nil), "CatalogsRequest")
	proto.RegisterType((*DatabasePropertyRequest)(nil), "DatabasePropertyRequest")
	proto.RegisterType((*SchemasRequest)(nil), "SchemasRequest")
	proto.RegisterType((*TablesRequest)(nil), "TablesRequest")
	proto.RegisterType((*TableTypesRequest)(nil), "TableTypesRequest")
	proto.RegisterType((*ColumnsRequest)(nil), "ColumnsRequest")
	proto.RegisterType((*TypeInfoRequest)(nil), "TypeInfoRequest")
	proto.RegisterType((*PrepareAndExecuteRequest)(nil), "PrepareAndExecuteRequest")
	proto.RegisterType((*PrepareRequest)(nil), "PrepareRequest")
	proto.RegisterType((*FetchRequest)(nil), "FetchRequest")
	proto.RegisterType((*CreateStatementRequest)(nil), "CreateStatementRequest")
	proto.RegisterType((*CloseStatementRequest)(nil), "CloseStatementRequest")
	proto.RegisterType((*OpenConnectionRequest)(nil), "OpenConnectionRequest")
	proto.RegisterType((*CloseConnectionRequest)(nil), "CloseConnectionRequest")
	proto.RegisterType((*ConnectionSyncRequest)(nil), "ConnectionSyncRequest")
	proto.RegisterType((*ExecuteRequest)(nil), "ExecuteRequest")
	proto.RegisterType((*SyncResultsRequest)(nil), "SyncResultsRequest")
	proto.RegisterType((*CommitRequest)(nil), "CommitRequest")
	proto.RegisterType((*RollbackRequest)(nil), "RollbackRequest")
	proto.RegisterType((*PrepareAndExecuteBatchRequest)(nil), "PrepareAndExecuteBatchRequest")
	proto.RegisterType((*UpdateBatch)(nil), "UpdateBatch")
	proto.RegisterType((*ExecuteBatchRequest)(nil), "ExecuteBatchRequest")
}

var fileDescriptor1 = []byte{
	// 879 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xc4, 0x56, 0xdd, 0x4e, 0xdb, 0x48,
	0x14, 0x96, 0x49, 0x02, 0xe4, 0x38, 0x09, 0xc1, 0xfc, 0x6c, 0xc4, 0x6a, 0xa5, 0xe0, 0xfd, 0x11,
	0x17, 0xbb, 0xde, 0x55, 0x16, 0xb1, 0x08, 0x69, 0x57, 0x5a, 0xbc, 0xa0, 0x45, 0x2a, 0x6d, 0xea,
	0xd0, 0xde, 0x5a, 0x13, 0x67, 0x02, 0x16, 0x8e, 0xc7, 0x78, 0xc6, 0x40, 0x7a, 0xdf, 0x9b, 0xde,
	0xf5, 0xaa, 0x0f, 0xd1, 0xc7, 0xe8, 0x2b, 0xf4, 0x11, 0x2a, 0xf5, 0xb2, 0xaf, 0xd0, 0x99, 0xb1,
	0x63, 0x92, 0x38, 0x55, 0x6b, 0x04, 0xed, 0x55, 0x72, 0xfe, 0xcf, 0xf7, 0xf9, 0x9c, 0x99, 0x81,
	0x5a, 0x88, 0x2f, 0x22, 0x4c, 0x19, 0x35, 0x82, 0x90, 0x30, 0xb2, 0x51, 0x71, 0xc8, 0x60, 0x40,
	0xfc, 0x58, 0xd2, 0x77, 0x60, 0xc9, 0x44, 0x0c, 0x79, 0xe4, 0x94, 0x5a, 0xb1, 0x9f, 0xf6, 0x23,
	0x54, 0x1d, 0xe2, 0xfb, 0xd8, 0x61, 0x2e, 0xf1, 0x6d, 0xb7, 0xd7, 0x50, 0x9a, 0xca, 0x56, 0xd9,
	0xaa, 0xdc, 0x28, 0x8f, 0x7a, 0xfa, 0x3f, 0xf0, 0xdd, 0x7f, 0x3c, 0xae, 0x8b, 0x28, 0x6e, 0x87,
	0x24, 0xc0, 0x21, 0x1b, 0xe6, 0x8a, 0xbf, 0x84, 0x5a, 0xc7, 0x39, 0xc3, 0x03, 0x94, 0x96, 0x6d,
	0xc0, 0x82, 0x13, 0x77, 0x92, 0x04, 0x8c, 0x44, 0xed, 0x67, 0xa8, 0x51, 0xe9, 0x6b, 0x07, 0x88,
	0x31, 0x1c, 0xfa, 0x8d, 0x39, 0xe9, 0x50, 0x8d, 0xb5, 0xed, 0x58, 0x99, 0xad, 0x5b, 0x98, 0x51,
	0xf7, 0x9d, 0x02, 0xd5, 0x13, 0xd4, 0xf5, 0xf0, 0xdd, 0xd5, 0xfd, 0x15, 0x34, 0x26, 0x32, 0xda,
	0x3e, 0x1a, 0xe0, 0xd4, 0x35, 0x2e, 0x5e, 0x97, 0x96, 0x87, 0xdc, 0x30, 0xf2, 0xfe, 0x1e, 0xca,
	0x6c, 0x18, 0x60, 0xdb, 0x73, 0x29, 0x6b, 0x14, 0x9b, 0x05, 0xee, 0xb4, 0x28, 0x14, 0x0f, 0xb8,
	0xac, 0xe9, 0x50, 0x3d, 0x43, 0xd4, 0xbe, 0x71, 0x98, 0xe7, 0x59, 0x16, 0x2d, 0x95, 0x2b, 0x4f,
	0x46, 0x3e, 0x19, 0x98, 0x0b, 0x33, 0x60, 0xee, 0xc2, 0xb2, 0x44, 0x29, 0xa2, 0xf2, 0x7d, 0xd8,
	0xb7, 0x0a, 0xd4, 0x4c, 0xe2, 0x45, 0x03, 0xff, 0x5b, 0x31, 0x64, 0xc0, 0x8a, 0x23, 0x1b, 0x98,
	0x74, 0x2f, 0x4a, 0xf7, 0xe5, 0xd8, 0x34, 0xee, 0x9f, 0x81, 0x55, 0x9a, 0x01, 0x8b, 0xcf, 0xb9,
	0xe0, 0xe2, 0xc8, 0xef, 0x93, 0x5c, 0x74, 0x7c, 0x50, 0xa0, 0xd1, 0x0e, 0x71, 0x80, 0x42, 0xfc,
	0xaf, 0xdf, 0x3b, 0xb8, 0xc6, 0x4e, 0xc4, 0x70, 0x9e, 0x0c, 0x5a, 0x1d, 0x0a, 0xf4, 0xc2, 0x4b,
	0x88, 0x11, 0x7f, 0xc5, 0x57, 0x1e, 0xa0, 0x6b, 0x3b, 0x24, 0x57, 0xb6, 0x43, 0x22, 0x9f, 0x49,
	0x26, 0x8a, 0x96, 0xca, 0x95, 0x16, 0xb9, 0x32, 0x85, 0x4a, 0xdb, 0x84, 0x0a, 0x65, 0x88, 0xe1,
	0x01, 0xf6, 0x99, 0xc8, 0x2c, 0xd0, 0x57, 0x2d, 0x35, 0xd5, 0xf1, 0xc4, 0x3f, 0x41, 0x2d, 0x49,
	0xc3, 0x27, 0x86, 0xf0, 0x0f, 0x22, 0x81, 0x17, 0xac, 0x4a, 0x9c, 0x87, 0x9e, 0x08, 0x9d, 0xf6,
	0x3b, 0xac, 0xf6, 0xdd, 0x90, 0x32, 0xbb, 0x1f, 0x0a, 0x36, 0x45, 0x04, 0x75, 0x9f, 0x61, 0x39,
	0x59, 0x25, 0x6b, 0x59, 0xda, 0x0e, 0x85, 0xe9, 0x18, 0x5d, 0x77, 0xb8, 0x41, 0x7f, 0xc9, 0x07,
	0x20, 0x41, 0xfc, 0x15, 0x70, 0x66, 0x41, 0x14, 0xb3, 0x20, 0xf4, 0x37, 0x0a, 0x54, 0x0e, 0x31,
	0x73, 0xce, 0x72, 0x75, 0x34, 0xcd, 0xe1, 0x5c, 0x96, 0xc3, 0x75, 0x98, 0x27, 0xfd, 0x3e, 0xc5,
	0xa3, 0xde, 0x12, 0x49, 0xfb, 0x0d, 0x56, 0xfa, 0xa2, 0x9e, 0x3d, 0x09, 0x20, 0xfe, 0x0a, 0x75,
	0x69, 0x3a, 0x9e, 0x44, 0x31, 0x45, 0x6f, 0x49, 0xd2, 0x5b, 0xe9, 0x8f, 0x33, 0xfb, 0x37, 0xac,
	0x9b, 0x21, 0xe6, 0xc5, 0x3b, 0xa3, 0x0e, 0x72, 0x8d, 0xa2, 0x0d, 0x6b, 0xa6, 0x47, 0xe8, 0xed,
	0xa2, 0xbf, 0x80, 0x0c, 0xfd, 0xb5, 0x02, 0x6b, 0x8f, 0x02, 0xec, 0x9b, 0x69, 0x5c, 0xae, 0x0a,
	0xdb, 0x50, 0x74, 0xf9, 0x7a, 0xf1, 0xcc, 0x85, 0x2d, 0xb5, 0xd5, 0x34, 0x66, 0xa6, 0x32, 0xc4,
	0x06, 0x1e, 0xf8, 0x2c, 0x1c, 0x5a, 0xd2, 0x7b, 0xe3, 0x2f, 0x28, 0xa7, 0x2a, 0x31, 0x43, 0xe7,
	0x78, 0x98, 0x64, 0x17, 0x7f, 0xb5, 0x55, 0x28, 0x5d, 0x22, 0x2f, 0xc2, 0xc9, 0x5c, 0xc5, 0xc2,
	0xde, 0xdc, 0xae, 0x22, 0xd9, 0x14, 0x74, 0xdc, 0xae, 0x5b, 0x3d, 0xe4, 0x6c, 0xa6, 0x72, 0x67,
	0xe8, 0x3b, 0x39, 0xb1, 0x82, 0x90, 0x6d, 0x7e, 0x89, 0x06, 0x54, 0xf6, 0xa6, 0xb6, 0xd6, 0x8c,
	0x9b, 0x84, 0xc9, 0x9d, 0xe8, 0xf2, 0xc3, 0xb7, 0x2c, 0x1c, 0x85, 0x4c, 0xf5, 0xf7, 0x7c, 0xb5,
	0xa6, 0x8e, 0x90, 0x3d, 0x58, 0x4a, 0x3f, 0xc1, 0xff, 0xc8, 0xef, 0x79, 0x58, 0xd6, 0x53, 0x5b,
	0x75, 0xa3, 0x33, 0xa9, 0xb7, 0xa6, 0x1d, 0xb5, 0x1d, 0xa8, 0xf3, 0x2d, 0xe5, 0x03, 0xc6, 0x8f,
	0x41, 0x5b, 0x12, 0x43, 0x13, 0xf2, 0x55, 0x43, 0x1c, 0x76, 0xbd, 0xa7, 0x42, 0x67, 0x2d, 0xa5,
	0x4e, 0x52, 0xa6, 0x9f, 0x3c, 0x12, 0xe2, 0x15, 0xc8, 0x1e, 0x09, 0xda, 0x1f, 0xb0, 0x2a, 0xae,
	0xa5, 0x4c, 0xb1, 0xa2, 0xbc, 0x9d, 0x34, 0x6e, 0x6b, 0x4f, 0x96, 0xd0, 0x5f, 0x29, 0xa0, 0xc5,
	0xa4, 0xd2, 0xc8, 0x63, 0xf4, 0xae, 0xd7, 0x76, 0x13, 0x4a, 0x52, 0x94, 0x2d, 0x0b, 0xb8, 0x8f,
	0x23, 0x1c, 0x0e, 0x25, 0x61, 0x56, 0x6c, 0x19, 0xdb, 0xec, 0xe2, 0xf8, 0x66, 0xeb, 0xdb, 0x50,
	0x35, 0xf9, 0x03, 0xc8, 0xcd, 0xb7, 0x7b, 0xfc, 0xfa, 0xb0, 0x88, 0xe7, 0x75, 0x91, 0x73, 0x9e,
	0x2b, 0xee, 0x85, 0x02, 0x3f, 0x64, 0xae, 0x8f, 0x7d, 0x74, 0x0f, 0x27, 0x99, 0x70, 0xb9, 0xf0,
	0x6c, 0xf1, 0xb8, 0xe3, 0xd3, 0x41, 0x39, 0x33, 0xe2, 0x69, 0xa1, 0x72, 0x9d, 0x99, 0xa8, 0xf4,
	0x03, 0x50, 0x9f, 0x04, 0x3d, 0x94, 0x34, 0x30, 0x73, 0x7c, 0x94, 0xcf, 0x8f, 0x8f, 0xfe, 0x5c,
	0x81, 0x95, 0xfb, 0x44, 0xf2, 0x0b, 0x2c, 0x44, 0xb2, 0xcd, 0x18, 0x84, 0xda, 0xaa, 0x18, 0x63,
	0x6d, 0x5b, 0x23, 0xe3, 0xbe, 0x0e, 0x4d, 0x12, 0x9e, 0x1a, 0x28, 0x40, 0xfc, 0xb1, 0x61, 0x38,
	0xc8, 0x73, 0x5c, 0x86, 0x0d, 0x74, 0x89, 0x98, 0xeb, 0xa0, 0xf8, 0x79, 0xdb, 0x9d, 0x97, 0x3f,
	0x7f, 0x7e, 0x0c, 0x00, 0x00, 0xff, 0xff, 0x4f, 0x49, 0xad, 0x3d, 0x05, 0x0b, 0x00, 0x00,
}
