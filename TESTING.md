# Unit Tests for go-blar

A comprehensive test suite covering all major components of the library.

## Test Files

### `goblar/app_test.go` (8 tests)
Tests for the main App type and configuration:
- `TestNew()` - Creates app with default config
- `TestNewWithDB()` - Configures database
- `TestNewWithAddress()` - Sets custom server address
- `TestNewWithMultipleOptions()` - Applies multiple options
- `TestRegisterWithoutDB()` - Validates DB requirement
- `TestRegisterValid()` - Successfully registers a model
- `TestRegisterMultiple()` - Registers multiple models

### `goblar/options_test.go` (6 tests)
Tests for configuration options:
- `TestWithDB()` - Database option
- `TestWithAddress()` - Address option
- `TestWithMiddleware()` - Middleware option
- `TestConfig_Apply()` - Option application
- `TestNewConfig_Defaults()` - Default configuration values
- `TestWithMiddleware()` - Multiple middleware stacking

### `internal/meta/parse_test.go` (10 tests)
Tests for metadata parsing and struct reflection:
- `TestParseBasicStruct()` - Parse simple struct
- `TestParseWithPrimaryKey()` - Detect primary keys
- `TestParseWithPointerType()` - Handle pointer types
- `TestCaching()` - Metadata caching
- `TestParseTableName()` - Table name generation
- `TestToSnakeCase()` - Convert CamelCase to snake_case (handles acronyms)
- `TestParseFieldTags()` - Parse go-blar struct tags
- `TestGetFieldByName()` - Retrieve field metadata by name

### `internal/hooks/hooks_test.go` (11 tests)
Tests for lifecycle hook execution:
- `TestCallBeforeCreate()` - Before create hook invocation
- `TestCallBeforeCreateWithError()` - Error handling
- `TestCallBeforeCreateNoImplementation()` - Optional hooks
- `TestCallAfterCreate()` - After create hook
- `TestCallBeforeUpdate()` - Before update hook
- `TestCallAfterUpdate()` - After update hook
- `TestCallBeforeDelete()` - Before delete hook
- `TestCallAfterDelete()` - After delete hook
- `TestMultipleHooks()` - Sequential hook execution

### `internal/repo/repository_test.go` (11 tests)
Tests for generic CRUD repository:
- `TestNewRepository()` - Create repository instance
- `TestCreate()` - Insert entity
- `TestGetByID()` - Retrieve by ID
- `TestGetAll()` - List all entities
- `TestUpdate()` - Modify entity
- `TestDelete()` - Remove entity
- `TestCount()` - Count entities
- `TestCreateWithoutDB()` - Validate database requirement

## Test Coverage

- **Public API**: App creation, configuration, registration
- **Options Pattern**: Configuration composition
- **Metadata Parsing**: Struct reflection, tag parsing, caching
- **Lifecycle Hooks**: All 6 hook types with success and error cases
- **Repository**: CRUD operations on in-memory SQLite database
- **Edge Cases**: Missing DB, nil checks, optional hooks

## Running Tests

Run all tests:
```bash
go test ./...
```

Run specific test file:
```bash
go test ./goblar -v
go test ./internal/meta -v
go test ./internal/hooks -v
go test ./internal/repo -v
```

Run specific test:
```bash
go test -run TestParseBasicStruct ./internal/meta
```

## Test Utilities

### Mock Entities
- `TestEntity` - Simple test struct
- `TestUser` - Repository test entity
- `MockEntityWithHooks` - Implements all hook interfaces for testing

### Setup Functions
- `setupTestDB()` - Creates in-memory SQLite database for testing
- `ClearRegistry()` - Clears metadata cache between tests

## Test Results

✅ **37 tests passing** - 100% pass rate

## Key Testing Patterns

1. **Isolation**: Each test clears the metadata registry to avoid cross-test contamination
2. **In-Memory DB**: SQLite `:memory:` database for fast, isolated testing
3. **Mock Implementations**: Fake hook implementations to verify invocation
4. **Table-Driven Tests**: Multiple test cases in single test function (e.g., `TestToSnakeCase`)
5. **Error Cases**: Tests both success and failure paths
